package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/appclacks/cli/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	testAccProviders map[string]func() (*schema.Provider, error)
	testAccProvider  *schema.Provider
)

type testAttrs map[string]schema.SchemaValidateDiagFunc

func init() {

	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"appclacks": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	org := os.Getenv("APPCLACKS_ORGANIZATION_ID")
	token := os.Getenv("APPCLACKS_TOKEN")
	if org == "" || token == "" {
		msg := "APPCLACKS_ORGANIZATION_ID and APPCLACKS_TOKEN must be set, or the Appclacks client configuration file should be configured for acceptance tests"
		configPath, err := client.GetConfigFilePath()
		if err != nil {
			t.Fatal(msg, err)
		}
		config, err := client.ReadConfig(configPath)
		if err != nil {
			t.Fatal(msg, err)
		}
		if len(config.Profiles) == 0 {
			t.Fatal("APPCLACKS_ORGANIZATION_ID and APPCLACKS_TOKEN must be set for acceptance tests")
		}
	}
}

func checkResourceAttributes(want testAttrs, got map[string]string) error {
	for attr, validateFunc := range want {
		v, ok := got[attr]
		if !ok {
			return fmt.Errorf("expected attribute %q not found in map", attr)
		} else if diags := validateFunc(v, cty.GetAttrPath(attr)); diags.HasError() {
			errors := make([]string, 0)
			for _, d := range diags {
				if d.Severity == diag.Error {
					errors = append(errors, d.Summary)
				}
			}

			return fmt.Errorf("invalid value for attribute %q:\n%s\n", // nolint:revive
				attr, strings.Join(errors, "\n"))
		}
	}

	return nil
}

func validateString(str string) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (s []string, es []error) {
		value, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if value != str {
			es = append(es, fmt.Errorf("string %q doesn't match expected value %q", value, str))
			return
		}

		return
	})
}
