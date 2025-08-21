package utils

import "fmt"

func Log(format string, level string) {
	switch level {
	case "low":
		fmt.Printf("\033[32m%s\033[0m\n", format)
	case "medium":
		fmt.Printf("\033[33m%s\033[0m\n", format)
	case "high":
		fmt.Printf("\033[31m%s\033[0m\n", format)
	}
}
