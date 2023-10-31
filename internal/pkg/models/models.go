package models

type Request struct {
	FullURL FullURL `json:"full_url"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	Result string `json:"result"`
}

type FullURL string

type ShortURL string
