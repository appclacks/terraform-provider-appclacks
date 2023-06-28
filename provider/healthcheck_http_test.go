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
	testAccResourceHealthcheckHTTPConfigCreate = `
resource "appclacks_healthcheck_http" "check_http" {
  name = "tf_acc_http"
  description = "bar"
  labels = {
    "check": "http",
  }
  interval = "35s"
  timeout = "7s"
  target = "google.com"
  port = 443
  enabled = false
  body = "foo"
  body_regexp = ["foo"]
  valid_status = [200]
  headers = {
    "a": "b"
  }
  query = {
    "p1": "p2"
  }
  protocol = "http"
  path = "/foo"
  cert = "/tmp/foocert"
  cacert = "/tmp/foocacert"
  key = "/tmp/fookey"
}
`

	testAccResourceHealthcheckHTTPConfigUpdate = `
resource "appclacks_healthcheck_http" "check_http" {
  name = "tf_acc_http2"
  description = "new_description"
  labels = {
    "check": "http2",
  }
  interval = "36s"
  timeout = "10s"
  target = "google.fr"
  port = 80
  enabled = false
  body = "bar"
  body_regexp = ["foo2"]
  valid_status = [201]
  headers = {
    "a": "b2"
  }
  query = {
    "p3": "p4"
  }
  protocol = "https"
  path = "/bar"
  cert = "/tmp/foocert2"
  cacert = "/tmp/foocacert2"
  key = "/tmp/fookey2"
}
`
)

func TestAccResourceHealthcheckHTTP(t *testing.T) {
	check := new(apitypes.Healthcheck)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckResourceHTTPCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHealthcheckHTTPConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_http.check_http", check),
					testAccCheckResourceCheck(check),
					testAccCheckResourceHTTPCheckAttributes(testAttrs{
						"name":           validateString("tf_acc_http"),
						"description":    validateString("bar"),
						"interval":       validateString("35s"),
						"body_regexp.0":  validateString("foo"),
						"valid_status.0": validateString("200"),
						"timeout":        validateString("7s"),
						"target":         validateString("google.com"),
						"protocol":       validateString("http"),
						"labels.check":   validateString("http"),
						"headers.a":      validateString("b"),
						"query.p1":       validateString("p2"),
						"path":           validateString("/foo"),
						"body":           validateString("foo"),
						"cert":           validateString("/tmp/foocert"),
						"key":            validateString("/tmp/fookey"),
						"cacert":         validateString("/tmp/foocacert"),
					}),
				),
			},
			{
				Config: testAccResourceHealthcheckHTTPConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_http.check_http", check),

					testAccCheckResourceCheck(check),
					testAccCheckResourceHTTPCheckAttributes(testAttrs{
						"name":           validateString("tf_acc_http2"),
						"description":    validateString("new_description"),
						"interval":       validateString("36s"),
						"valid_status.0": validateString("201"),
						"timeout":        validateString("10s"),
						"body_regexp.0":  validateString("foo2"),
						"target":         validateString("google.fr"),
						"labels.check":   validateString("http2"),
						"headers.a":      validateString("b2"),
						"query.p3":       validateString("p4"),
						"protocol":       validateString("https"),
						"path":           validateString("/bar"),
						"body":           validateString("bar"),
						"cert":           validateString("/tmp/foocert2"),
						"key":            validateString("/tmp/fookey2"),
						"cacert":         validateString("/tmp/foocacert2"),
					}),
				),
			},
			{
				ResourceName:      "appclacks_healthcheck_http.check_http",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					return checkResourceAttributes(
						testAttrs{
							"name":           validateString("tf_acc_http2"),
							"description":    validateString("new_description"),
							"interval":       validateString("36s"),
							"timeout":        validateString("10s"),
							"valid_status.0": validateString("201"),
							"body_regexp.0":  validateString("foo2"),
							"labels.check":   validateString("http2"),
							"target":         validateString("google.fr"),
							"headers.a":      validateString("b2"),
							"query.p3":       validateString("p4"),
							"protocol":       validateString("https"),
							"path":           validateString("/bar"),
							"body":           validateString("bar"),
							"cert":           validateString("/tmp/foocert2"),
							"key":            validateString("/tmp/fookey2"),
							"cacert":         validateString("/tmp/foocacert2"),
						},
						s[0].Attributes)
				},
			},
		},
	})
}

func testAccCheckResourceHTTPCheckAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "appclacks_healthcheck_http" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("resource not found in the state")
	}
}

func testAccCheckResourceHTTPCheckDestroy(s *terraform.State) error {
	client := GetAppclacksClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appclacks_healthcheck_http" {
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

	return errors.New("HTTP Health check still exists")
}
