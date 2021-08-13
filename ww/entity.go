package ww

const (
	WechatWorkBotWebhookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	WechatWorkColorGreen    = "info"
	WechatWorkColorGray     = "comment"
	WechatWorkColorRed      = "warning"
)

type WechatWorkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
