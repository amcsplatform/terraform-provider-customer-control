package customercontrol

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHAProxyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceHAProxyRuleCreate,
		Read:   resourceHAProxyRuleRead,
		Delete: resourceHAProxyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"virtual_host_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"setup_kind": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateValueFunc([]string{
					"simple-forward",
					"multi-forward",
				}),
			},
			"setup_configuration": &schema.Schema{
				Type:         schema.TypeList,
				MaxItems:     1,
				ExactlyOneOf: []string{"setup_configuration", "setup_configuration_multi_forward"},
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
				ForceNew: true,
			},
			"setup_configuration_multi_forward": &schema.Schema{
				Type:         schema.TypeList,
				MaxItems:     1,
				ExactlyOneOf: []string{"setup_configuration", "setup_configuration_multi_forward"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
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
				ForceNew: true,
			},
			"valid_until": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"manage_certificate": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceHAProxyRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cc.CustomerControlClient)

	virtualHostId := d.Get("virtual_host_id").(int)
	virtualHost, err := client.GetVirtualHostById(&virtualHostId)

	if err != nil {
		return err
	}

	domainId := d.Get("domain_id").(int)
	domain, err := client.GetDomainById(domainId)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(domain.DomainNameId))
	d.Set("valid_until", domain.ValidUntil)
	d.Set("manage_certificate", domain.ManageCertificate)
	d.Set("domain_name", domain.DomainName)
	d.Set("setup_kind", virtualHost.SetupKind)

	if virtualHost.SetupKind == "simple-forward" {
		setupConfigurationMap := map[string]interface{}{
			"backend":      virtualHost.Configuration.Backend,
			"backend_port": virtualHost.Configuration.BackendPort,
			"is_ssl":       virtualHost.Configuration.IsSsl,
			"set_host":     virtualHost.Configuration.SetHost,
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

	return nil
}

func resourceHAProxyRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cc.CustomerControlClient)
	domainName := d.Get("domain_name").(string)
	// manageCertificate := d.Get("manage_certificate").(bool)

	domains, err := client.GetDomains()

	if err != nil {
		return err
	}

	// Check if domain already exists
	var matchingDomain *cc.Domain

	for _, domain := range *domains {
		if strings.Compare(domain.DomainName, domainName) == 0 {
			matchingDomain = &domain
			break
		}
	}

	if matchingDomain != nil {
		return fmt.Errorf("Domain %s already exists", domainName)
	}

	// Create domain
	log.Printf("[INFO] Creating new domain %s", domainName)
	domainId, err := client.CreateDomain(domainName)
	log.Printf("[INFO] Created new domain %s", strconv.Itoa(domainId))

	if err != nil {
		return err
	}

	// Create virtual host
	setupKind := d.Get("setup_kind").(string)
	setupKindType := cc.SimpleForward
	var setupConfiguration interface{}

	if setupKind == "simple-forward" {
		sc := d.Get("setup_configuration").([]interface{})
		setupConfigurationMap := sc[0].(map[string]interface{})
		setupConfiguration = cc.VirtualHostConfiguration{
			Backend:     setupConfigurationMap["backend"].(string),
			BackendPort: setupConfigurationMap["backend_port"].(int),
			IsSsl:       setupConfigurationMap["is_ssl"].(bool),
			SetHost:     setupConfigurationMap["set_host"].(bool),
		}
	} else if setupKind == "multi-forward" {
		// TODO: implement
		// setupKindType := cc.MultiForward
		// setupConfigurationMap := d.Get("setup_configuration_multi_forward").(map[string]interface{})
	}

	configurationBytes, err := json.Marshal(setupConfiguration)

	virtualHostModel := cc.VirtualHostPostModel{
		DomainNameId:       domainId,
		SetupKind:          setupKindType,
		SetupConfiguration: string(configurationBytes),
	}

	log.Printf("[INFO] Creating virtual host")
	virtualHostId, err := client.CreateVirtualHost(&virtualHostModel)
	log.Printf("[INFO] Created virtual host %s", strconv.Itoa(virtualHostId))

	if err != nil {
		return err
	}

	// Write configuration
	log.Printf("[INFO] Writing configration")
	err = client.WriteConfiguration()

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(virtualHostId))
	d.Set("virtual_host_id", virtualHostId)
	d.Set("domain_id", domainId)

	log.Printf("[INFO] Finished creating HAProxy rule")

	return resourceHAProxyRuleRead(d, m)
}

func resourceHAProxyRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*cc.CustomerControlClient)

	// Delete virtual host
	virtualHostId := d.Get("virtual_host_id").(int)
	_, err := client.DeleteVirtualHostById(&virtualHostId)
	if err != nil {
		return err
	}

	// Delete doman
	domainId := d.Get("domain_id").(int)
	err = client.DeleteDomainById(domainId)
	if err != nil {
		return err
	}

	// Write configuration after delete
	err = client.WriteConfiguration()
	if err != nil {
		return err
	}

	return nil
}
