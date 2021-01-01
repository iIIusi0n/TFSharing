package main

import "fmt"

func AskAnswer(question string) string {
	var answer string
	fmt.Print(question)
	fmt.Scanln(&answer)
	return answer
}
