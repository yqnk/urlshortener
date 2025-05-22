package server

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"github.com/yqnk/urlshortener/internal/model"
)

type Service struct {
	storage *Storage
}

func NewService(dataSource string) (*Service, error) {
	storage, err := NewStorage(dataSource)
	if err != nil {
		return nil, err
	}

	return &Service{storage: storage}, nil
}

func (s *Service) GenerateURL(originalURL string) string {
	h := sha1.New()
	h.Write([]byte(originalURL))
	bs := h.Sum(nil)
	shortURL := base64.URLEncoding.EncodeToString(bs)
	return strings.TrimRight(shortURL[:6], "=") // only keep 6 chars
}

func (s *Service) ShortenURL(originalURL string) (*model.URL, error) {
	shortURL := s.GenerateURL(originalURL)

	// check if url already present and return it if so
	if m, err := s.storage.GetURL(shortURL); err == nil {
		return m, nil
	}

	// otherwise, create it
	newURL := model.URL{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		Clicks:      0,
	}

	if err := s.storage.SaveURL(newURL); err != nil {
		return nil, err
	}

	return &newURL, nil
}

func (s *Service) GetOriginalURL(shortURL string) (*model.URL, error) {
	return s.storage.GetURL(shortURL)
}

func (s *Service) AddClick(shortURL string) error {
	return s.storage.AddClick(shortURL)
}

func (s *Service) GetRandomURL() (string, error) {
	return s.storage.GetRandomURL()
}
