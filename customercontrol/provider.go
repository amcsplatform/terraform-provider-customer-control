package customercontrol

import (
	"context"

	"dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("CUSTOMERCONTROL_URL", "https://customercontrol-dev.amcsgroup.io"),
			},
			"privateKey": &schema.Schema{
				Type: schema.TypeString,
				Optional: false,
			},
		},
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	privateKey := d.Get("privateKey").(string)

	var diags diag.Diagnostics

	if (url != "") && (privateKey != "") {
		client, err := customercontrolclient.NewClient(&url, &privateKey)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary: "Unable to create CustomerControl client",
			})

			return nil, diags
		}

		return client, diags
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: "CustomerControl URL and privateKey must be provided",
		})
	}

	return nil, diags
}
