package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/internal/storage"
)

type Service interface {
	CreateShortenURL(ctx context.Context, fURL models.FullURL) (*models.ShortURL, error)
	GetAllURLs(ctx context.Context, user models.User) ([]*models.ShortAndFullURLs, error)
	GetStore() storage.Storage
	GetFullURL(ctx context.Context, sID models.ShortenID) (*models.FullURL, error)
	GetConfigs() *config.Config
}

type service struct {
	store storage.Storage
	cfg   *config.Config
}

func NewService(s storage.Storage, cfg *config.Config) Service {
	return &service{s, cfg}
}

func (s service) CreateShortenURL(ctx context.Context, fURL models.FullURL) (*models.ShortURL, error) {
	var shortURL models.ShortURL
	sID, err := s.store.Set(ctx, fURL)
	if err != nil {
		existsError := storage.URLExistsError{}
		if errors.As(err, &existsError) {
			shortURL = models.ShortURL(s.cfg.BaseURL + "/" + string(existsError.SID))
			return &shortURL, existsError
		}
		return nil, fmt.Errorf("error while creating shorten url: %w", err)
	}
	shortURL = models.ShortURL(s.cfg.BaseURL + "/" + string(*sID))
	return &shortURL, nil
}

func (s service) GetAllURLs(ctx context.Context, user models.User) ([]*models.ShortAndFullURLs, error) {
	data, err := s.store.ListURLs(ctx, user)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s service) GetStore() storage.Storage {
	return s.store
}

func (s service) GetFullURL(ctx context.Context, sID models.ShortenID) (*models.FullURL, error) {
	fURL, err := s.store.Get(ctx, sID)
	if err != nil {
		return nil, err
	}
	return fURL, nil
}

func (s service) GetConfigs() *config.Config {
	return s.cfg
}
