package models

import "github.com/google/uuid"

type FullURL string
type ShortURL string

type Request struct {
	FullURL FullURL `json:"url"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	Result string `json:"result"`
}

type JSONEntry struct {
	UUID     *uuid.UUID `json:"uuid"`
	ShortURL ShortURL   `json:"short_url"`
	FullURL  FullURL    `json:"full_url"`
}