package chat

import (
	"context"
	"io"

	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
)

// Request requesta OpenAI Chat completion, and returns response message.
func (cctx *ChatContext) Request(ctx context.Context, rest bool, msgs []string, w io.Writer) error {
	if cctx == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if err := cctx.AppendUserMessages(msgs); err != nil {
		return errs.Wrap(err)
	}
	var err error
	var resText string
	if rest {
		resText, err = cctx.rest(ctx, cctx.Client(), w)
	} else {
		resText, err = cctx.stream(ctx, cctx.Client(), w)
	}
	if err != nil {
		return errs.Wrap(err)
	}
	_ = cctx.AppendAssistantMessages([]string{resText})
	return cctx.Save()
}
