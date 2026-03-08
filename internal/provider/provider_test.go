package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	provider "github.com/cberndt/tribal-tf-provider/internal/provider"
)

const (
	testAPIKey = "tribal_sk_7aaae0c48fb0c2cc119022a57d3229bc06a50583e097bc3cb08dcb0ba19ccd4f"
	testHost   = "http://localhost:8000"
)

func testProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"tribal": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

func providerConfig() string {
	return `
provider "tribal" {
  host    = "` + testHost + `"
  api_key = "` + testAPIKey + `"
}
`
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestAccProvider verifies the provider initializes and can manage admin settings.
func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "tribal_admin_settings" "smoke" {
  reminder_days    = [30, 7]
  notify_hour      = 9
  alert_on_overdue = false
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tribal_admin_settings.smoke", "notify_hour", "9"),
				),
			},
		},
	})
}
