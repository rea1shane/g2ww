package old

import (
	"bytes"
	"encoding/json"
	"fmt"
	"g2ww/common"
	"g2ww/grafana/general"
	"g2ww/ww"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

var counter = common.Counter{
	SentSuccessCount: 0,
	SentFailureCount: 0,
}

func GetSendCount(c *gin.Context) {
	general.GetSendCount(c, counter)
}

func SendMsg(c *gin.Context) {
	common.PrintCutOffRule()
	status := common.InternalError

	// 将 webhook 数据装载为 struct 对象
	h := Hook{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &h); err != nil {
		status = general.GrafanaWebhookUnmarshalJsonError(c, err)
	} else {
		// 组装 url
		url := ww.WechatWorkBotWebhookUrl + c.Query("key")

		// 消息体
		var msgType, msgStr string
		if c.Query("type") == "news" {
			msgType = "news"
			msgStr = MsgNews(&h)
		} else {
			msgType = "markdown"
			msgStr = MsgMarkdown(&h)
		}

		// 发送请求
		jsonStr := []byte(msgStr)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			status = general.ClientCallAPIError(c, err)
		} else {
			status = ww.CheckWechatWorkResponse(resp)

			// 关闭 response body
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			// 日志记录
			fmt.Println("MsgType  :", msgType)
			PrintAlertLog(&h)
			fmt.Println()
		}
	}
	common.CheckStatus(status, &counter)
	fmt.Println("Status   :", status)
	fmt.Println()
	return
}

// MsgNews 发送消息类型 news
func MsgNews(h *Hook) string {
	return fmt.Sprintf(`
		{
			"msgtype": "news",
			"news": {
			  	"articles": [
					{
				  		"title": "%s",
				  		"description": "%s",
				  		"url": "%s",
				  		"picurl": "%s"
					}
			  	]
			}
		}`, h.Title, h.Message, h.RuleUrl, h.ImageUrl)
}

// MsgMarkdown 发送消息类型
func MsgMarkdown(h *Hook) string {
	var color, stateMsg string
	if h.State == general.OK {
		color = ww.WechatWorkColorGreen
		stateMsg = general.OKMsg
	} else if h.State == general.Alerting {
		color = ww.WechatWorkColorRed
		stateMsg = general.AlertingMsg
	} else {
		color = ww.WechatWorkColorRed
		stateMsg = h.State
	}
	return fmt.Sprintf(`
		{
		   	"msgtype": "markdown",
		   	"markdown": {
				"content": "# <font color=\"%s\">[%s]</font> %s\n\n\n**Message**: <font color=\"comment\">%s</font>\n**Grafana**: [%s](%s)\n**AlertImage**: [%s](%s)"
		   	}
		}`, color, stateMsg, h.RuleName, h.Message, h.RuleUrl, h.RuleUrl, h.ImageUrl, h.ImageUrl)
}

// PrintAlertLog 打印告警信息
func PrintAlertLog(h *Hook) {
	fmt.Println("Title    :", h.Title)
	fmt.Println("RuleName :", h.RuleName)
	fmt.Println("State    :", h.State)
	fmt.Println("Message  :", h.Message)
	fmt.Println("RuleUrl  :", h.RuleUrl)
	fmt.Println("ImageUrl :", h.ImageUrl)
}
