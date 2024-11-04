package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Brom95/duckllm/duckduckgo"
	"gopkg.in/AlecAivazis/survey.v1"
)

func multilineInput() (string, error) {
	text := ""
	prompt := &survey.Multiline{
		Message: "Q: ",
	}
	survey.AskOne(prompt, &text, nil)
	return text, nil

}
func main() {
	// reader := bufio.NewReader(os.Stdin)
	client := &http.Client{}
	initialRequst, err := http.NewRequest("GET", "https://duckduckgo.com/duckchat/v1/status", nil)
	if err != nil {
		panic(err)
	}
	initialRequst.Header.Set("x-vqd-accept", "1")
	initialResp, err := client.Do(initialRequst)

	if err != nil {
		panic(err)
	}
	vqd := initialResp.Header.Get("x-vqd-4")
	initialResp.Body.Close()
	context := duckduckgo.NewContext()

	for {

		query, err := multilineInput()
		if err != nil {
			panic(err)
		}
		context.Messages = append(context.Messages, duckduckgo.NewMessage("user", query))
		fmt.Print("A: ")
		contextJson, err := json.Marshal(context)
		if err != nil {
			panic(err)
		}

		queryRequest, err := http.NewRequest("POST", "https://duckduckgo.com/duckchat/v1/chat", bytes.NewBuffer(contextJson))
		if err != nil {
			panic(err)
		}
		queryRequest.Header.Set("x-vqd-4", vqd)
		queryRequest.Header.Set("Content-Type", "application/json")
		queryResponse, err := client.Do(queryRequest)
		if err != nil {
			panic(err)

		}
		vqd = queryResponse.Header.Get("x-vqd-4")
		scanner := bufio.NewScanner(queryResponse.Body)

		message := strings.Builder{}
		for scanner.Scan() {
			line := string(scanner.Text())
			spitedline := strings.Split(line, "data: ")
			if len(spitedline) > 1 {

				var answerContent map[string]any
				err := json.Unmarshal([]byte(spitedline[1]), &answerContent)
				if err != nil {
					continue
				}
				if chank, ok := answerContent["message"]; ok {
					fmt.Print(chank.(string))
					message.WriteString(chank.(string))
				}
			}
		}
		fmt.Println()
		context.Messages = append(context.Messages, duckduckgo.NewMessage("assistant", message.String()))

	}

}
