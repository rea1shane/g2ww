package ngalert

import (
	"fmt"
	"g2ww/common"
	"g2ww/ww"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TimeLayout = "2006-01-02 Mon 15:04:05"

	RESOLVED = "resolved"
	FIRING   = "firing"
)

type Labels struct {
	Alertname string `json:"alertname"`
}

type Annotations struct {
}

type Alert struct {
	Status       string      `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
	SilenceURL   string      `json:"silenceURL"`
	DashboardURL string      `json:"dashboardURL"`
	PanelURL     string      `json:"panelURL"`
	ValueString  string      `json:"valueString"`
}

type GroupLabels struct {
}

type CommonLabels struct {
	Alertname string `json:"alertname"`
}

type CommonAnnotations struct {
}

type Hook struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	GroupLabels       GroupLabels       `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Title             string            `json:"title"`
	State             string            `json:"state"`
	Message           string            `json:"message"`
}

// MsgNews TODO
func (h Hook) MsgNews() string {
	return ""
}

func (h Hook) MsgMarkdown() string {
	var color string
	if h.Status == RESOLVED {
		color = ww.WechatWorkColorGreen
	} else {
		color = ww.WechatWorkColorRed
	}
	return fmt.Sprintf(`
		{
		   	"msgtype": "markdown",
		   	"markdown": {
				"content": "# <font color=\"%s\">%s</font>%s"
		   	}
		}`, color, h.GetTitleStatus(), h.GetAlertDetailList())
}

// PrintAlertLog TODO
func (h Hook) PrintAlertLog() {

}

func (h Hook) GetTitleStatus() string {
	return "[" + strings.ToUpper(h.Status) + ":" + strconv.Itoa(len(h.Alerts)) + "]"
}

// GetTitleAlerts 暂时去掉 显示的感觉些许多余 本应跟在 TitleStatus 后面
func (h Hook) GetTitleAlerts() string {
	var alertNames = ""
	for _, alert := range h.Alerts {
		alertNames += "「" + alert.Labels.Alertname + "」"
	}
	return alertNames
}

func (h Hook) GetAlertDetailList() string {
	alertDetailList := ""
	for _, a := range h.Alerts {
		alertDetailList += "\n"
		alertDetailList += a.GetAlertDetail()
	}
	return alertDetailList
}

func (a Alert) GetAlertDetail() string {
	var color, endTimeString string
	if a.Status == RESOLVED {
		color = ww.WechatWorkColorGreen
		endTimeString = fmt.Sprintf(`\n><font color=\"%s\">恢复时间: </font>%s`, ww.WechatWorkColorGray, a.EndsAt.Format(TimeLayout))
	} else {
		color = ww.WechatWorkColorRed
	}
	return fmt.Sprintf(
		`
><font color=\"%s\">告警名称: </font><font color=\"%s\">**%s**</font>
><font color=\"%s\">信息: </font>%s
><font color=\"%s\">触发时间: </font>%s%s
><font color=\"%s\">图表: </font>[%s](%s)
><font color=\"%s\">仪表盘: </font>[%s](%s)
`,
		ww.WechatWorkColorGray, color, a.Labels.Alertname,
		ww.WechatWorkColorGray, a.GetMessage(),
		ww.WechatWorkColorGray, a.StartsAt.Format(TimeLayout),
		endTimeString,
		ww.WechatWorkColorGray, a.DashboardURL, a.DashboardURL,
		ww.WechatWorkColorGray, a.PanelURL, a.PanelURL+"&kiosk",
	)
}

func (a Alert) GetMessage() string {
	var color string
	if a.Status == RESOLVED {
		color = ww.WechatWorkColorGreen
	} else {
		color = ww.WechatWorkColorRed
	}
	message := ""
	messageRegexp := regexp.MustCompile(`^\[ metric='(.*)' labels=\{\} value=([-+]?\d+\.?\d*) ]$`)
	params := messageRegexp.FindStringSubmatch(a.ValueString)
	metric := params[1]
	value, err := strconv.ParseFloat(params[2], 64)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("%s", common.ConvertFailureWarning)
		fmt.Println()
		fmt.Printf("value: %s\n", params[2])
		fmt.Println()
		fmt.Println()
	}
	message += fmt.Sprintf(`\n\t<font color=\"%s\">指标: </font><font color=\"%s\">**%s**</font>`, ww.WechatWorkColorGray, color, metric)
	message += fmt.Sprintf(`\n\t<font color=\"%s\">值    : </font><font color=\"%s\">**%.2f**</font>`, ww.WechatWorkColorGray, color, value)
	return message
}
