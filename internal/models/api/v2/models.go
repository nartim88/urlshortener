package v2

import (
	"github.com/nartim88/urlshortener/internal/models"
)

type Request struct {
	Data []RequestData
}

type RequestData struct {
	CorrelationID models.CorrelationID `json:"correlation_id"`
	FullURL       models.FullURL       `json:"original_url"`
}

type Response struct {
	Response []ResponsePayload
}

type ResponsePayload struct {
	CorrelationID models.CorrelationID `json:"correlation_id"`
	ShortURL      string               `json:"short_url"`
}
