package neonApi

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/imroc/req/v3"
)

func NewNeonApiClientFixture() NeonApiClient {
	var authToken string
	if authToken = os.Getenv("NEON_API_KEY"); authToken == "" {
		panic("No Neon API authentication token provided. Please set environment variable `NEON_API_KEY`.")
	}

	return NewNeonApiClient(req.C(), authToken)
}

func NewDefaultNeonApiClientOptionsFixture() NeonApiClientOptions {

	return NeonApiClientOptions{
		NumRetries: 0,
	}
}

func NewProjectFixture(t *testing.T, cleanup bool) NeonProject {
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

	if cleanup {
		t.Cleanup(func() {
			ProjectFixtureDelete(t, result.Project.ID)
		})
	}

	return result.Project
}
