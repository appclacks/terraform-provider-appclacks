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
	testAccResourceHealthcheckTLSConfigCreate = `
resource "appclacks_healthcheck_tls" "check_tls" {
  name = "tf_acc_tls"
  description = "bar"
  labels = {
    "check": "tls",
  }
  interval = "35s"
  timeout = "7s"
  target = "google.com"
  port = 443
  enabled = false
  cert = "/tmp/foocert"
  cacert = "/tmp/foocacert"
  key = "/tmp/fookey"
  insecure = true
  server_name = "google.com"
  expiration_delay = "10s"
}
`

	testAccResourceHealthcheckTLSConfigUpdate = `
resource "appclacks_healthcheck_tls" "check_tls" {
  name = "tf_acc_tls2"
  description = "new_description"
  labels = {
    "check": "tls2",
  }
  interval = "36s"
  timeout = "10s"
  target = "google.fr"
  port = 80
  enabled = false
  cert = "/tmp/foocert2"
  cacert = "/tmp/foocacert2"
  key = "/tmp/fookey2"
  insecure = false
  server_name = "google.fr"
  expiration_delay = "20s"
}
`
)

func TestAccResourceHealthcheckTLS(t *testing.T) {
	check := new(goclient.Healthcheck)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckResourceTLSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHealthcheckTLSConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_tls.check_tls", check),
					testAccCheckResourceCheck(check),
					testAccCheckResourceTLSCheckAttributes(testAttrs{
						"insecure":         validateString("true"),
						"name":             validateString("tf_acc_tls"),
						"port":             validateString("443"),
						"description":      validateString("bar"),
						"interval":         validateString("35s"),
						"timeout":          validateString("7s"),
						"target":           validateString("google.com"),
						"labels.check":     validateString("tls"),
						"cert":             validateString("/tmp/foocert"),
						"key":              validateString("/tmp/fookey"),
						"cacert":           validateString("/tmp/foocacert"),
						"server_name":      validateString("google.com"),
						"expiration_delay": validateString("10s"),
					}),
				),
			},
			{
				Config: testAccResourceHealthcheckTLSConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_tls.check_tls", check),

					testAccCheckResourceCheck(check),
					testAccCheckResourceTLSCheckAttributes(testAttrs{
						"insecure":         validateString("false"),
						"name":             validateString("tf_acc_tls2"),
						"description":      validateString("new_description"),
						"port":             validateString("80"),
						"interval":         validateString("36s"),
						"timeout":          validateString("10s"),
						"target":           validateString("google.fr"),
						"labels.check":     validateString("tls2"),
						"cert":             validateString("/tmp/foocert2"),
						"key":              validateString("/tmp/fookey2"),
						"cacert":           validateString("/tmp/foocacert2"),
						"server_name":      validateString("google.fr"),
						"expiration_delay": validateString("20s"),
					}),
				),
			},
			{
				ResourceName:      "appclacks_healthcheck_tls.check_tls",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					return checkResourceAttributes(
						testAttrs{
							"insecure":         validateString("false"),
							"name":             validateString("tf_acc_tls2"),
							"description":      validateString("new_description"),
							"port":             validateString("80"),
							"interval":         validateString("36s"),
							"timeout":          validateString("10s"),
							"target":           validateString("google.fr"),
							"labels.check":     validateString("tls2"),
							"cert":             validateString("/tmp/foocert2"),
							"key":              validateString("/tmp/fookey2"),
							"cacert":           validateString("/tmp/foocacert2"),
							"server_name":      validateString("google.fr"),
							"expiration_delay": validateString("20s"),
						},
						s[0].Attributes)
				},
			},
		},
	})
}

func testAccCheckResourceTLSCheckAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "appclacks_healthcheck_tls" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("resource not found in the state")
	}
}

func testAccCheckResourceTLSCheckDestroy(s *terraform.State) error {
	client := GetAppclacksClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appclacks_healthcheck_tls" {
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

	return errors.New("TLS Health check still exists")
}
