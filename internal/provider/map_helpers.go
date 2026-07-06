package provider

import (
	"encoding/json"
	"strconv"
	"strings"
)

// aliasStr cari value string pake beberapa kandidat key, case-insensitive.
func aliasStr(v map[string]any, keys ...string) string {
	for k, x := range v {
		for _, want := range keys {
			if strings.EqualFold(k, want) {
				if s, ok := x.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}

// aliasInt cari value int64 pake beberapa kandidat key (number atau string numerik).
func aliasInt(v map[string]any, keys ...string) int64 {
	for k, x := range v {
		for _, want := range keys {
			if !strings.EqualFold(k, want) {
				continue
			}
			switch n := x.(type) {
			case json.Number:
				i, err := n.Int64()
				if err == nil {
					return i
				}
			case float64:
				return int64(n)
			case string:
				i, err := strconv.ParseInt(n, 10, 64)
				if err == nil {
					return i
				}
			}
		}
	}
	return 0
}

// rawJSON marshal map ke string JSON pake redactJSON, fallback empty string.
func rawJSON(v map[string]any) string {
	if v == nil {
		return ""
	}
	return string(redactJSON(v))
}

// lowerStatus ambil status/state lalu lower-case biar polling case-insensitive.
func lowerStatus(v map[string]any) string {
	return strings.ToLower(aliasStr(v, "status", "state"))
}
// wip 884
