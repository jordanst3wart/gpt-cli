package chat

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goark/errs"
	"github.com/sashabaranov/go-openai"
)

// OutputHistory function converts markdown like text data from history data.
func OutputHistory(r io.Reader, w io.Writer, userName, assistantName string) error {
	hist := openai.ChatCompletionRequest{}
	if err := json.NewDecoder(r).Decode(&hist); err != nil {
		return errs.Wrap(err)
	}

	// Output
	fmt.Fprintln(w, "# Chat with GPT")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "- `model`:", hist.Model)
	if hist.MaxTokens != 0 {
		fmt.Fprintln(w, "- `max_tokens`:", hist.MaxTokens)
	}
	if hist.Temperature != 0 {
		fmt.Fprintln(w, "- `temperature`:", hist.Temperature)
	}
	if hist.TopP != 0 {
		fmt.Fprintln(w, "- `top_p`:", hist.TopP)
	}
	if hist.N != 0 {
		fmt.Fprintln(w, "- `n`:", hist.N)
	}
	for _, msg := range hist.Messages {
		role := msg.Role
		switch {
		case role == openai.ChatMessageRoleUser && len(userName) > 0:
			role = userName
		case role == openai.ChatMessageRoleAssistant && len(assistantName) > 0:
			role = assistantName
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, "##", role)
		fmt.Fprintln(w)
		fmt.Fprintln(w, msg.Content)
	}
	return nil
}
