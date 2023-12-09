package v1

import (
	"github.com/nartim88/urlshortener/internal/models"
)

type Request struct {
	FullURL models.FullURL `json:"url"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	Result string `json:"result"`
}
