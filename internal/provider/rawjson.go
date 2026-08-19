package provider

import (
	"encoding/json"
	"strings"
)

// secretJSONKeys key yang nilainya harus di-mask sebelum masuk state/log.
var secretJSONKeys = []string{
	"cipassword", "console_password", "consolepassword", "password", "passwd",
	"private_key", "privatekey", "private", "secret_key", "secretkey", "secret",
	"pem", "token",
}

func isSecretJSONKey(k string) bool {
	for _, s := range secretJSONKeys {
		if strings.EqualFold(k, s) {
			return true
		}
	}
	return false
}

// redactMap salin map, mask key rahasia (case-insensitive), rekursif ke nested map.
func redactMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if isSecretJSONKey(k) {
			out[k] = "***"
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			out[k] = redactMap(nested)
			continue
		}
		out[k] = v
	}
	return out
}

// redactJSON marshal map dengan nilai rahasia di-mask, fallback nil.
func redactJSON(m map[string]any) []byte {
	b, err := json.Marshal(redactMap(m))
	if err != nil {
		return nil
	}
	return b
}
