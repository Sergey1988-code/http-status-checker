package main

import (
	"log"
	"os"

	"github.com/Sergey1988-code/http-status-checker/internal/checker"
	"github.com/Sergey1988-code/http-status-checker/internal/cli"
)

func main() {
	// Парсим аргументы командной строки
	config, err := cli.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем проверку
	results, err := checker.CheckURLs(config)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим результаты
	checker.PrintResults(results)

	// Если есть ошибки - выходим с кодом 1
	if checker.HasErrors(results) {
		os.Exit(1)
	}
}
