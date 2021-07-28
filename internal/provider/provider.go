package provider

import (
	"context"
	"fmt"
	"strings"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func Provider() func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("CUSTOMERCONTROL_URL", "https://customercontrol-dev.amcsgroup.io"),
					Description: "Url to CustomerControl API",
				},
				"private_key": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "CustomerControl private key for authentication",
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

		return p
	}

}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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
				Detail:   err.Error(),
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
