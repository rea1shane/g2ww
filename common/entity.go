package common

const (
	InternalError                           int = -1
	OK                                      int = 0
	GrafanaWebhookUnmarshalJsonError        int = 100001
	ClientCallAPIError                      int = 100002
	WechatWorkCallAPIError                  int = 200001
	WechatWorkCallAPIWrongJsonFormatWarning int = 200002
	WechatWorkParseResponseBodyFailure      int = 200003
)

type Counter struct {
	SentSuccessCount int `json:"sentSuccessCount"`
	SentFailureCount int `json:"sentFailureCount"`
}
