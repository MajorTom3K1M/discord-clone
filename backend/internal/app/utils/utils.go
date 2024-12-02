package utils

import "strings"

func ExtractBaseDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".") // e.g., jkrn.me
	}
	return host
}
