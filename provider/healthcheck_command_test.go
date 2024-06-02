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
	testAccResourceHealthcheckCommandConfigCreate = `
resource "appclacks_healthcheck_command" "check_command" {
  name = "tf_acc_command"
  description = "bar"
  labels = {
    "check": "command",
  }
  interval = "35s"
  timeout = "7s"
  command = "ls"
  arguments = ["/"]
}
`

	testAccResourceHealthcheckCommandConfigUpdate = `
resource "appclacks_healthcheck_command" "check_command" {
  name = "tf_acc_command2"
  description = "new_description"
  labels = {
    "check": "command2",
  }
  interval = "36s"
  timeout = "10s"
  command = "cat"
  arguments = ["/foo"]
}
`
)

func TestAccResourceHealthcheckCommand(t *testing.T) {
	check := new(goclient.Healthcheck)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckResourceCommandCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHealthcheckCommandConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_command.check_command", check),
					testAccCheckResourceCheck(check),
					testAccCheckResourceCommandCheckAttributes(testAttrs{
						"name":         validateString("tf_acc_command"),
						"description":  validateString("bar"),
						"interval":     validateString("35s"),
						"timeout":      validateString("7s"),
						"command":      validateString("ls"),
						"labels.check": validateString("command"),
						"arguments.0":  validateString("/"),
					}),
				),
			},
			{
				Config: testAccResourceHealthcheckCommandConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceCheckExists("appclacks_healthcheck_command.check_command", check),

					testAccCheckResourceCheck(check),
					testAccCheckResourceCommandCheckAttributes(testAttrs{
						"name":         validateString("tf_acc_command2"),
						"description":  validateString("new_description"),
						"interval":     validateString("36s"),
						"timeout":      validateString("10s"),
						"command":      validateString("cat"),
						"labels.check": validateString("command2"),
						"arguments.0":  validateString("/foo"),
					}),
				),
			},
			{
				ResourceName:      "appclacks_healthcheck_command.check_command",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					return checkResourceAttributes(
						testAttrs{
							"name":         validateString("tf_acc_command2"),
							"description":  validateString("new_description"),
							"interval":     validateString("36s"),
							"timeout":      validateString("10s"),
							"command":      validateString("cat"),
							"labels.check": validateString("command2"),
							"arguments.0":  validateString("/foo"),
						},
						s[0].Attributes)
				},
			},
		},
	})
}

func testAccCheckResourceCheckExists(name string, check *goclient.Healthcheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return errors.New("resource not found in the state")
		}

		if rs.Primary.ID == "" {
			return errors.New("resource ID not set")
		}

		client := GetAppclacksClient(testAccProvider.Meta())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, err := client.GetHealthcheck(ctx, goclient.GetHealthcheckInput{
			Identifier: rs.Primary.ID,
		})
		if err != nil {
			return err
		}
		check.ID = result.ID

		return nil
	}
}

func testAccCheckResourceCheck(check *goclient.Healthcheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if check.ID == "" {
			return errors.New("Health check ID is nil")
		}

		return nil
	}
}

func testAccCheckResourceCommandCheckAttributes(expected testAttrs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "appclacks_healthcheck_command" {
				continue
			}

			return checkResourceAttributes(expected, rs.Primary.Attributes)
		}

		return errors.New("resource not found in the state")
	}
}

func testAccCheckResourceCommandCheckDestroy(s *terraform.State) error {
	client := GetAppclacksClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appclacks_healthcheck_command" {
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

	return errors.New("Command Health check still exists")
}
