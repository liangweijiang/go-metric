package utils

import (
	"strings"
	"unicode"
)

// SanitizeMetricName 将非标准化的指标名称转换为符合OpenTelemetry规范的格式
func SanitizeMetricName(name string) string {
	name = strings.ToLower(name)
	var sb strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	name = strings.Trim(sb.String(), "_")
	if !unicode.IsLetter(rune(name[0])) {
		name = "o_" + name
	}
	return name

}
