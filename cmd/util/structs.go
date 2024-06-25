package structs

type MetroStop struct {
	Id              int    `json:"Id"`
	StationLocation string `json:"StationLocation"`
	AtcoCode        string `json:"AtcoCode"`
	Direction       string `json:"Direction"`
	Dest0           string `json:"Dest0"`
}

type Result struct {
	ODataContext string      `json:"@odata.context"`
	Value        []MetroStop `json:"value"`
}

type EmailContent struct {
	From    map[string]string `json:"from"`
	Subject string            `json:"subject"`
	Html    string            `json:"html"`
}

type EmailBody struct {
	Recipients [1]map[string]map[string]string `json:"recipients"`
	Content    EmailContent                    `json:"content"`
}