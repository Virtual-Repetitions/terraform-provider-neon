package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNeonProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNeonProjectResourceConfig(randomProjectName()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neon_project.test", "id"),
				),
			},

			// Update and Read testing
			{
				Config: testAccNeonProjectResourceConfig("updated-project-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("neon_project.test", "name", "updated-project-name"),
				),
			},
		},
	})
}

func testAccNeonProjectResourceConfig(projectName string) string {

	config := fmt.Sprintf(`
	provider "neon" { }
	resource "neon_project" "test" {
		name = "%s"
		instance_handle = "scalable"
		platform_id = "aws"
		region_id = "aws-us-west-2"
	}
`, projectName)
	return config
}

func randomProjectName() string {
	return fmt.Sprintf("Test Project %d", rand.Intn(10000))
}
