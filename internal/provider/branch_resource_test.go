package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNeonBranchResourceIsCreated(t *testing.T) {
	var branchID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccNeonBranchResourceConfig("neon_project.test_parent_initial.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccNeonBranchID(&branchID),
					resource.TestCheckResourceAttrSet("neon_branch.test", "id"),
				),
			},

			// Tests that when parent_project_id changes the previous branch is destroyed and a new one exists
			{
				Config: testAccNeonBranchResourceConfig("neon_project.test_parent_updated.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBranchRecreated(&branchID),
					// Ideally we would check that the ID has changed but I couldn't figure out how to do that with the testing framework
					resource.TestCheckResourceAttrSet("neon_branch.test", "id"),
				),
			},
		},
	})
}

func testAccNeonBranchResourceConfig(parentProjectIDReference string) string {

	config := fmt.Sprintf(`
	provider "neon" { }
	resource "neon_project" "test_parent_initial" {
		name = "test-branches-parent-project-initial"
		instance_handle = "scalable"
		platform_id = "aws"
		region_id = "aws-us-west-2"
	}

	resource "neon_project" "test_parent_updated" {
		name = "test-branches-parent-project-updated"
		instance_handle = "scalable"
		platform_id = "aws"
		region_id = "aws-us-west-2"
	}

	resource "neon_branch" "test" {
		parent_project_id = %s
	}
`, parentProjectIDReference)
	return config
}

func testAccNeonBranchID(branchID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fetchedBranchID, err := testAccGetNeonBranchID("neon_branch.test", s)
		if err != nil {
			return err
		}

		*branchID = fetchedBranchID

		return nil
	}
}

func testAccCheckBranchRecreated(previousBranchID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		updatedBranchID, err := testAccGetNeonBranchID("neon_branch.test", s)
		if err != nil {
			return err
		}

		if updatedBranchID == *previousBranchID {
			return fmt.Errorf("Branch ID has not changed. resource: neon_branch.test")
		}

		return nil
	}
}

func testAccGetNeonBranchID(resourceName string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Resource not found. resource: %s", resourceName)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("Branch ID not set. resource: %s", resourceName)
	}

	return rs.Primary.ID, nil
}
