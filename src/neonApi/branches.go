package neonApi

import (
	"fmt"
	"time"
)

type NeonBranchCreateResult struct {
	Project  NeonProject
	Response NeonBranchCreateSuccessResponse
}

type NeonBranchCreateSuccessResponse struct {
	CreatedAt       time.Time         `json:"created_at"`
	CurrentState    string            `json:"current_state"`
	Deleted         bool              `json:"deleted"`
	ID              string            `json:"id"`
	InstanceHandle  string            `json:"instance_handle"`
	InstanceTypeID  string            `json:"instance_type_id"`
	MaxProjectSize  int               `json:"max_project_size"`
	Name            string            `json:"name"`
	ParentID        string            `json:"parent_id"`
	PendingState    string            `json:"pending_state"`
	PlatformID      string            `json:"platform_id"`
	PlatformName    string            `json:"platform_name"`
	PoolerEnabled   bool              `json:"pooler_enabled"`
	RegionID        string            `json:"region_id"`
	RegionName      string            `json:"region_name"`
	Settings        map[string]string `json:"settings"`
	AdditionalProp1 string            `json:"additionalProp1"`
	AdditionalProp2 string            `json:"additionalProp2"`
	AdditionalProp3 string            `json:"additionalProp3"`
	Size            int               `json:"size"`
	UpdatedAt       string            `json:"updated_at"`
}

func (client *NeonApiClient) BranchCreate(parentProjectID string, options NeonApiClientOptions) (NeonBranchCreateResult, error) {
	var response NeonBranchCreateSuccessResponse

	_, err := client.NewRequest().SetRetryCount(options.NumRetries).SetResult(&response).Post(fmt.Sprintf("/api/v1/projects/%s/branches", parentProjectID))

	if err != nil {
		return NeonBranchCreateResult{}, err
	}

	project := NeonProject{
		ID:             response.ID,
		Name:           response.Name,
		InstanceTypeID: response.InstanceTypeID,
		ParentID:       response.ParentID,
		PlatformID:     response.PlatformID,
		RegionID:       response.RegionID,
		Settings:       response.Settings,
	}

	return NeonBranchCreateResult{
		Project:  project,
		Response: response,
	}, err
}
