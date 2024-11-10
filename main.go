package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Brom95/duckllm/duckduckgo"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/yorukot/ansichroma"
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

	for chunk := range c {
		message.WriteString(chunk)
	}
	fmt.Println(string(markdown.Render(message.String(), 80, 0)))
}
func codeOutput(c <-chan string, render bool) {
	message := strings.Builder{}
	line := strings.Builder{}
	lang := ""
	for chank := range c {
		line.WriteString(chank)
		if strings.Contains(chank, "\n") {
			strLine := line.String()
			if strings.Contains(strLine, "```") {
				if len(strLine) > 3 {
					lang = strLine[3 : len(strLine)-1]
				}
			} else {
				message.WriteString(line.String())
			}
			line.Reset()
		}
	}

	if render {
		format := lang
		style := "github"
		background := "#0d1117"
		resultString, err := ansichroma.HightlightString(message.String(), format, style, background)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(resultString)
	} else {
		fmt.Println(message.String())
	}

}

var (
	model     *string
	render    *bool
	query     *string
	only_code *bool
)

func main() {
	model = flag.String("m", "gpt-4o-mini", "Model to use [gpt-4o-mini, mistralai/Mixtral-8x7B-Instruct-v0.1, meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo, claude-3-haiku-20240307]")
	render = flag.Bool("r", false, "Render output as markdowm")
	query = flag.String("q", "", "query to send in non interactive  mode")
	only_code = flag.Bool("c", false, "Try provide  only code in non interactive mode")
	flag.Parse()
	session := duckduckgo.NewSession(*model)
	session.Init()
	if *only_code {
		*query += "\nprovide only code. No additional comments"
	}
	if query != nil && len(*query) > 0 {
		responseStream := session.Send(*query)
		// message := strings.Builder{}
		if *only_code {
			codeOutput(responseStream, *render)
		} else {

			if *render {
				renderedOutput(responseStream)
			} else {
				for c := range responseStream {
					fmt.Print(c)
				}
			}
		}

		fmt.Println()

		return
	}
	fmt.Println("Use return twice to send message!!!")
	for {
		fmt.Print("Q: ")
		query, err := multilineInput()
		responseStream := session.Send(query)

		if err != nil {
			panic(err)
		}
		fmt.Print("A: ")

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
