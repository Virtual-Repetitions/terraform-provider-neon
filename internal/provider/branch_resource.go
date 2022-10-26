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
var _ resource.Resource = &NeonBranchResource{}
var _ resource.ResourceWithImportState = &NeonBranchResource{}

func NewNeonBranchResource() resource.Resource {
	return &NeonBranchResource{}
}

// NeonBranchResource defines the resource implementation.
type NeonBranchResource struct {
	client neonApi.NeonApiClient
}

// neonBranchResourceModel describes the resource data model.
type neonBranchResourceModel struct {
	ID              types.String `tfsdk:"id"`
	ParentProjectID types.String `tfsdk:"parent_project_id"`
}

func (r *NeonBranchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (r *NeonBranchResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Neon branch resource",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Branch ID",
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"parent_project_id": {
				Required: true,
				Type:     types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					// When parent_project_id changes, force recreation
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *NeonBranchResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NeonBranchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan neonBranchResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Neon branch resource.")

	result, err := r.client.BranchCreate(plan.ParentProjectID.Value, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating branch",
			"Could not create branch, unexpected error: "+err.Error(),
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

func (r *NeonBranchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state neonBranchResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	branch, err := r.client.ProjectRead(state.ID.Value, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading branch",
			"Could not read branch, unexpected error: "+err.Error(),
		)
		return
	}

	state = neonBranchResourceModel{
		ID:              state.ID,
		ParentProjectID: types.String{Value: branch.ParentID},
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// No updates allowed. See `parent_project_id` attribute.
func (r *NeonBranchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data neonBranchResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *NeonBranchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state neonBranchResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ProjectDelete(state.ID.Value, neonApi.NeonApiClientOptions{
		NumRetries: 0,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete branch, got error: %s", err))
		return
	}
}

func (r *NeonBranchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
