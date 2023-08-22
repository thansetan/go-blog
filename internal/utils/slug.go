package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func GenerateSlug(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	return fmt.Sprintf("%s-%d", s, time.Now().Unix())
}
