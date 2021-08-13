package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"g2ww/common"
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

func SendMsgOld(c *gin.Context) {
	common.PrintCutOffRule()
	status := common.InternalError

	// 将 webhook 数据装载为 struct 对象
	h := old.Hook{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &h); err != nil {
		status = GrafanaWebhookUnmarshalJsonError(c, err)
	} else {
		// 组装 url
		url := ww.WechatWorkBotWebhookUrl + c.Query("key")

		// 消息体
		var msgType, msgStr string
		if c.Query("type") == "news" {
			msgType = "news"
			msgStr = MsgNewsOld(&h)
		} else {
			msgType = "markdown"
			msgStr = MsgMarkdownOld(&h)
		}

		// 发送请求
		jsonStr := []byte(msgStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			status = ClientCallAPIError(c, err)
		} else {
			status = ww.CheckWechatWorkResponse(resp)

			// 关闭 response body
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			// 日志记录
			fmt.Println("MsgType  :", msgType)
			PrintAlertLogOld(&h)
			fmt.Println()
		}
	}
	common.CheckStatus(status, &counter)
	fmt.Println("Status   :", status)
	fmt.Println()
	return
}

func SendMsgNgalert(c *gin.Context) {
	common.PrintCutOffRule()
	status := common.InternalError

	// 将 webhook 数据装载为 struct 对象
	h := ngalert.Hook{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &h); err != nil {
		status = GrafanaWebhookUnmarshalJsonError(c, err)
	} else {
		// 组装 url
		url := ww.WechatWorkBotWebhookUrl + c.Query("key")

		// 消息体
		var msgType, msgStr string
		if c.Query("type") == "news" {
			msgType = "news"
			msgStr = MsgNewsNgalert(&h)
		} else {
			msgType = "markdown"
			msgStr = MsgMarkdownNgalert(&h)
		}

		// 发送请求
		jsonStr := []byte(msgStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			status = ClientCallAPIError(c, err)
		} else {
			status = ww.CheckWechatWorkResponse(resp)

			// 关闭 response body
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			// 日志记录
			fmt.Println("MsgType  :", msgType)
			PrintAlertLogNgalert(&h)
			fmt.Println()
		}
	}
	common.CheckStatus(status, &counter)
	fmt.Println("Status   :", status)
	fmt.Println()
	return
}

func DealError(c *gin.Context, err error, errMsg string) {
	fmt.Println(errMsg)
	fmt.Println(err.Error())
	_, _ = c.Writer.WriteString(errMsg)
}

func GrafanaWebhookUnmarshalJsonError(c *gin.Context, err error) int {
	errMsg := `[ERROR] JSON Unmarshal failure`
	DealError(c, err, errMsg)
	return common.GrafanaWebhookUnmarshalJsonError
}

func ClientCallAPIError(c *gin.Context, err error) int {
	errMsg := `[ERROR] Client call API failure`
	DealError(c, err, errMsg)
	return common.ClientCallAPIError
}
