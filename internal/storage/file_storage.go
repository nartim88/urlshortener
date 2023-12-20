package storage

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/pkg/logger"
)

type FileStorage struct {
	// FilePath абсолютный путь к файлу для хранения данных
	FilePath string
	FilePerm os.FileMode
}

// NewFileStorage инициализация конкретного Storage
func NewFileStorage(path string) (Storage, error) {
	s := FileStorage{path, 0666}
	if !s.fileExists() {
		err := s.createFile()
		if err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (s *FileStorage) Get(ctx context.Context, sID models.ShortenID) (*models.FullURL, error) {
	entry, err := s.getByShortURL(sID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &entry.FullURL, nil
}

func (s *FileStorage) Set(ctx context.Context, fURL models.FullURL) (*models.ShortenID, error) {
	randChars := GenerateRandChars(shortURLLen)
	sID := models.ShortenID(randChars)

	_, err := s.getByShortURL(sID)
	if err != nil {
		return nil, err
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	newEntry := models.FileJSONEntry{
		ID:        &newUUID,
		ShortenID: sID,
		FullURL:   fURL,
	}

	if err = s.saveToFile(newEntry); err != nil {
		return nil, err
	}

	return &sID, nil
}

func (s *FileStorage) ListURLs(ctx context.Context, u models.User) ([]SIDAndFullURL, error) {
	return nil, nil
}

func (s *FileStorage) MarkAsDeletedByID(ctx context.Context, IDs []models.ShortenID) error {
	return nil
}

// getByShortURL возвращает полный урл по короткому
func (s *FileStorage) getByShortURL(sID models.ShortenID) (*models.FileJSONEntry, error) {
	var entry models.FileJSONEntry

	scanner, err := s.newScanner()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		if entry.ShortenID == sID {
			return &entry, nil
		}
	}

	return nil, nil
}

// newScanner для чтения из файла без буфера
func (s *FileStorage) newScanner() (*bufio.Scanner, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(data)
	return bufio.NewScanner(buf), nil
}

func (s *FileStorage) saveToFile(entry models.FileJSONEntry) error {
	file, err := os.OpenFile(s.FilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, s.FilePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logger.Log.Info().Err(err).Send()
		}
	}(file)

	w := json.NewEncoder(file)

	if err = w.Encode(entry); err != nil {
		return err
	}

	return nil
}

// fileExists проверяет существует ли FilePath
func (s *FileStorage) fileExists() bool {
	_, err := os.Stat(s.FilePath)
	return !os.IsNotExist(err)
}

func (s *FileStorage) createFile() error {
	f, err := os.Create(s.FilePath)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) SetBatch(ctx context.Context, fURLs []models.FullURL) (map[models.FullURL]models.ShortenID, error) {
	return nil, nil
}
