package customercontrol

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"

	cc "dev.azure.com/amcsgroup/DevOps/_git/CustomerControlClientGo.git"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHAProxy_SimpleForward(t *testing.T) {
	var domainId int
	var virtualHostId int

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccHAProxyRuleCheckDestroy(&virtualHostId, &domainId),
		Steps: []resource.TestStep{
			{
				// Test resource creation
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_simple_forward.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.simple-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_kind", "simple-forward"),
				),
			},
			{
				// Update domain name, port and SSL
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_simple_forward_update_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.simple-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "domain_name", "terraform-provider-test-2.amcsgroup.io"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_kind", "simple-forward"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration.0.backend_port", "80"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration.0.is_ssl", "false"),
				),
			},
		},
	})
}

func TestAccHAProxy_MultiForward(t *testing.T) {
	var domainId int
	var virtualHostId int

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccHAProxyRuleCheckDestroy(&virtualHostId, &domainId),
		Steps: []resource.TestStep{
			{
				// Test resource creation
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_multi_forward.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.multi-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_kind", "multi-forward"),
				),
			},
			{
				// Update domain name, port and SSL
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_multi_forward_update_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.multi-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "domain_name", "terraform-provider-test-2.amcsgroup.io"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_kind", "multi-forward"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_configuration_multi_forward.0.port", "443"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_configuration_multi_forward.1.port", "443"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_configuration_multi_forward.0.is_ssl", "true"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.multi-forward", "setup_configuration_multi_forward.1.is_ssl", "true"),
				),
			},
		},
	})
}

func testAccHAProxyRuleCheckExists(rn string, domainId *int, virtualHostId *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]

		r, _ := json.Marshal(rs.Primary.Attributes)
		fmt.Println(string(r))

		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		domainIdAttr, _ := strconv.Atoi(rs.Primary.Attributes["domain_id"])
		if domainIdAttr <= 0 {
			return fmt.Errorf("domainId is not set")
		}

		*domainId = domainIdAttr

		virtualHostIdAttr, _ := strconv.Atoi(rs.Primary.Attributes["virtual_host_id"])
		if virtualHostIdAttr <= 0 {
			return fmt.Errorf("virtualHostId is not set")
		}

		*virtualHostId = virtualHostIdAttr

		client := testAccProvider.Meta().(*cc.CustomerControlClient)
		_, err := client.GetDomainById(*domainId)

		if err != nil {
			return fmt.Errorf("error getting domain: %s", err)
		}

		_, err = client.GetVirtualHostById(virtualHostId)

		if err != nil {
			return fmt.Errorf("error getting virtual host: %s", err)
		}

		return nil
	}
}

func testAccHAProxyRuleCheckDestroy(virtualHostId *int, domainId *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*cc.CustomerControlClient)

		virtualHost, err := client.GetVirtualHostById(virtualHostId)
		if err == nil && virtualHost.VirtualHostId > 0 {
			return fmt.Errorf("virtual host still exists, id: %s", strconv.Itoa(virtualHost.VirtualHostId))
		}

		domain, err := client.GetDomainById(*domainId)
		if err == nil && domain.DomainNameId > 0 {
			return fmt.Errorf("domain still exists, id: %s", strconv.Itoa(domain.DomainNameId))
		}

		return nil
	}
}
