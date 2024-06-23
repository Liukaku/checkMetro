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