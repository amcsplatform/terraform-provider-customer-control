terraform {
  required_providers {
    customercontrol = {
      version = "0.0.1"
      source  = "amcsgroup.com/amcs/customercontrol"
    }
  }
}

provider "customercontrol" {
  url = "https://customercontrol-dev.amcsgroup.com"
  privateKey = ""
}

data "customercontrol_domain" "test" {
  domain_name = "d1-p83-svc-publisher-proxy.amcsplatform.com"
}

output "domain" {
  value = data.customercontrol_domain.id
}

