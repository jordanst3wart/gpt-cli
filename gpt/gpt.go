package gpt

import (
	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
)

const (
	ENV_API_KEY = "OPENAI_API_KEY"
)

type GPTContext struct {
	apiKey   string
	cacheDir string
	logger   *zerolog.Logger
}

// New function creates APIContext instance.
func New(apiKey, cacheDir string, logger *zerolog.Logger) (*GPTContext, error) {
	if len(apiKey) == 0 {
		return nil, errs.Wrap(ecode.ErrAPIKey)
	}
	if len(cacheDir) == 0 {
		cacheDir = "."
	}
	return &GPTContext{
		apiKey:   apiKey,
		cacheDir: cacheDir,
		logger:   logger,
	}, nil
}

// Client method creates new openai.Client.
func (gctx *GPTContext) Client() *openai.Client {
	return openai.NewClient(gctx.apiKey)
}

// CacheDir method returns cache directory.
func (gctx *GPTContext) CacheDir() string {
	return gctx.cacheDir
}

// Logger method returns logger.
func (gctx *GPTContext) Logger() *zerolog.Logger {
	return gctx.logger
}
