package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	goclient "github.com/appclacks/cli/client"
	apitypes "github.com/appclacks/go-types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccResourceHealthcheckDNSConfigCreate = `
resource "appclacks_healthcheck_dns" "check_dns" {
  name = "tf_acc_dns"
  description = "bar"
  labels = {
    "check": "dns",
  }
  interval = "35s"
  timeout = "7s"
  domain = "google.com"
  expected_ips = ["10.0.0.1"]
}
`

	testAccResourceHealthcheckDNSConfigUpdate = `
resource "appclacks_healthcheck_dns" "check_dns" {
  name = "tf_acc_dns2"
  description = "new_description"
  labels = {
    "check": "dns2",
  }
  interval = "36s"
  timeout = "10s"
  domain = "google.fr"
  expected_ips = ["10.0.0.2"]
}
`
)

func TestAccResourceHealthcheckDNS(t *testing.T) {
	check := new(apitypes.Healthcheck)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckResourceDNSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHealthcheckDNSConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_dns.check_dns", check),
					testAccCheckResourceCheck(check),
					testAccCheckResourceDNSCheckAttributes(testAttrs{
						"name":           validateString("tf_acc_dns"),
						"labels.check":   validateString("dns"),
						"description":    validateString("bar"),
						"expected_ips.0": validateString("10.0.0.1"),
						"interval":       validateString("35s"),
						"timeout":        validateString("7s"),
						"domain":         validateString("google.com"),
					}),
				),
			},
			{
				Config: testAccResourceHealthcheckDNSConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_dns.check_dns", check),

					testAccCheckResourceCheck(check),
					testAccCheckResourceDNSCheckAttributes(testAttrs{
						"name":           validateString("tf_acc_dns2"),
						"description":    validateString("new_description"),
						"expected_ips.0": validateString("10.0.0.2"),
						"labels.check":   validateString("dns2"),

						"interval": validateString("36s"),
						"timeout":  validateString("10s"),
						"domain":   validateString("google.fr"),
					}),
				),
			},
			{
				ResourceName:      "appclacks_healthcheck_dns.check_dns",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					return checkResourceAttributes(
						testAttrs{
							"name":           validateString("tf_acc_dns2"),
							"description":    validateString("new_description"),
							"labels.check":   validateString("dns2"),
							"expected_ips.0": validateString("10.0.0.2"),
							"interval":       validateString("36s"),
							"timeout":        validateString("10s"),
							"domain":         validateString("google.fr"),
						},
						s[0].Attributes)
				},
			},
		},
	})
}

func testAccCheckResourceDNSCheckAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "appclacks_healthcheck_dns" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("resource not found in the state")
	}
}

func testAccCheckResourceDNSCheckDestroy(s *terraform.State) error {
	client := GetAppclacksClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appclacks_healthcheck_dns" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := client.GetHealthcheck(ctx, apitypes.GetHealthcheckInput{
			Identifier: rs.Primary.ID,
		})

		if err != nil {
			if errors.Is(err, goclient.ErrNotFound) {
				return nil
			}
			return err
		}
	}

	return errors.New("DNS Health check still exists")
}
