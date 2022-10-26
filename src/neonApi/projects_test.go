package neonApi

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestProjectCreate verifies Neon project can be created
func TestProjectCreate(t *testing.T) {

	projectName := fmt.Sprintf("Test Project %d", rand.Intn(10000))
	neonApiClient := NewNeonApiClientFixture()

	createData := NeonProjectCreateData{
		Project: NeonProjectCreateProjectAttributes{
			InstanceHandle: "scalable",
			Name:           projectName,
			PlatformID:     "aws",
			RegionID:       "aws-us-west-2",
			Settings:       map[string]string{},
		},
	}

	result, err := neonApiClient.SetDebug(false).ProjectCreate(createData, NewDefaultNeonApiClientOptionsFixture())

	if err != nil {
		t.Error(err)
	}

	ProjectFixtureDelete(t, result.Project.ID)

	if createData.Project.Name != result.Project.Name {
		t.Errorf("Expected fields to be set during creation. expected field values of %+v, got %+v", createData, result)
	}

	if createData.Project.RegionID != result.Project.RegionID {
		t.Errorf("Expected input region ID to be normalized during creation. expected %s, got %s", createData.Project.RegionID, result.Project.RegionID)
	}
}

// TestProjectRead verifies Neon project can be read
func TestProjectRead(t *testing.T) {
	projectFixture := NewProjectFixture(t, true)

	neonApiClient := NewNeonApiClientFixture()

	project, err := neonApiClient.SetDebug(false).ProjectRead(projectFixture.ID, NewDefaultNeonApiClientOptionsFixture())
	if err != nil {
		t.Error(err)
	}

	if project.ID != projectFixture.ID {
		t.Errorf("Expected project ID %s, got %s", projectFixture.ID, project.ID)
	}
}

// TestProjectUpdate verifies Neon project can be read
func TestProjectUpdate(t *testing.T) {
	projectFixture := NewProjectFixture(t, true)

	neonApiClient := NewNeonApiClientFixture()

	updateData := NeonProjectUpdateData{
		Project: NeonProjectUpdateProjectAttributes{
			InstanceTypeID: "1",
			Name:           "updated-project-name",
			Settings:       map[string]string{},
		},
	}
	result, err := neonApiClient.SetDebug(false).ProjectUpdate(projectFixture.ID, updateData, NewDefaultNeonApiClientOptionsFixture())
	if err != nil {
		t.Error(err)
	}

	if result.Project.ID != projectFixture.ID {
		t.Errorf("Expected project ID %s, got %s", projectFixture.ID, result.Project.ID)
	}
}

// TestProjectRead verifies Neon project can be read
func TestProjectReadOnInvalidProject(t *testing.T) {
	projectID := "some-nonexistent-project"
	neonApiClient := NewNeonApiClientFixture()

	_, err := neonApiClient.ProjectRead(projectID, NewDefaultNeonApiClientOptionsFixture())
	if err == nil {
		t.Error("Expected to receive error, got nil error instead.")
	}
}

// TestProjectDelete verifies Neon project is deleted
func TestProjectDelete(t *testing.T) {
	projectFixture := NewProjectFixture(t, false)
	neonApiClient := NewNeonApiClientFixture()

	err := neonApiClient.ProjectDelete(projectFixture.ID, NewDefaultNeonApiClientOptionsFixture())
	if err != nil {
		t.Errorf("Could not delete project. err: %s", err)
	}

	project, err := neonApiClient.SetDebug(false).ProjectRead(projectFixture.ID, NewDefaultNeonApiClientOptionsFixture())
	if project.ID != "" {
		t.Errorf("Project was not deleted.")
	}
}
