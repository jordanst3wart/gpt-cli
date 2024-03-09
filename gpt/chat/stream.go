package chat

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
	"github.com/sashabaranov/go-openai"
)

func (cctx *ChatContext) stream(ctx context.Context, client *openai.Client, w io.Writer) (string, error) {
	cctx.Logger().Info().Interface("request", cctx.prepare).Send()
	stream, err := client.CreateChatCompletionStream(ctx, cctx.prepare)
	if err != nil {
		err = errs.Wrap(err)
		cctx.Logger().Error().Interface("error", err).Send()
		return "", err
	}
	defer stream.Close()

	builder := strings.Builder{}
	fmt.Fprintln(w)
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errs.Is(err, io.EOF) {
				break
			}
			err = errs.Wrap(ecode.ErrStream, errs.WithCause(err))
			cctx.Logger().Error().Interface("error", err).Send()
			return "", err
		}
		cctx.Logger().Info().Interface("response", resp).Send()
		if len(resp.Choices) > 0 {
			if delta := resp.Choices[0].Delta.Content; len(delta) > 0 {
				if _, err := builder.WriteString(delta); err != nil {
					err = errs.Wrap(err)
					cctx.Logger().Error().Interface("error", err).Send()
					return "", err
				}
				if _, err := io.WriteString(w, delta); err != nil {
					err = errs.Wrap(err)
					cctx.Logger().Error().Interface("error", err).Send()
					return "", err
				}
			}
		}
	}
	fmt.Fprintln(w)
	return builder.String(), nil
}
