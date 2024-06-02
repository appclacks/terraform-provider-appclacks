package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	goclient "github.com/appclacks/go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccResourceHealthcheckTCPConfigCreate = `
resource "appclacks_healthcheck_tcp" "check_tcp" {
  name = "tf_acc_tcp"
  description = "bar"
  labels = {
    "check": "tcp",
  }
  interval = "35s"
  timeout = "7s"
  target = "google.com"
  port = 443
}
`

	testAccResourceHealthcheckTCPConfigUpdate = `
resource "appclacks_healthcheck_tcp" "check_tcp" {
  name = "tf_acc_tcp2"
  description = "new_description"
  labels = {
    "check": "tcp2",
  }
  interval = "36s"
  timeout = "10s"
  target = "google.fr"
  port = 80
}
`
)

func TestAccResourceHealthcheckTCP(t *testing.T) {
	check := new(goclient.Healthcheck)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckResourceTCPCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHealthcheckTCPConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_tcp.check_tcp", check),
					testAccCheckResourceCheck(check),
					testAccCheckResourceTCPCheckAttributes(testAttrs{
						"name":         validateString("tf_acc_tcp"),
						"description":  validateString("bar"),
						"interval":     validateString("35s"),
						"port":         validateString("443"),
						"timeout":      validateString("7s"),
						"labels.check": validateString("tcp"),
						"target":       validateString("google.com"),
					}),
				),
			},
			{
				Config: testAccResourceHealthcheckTCPConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_tcp.check_tcp", check),

					testAccCheckResourceCheck(check),
					testAccCheckResourceTCPCheckAttributes(testAttrs{
						"name":         validateString("tf_acc_tcp2"),
						"port":         validateString("80"),
						"description":  validateString("new_description"),
						"interval":     validateString("36s"),
						"timeout":      validateString("10s"),
						"labels.check": validateString("tcp2"),
						"target":       validateString("google.fr"),
					}),
				),
			},
			{
				ResourceName:      "appclacks_healthcheck_tcp.check_tcp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					return checkResourceAttributes(
						testAttrs{
							"name":         validateString("tf_acc_tcp2"),
							"port":         validateString("80"),
							"description":  validateString("new_description"),
							"interval":     validateString("36s"),
							"timeout":      validateString("10s"),
							"labels.check": validateString("tcp2"),
							"target":       validateString("google.fr"),
						},
						s[0].Attributes)
				},
			},
		},
	})
}

func testAccCheckResourceTCPCheckAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "appclacks_healthcheck_tcp" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("resource not found in the state")
	}
}

func testAccCheckResourceTCPCheckDestroy(s *terraform.State) error {
	client := GetAppclacksClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appclacks_healthcheck_tcp" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := client.GetHealthcheck(ctx, goclient.GetHealthcheckInput{
			Identifier: rs.Primary.ID,
		})

		if err != nil {
			if errors.Is(err, goclient.ErrNotFound) {
				return nil
			}
			return err
		}
	}

	return errors.New("TCP Health check still exists")
}
