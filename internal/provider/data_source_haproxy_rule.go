package provider

import (
	"context"
	"strconv"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHAProxyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data resource to access information about an existing HAProxy rule",
		ReadContext: dataSourceHAProxyRuleRead,
		Schema: map[string]*schema.Schema{
			"virtual_host_id": {
				Description: "VirtualHost ID registered in CustomerControl",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"domain_id": {
				Description: "Domain ID registered in CustomerControl",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"domain_name": {
				Description: "Domain name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"setup_kind": {
				Description: "Rule kind",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"setup_configuration": {
				Description: "Rule configuration for simple-forward kind",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Description: "Backend address or IP to redirect requests to",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"is_ssl": {
							Description: "Enables SSL if true; terminates SSL if false",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"backend_port": {
							Description: "Backend port",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"set_host": {
							Description: "Passes host name in the request header to target backends if true",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"setup_configuration_multi_forward": {
				Description: "Rule configuration for multi-forward kind",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"servers": {
							Description: "List of backends",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Description: "Backend address or IP to redirect requests to",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"is_ssl": {
										Description: "Enables SSL if true; terminates SSL if false",
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
									},
									"port": {
										Description: "Backend port",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"valid_until": {
				Description: "SSL certificate validity if manage_certificate was set to true",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"manage_certificate": {
				Description: "Generates new SSL certificate for custom domain via LetsEncrypt and auto-renews it if true",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceHAProxyRuleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cc.CustomerControlClient)

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
		var virtualHostConfiguration = (virtualHost.Configuration).(cc.VirtualHostConfiguration)
		setupConfigurationMap := map[string]interface{}{
			"backend":      virtualHostConfiguration.Backend,
			"backend_port": virtualHostConfiguration.BackendPort,
			"is_ssl":       virtualHostConfiguration.IsSsl,
			"set_host":     virtualHostConfiguration.SetHost,
		}
		d.Set("setup_configuration", setupConfigurationMap)
	} else if virtualHost.SetupKind == "multi-forward" {
		var virtualHostConfiguration = (virtualHost.Configuration).(cc.VirtualHostConfigurationMultiBackends)
		var servers []map[string]interface{}

		for _, s := range virtualHostConfiguration.Servers {
			var server = map[string]interface{}{
				"url":    s.Url,
				"port":   s.Port,
				"is_ssl": s.IsSsl,
			}

			servers = append(servers, server)
		}

		setupConfigurationMap := servers
		d.Set("setup_configuration_multi_forward", setupConfigurationMap)
	}

	return diag.Diagnostics{}
}
