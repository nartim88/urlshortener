package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	URL string `json:"url"`
}
