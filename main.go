package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/josuablejeru/terraform-provider-epilot/epilot"
)

func main() {
	tfsdk.Serve(context.Background(), epilot.New, tfsdk.ServeOpts{
		Name: "epilot",
})
}