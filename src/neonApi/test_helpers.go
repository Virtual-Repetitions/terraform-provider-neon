package neonApi

import (
	"testing"
)

func ProjectFixtureDelete(t *testing.T, projectID string) {
	neonApiClient := NewNeonApiClientFixture()

	err := neonApiClient.ProjectDelete(projectID, NewDefaultNeonApiClientOptionsFixture())

	if err != nil {
		t.Errorf("Could not delete project. err: %s", err)
	}
}
