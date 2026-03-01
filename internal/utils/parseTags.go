package utils

import "strings"

func ParseTags(tagFlags string) []string {
	var tags []string
	tags = append(tags, strings.Split(tagFlags, ",")...)

	return tags
}
