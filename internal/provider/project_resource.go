package provider

import (
	"context"
	"fmt"
	"virtual-repetitions/terraform-provider-neon/src/neonApi"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NeonProjectResource{}
var _ resource.ResourceWithImportState = &NeonProjectResource{}

func NewNeonProjectResource() resource.Resource {
	return &NeonProjectResource{}
}

// NeonProjectResource defines the resource implementation.
type NeonProjectResource struct {
	client neonApi.NeonApiClient
}

// neonProjectResourceModel describes the resource data model.
type neonProjectResourceModel struct {
	ID             types.String `tfsdk:"id"`
	InstanceHandle types.String `tfsdk:"instance_handle"`
	Name           types.String `tfsdk:"name"`
	PlatformID     types.String `tfsdk:"platform_id"`
	RegionID       types.String `tfsdk:"region_id"`
}

func (r *NeonProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *NeonProjectResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Neon project resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Project ID",
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"instance_handle": {
				Required: true,
				Type:     types.StringType,
			},
			"name": {
				Required: true,
				Type:     types.StringType,
			},
			"platform_id": {
				Required: true,
				Type:     types.StringType,
			},
			"region_id": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (r *NeonProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(neonApi.NeonApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected neonApi.NeonApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NeonProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan neonProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Neon project resource.")

	result, err := r.client.ProjectCreate(neonApi.NeonProjectCreateData{
		Project: neonApi.NeonProjectCreateProjectAttributes{
			InstanceHandle: plan.InstanceHandle.Value,
			Name:           plan.Name.Value,
			PlatformID:     plan.PlatformID.Value,
			RegionID:       plan.RegionID.Value,
			Settings:       map[string]string{},
		},
	}, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.String{Value: result.Project.ID}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NeonProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state neonProjectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.ProjectRead(state.ID.Value, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			"Could not read project, unexpected error: "+err.Error(),
		)
		return
	}

	state = neonProjectResourceModel{
		ID:             state.ID,
		InstanceHandle: types.String{Value: project.InstanceHandle},
		Name:           types.String{Value: project.Name},
		PlatformID:     types.String{Value: project.PlatformID},
		RegionID:       types.String{Value: project.RegionID},
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NeonProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data neonProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.ProjectUpdate(data.ID.Value, neonApi.NeonProjectUpdateData{
		Project: neonApi.NeonProjectUpdateProjectAttributes{
			Name: data.Name.Value,
		},
	}, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project, unexpected error: "+err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NeonProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state neonProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ProjectDelete(state.ID.Value, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
		return
	}
}

func (r *NeonProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
