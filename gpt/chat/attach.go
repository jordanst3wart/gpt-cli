package chat

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gpt-cli/ecode"
)

// AttachFile function converts markdown like text from text file data
func AttachFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("path", path))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("path", path))
	}
	if bytes.Contains(b, []byte{0}) {
		return "", errs.Wrap(ecode.ErrBinary, errs.WithContext("path", path))
	}
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "Path: %s\n\n", path)
	fmt.Fprintln(builder, "```")
	if _, err := io.Copy(builder, bytes.NewReader(b)); err != nil {
		return "", errs.Wrap(err, errs.WithContext("path", path))
	}
	fmt.Fprintln(builder, "\n```")
	return builder.String(), nil
}
