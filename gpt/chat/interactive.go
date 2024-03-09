package chat

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
	"github.com/hymkor/go-multiline-ny"
	"github.com/nyaosorg/go-readline-ny"
)

// Interactive method is chatting in interactive mode (stream access).
func (cctx *ChatContext) Interactive(ctx context.Context, writer io.Writer) error {
	if cctx == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	client := cctx.Client()
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) { return fmt.Fprint(w, "\nChat>") },
	}
	fmt.Fprintln(writer, "Input 'q' or 'quit' to stop")
	cctx.prepare.Stream = true
	for {
		text, err := editor.ReadLine(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		if strings.EqualFold(text, "q") || strings.EqualFold(text, "quit") {
			break
		}
		_ = cctx.AppendUserMessages([]string{text})
		resText, err := cctx.stream(ctx, client, writer)
		if err != nil {
			return errs.Wrap(err)
		}
		_ = cctx.AppendAssistantMessages([]string{resText})
	}
	return cctx.Save()
}

// InteractiveMulti method is chatting in interactive mode with multiline editing (stream access).
func (cctx *ChatContext) InteractiveMulti(ctx context.Context, writer io.Writer) error {
	if cctx == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	client := cctx.Client()
	var editor multiline.Editor
	editor.SetPrompt(func(w io.Writer, lnum int) (int, error) {
		return fmt.Fprintf(w, "Chat:%2d>", lnum+1)
	})

	fmt.Fprintln(writer, "Input 'Ctrl+J' or 'Ctrl+Enter' to submit message")
	fmt.Fprintln(writer, "Input 'Ctrl+D' with no chars to stop")
	fmt.Fprintln(writer, "      or input text \"q\" or \"quit\" and submit to stop")
	cctx.prepare.Stream = true
	for {
		lines, err := editor.Read(ctx)
		if err != nil {
			if errs.Is(err, io.EOF) {
				break
			}
			return errs.Wrap(err)
		}
		if len(lines) == 0 {
			fmt.Fprintln(writer)
			continue
		}
		text := strings.TrimSpace(strings.Join(lines, "\n"))
		if len(text) == 0 {
			fmt.Fprintln(writer)
			continue
		}
		if strings.EqualFold(text, "q") || strings.EqualFold(text, "quit") {
			break
		}

		_ = cctx.AppendUserMessages([]string{text})
		resText, err := cctx.stream(ctx, client, writer)
		if err != nil {
			return errs.Wrap(err)
		}
		fmt.Fprintln(writer)
		_ = cctx.AppendAssistantMessages([]string{resText})
	}
	return cctx.Save()
}
