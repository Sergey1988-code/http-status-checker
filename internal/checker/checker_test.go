package checker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Sergey1988-code/http-status-checker/pkg/types"
)

// TestCheckSingleURL тестирует проверку одного URL
func TestCheckSingleURL(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/200":
			w.WriteHeader(http.StatusOK)
		case "/404":
			w.WriteHeader(http.StatusNotFound)
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
		case "/slow":
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "Status 200",
			url:        server.URL + "/200",
			wantStatus: 200,
			wantError:  false,
		},
		{
			name:       "Status 404",
			url:        server.URL + "/404",
			wantStatus: 404,
			wantError:  false,
		},
		{
			name:       "Status 500",
			url:        server.URL + "/500",
			wantStatus: 500,
			wantError:  false,
		},
		{
			name:       "Invalid URL",
			url:        "http://invalid-url-that-does-not-exist.local",
			wantStatus: 0,
			wantError:  true,
		},
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			status, err := checkSingleURL(ctx, client, tt.url)

			if (err != nil) != tt.wantError {
				t.Errorf("checkSingleURL() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("checkSingleURL() = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

// TestGetURLs тестирует чтение URL из файла
func TestGetURLs(t *testing.T) {
	// Создаем временный файл с URL
	content := `https://example.com/1
https://example.com/2

# Это комментарий
https://example.com/3

https://example.com/4`

	tmpfile, err := os.CreateTemp("", "urls.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	tests := []struct {
		name     string
		config   *types.Config
		wantURLs []string
		wantErr  bool
	}{
		{
			name: "Read from file",
			config: &types.Config{
				FilePath: tmpfile.Name(),
			},
			wantURLs: []string{
				"https://example.com/1",
				"https://example.com/2",
				"https://example.com/3",
				"https://example.com/4",
			},
			wantErr: false,
		},
		{
			name: "Single URL",
			config: &types.Config{
				URL: "https://google.com",
			},
			wantURLs: []string{"https://google.com"},
			wantErr:  false,
		},
		{
			name: "No URLs provided",
			config: &types.Config{
				FilePath: "",
				URL:      "",
			},
			wantURLs: nil,
			wantErr:  true,
		},
		{
			name: "File not found",
			config: &types.Config{
				FilePath: "/nonexistent/file.txt",
			},
			wantURLs: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := getURLs(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("getURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(urls) != len(tt.wantURLs) {
				t.Errorf("getURLs() = %v, want %v", urls, tt.wantURLs)
				return
			}

			for i, url := range urls {
				if url != tt.wantURLs[i] {
					t.Errorf("getURLs()[%d] = %v, want %v", i, url, tt.wantURLs[i])
				}
			}
		})
	}
}

// TestCheckURLsConcurrently тестирует параллельную проверку URL
func TestCheckURLsConcurrently(t *testing.T) {
	// Создаем тестовый сервер с небольшой задержкой
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond) // Небольшая задержка
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	urls := []string{
		server.URL + "/1",
		server.URL + "/2",
		server.URL + "/3",
		server.URL + "/4",
		server.URL + "/5",
	}

	ctx := context.Background()
	results := checkURLsConcurrently(ctx, urls)

	if len(results) != len(urls) {
		t.Errorf("checkURLsConcurrently() returned %d results, want %d", len(results), len(urls))
	}

	for i, result := range results {
		if result.URL != urls[i] {
			t.Errorf("Result %d URL = %v, want %v", i, result.URL, urls[i])
		}
		if result.StatusCode != 200 {
			t.Errorf("Result %d StatusCode = %v, want 200", i, result.StatusCode)
		}
		if result.Error != nil {
			t.Errorf("Result %d Error = %v, want nil", i, result.Error)
		}
		// Проверяем что Duration не отрицательный (может быть 0 для очень быстрых запросов)
		if result.Duration < 0 {
			t.Errorf("Result %d Duration should not be negative, got %v", i, result.Duration)
		}
	}
}

// TestFilterOnlyErrors тестирует фильтрацию только ошибок
func TestFilterOnlyErrors(t *testing.T) {
	results := []types.Result{
		{URL: "https://example.com/1", StatusCode: 200, Error: nil},
		{URL: "https://example.com/2", StatusCode: 404, Error: nil},
		{URL: "https://example.com/3", StatusCode: 500, Error: nil},
		{URL: "https://example.com/4", StatusCode: 0, Error: fmt.Errorf("connection failed")},
		{URL: "https://example.com/5", StatusCode: 301, Error: nil},
	}

	filtered := filterOnlyErrors(results)

	if len(filtered) != 3 {
		t.Errorf("filterOnlyErrors() returned %d results, want 3", len(filtered))
	}

	expectedURLs := []string{
		"https://example.com/2",
		"https://example.com/3",
		"https://example.com/4",
	}

	for i, result := range filtered {
		if result.URL != expectedURLs[i] {
			t.Errorf("Filtered result %d = %v, want %v", i, result.URL, expectedURLs[i])
		}
	}
}

// TestHasErrors тестирует проверку наличия ошибок
func TestHasErrors(t *testing.T) {
	tests := []struct {
		name    string
		results []types.Result
		wantHas bool
	}{
		{
			name: "No errors",
			results: []types.Result{
				{StatusCode: 200},
				{StatusCode: 301},
			},
			wantHas: false,
		},
		{
			name: "With client errors",
			results: []types.Result{
				{StatusCode: 200},
				{StatusCode: 404},
			},
			wantHas: true,
		},
		{
			name: "With server errors",
			results: []types.Result{
				{StatusCode: 200},
				{StatusCode: 500},
			},
			wantHas: true,
		},
		{
			name: "With connection errors",
			results: []types.Result{
				{StatusCode: 200},
				{Error: fmt.Errorf("connection failed")},
			},
			wantHas: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasErrors(tt.results); got != tt.wantHas {
				t.Errorf("HasErrors() = %v, want %v", got, tt.wantHas)
			}
		})
	}
}

// TestStatusHelpers тестирует вспомогательные функции
func TestStatusHelpers(t *testing.T) {
	tests := []struct {
		statusCode int
		wantColor  string
		wantEmoji  string
	}{
		{200, "🟢", "✅"},
		{301, "🔵", "🔄"},
		{404, "🟡", "❌"},
		{500, "🔴", "💥"},
		{999, "⚫", "❓"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Status%d", tt.statusCode), func(t *testing.T) {
			if got := getColorForStatus(tt.statusCode); got != tt.wantColor {
				t.Errorf("getColorForStatus(%d) = %v, want %v", tt.statusCode, got, tt.wantColor)
			}
			if got := getEmojiForStatus(tt.statusCode); got != tt.wantEmoji {
				t.Errorf("getEmojiForStatus(%d) = %v, want %v", tt.statusCode, got, tt.wantEmoji)
			}
		})
	}
}
