package v2

import "github.com/nartim88/urlshortener/internal/pkg/models"

type Request struct {
	CorrelationID models.CorrelationID `json:"correlation_id"`
	FullURL       models.FullURL       `json:"original_url"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
}

type ResponsePayload struct {
	CorrelationID models.CorrelationID `json:"correlation_id"`
	ShortURL      string               `json:"short_url"`
}
