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
	RESOLVED = "resolved"
	FIRING   = "firing"
)

type Labels struct {
	Alertname string `json:"alertname"`
}

type Annotations struct {
	Unit string `json:"Unit"`
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
	return fmt.Sprintf(`
		{
		   	"msgtype": "markdown",
		   	"markdown": {
				"content": "%s\n%s"
		   	}
		}`, h.GetStatusCount(), h.GetAlertDetailList())
}

// PrintAlertLog TODO
func (h Hook) PrintAlertLog() {

}

func (h Hook) GetStatusCount() string {
	firingCount, resolvedCount := 0, 0
	for _, a := range h.Alerts {
		if a.Status == FIRING {
			firingCount++
		} else if a.Status == RESOLVED {
			resolvedCount++
		} else {
			fmt.Println(common.GrafanaUnknownStatusWarning)
		}
	}
	return fmt.Sprintf(`## 新增告警：<font color=\"%s\">%d</font> 例
\n ## 恢复正常：<font color=\"%s\">%d</font> 例`,
		ww.WechatWorkColorRed, firingCount,
		ww.WechatWorkColorGreen, resolvedCount,
	)
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
	var duringTimeString time.Duration
	if a.Status == RESOLVED {
		color = ww.WechatWorkColorGreen
		endTimeString = fmt.Sprintf(`\n>恢复时间：<font color=\"%s\">%s</font>`, ww.WechatWorkColorGreen, a.EndsAt.Format(common.TimeLayout))
		duringTimeString = a.EndsAt.Sub(a.StartsAt)
	} else {
		color = ww.WechatWorkColorRed
		duringTimeString = time.Now().Sub(a.StartsAt)
	}
	if duringTimeString < 0 {
		fmt.Println(common.GrafanaWrongTimeSynchronizationError)
	}
	return fmt.Sprintf(
		`
>告警名称：<font color=\"%s\">**%s**</font>
>状态：<font color=\"%s\">**%s**</font>
>信息：{%s
>}
>触发时间：<font color=\"%s\">%s</font>%s
>持续时长：%v
>图表：[%s](%s)
>仪表盘：[%s](%s)
`,
		color, a.Labels.Alertname,
		color, strings.ToUpper(a.Status),
		a.GetMessage(),

		ww.WechatWorkColorRed, a.StartsAt.Format(common.TimeLayout), endTimeString,
		common.FormatDuration(duringTimeString),
		a.DashboardURL, a.DashboardURL,
		a.PanelURL, a.PanelURL+"&kiosk",
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
	metricArray := strings.Split(a.ValueString, "], [")
	for _, metric := range metricArray {
		message += "\n"
		messageRegexp := regexp.MustCompile(`metric='(.*)' labels=\{(.*)\} value=([-+]?\d+\.?\d*)`)
		params := messageRegexp.FindStringSubmatch(metric)
		metric := params[1]
		// 暂时不使用 labels
		_ = params[2]
		value, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("%s", common.ConvertFailureWarning)
			fmt.Println()
			fmt.Printf("value: %s\n", params[3])
			fmt.Println()
			fmt.Println()
		}
		message += fmt.Sprintf(`\t\t%s：<font color=\"%s\">%.2f%s</font>`, metric, color, value, a.Annotations.Unit)
	}
	return message
}
