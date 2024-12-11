package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_interpreter/object"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

func openaiCall(prompt string) (string, bool) {
	key := "<your_key>"
	if key == "<your_key>" {
		return "missing api key in utils.openai..", false
	}

	reqBody := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", false
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	var chatResp ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", false
	}

	return chatResp.Choices[0].Message.Content, true
}

func Cattsort(args *object.Array) (*object.Array, string) {
	tmp := make(map[string]object.Object)

	var sb strings.Builder
	sb.WriteString("[")
	for _, value := range args.Elements {
		val := fmt.Sprintf("%s", value.Inspect())

		tmp[val] = value
		sb.WriteString(fmt.Sprintf("'%s', ", value.Inspect()))
	}
	sb.WriteString("]")

	prompt := fmt.Sprintf(`sort this array '%s', you should sort this in a way that makes sense to humans but maintain the single quotes around all values when returning the array`, sb.String())

	input, ok := openaiCall(prompt)
	if !ok {
		return nil, input
	}

	trimmed := strings.Trim(input, "[]")

	var result []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(trimmed); i++ {
		char := trimmed[i]

		switch char {
		case '\'':
			inQuote = !inQuote
			continue
		case ',':
			if !inQuote {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
				continue
			}
		case ' ':
			if !inQuote {
				continue
			}
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	var elements []object.Object
	for _, val := range result {
		value, ok := tmp[val]
		if !ok {
			return nil, ""
		}
		elements = append(elements, value)
	}

	return &object.Array{Elements: elements}, ""
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func Cattify(message string) (string, string) {
	prompt := fmt.Sprintf(`take this message but cattify it, you should respond 
        mainly with noises that cats make, like Mrrp or meow, 
        but make sure to include punctuation. Only reply with what 
        you would expect this phrase to be translated in cat-speak, only reply with the cat
        translation of '%s'`, message)

	result, ok := openaiCall(prompt)
	if !ok {
		return "", "A catastrophe occured.."
	}

	return result, ""
}
