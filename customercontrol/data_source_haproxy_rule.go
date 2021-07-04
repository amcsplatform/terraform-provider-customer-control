package customercontrol

import (
	"context"
	"strconv"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHAProxyRule() *schema.Resource {
	return &schema.Resource{
		Description: "HAProxy rule resource, creates domain and virtual host",
		ReadContext: dataSourceHAProxyRuleRead,
		Schema: map[string]*schema.Schema{
			"virtual_host_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"setup_kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"setup_configuration": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"backend_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"set_host": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"setup_configuration_multi_forward": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"is_ssl": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"valid_until": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_certificate": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceHAProxyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cc.CustomerControlClient)
	var diags diag.Diagnostics

	virtualHostId := d.Get("virtual_host_id").(int)
	virtualHost, err := client.GetVirtualHostById(&virtualHostId)

	if err != nil {
		return diag.FromErr(err)
	}

	domainId := d.Get("domain_id").(int)
	domain, err := client.GetDomainById(domainId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(domain.DomainNameId))
	d.Set("valid_until", domain.ValidUntil)
	d.Set("manage_certificate", domain.ManageCertificate)
	d.Set("domain_name", domain.DomainName)
	d.Set("setup_kind", virtualHost.SetupKind)

	if virtualHost.SetupKind == "simple-forward" {
		setupConfigurationMap := map[string]interface{}{
			"backend": virtualHost.Configuration.Backend,
			"backend_port": virtualHost.Configuration.BackendPort,
			"is_ssl": virtualHost.Configuration.IsSsl,
			"set_host": virtualHost.Configuration.SetHost,
		}
		d.Set("setup_configuration", setupConfigurationMap)
	} else if virtualHost.SetupKind == "multi-forward" {
		// TODO: implement
		// setupConfigurationMap := map[string]interface{}{
		// 	"servers": []map[string]interface{}{

		// 	},
		// }
		// d.Set("setup_configuration_multi_forward", setupConfigurationMap)
	}

	return diags
}
