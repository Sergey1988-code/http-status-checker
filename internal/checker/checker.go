package checker

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Sergey1988-code/http-status-checker/pkg/types"
)

// CheckURLs - основная функция проверки URL
func CheckURLs(config *types.Config) ([]types.Result, error) {
	urls, err := getURLs(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Timeout)*time.Second)
	defer cancel()

	results := checkURLsConcurrently(ctx, urls)

	if config.OnlyErrors {
		results = filterOnlyErrors(results)
	}

	return results, nil
}

// getURLs - читает URL из файла или возвращает одиночный URL
func getURLs(config *types.Config) ([]string, error) {
	if config.URL != "" {
		return []string{config.URL}, nil
	}

	if config.FilePath == "" {
		return nil, fmt.Errorf("no URLs provided")
	}

	file, err := os.Open(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && !strings.HasPrefix(url, "#") { // Игнорируем пустые строки и комментарии
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return urls, nil
}

// checkURLsConcurrently - проверяет URL параллельно с помощью горутин
func checkURLsConcurrently(ctx context.Context, urls []string) []types.Result {
	var wg sync.WaitGroup
	results := make([]types.Result, len(urls))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()

			start := time.Now()
			statusCode, err := checkSingleURL(ctx, client, url)
			duration := time.Since(start)

			results[i] = types.Result{
				URL:        url,
				StatusCode: statusCode,
				Error:      err,
				Duration:   duration,
			}
		}(i, url)
	}

	wg.Wait()
	return results
}

// checkSingleURL - проверяет один URL
func checkSingleURL(ctx context.Context, client *http.Client, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// filterOnlyErrors - фильтрует только URL с ошибками
func filterOnlyErrors(results []types.Result) []types.Result {
	var filtered []types.Result
	for _, result := range results {
		if result.Error != nil || result.StatusCode >= 400 {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// HasErrors - проверяет, есть ли в результатах ошибки
func HasErrors(results []types.Result) bool {
	for _, result := range results {
		if result.Error != nil || result.StatusCode >= 400 {
			return true
		}
	}
	return false
}

// PrintResults - красивый вывод результатов
func PrintResults(results []types.Result) {
	fmt.Printf("\n🚀 Результаты проверки (%d URL):\n\n", len(results))

	successCount := 0
	errorCount := 0

	for _, result := range results {
		color := getColorForStatus(result.StatusCode)
		emoji := getEmojiForStatus(result.StatusCode)
		statusText := getStatusText(result)

		fmt.Printf("%s %s %s\n", emoji, color, result.URL)
		fmt.Printf("   ↳ %s (%.2f сек)\n\n", statusText, result.Duration.Seconds())

		if result.StatusCode >= 200 && result.StatusCode < 400 {
			successCount++
		} else {
			errorCount++
		}
	}

	fmt.Printf("📊 Статистика:\n")
	fmt.Printf("   • Успешных: %d\n", successCount)
	fmt.Printf("   • Ошибок: %d\n", errorCount)
	fmt.Printf("   • Всего: %d\n", len(results))
}

// Вспомогательные функции для оформления
func getColorForStatus(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "🟢" // Зеленый для успешных
	case statusCode >= 300 && statusCode < 400:
		return "🔵" // Синий для перенаправлений
	case statusCode >= 400 && statusCode < 500:
		return "🟡" // Желтый для клиентских ошибок
	case statusCode >= 500 && statusCode < 600:
		return "🔴" // Красный для серверных ошибок
	default:
		return "⚫" // Черный для неизвестных
	}
}

func getEmojiForStatus(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "✅"
	case statusCode >= 300 && statusCode < 400:
		return "🔄"
	case statusCode == 404:
		return "❌"
	case statusCode >= 400 && statusCode < 500:
		return "⚠️"
	case statusCode >= 500 && statusCode < 600:
		return "💥"
	default:
		return "❓"
	}
}

func getStatusText(result types.Result) string {
	if result.Error != nil {
		return fmt.Sprintf("Ошибка: %v", result.Error)
	}

	switch result.StatusCode {
	case 200:
		return "200 OK"
	case 404:
		return "404 Not Found"
	case 500:
		return "500 Internal Server Error"
	default:
		return fmt.Sprintf("%d", result.StatusCode)
	}
}
