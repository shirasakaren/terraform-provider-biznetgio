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
// wip 728
