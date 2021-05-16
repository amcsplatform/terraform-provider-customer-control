package customercontrol

import (
	"context"
	"strconv"
	"strings"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHAProxyDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHAProxyDomainRead,
		Schema: map[string]*schema.Schema{
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
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

func dataSourceHAProxyDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*cc.CustomerControlClient)
	var diags diag.Diagnostics

	domainName := d.Get("domain_name").(string)
	domains, err := client.GetDomains()

	if err != nil {
		return diag.FromErr(err)
	}

	var matchingDomain *cc.Domain

	for _, domain := range *domains {
		if strings.Compare(domain.DomainName, domainName) == 0 {
			matchingDomain = &domain
			break
		}
	}

	if matchingDomain == nil {
		return diag.Errorf("Could not find domain: %s", domainName)
	}

	d.SetId(strconv.Itoa(matchingDomain.DomainNameId))
	d.Set("valid_until", matchingDomain.ValidUntil)
	d.Set("manage_certificate", matchingDomain.ManageCertificate)

	return diags
}
