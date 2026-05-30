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
