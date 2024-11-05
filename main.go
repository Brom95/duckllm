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
func main() {
	model := flag.String("m", "gpt-4o-mini", "Model to use [gpt-4o-mini, mistralai/Mixtral-8x7B-Instruct-v0.1, meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo, claude-3-haiku-20240307]")
	render := flag.Bool("r", false, "Render output as markdowm")
	flag.Parse()
	fmt.Println("Use return twice to send message!!!")
	session := duckduckgo.NewSession(*model)
	session.Init()

	for {
		fmt.Print("Q: ")
		query, err := multilineInput()
		if err != nil {
			panic(err)
		}
		fmt.Print("A: ")
		if *render {
			message := strings.Builder{}
			for chank := range session.Send(query) {
				message.WriteString(chank)
			}
			fmt.Println(string(markdown.Render(message.String(), 80, 6)))
		} else {

			for chank := range session.Send(query) {

				fmt.Print(chank)
			}
		}
		fmt.Println()
		fmt.Println()

	}

}
