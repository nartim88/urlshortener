package config

// DB constants
const (
	DBTableName   = "shortener"
	DBFullURLRow  = "full_url"
	DBShortURLRow = "short_url"
)

const (
	RunAddr         = "localhost:8080"
	BaseURL         = "http://localhost:8080"
	LogLevel        = "info"
	FileStoragePath = "/tmp/short-url-db.json"
	DatabaseDSN     = "host=localhost user=videos password=videos dbname=videos"
)

const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
