package facade

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/goark/gpt-cli/gpt"
	"github.com/goark/gpt-cli/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type options struct {
	APIKey   string
	CacheDir string
	Logger   *zerolog.Logger
}

func getOptions() (*options, error) {
	apiKey := os.Getenv(gpt.ENV_API_KEY)
	if s := viper.GetString("api-key"); len(s) > 0 {
		apiKey = s
	}
	log, err := logger.New(
		logger.LevelFrom(viper.GetString("log-level")),
		viper.GetString("log-dir"),
	)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &options{
		APIKey:   apiKey,
		CacheDir: cache.Dir(Name),
		Logger:   log,
	}, nil
}

func getFiles(ss []string) ([]string, error) {
	paths := map[string]bool{}
	for _, s := range ss {
		pp, err := filepath.Glob(s)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("path", s))
		}
		for _, p := range pp {
			paths[p] = true
		}
	}
	if len(paths) > 0 {
		plist := make([]string, 0, len(paths))
		for k := range paths {
			plist = append(plist, k)
		}
		sort.Strings(plist)
		return plist, nil
	}
	return []string{}, nil
}
