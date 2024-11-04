package duckduckgo

type Context struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func NewContext() *Context {
	return &Context{
		Model:    "gpt-4o-mini",
		Messages: make([]Message, 0),
	}

}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role string, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}
