package duckduckgo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Session struct {
	Context *Context
	VQD     string
	init    bool
	client  *http.Client
}

func NewSession(model string) *Session {
	return &Session{
		Context: NewContext(model),
		init:    false,
		client:  &http.Client{},
	}
}
func (s *Session) Init() {
	defer func() {
		s.init = true
	}()
	initialRequst, err := http.NewRequest("GET", "https://duckduckgo.com/duckchat/v1/status", nil)
	if err != nil {
		panic(err)
	}
	initialRequst.Header.Set("x-vqd-accept", "1")
	initialResp, err := s.client.Do(initialRequst)

	if err != nil {
		panic(err)
	}
	s.VQD = initialResp.Header.Get("x-vqd-4")
	initialResp.Body.Close()
}
func (s *Session) Send(msg string) <-chan string {
	if !s.init {
		panic(fmt.Errorf("session not initialized"))
	}
	s.Context.Messages = append(s.Context.Messages, NewMessage("user", msg))
	contextJson, err := json.Marshal(s.Context)
	if err != nil {
		panic(err)
	}

	queryRequest, err := http.NewRequest("POST", "https://duckduckgo.com/duckchat/v1/chat", bytes.NewBuffer(contextJson))
	if err != nil {
		panic(err)
	}
	queryRequest.Header.Set("x-vqd-4", s.VQD)
	queryRequest.Header.Set("Content-Type", "application/json")
	queryResponse, err := s.client.Do(queryRequest)
	if err != nil {
		panic(err)

	}
	s.VQD = queryResponse.Header.Get("x-vqd-4")
	result := make(chan string)
	go func() {
		scanner := bufio.NewScanner(queryResponse.Body)
		message := strings.Builder{}
		defer queryResponse.Body.Close()
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

					result <- chank.(string)
					message.WriteString(chank.(string))
				}
			}
		}
		s.Context.Messages = append(s.Context.Messages, NewMessage("assistant", message.String()))
		close(result)

	}()
	return result
}
