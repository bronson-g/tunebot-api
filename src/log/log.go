package log

import "fmt"

func Red(txt string) string {
	return "\\e[31m" + txt + "\\e[39m"
}

func Green(txt string) string {
	return "\\e[32m" + txt + "\\e[39m"
}

func Println(message string) {
	fmt.Println(message)
}
