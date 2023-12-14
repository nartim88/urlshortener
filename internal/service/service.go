package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/internal/storage"
	"github.com/nartim88/urlshortener/pkg/logger"
)

type Service interface {
	CreateShortenURL(ctx context.Context, fURL models.FullURL) (*models.ShortURL, error)
	GetAllURLs(ctx context.Context, user models.User) ([]models.ShortAndFullURLs, error)
	GetStore() storage.Storage
	GetFullURL(ctx context.Context, sID models.ShortenID) (*models.FullURL, error)
	GetConfigs() *config.Config
	DeleteURLs(ctx context.Context, IDs []models.ShortenID) error
	markAsDeletedListener()
	MarkAsDeletedCh() chan models.ShortenID
}

type service struct {
	store                 storage.Storage
	cfg                   *config.Config
	markAsDeletedCh       chan models.ShortenID
	markAsDeletedResultCh chan models.ShortenID
}

func NewService(s storage.Storage, cfg *config.Config) Service {
	svc := &service{
		store:                 s,
		cfg:                   cfg,
		markAsDeletedCh:       make(chan models.ShortenID, 1),
		markAsDeletedResultCh: make(chan models.ShortenID, 64),
	}
	go svc.markAsDeletedListener()
	return svc
}

func (s service) CreateShortenURL(ctx context.Context, fURL models.FullURL) (*models.ShortURL, error) {
	var shortURL models.ShortURL
	sID, err := s.store.Set(ctx, fURL)
	if err != nil {
		existsError := storage.ErrURLExists{}
		if errors.As(err, &existsError) {
			shortURL = models.ShortURL(s.cfg.BaseURL + "/" + string(existsError.SID))
			return &shortURL, existsError
		}
		return nil, fmt.Errorf("error while creating shorten url: %w", err)
	}
	shortURL = models.ShortURL(s.cfg.BaseURL + "/" + string(*sID))
	return &shortURL, nil
}

func (s service) GetAllURLs(ctx context.Context, user models.User) ([]models.ShortAndFullURLs, error) {
	data, err := s.store.ListURLs(ctx, user)
	if err != nil {
		return nil, err
	}
	URLs := make([]models.ShortAndFullURLs, 0)
	logger.Log.Info().Any("data", data).Send()
	for _, v := range data {
		shortURL := models.ShortURL(s.cfg.BaseURL + "/" + string(v.ShortenID))
		URLs = append(URLs, models.ShortAndFullURLs{ShortURL: shortURL, FullURL: v.FullURL})
	}
	return URLs, nil
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

func (s service) DeleteURLs(ctx context.Context, IDs []models.ShortenID) error {
	if err := s.store.MarkAsDeletedByID(ctx, IDs); err != nil {
		return err
	}
	return nil
}

func (s service) markAsDeletedListener() {
	logger.Log.Info().Msg("MarkAsDeletedListener is active")
	defer close(s.markAsDeletedResultCh)

	for sID := range s.markAsDeletedCh {
		s.markAsDeletedResultCh <- sID
	}

	logger.Log.Info().Msg("MarkAsDeletedListener is closed")
}

func (s service) MarkAsDeletedCh() chan models.ShortenID {
	return s.markAsDeletedCh
}
