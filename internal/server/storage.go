package server

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yqnk/urlshortener/internal/model"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dataSource string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		short_url TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0
	);
	`

	_, err := db.Exec(query)
	return err
}

func (s *Storage) SaveURL(m model.URL) error {
	_, err := s.db.Exec(`
		INSERT INTO urls (short_url, original_url, clicks) VALUES (?, ?, ?)
		`, m.ShortURL, m.OriginalURL, m.Clicks)

	return err
}

func (s *Storage) GetURL(shortURL string) (*model.URL, error) {
	row := s.db.QueryRow(`
		SELECT short_url, original_url, clicks
		FROM urls
		WHERE short_url = ?
		`, shortURL)

	var m model.URL
	err := row.Scan(&m.ShortURL, &m.OriginalURL, &m.Clicks)
	if err == sql.ErrNoRows {
		return nil, errors.New("no row found")
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *Storage) AddClick(shortURL string) error {
	_, err := s.db.Exec(`
		UPDATE urls
		SET clicks = clicks + 1
		WHERE short_url = ?
		`, shortURL)
	return err
}

func (s *Storage) GetRandomURL() (string, error) {
	row := s.db.QueryRow(`
		SELECT short_url
		FROM urls
		ORDER BY RANDOM() LIMIT 1
		`)

	var shortURL string
	err := row.Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", errors.New("no row found")
	}

	if err != nil {
		return "", err
	}

	return shortURL, nil
}
