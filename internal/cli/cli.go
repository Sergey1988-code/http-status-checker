package cli

import (
	"flag"
	"fmt"
	"github.com/Sergey1988-code/http-status-checker/pkg/types"
)

func ParseFlags() (*types.Config, error) {
	var config types.Config

	flag.StringVar(&config.FilePath, "file", "", "Path to file with URLs")
	flag.StringVar(&config.URL, "url", "", "Single URL to check")
	flag.IntVar(&config.Timeout, "timeout", 30, "Request timeout in seconds")
	flag.BoolVar(&config.OnlyErrors, "only-errors", false, "Show only errors")

	flag.Usage = func() {
		fmt.Printf("Usage: http-status-checker [options]\n\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  http-status-checker -file urls.txt\n")
		fmt.Printf("  http-status-checker -url https://example.com\n")
		fmt.Printf("  http-status-checker -file urls.txt -only-errors\n")
	}

	flag.Parse()

	if config.FilePath == "" && config.URL == "" {
		flag.Usage()
		return nil, fmt.Errorf("must provide either -file or -url")
	}

	return &config, nil
}
