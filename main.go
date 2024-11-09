package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Brom95/duckllm/duckduckgo"
	markdown "github.com/MichaelMure/go-term-markdown"
)

func multilineInput() (string, error) {
	text := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {

			return "", err
		}

		// Check for double newline
		if line == "\n" {
			// If the last line was also a newline, break the loop
			break
		} else {
			text += line
		}
	}
	return text, nil

}
func renderedOutput(c <-chan string) {
	message := strings.Builder{}
	for chank := range c {
		message.WriteString(chank)
	}
	fmt.Print(string(markdown.Render(message.String(), 80, 6)))
}

func main() {
	model := flag.String("m", "gpt-4o-mini", "Model to use [gpt-4o-mini, mistralai/Mixtral-8x7B-Instruct-v0.1, meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo, claude-3-haiku-20240307]")
	render := flag.Bool("r", false, "Render output as markdowm")
	query := flag.String("q", "", "query to send in non interactive  mode")
	only_code := flag.Bool("c", false, "Try provide  only code in non interactive mode")
	_ = only_code
	flag.Parse()
	session := duckduckgo.NewSession(*model)
	session.Init()
	if query != nil && len(*query) > 0 {
		// message := strings.Builder{}
		if *only_code {
			*query += "\nprovide only code. No additional comments"
		}
		prev := false
		responseStream := session.Send(*query)
		if *render {
			renderedOutput(responseStream)
		} else {
			for chank := range responseStream {

				if *only_code && (prev || strings.Contains(chank, "```")) {
					prev = !prev
					continue
				}
				fmt.Print(chank)
			}

		}
		fmt.Println()

		return
	}
	fmt.Println("Use return twice to send message!!!")
	for {
		fmt.Print("Q: ")
		query, err := multilineInput()
		if err != nil {
			panic(err)
		}
		fmt.Print("A: ")
		responseStream := session.Send(query)
		if *render {
			renderedOutput(responseStream)
		} else {

			for chank := range responseStream {

				fmt.Print(chank)
			}
		}
		fmt.Println()

	}

}
