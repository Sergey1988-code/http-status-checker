package types

import "time"

// Config хранит настройки из командной строки
type Config struct {
	FilePath   string
	URL        string
	Timeout    int
	OnlyErrors bool
}

// Result хранит результат проверки одного URL
type Result struct {
	URL        string
	StatusCode int
	Error      error
	Duration   time.Duration
}
