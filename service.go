package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"g2ww/common"
	"g2ww/grafana/general"
	"g2ww/grafana/ngalert"
	"g2ww/grafana/old"
	"g2ww/ww"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var counter = common.Counter{
	SentSuccessCount: 0,
	SentFailureCount: 0,
}

func GetSendCount(c *gin.Context) {
	common.PrintCutOffRule()

	_, _ = c.Writer.WriteString("G2WW Server is running! \nParsed & forwarded " + strconv.Itoa(counter.SentSuccessCount) + " messages to WeChat Work! \nParsed & forwarded Failure " + strconv.Itoa(counter.SentFailureCount) + " messages!")

	fmt.Println("Sent Success Count:", counter.SentSuccessCount)
	fmt.Println("Sent Failure Count:", counter.SentFailureCount)
	fmt.Println()

	return
}

func SendMsg(c *gin.Context, h general.Hook) {
	common.PrintCutOffRule()
	status := common.InternalError

	// 将 webhook 数据装载为 struct 对象
	data, _ := ioutil.ReadAll(c.Request.Body)
	// 用于 debug 输出消息体, 测试时打开
	// fmt.Println(string(data))
	if err := json.Unmarshal(data, &h); err != nil {
		fmt.Println(err.Error())
		fmt.Println()
		status = common.GrafanaWebhookUnmarshalJsonError
		_, _ = c.Writer.WriteString(status.String())
	} else {
		// 组装 url
		// 如果想要进行 debug, 需要在链接末尾  + "&debug=1", 然后在 https://open.work.weixin.qq.com/devtool/query 进行 hint 查询
		url := ww.WechatWorkBotWebhookUrl + c.Query("key")

		// 消息体
		var msgType, msgStr string
		if c.Query("type") == "news" {
			msgType = "news"
			msgStr = h.MsgNews()
		} else {
			msgType = "markdown"
			msgStr = h.MsgMarkdown()
		}

		// 发送请求
		jsonStr := []byte(msgStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println()
			status = common.ClientCallAPIError
			_, _ = c.Writer.WriteString(status.String())
		} else {
			status = ww.CheckWechatWorkResponse(resp)

			// 关闭 response body
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			// 日志记录
			fmt.Println("MsgType  :", msgType)
			h.PrintAlertLog()
			fmt.Println()
		}
	}
	common.CheckStatus(status, &counter)
	fmt.Println()
	return
}

func SendMsgOld(c *gin.Context) {
	h := old.Hook{}
	SendMsg(c, &h)
}

func SendMsgNgalert(c *gin.Context) {
	h := ngalert.Hook{}
	SendMsg(c, &h)
}
