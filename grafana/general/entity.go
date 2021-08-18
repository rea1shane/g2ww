package general

type Hook interface {
	MsgNews() string
	MsgMarkdown() string
	PrintAlertLog()
}
