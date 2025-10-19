package cli

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantFile    string
		wantURL     string
		wantTimeout int
		wantErr     bool
	}{
		{
			name:        "File provided",
			args:        []string{"-file", "urls.txt"},
			wantFile:    "urls.txt",
			wantURL:     "",
			wantTimeout: 30,
			wantErr:     false,
		},
		{
			name:        "URL provided",
			args:        []string{"-url", "https://example.com"},
			wantFile:    "",
			wantURL:     "https://example.com",
			wantTimeout: 30,
			wantErr:     false,
		},
		{
			name:        "Custom timeout",
			args:        []string{"-url", "https://example.com", "-timeout", "60"},
			wantFile:    "",
			wantURL:     "https://example.com",
			wantTimeout: 60,
			wantErr:     false,
		},
		{
			name:        "Only errors flag",
			args:        []string{"-file", "urls.txt", "-only-errors"},
			wantFile:    "urls.txt",
			wantURL:     "",
			wantTimeout: 30,
			wantErr:     false,
		},
		{
			name:    "No arguments",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальные аргументы
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Устанавливаем тестовые аргументы
			os.Args = append([]string{"test"}, tt.args...)

			// Сбрасываем флаги для чистого теста
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			config, err := ParseFlags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Ожидаемая ошибка
			}

			if config.FilePath != tt.wantFile {
				t.Errorf("ParseFlags().FilePath = %v, want %v", config.FilePath, tt.wantFile)
			}
			if config.URL != tt.wantURL {
				t.Errorf("ParseFlags().URL = %v, want %v", config.URL, tt.wantURL)
			}
			if config.Timeout != tt.wantTimeout {
				t.Errorf("ParseFlags().Timeout = %v, want %v", config.Timeout, tt.wantTimeout)
			}
		})
	}
}
