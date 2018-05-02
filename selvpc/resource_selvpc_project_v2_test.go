package selvpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/projects"
)

func TestAccResellV2ProjectBasic(t *testing.T) {
	var project projects.Project
	projectName := acctest.RandomWithPrefix("tf-acc")
	projectNameUpdated := acctest.RandomWithPrefix("tf-acc-updated")
	projectCustomURL := acctest.RandomWithPrefix("tf-acc-url") + ".selvpc.ru"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccSelVPCPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResellV2ProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResellV2ProjectBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResellV2ProjectExists("selvpc_project_v2.project_tf_acc_test_1", &project),
					resource.TestCheckResourceAttr("selvpc_project_v2.project_tf_acc_test_1", "name", projectName),
				),
			},
			resource.TestStep{
				Config: testAccResellV2ProjectUpdate1(projectName, projectCustomURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "name", projectName),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "custom_url", projectCustomURL),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "theme.color", "000000"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "theme.logo", "fake.png"),
				),
			},
			resource.TestStep{
				Config: testAccResellV2ProjectUpdate2(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "name", projectName),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "custom_url", ""),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "theme.color", "FF0000"),
				),
			},
			resource.TestStep{
				Config: testAccResellV2ProjectUpdate3(projectNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "name", projectNameUpdated),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "custom_url", ""),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "theme.color", "5D6D7E"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.0.resource_name", "image_gigabytes"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.0.resource_quotas.0.region", "ru-1"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.0.resource_quotas.0.zone", ""),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.0.resource_quotas.0.value", "1"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_name", "volume_gigabytes_basic"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.0.region", "ru-1"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.0.zone", "ru-1a"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.0.value", "1"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.1.region", "ru-2"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.1.zone", "ru-2a"),
					resource.TestCheckResourceAttr(
						"selvpc_project_v2.project_tf_acc_test_1", "quotas.1.resource_quotas.1.value", "2"),
				),
			},
		},
	})
}

func testAccCheckResellV2ProjectDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	resellV2Client := config.resellV2Client()
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "selvpc_project_v2" {
			continue
		}

		_, _, err := projects.Get(ctx, resellV2Client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("project still exists")
		}
	}

	return nil
}

func testAccCheckResellV2ProjectExists(n string, project *projects.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		resellV2Client := config.resellV2Client()
		ctx := context.Background()

		foundProject, _, err := projects.Get(ctx, resellV2Client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundProject.ID != rs.Primary.ID {
			return fmt.Errorf("project not found")
		}

		*project = *foundProject

		return nil
	}
}

func testAccResellV2ProjectBasic(name string) string {
	return fmt.Sprintf(`
resource "selvpc_project_v2" "project_tf_acc_test_1" {
  name = "%s"
}`, name)
}

func testAccResellV2ProjectUpdate1(name, customURL string) string {
	return fmt.Sprintf(`
resource "selvpc_project_v2" "project_tf_acc_test_1" {
  name       = "%s"
  custom_url = "%s"
  theme {
    color = "000000"
    logo  = "fake.png"
  }
}`, name, customURL)
}

func testAccResellV2ProjectUpdate2(name string) string {
	return fmt.Sprintf(`
resource "selvpc_project_v2" "project_tf_acc_test_1" {
  name       = "%s"
  theme {
    color = "FF0000"
  }
}`, name)
}

func testAccResellV2ProjectUpdate3(name string) string {
	return fmt.Sprintf(`
resource "selvpc_project_v2" "project_tf_acc_test_1" {
  name = "%s"
  theme {
    color = "5D6D7E"
  }
  quotas = [
    {
      resource_name = "image_gigabytes"
      resource_quotas = [
        {
          region = "ru-1"
          value = 1
        }
      ]
    },
    {
      resource_name = "volume_gigabytes_basic"
      resource_quotas = [
        {
          region = "ru-1"
          zone = "ru-1a"
          value = 1
        },
        {
          region = "ru-2"
          zone = "ru-2a"
          value = 2
        }
      ]
    }
  ]
}`, name)
}
