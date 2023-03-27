package chat

import (
	"encoding/json"
	"os"

	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
	"github.com/goark/gpt-cli/gpt"
	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
)

// ChatContext is context data for chat
type ChatContext struct {
	*gpt.GPTContext
	profile  openai.ChatCompletionRequest
	savePath string
}

// New function create new ChatContext instance.
func New(apiKey, cacheDir string, logger *zerolog.Logger, profilePath, savePath string) (*ChatContext, error) {
	gctx, err := gpt.New(apiKey, cacheDir, logger)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	profile := openai.ChatCompletionRequest{}
	if len(profilePath) > 0 {
		file, err := os.Open(profilePath)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("profilePath", profilePath))
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&profile); err != nil {
			return nil, errs.Wrap(err)
		}
	}
	if len(profile.Model) == 0 {
		profile.Model = openai.GPT3Dot5Turbo0301
	}
	if profile.Messages == nil {
		profile.Messages = []openai.ChatCompletionMessage{}
	}
	return &ChatContext{
		GPTContext: gctx,
		profile:    profile,
		savePath:   savePath,
	}, nil
}

func (cctx *ChatContext) SavePath() string {
	if cctx == nil {
		return ""
	}
	return cctx.savePath
}

// Save method saves openai.ChatCompletionRequest data.
func (cctx *ChatContext) Save() error {
	if cctx == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	var file *os.File
	var err error
	if len(cctx.savePath) == 0 {
		file, err = os.CreateTemp(cctx.CacheDir(), "chat.*.json")
	} else {
		file, err = os.Create(cctx.savePath)
	}
	if err != nil {
		return errs.Wrap(err, errs.WithContext("savePath", cctx.savePath))
	}
	defer file.Close()
	cctx.savePath = file.Name()
	return json.NewEncoder(file).Encode(cctx.profile)
}

/* MIT License
 *
 * Copyright 2023 Spiegel
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
