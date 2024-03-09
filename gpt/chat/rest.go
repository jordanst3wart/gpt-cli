package chat

import (
	"context"
	"io"

	"github.com/goark/errs"
	"github.com/sashabaranov/go-openai"
)

func (cctx *ChatContext) rest(ctx context.Context, client *openai.Client, w io.Writer) (string, error) {
	cctx.prepare.Stream = false
	cctx.Logger().Info().Interface("request", cctx.prepare).Send()
	resp, err := client.CreateChatCompletion(ctx, cctx.prepare)
	if err != nil {
		err = errs.Wrap(err, errs.WithContext("request", cctx.prepare))
		cctx.Logger().Error().Interface("error", err).Send()
		return "", err
	}
	cctx.Logger().Info().Interface("response", resp).Send()

	if len(resp.Choices) == 0 {
		return "", nil
	}
	resText := resp.Choices[0].Message.Content
	if _, err := io.WriteString(w, resText); err != nil {
		err = errs.Wrap(err)
		cctx.Logger().Error().Interface("error", err).Send()
		return "", err
	}
	return resText, nil
}
