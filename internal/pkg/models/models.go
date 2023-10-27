package models

type Request struct {
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	Text string `json:"text"`
}
