package neonApi

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type NeonProjectMutationResult struct {
	Project  NeonProject
	Response NeonProjectMutationSuccessResponse
}

type NeonProjectMutationSuccessResponseDatabase struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OwnerID int    `json:"owner_id"`
}

type NeonProjectMutationSuccessResponseRole struct {
	Dsn      string `json:"dsn"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type NeonProjectMutationSuccessResponse struct {
	CreatedAt      time.Time                                    `json:"created_at"`
	CurrentState   string                                       `json:"current_state"`
	Databases      []NeonProjectMutationSuccessResponseDatabase `json:"databases"`
	Deleted        bool                                         `json:"deleted"`
	ID             string                                       `json:"id"`
	InstanceHandle string                                       `json:"instance_handle"`
	InstanceTypeID string                                       `json:"instance_type_id"`
	MaxProjectSize int                                          `json:"max_project_size"`
	Name           string                                       `json:"name"`
	ParentID       string                                       `json:"parent_id"`
	PendingState   string                                       `json:"pending_state"`
	PlatformID     string                                       `json:"platform_id"`
	PlatformName   string                                       `json:"platform_name"`
	PoolerEnabled  bool                                         `json:"pooler_enabled"`
	RegionID       string                                       `json:"region_id"`
	RegionName     string                                       `json:"region_name"`
	Roles          []NeonProjectMutationSuccessResponseRole     `json:"roles"`
	Settings       map[string]string                            `json:"settings" `
	Size           int                                          `json:"size"`
	UpdatedAt      time.Time                                    `json:"updated_at"`
}

type NeonProject struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	InstanceHandle string            `json:"instance_handle"`
	InstanceTypeID string            `json:"instance_type_id"`
	ParentID       string            `json:"parent_id"`
	PlatformID     string            `json:"platform_id"`
	RegionID       string            `json:"region_id"`
	Settings       map[string]string `json:"settings"`
}

type NeonProjectCreateData struct {
	Project NeonProjectCreateProjectAttributes `json:"project"`
}

type NeonProjectCreateProjectAttributes struct {
	InstanceHandle string            `json:"instance_handle"`
	Name           string            `json:"name"`
	PlatformID     string            `json:"platform_id"`
	RegionID       string            `json:"region_id"`
	Settings       map[string]string `json:"settings"`
}

type NeonProjectUpdateData struct {
	Project NeonProjectUpdateProjectAttributes `json:"project"`
}

type NeonProjectUpdateProjectAttributes struct {
	InstanceTypeID string            `json:"instance_type_id"`
	Name           string            `json:"name"`
	PoolerEnabled  bool              `json:"pooler_enabled"`
	Settings       map[string]string `json:"settings"`
}

func (client *NeonApiClient) ProjectCreate(data NeonProjectCreateData, options NeonApiClientOptions) (NeonProjectMutationResult, error) {
	var response NeonProjectMutationSuccessResponse

	normalizedRegionID, err := normalizeRegionID(data.Project.RegionID)
	if err != nil {
		return NeonProjectMutationResult{}, err
	}

	data.Project.RegionID = normalizedRegionID

	_, err = client.NewRequest().SetRetryCount(options.NumRetries).SetResult(&response).SetBody(data).Post(fmt.Sprintf("/api/v1/projects"))

	result := NeonProjectMutationResult{
		Project: NeonProject{
			ID:             response.ID,
			Name:           response.Name,
			InstanceHandle: response.InstanceHandle,
			InstanceTypeID: response.InstanceTypeID,
			ParentID:       response.ParentID,
			PlatformID:     response.PlatformID,
			RegionID:       response.RegionID,
			Settings:       response.Settings,
		},
		Response: response,
	}
	return result, err
}

func (client *NeonApiClient) ProjectUpdate(projectID string, data NeonProjectUpdateData, options NeonApiClientOptions) (NeonProjectMutationResult, error) {
	var response NeonProjectMutationSuccessResponse

	_, err := client.NewRequest().SetRetryCount(options.NumRetries).SetResult(&response).SetBody(data).Patch(fmt.Sprintf("/api/v1/projects/%s", projectID))

	if err != nil {
		fmt.Printf("Error updating project. body data: %+v err: %s projectID: %s", data, err, projectID)
		return NeonProjectMutationResult{}, err
	}

	result := NeonProjectMutationResult{
		Project: NeonProject{
			ID:             response.ID,
			Name:           response.Name,
			InstanceHandle: response.InstanceHandle,
			InstanceTypeID: response.InstanceTypeID,
			ParentID:       response.ParentID,
			PlatformID:     response.PlatformID,
			RegionID:       response.RegionID,
			Settings:       response.Settings,
		},
		Response: response,
	}
	return result, err
}

func (client *NeonApiClient) ProjectDelete(projectID string, options NeonApiClientOptions) error {

	var project NeonProject
	_, err := client.NewRequest().SetRetryCount(options.NumRetries).SetResult(&project).Post(fmt.Sprintf("/api/v1/projects/%s/delete", projectID))

	return err
}

func (client *NeonApiClient) ProjectRead(projectID string, options NeonApiClientOptions) (NeonProject, error) {

	var project NeonProject
	_, err := client.NewRequest().SetRetryCount(options.NumRetries).SetResult(&project).Get(fmt.Sprintf("/api/v1/projects/%s", projectID))

	return project, err
}

// Neon API expects region IDs for input data to be of form `us-west-2` but returns `aws-us-west-2` in response data
// This supports `aws-us-west-2` for consistency
func normalizeRegionID(regionID string) (string, error) {
	re := regexp.MustCompile("[A-Za-z]+-(?P<regionID>.*)")
	match := re.FindStringSubmatch(regionID)

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	normalizedRegionID, valid := result["regionID"]

	if !valid {
		return "", errors.New(fmt.Sprintf("Could not parse region ID. Expected to be of form `aws-us-west-2`. given: %s", regionID))
	}

	return normalizedRegionID, nil
}
