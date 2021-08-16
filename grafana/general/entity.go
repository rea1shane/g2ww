package general

const (
	OK          = "ok"
	Alerting    = "alerting"
	OKMsg       = "OK"
	AlertingMsg = "Alerting"
)

type Hook interface {
	MsgNews() string
	MsgMarkdown() string
	PrintAlertLog()
}
