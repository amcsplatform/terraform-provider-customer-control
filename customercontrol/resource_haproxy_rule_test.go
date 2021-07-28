package customercontrol

import (
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
		CheckDestroy:      testAccHAProxyRuleCheckDestroy(&virtualHostId, domainId),
		Steps: []resource.TestStep{
			{
				// Test resource creation
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_simple_forward.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.simple-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration[0].setup_kind", "simple-forward"),
				),
			},
			{
				// Update domain name, port and SSL
				Config: testAccExample(t, "resources/customercontrol_haproxy_rule/_acc_simple_forward_update_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccHAProxyRuleCheckExists("customercontrol_haproxy_rule.simple-forward", &domainId, &virtualHostId),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "domain_name", "terraform-provider-test.amcsgroup.io"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration[0].setup_kind", "simple-forward"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration[0].backend_port", "80"),
					resource.TestCheckResourceAttr("customercontrol_haproxy_rule.simple-forward", "setup_configuration[0].is_ssl", "false"),
				),
			},
			{
				// Importing matches the state of the previous step.
				ResourceName:      "customercontrol_haproxy_rule.simple-forward",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHAProxyRuleCheckExists(rn string, domainId *int, virtualHostId *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]

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

func testAccHAProxyRuleCheckDestroy(virtualHostId *int, domainId int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*cc.CustomerControlClient)

		virtualHost, err := client.GetVirtualHostById(virtualHostId)
		if err == nil && virtualHost.VirtualHostId > 0 {
			return fmt.Errorf("virtual host still exists, id: %s", strconv.Itoa(virtualHost.VirtualHostId))
		}

		domain, err := client.GetDomainById(domainId)
		if err == nil && domain.DomainNameId > 0 {
			return fmt.Errorf("domain still exists, id: %s", strconv.Itoa(domain.DomainNameId))
		}

		return nil
	}
}
