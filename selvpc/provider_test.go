package selvpc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"selvpc": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(_ *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccSelVPCPreCheck(t *testing.T) {
	if v := os.Getenv("SEL_TOKEN"); v == "" {
		t.Fatal("SEL_TOKEN must be set for acceptance tests")
	}
}
