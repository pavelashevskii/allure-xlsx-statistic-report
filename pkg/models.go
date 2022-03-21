package pkg

type TestObject struct {
	Uid             string   `json:"uid"`
	Name            string   `json:"name"`
	FullName        string   `json:"fullName"`
	HistoryId       string   `json:"historyId"`
	Extra           Extra    `json:"extra"`
	ParameterValues []string `json:"parameterValues"`
}

type Extra struct {
	History History `json:"history"`
}

type History struct {
	Statistic GlobalStatistic `json:"statistic"`
	Items     []TestRunItem   `json:"items"`
}

type GlobalStatistic struct {
	Failed  int `json:"failed"`
	Broken  int `json:"broken"`
	Skipped int `json:"skipped"`
	Passed  int `json:"passed"`
	Unknown int `json:"unknown"`
	Total   int `json:"total"`
}

type TestRunItem struct {
	Uid           string   `json:"uid"`
	ReportUrl     string   `json:"reportUrl"`
	Status        string   `json:"status"`
	StatusDetails string   `json:"statusDetails"`
	Time          TimeItem `json:"time"`
}

type TimeItem struct {
	Start    uint64 `json:"start"`
	End      uint64 `json:"end"`
	Duration uint   `json:"duration"`
}
