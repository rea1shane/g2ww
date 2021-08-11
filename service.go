package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// TODO 筛选 Hook 中有用的字段
// TODO 更新 MsgMarkdown 消息结构

// Hook 注释的为数据类型对不上的结构
type Hook struct {
	ImageUrl string `json:"imageUrl"`
	Message  string `json:"message"`
	RuleName string `json:"ruleName"`
	RuleUrl  string `json:"ruleUrl"`
	State    string `json:"state"`
	Title    string `json:"title"`
	// OrgId       string `json:"orgId"`
	// PanelId     string `json:"panelId"`
	// RuleId      string `json:"ruleId,omitempty"`
	// Tags        string `json:"tags"`
	// DashboardId string `json:"dashboardId"`
	// EvalMatches string `json:"evalMatches,omitempty"`
}

type WechatWorkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var sentSuccessCount = 0
var sentFailureCount = 0

const (
	Url         = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	OK          = "ok"
	Alerting    = "alerting"
	OKMsg       = "OK"
	AlertingMsg = "Alerting"
	ColorGreen  = "info"
	ColorGray   = "comment"
	ColorRed    = "warning"
)

// GetSendCount 记录发送次数
func GetSendCount(c *gin.Context) {
	PrintCutOffRule()

	_, _ = c.Writer.WriteString("G2WW Server is running! \nParsed & forwarded " + strconv.Itoa(sentSuccessCount) + " messages to WeChat Work! \nParsed & forwarded Failure " + strconv.Itoa(sentFailureCount) + " messages!")

	fmt.Println("Sent Success Count:", sentSuccessCount)
	fmt.Println("Sent Failure Count:", sentFailureCount)
	fmt.Println()

	return
}

// SendMsg 发送消息
func SendMsg(c *gin.Context) {
	PrintCutOffRule()

	h := Hook{}
	data, _ := ioutil.ReadAll(c.Request.Body)

	if err := json.Unmarshal(data, &h); err != nil {
		fmt.Println("[ERROR]", err.Error())
		_, _ = c.Writer.WriteString("Error on JSON format")
		sentFailureCount++
		return
	}

	// Send to WeChat Work
	url := Url + c.Query("key")

	// 消息体
	var msgType, msgStr string
	if c.Query("type") == "news" {
		msgType = "news"
		msgStr = MsgNews(&h)
	} else {
		msgType = "markdown"
		msgStr = MsgMarkdown(&h)
	}

	// 发送http请求
	jsonStr := []byte(msgStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		_, _ = c.Writer.WriteString("Error sending to WeChat Work API")
		sentFailureCount++
		return
	}

	// 检测是否发送失败
	// 企业微信有 20条/min 的发送速率限制
	r := WechatWorkResponse{}
	buffer := new(bytes.Buffer)
	_, _ = buffer.ReadFrom(resp.Body)
	_ = json.Unmarshal(buffer.Bytes(), &r)
	if r.ErrCode != 0 {
		fmt.Println("[ERROR] Error sending to WeChat Work API")
		fmt.Printf("ErrorCode: [%v]", r.ErrCode)
		fmt.Println()
		fmt.Printf("ErrorMsg : [%v]", r.ErrMsg)
		fmt.Println()
		fmt.Println()
		sentFailureCount++
	} else if r.ErrMsg != "ok" {
		fmt.Println("[ERROR] Parse response body failure")
		fmt.Println()
		sentFailureCount++
	} else {
		sentSuccessCount++
	}

	// 关闭 response body
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 日志记录
	fmt.Println("MsgType  :", msgType)
	fmt.Println("Title    :", h.Title)
	fmt.Println("RuleName :", h.RuleName)
	fmt.Println("State    :", h.State)
	fmt.Println("Message  :", h.Message)
	fmt.Println("RuleUrl  :", h.RuleUrl)
	fmt.Println("ImageUrl :", h.ImageUrl)
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
	if h.State == OK {
		color = ColorGreen
		stateMsg = OKMsg
	} else if h.State == Alerting {
		color = ColorRed
		stateMsg = AlertingMsg
	} else {
		color = ColorRed
		stateMsg = h.State
	}
	return fmt.Sprintf(`
		{
		   	"msgtype": "markdown",
		   	"markdown": {
				"content": "<font color=\"%s\">[%s]</font> <font>%s</font>\r\n<font color=\"comment\">%s\r\n[点击查看详情](%s)![](%s)</font>"
		   	}
		}`, color, stateMsg, h.RuleName, h.Message, h.RuleUrl, h.ImageUrl)
}

func PrintCutOffRule() {
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println()
}
