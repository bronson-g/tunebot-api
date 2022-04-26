package log

import "fmt"

func Red(txt string) string {
	return "\033[31m" + txt + "\033[0m"
}

func Green(txt string) string {
	return "\033[32m" + txt + "\033[39m"
}

func Yellow(txt string) string {
	return "\033[33m" + txt + "\033[39m"
}

func Blue(txt string) string {
	return "\033[34m" + txt + "\033[39m"
}

func Purple(txt string) string {
	return "\033[35m" + txt + "\033[39m"
}

func Cyan(txt string) string {
	return "\033[36m" + txt + "\033[39m"
}

func Println(message string) {
	fmt.Println(message)
}
