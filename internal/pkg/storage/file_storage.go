package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"

	"github.com/google/uuid"

	"github.com/nartim88/urlshortener/internal/pkg/logger"
	"github.com/nartim88/urlshortener/internal/pkg/models"
	"github.com/nartim88/urlshortener/internal/pkg/service"
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

// Get возвращает полный урл по сокращенному.
func (s *FileStorage) Get(sURL models.ShortURL) (*models.FullURL, error) {
	entry, err := s.getByShortURL(sURL)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &entry.FullURL, nil
}

// Set сохраняет в базу полный УРЛ и соответствующий ему короткий УРЛ
func (s *FileStorage) Set(fURL models.FullURL) (*models.ShortURL, error) {
	randChars := service.GenerateRandChars(shortURLLen)
	shortURL := models.ShortURL(randChars)

	_, err := s.getByShortURL(shortURL)
	if err != nil {
		return nil, err
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	newEntry := models.JSONEntry{
		UUID:     &newUUID,
		ShortURL: shortURL,
		FullURL:  fURL,
	}

	if err = s.saveToFile(newEntry); err != nil {
		return nil, err
	}

	return &shortURL, nil
}

// getByShortURL возвращает полный урл по короткому
func (s *FileStorage) getByShortURL(sURL models.ShortURL) (*models.JSONEntry, error) {
	var entry models.JSONEntry

	scanner, err := s.newScanner()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		if entry.ShortURL == sURL {
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

func (s *FileStorage) saveToFile(entry models.JSONEntry) error {
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
