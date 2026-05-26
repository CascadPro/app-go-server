package utils

import "strings"

func ValidateTagParam(tag string, cfg *Config) (string, error) {
	allowedTags := map[string]struct{}{}

	for allowedTag := range strings.SplitSeq(cfg.S3AllowedTags, ",") {
		allowedTags[strings.TrimSpace(allowedTag)] = struct{}{}
	}

	if _, exists := allowedTags[tag]; !exists {
		return "Неверный тэг!", nil
	}

	return tag, nil
}

func GetBucketFolder(tag string) string {
	var folder string

	switch tag {
	case "images":
		folder = "images/"
	case "avatars":
		folder = "avatars/"
	case "docs":
		folder = "documents/"
	default:
		folder = "misc/"
	}

	return folder
}
