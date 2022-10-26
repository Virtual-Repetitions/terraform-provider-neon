package provider

import (
	"context"
	"os"
	"virtual-repetitions/terraform-provider-neon/src/neonApi"

	reqPkg "github.com/imroc/req/v3"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure NeonProvider satisfies various provider interfaces.
var _ provider.Provider = &NeonProvider{}
var _ provider.ProviderWithMetadata = &NeonProvider{}

// NeonProvider defines the provider implementation.
type NeonProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NeonProviderModel describes the provider data model.
type providerModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *NeonProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "neon"
	resp.Version = p.version
}

func (p *NeonProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Neon API key. This can be generated at https://console.neon.tech/app/settings/account",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *NeonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Neon API Key",
			"The provider cannot create the Neon API client as there is an unknown configuration value for the Neon API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NEON_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	neonApiKey := os.Getenv("NEON_API_KEY")

	if !config.ApiKey.IsNull() {
		neonApiKey = config.ApiKey.Value
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if neonApiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Neon API key",
			"The provider cannot create the Neon API client as there is a missing or empty value for the Neon API key. "+
				"Set the api_key value in the configuration or use the NEON_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Neon client using the configuration values
	client := neonApi.NewNeonApiClient(reqPkg.C(), neonApiKey)

	// Make the Neon client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *NeonProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNeonBranchResource,
		NewNeonProjectResource,
	}
}

func (p *NeonProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NeonProvider{
			version: version,
		}
	}
}
