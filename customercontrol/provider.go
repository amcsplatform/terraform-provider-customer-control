package customercontrol

import (
	"context"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CUSTOMERCONTROL_URL", "https://customercontrol-dev.amcsgroup.io"),
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"customercontrol_haproxy_rule": resourceHAProxyRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"customercontrol_haproxy_rule": dataSourceHAProxyRule(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	privateKey := d.Get("private_key").(string)

	var diags diag.Diagnostics
	var client *cc.CustomerControlClient
	var err error

	if (url != "") && (privateKey != "") {
		client, err = cc.NewClient(&url, &privateKey)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create CustomerControl client",
				Detail: err.Error(),
			})

			return nil, diags
		}

		return client, diags
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "CustomerControl URL and privateKey must be provided",
		})
	}

	return nil, diags
}
