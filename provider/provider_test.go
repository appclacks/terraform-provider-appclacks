package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
	endpoint := os.Getenv("APPCLACKS_API_ENDPOINT")
	if endpoint == "" {
		msg := "APPCLACKS_API_ENDPOINT must be set to launch tests"
		t.Fatal(msg)
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
