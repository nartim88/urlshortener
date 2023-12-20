package service

import "github.com/nartim88/urlshortener/internal/models"

func (s service) makeShortenURL(sID models.ShortenID) models.ShortURL {
	return models.ShortURL(s.cfg.BaseURL + "/" + string(sID))
}
