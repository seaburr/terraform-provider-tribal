package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdminSettings_basic(t *testing.T) {
	resourceName := "tribal_admin_settings.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			// Create (apply initial settings)
			{
				Config: providerConfig() + testAdminSettingsConfig("[30, 14, 7]", 9, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "reminder_days.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "reminder_days.0", "30"),
					resource.TestCheckResourceAttr(resourceName, "reminder_days.1", "14"),
					resource.TestCheckResourceAttr(resourceName, "reminder_days.2", "7"),
					resource.TestCheckResourceAttr(resourceName, "notify_hour", "9"),
					resource.TestCheckResourceAttr(resourceName, "alert_on_overdue", "false"),
				),
			},
			// Update
			{
				Config: providerConfig() + testAdminSettingsConfig("[60, 30, 7, 1]", 12, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "reminder_days.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "reminder_days.0", "60"),
					resource.TestCheckResourceAttr(resourceName, "notify_hour", "12"),
					resource.TestCheckResourceAttr(resourceName, "alert_on_overdue", "true"),
				),
			},
		},
	})
}

func TestAccAdminSettings_withSlackWebhook(t *testing.T) {
	resourceName := "tribal_admin_settings.with_webhook"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "tribal_admin_settings" "with_webhook" {
  reminder_days   = [30, 7]
  notify_hour     = 8
  slack_webhook   = "https://hooks.slack.com/services/fake/org/webhook"
  alert_on_overdue = false
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "slack_webhook", "https://hooks.slack.com/services/fake/org/webhook"),
				),
			},
		},
	})
}

func testAdminSettingsConfig(reminderDays string, notifyHour int, alertOnOverdue bool) string {
	alertStr := "false"
	if alertOnOverdue {
		alertStr = "true"
	}
	return `
resource "tribal_admin_settings" "test" {
  reminder_days    = ` + reminderDays + `
  notify_hour      = ` + itoa(notifyHour) + `
  alert_on_overdue = ` + alertStr + `
}
`
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
