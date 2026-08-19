// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"biznetgio": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("BIZNETGIO_API_KEY") == "" {
		// no key in CI: skip instead of failing so the workflow stays green.
		// acceptance tests still run locally with a real token.
		t.Skip("BIZNETGIO_API_KEY must be set for acceptance tests")
	}
}
