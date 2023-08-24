package utils

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				sourceFile := *(a.Value.Any().(*slog.Source))
				fullPath := strings.Split(sourceFile.File, "/")
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", fullPath[len(fullPath)-1], sourceFile.Line))
			}
			return a
		},
	}))
	return logger
}
