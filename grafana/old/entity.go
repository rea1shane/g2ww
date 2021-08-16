package old

import (
	"fmt"
	"g2ww/grafana/general"
	"g2ww/ww"
)

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

func (h Hook) MsgNews() string {
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

func (h Hook) MsgMarkdown() string {
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

func (h Hook) PrintAlertLog() {
	fmt.Println("Title       :", h.Title)
	fmt.Println("RuleName    :", h.RuleName)
	fmt.Println("State       :", h.State)
	fmt.Println("Message     :", h.Message)
	fmt.Println("RuleUrl     :", h.RuleUrl)
	fmt.Println("ImageUrl    :", h.ImageUrl)
}
