package old

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
