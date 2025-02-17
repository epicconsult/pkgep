package pkgep

import (
	"strings"
)

func StringPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func PtrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
