package utils

import "strings"

const (
	TagImages  = "images"
	TagAvatars = "avatars"
	TagDocs    = "docs"
)

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
	case TagImages:
		folder = "images/"
	case TagAvatars:
		folder = "avatars/"
	case TagDocs:
		folder = "documents/"
	default:
		folder = "misc/"
	}

	return folder
}
