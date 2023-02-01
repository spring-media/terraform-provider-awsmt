package awsmt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &awsmtProvider{}
)

func New() provider.Provider {
	return &awsmtProvider{}
}

type awsmtProvider struct{}

func (p *awsmtProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awsmt"
}

func (p *awsmtProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *awsmtProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *awsmtProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *awsmtProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
