package epilot

import (
	"context"
	"os"

	"github.com/josuablejeru/terraform-provider-epilot/epilot-webhooks-client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const server = "https://webhooks.staging.sls.epilot.io"

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *webhooks.ClientWithResponses
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"orgId": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"apiKey": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	OrgId types.String `tfsdk:"orgId"`
	ApiKey types.String `tfsdk:"apiKey"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a orgId to the provider
	var orgId string
	if config.OrgId.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as OrgId",
		)
		return
	}

	if config.OrgId.Null {
		orgId = os.Getenv("EPILOT_ORG_ID")
	} else {
		orgId = config.OrgId.Value
	}

	if orgId == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find OrgId",
			"Username cannot be an empty string",
		)
		return
	}

	// User must provide a password to the provider
	var apiKey string
	if config.ApiKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as ApiKey",
		)
		return
	}

	if config.ApiKey.Null {
		apiKey = os.Getenv("EPILOT_API_KEY")
	} else {
		apiKey = config.ApiKey.Value
	}

	if apiKey == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find apiKey",
			"apiKey cannot be an empty string",
		)
		return
	}

	// Create a new epilot client and set it to the provider client
	c, err := webhooks.NewClientWithResponses(server)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create epilot client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"epilot_webhook": resourceWebhookType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"hashicups_coffees": dataSourceWebhookType{},
	}, nil
}