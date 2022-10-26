package neonApi

import (
	"testing"
)

// TestBranchCreate verifies Neon project branch can be created
func TestBranchCreate(t *testing.T) {
	projectFixture := NewProjectFixture(t, true)

	neonApiClient := NewNeonApiClientFixture()

	result, err := neonApiClient.SetDebug(false).BranchCreate(projectFixture.ID, NewDefaultNeonApiClientOptionsFixture())
	if err != nil {
		t.Error(err)
	}

	if projectFixture.ID != result.Project.ParentID {
		t.Errorf("Expected parent project ID %s, got %s", projectFixture.ID, result.Project.ParentID)
	}

	if result.Project.ID != "" {
		ProjectFixtureDelete(t, result.Project.ID)
	}
}
