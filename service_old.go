package main

import (
	"fmt"
	"g2ww/grafana/general"
	"g2ww/grafana/old"
	"g2ww/ww"
)

// MsgNewsOld 发送消息类型 news
func MsgNewsOld(h *old.Hook) string {
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

// MsgMarkdownOld 发送消息类型
func MsgMarkdownOld(h *old.Hook) string {
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

// PrintAlertLogOld 打印告警信息
func PrintAlertLogOld(h *old.Hook) {
	fmt.Println("Title    :", h.Title)
	fmt.Println("RuleName :", h.RuleName)
	fmt.Println("State    :", h.State)
	fmt.Println("Message  :", h.Message)
	fmt.Println("RuleUrl  :", h.RuleUrl)
	fmt.Println("ImageUrl :", h.ImageUrl)
}
