package minerapi

type MultipleResponse struct {
	ID int `json:"id"`
}

type APIStatus struct {
	Status      string `json:"STATUS"`
	When        int    `json:"When"`
	Code        int    `json:"Code"`
	Msg         string `json:"Msg"`
	Description string `json:"Description"`
}

type ResponseCommon struct {
	Status []APIStatus `json:"STATUS"`
	ID     int         `json:"id"`
}
