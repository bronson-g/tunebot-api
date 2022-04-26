package log

import "fmt"

func Red(txt string) string {
	return "\033[31m" + txt + "\033[0m"
}

func Green(txt string) string {
	return "\033[32m" + txt + "\033[39m"
}

func Println(message string) {
	fmt.Println(message)
}
