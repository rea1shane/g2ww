package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 注释的为数据类型对不上的结构
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

var sentCount = 0

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

// 记录发送次数
func GetSendCount(c *gin.Context) {
	fmt.Println("\n----------------------------------------------------------------------\n")

	_, _ = c.Writer.WriteString("G2WW Server is running! \nParsed & forwarded " + strconv.Itoa(sentCount) + " messages to WeChat Work!")

	fmt.Println("sentCount:", sentCount)
	fmt.Println()

	return
}

// 发送消息
func SendMsg(c *gin.Context) {
	fmt.Println("\n----------------------------------------------------------------------\n")

	h := &Hook{}
	data, _ := ioutil.ReadAll(c.Request.Body)

	if err := json.Unmarshal(data, &h); err != nil {
		fmt.Println("err:", err.Error())
		_, _ = c.Writer.WriteString("Error on JSON format")
		return
	}

	// Send to WeChat Work
	url := Url + c.Query("key")

	// 消息体
	var msgType, msgStr string
	if c.Query("type") == "news" {
		msgType = "news"
		msgStr = MsgNews(h)
	} else {
		msgType = "markdown"
		msgStr = MsgMarkdown(h)
	}

	jsonStr := []byte(msgStr)
	// 发送http请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", err.Error())
		_, _ = c.Writer.WriteString("Error sending to WeChat Work API")
		return
	}
	defer resp.Body.Close()
	sentCount++

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

// 发送消息类型 news
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
		  }
		`, h.Title, h.Message, h.RuleUrl, h.ImageUrl)
}

// 发送消息类型
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
