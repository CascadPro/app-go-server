package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type UserAgent struct {
	AppName    string
	AppVersion string

	OS    string
	Model string
	Type  string
}

func ParseUserAgent(ua string) (*UserAgent, error) {
	ua = strings.TrimSpace(ua)

	re := regexp.MustCompile(`^([^/]+)/([^(]+)\s*\(([^;]+);\s*([^;]+);\s*([^)]+)\)$`)

	matches := re.FindStringSubmatch(ua)
	if matches == nil {
		return nil, fmt.Errorf("неверный формат заголовка user agent")
	}

	deviceType := strings.TrimSpace(matches[5])

	switch deviceType {
	case "1":
		deviceType = "phone"
	case "2":
		deviceType = "tablet"
	case "3":
		deviceType = "desktop"
	default:
		deviceType = "unknown"
	}

	return &UserAgent{
		AppName:    strings.TrimSpace(matches[1]),
		AppVersion: strings.TrimSpace(matches[2]),
		OS:         strings.TrimSpace(matches[3]),
		Model:      strings.TrimSpace(matches[4]),
		Type:       deviceType,
	}, nil
}
