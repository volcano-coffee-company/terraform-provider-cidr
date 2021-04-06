package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetwork_IPv4Mask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkIPMask("192.168.0.1", "255.255.255.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/24"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.255.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.0.254"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "broadcast", "192.168.0.255"),
				),
			},
		},
	})
}

func TestAccDataSourceNetwork_IPv4Prefix(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkPrefix("192.168.0.1/22"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/22"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/22"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.252.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "22"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.3.254"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "broadcast", "192.168.3.255"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("192.168.0.1/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/24"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.255.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "24"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.0.254"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "broadcast", "192.168.0.255"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("192.168.0.1/30"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/30"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/30"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.255.252"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "30"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "broadcast", "192.168.0.3"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("192.168.0.1/31"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/31"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/31"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.255.254"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "31"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.0.1"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("192.168.0.1/32"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.1/32"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.1/32"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "255.255.255.255"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "32"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.0.1"),
				),
			},
		},
	})
}

func TestAccDataSourceNetwork_IPv6Mask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkIPMask("2a03:2880:20:4f06::", "ffff:ffff:ffff:ffff::"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "2a03:2880:20:4f06::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "2a03:2880:20:4f06::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff:ffff::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "2a03:2880:20:4f06::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "2a03:2880:20:4f06::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "2a03:2880:20:4f06:ffff:ffff:ffff:ffff"),
				),
			},
		},
	})
}

func TestAccDataSourceNetwork_IPv6Prefix(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkPrefix("2a03:2880:20:4f06::/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "2a03:2880:20:4f06::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "2a03:2880:20:4f06::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff:ffff::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "2a03:2880:20:4f06::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "2a03:2880:20:4f06::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "2a03:2880:20:4f06:ffff:ffff:ffff:ffff"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("1::/64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "1::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "1::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff:ffff::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "1::/64"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "1::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "1::ffff:ffff:ffff:ffff"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("::ffff:192.168.0.0/112"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "192.168.0.0/112"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "192.168.0.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "112"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "192.168.0.0/16"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "192.168.0.0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "192.168.255.255"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("fe80::/48"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "fe80::/48"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "fe80::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "48"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "fe80::/48"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "fe80::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "fe80::ffff:ffff:ffff:ffff:ffff"),
				),
			},
			{
				Config: testAccDataSourceNetworkPrefix("fe80::3:0:0:0/81"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_network.test", "prefix", "fe80::3:0:0:0/81"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "ip", "fe80::3:0:0:0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask", "ffff:ffff:ffff:ffff:ffff:8000::"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "mask_bits", "81"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "network", "fe80::3:0:0:0/81"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "first_ip", "fe80::3:0:0:0"),
					resource.TestCheckResourceAttr("data.cidr_network.test", "last_ip", "fe80::3:7fff:ffff:ffff"),
				),
			},
		},
	})
}

func TestAccDataSourceNetwork_MissingArgs(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNetworkMissingArgs,
				ExpectError: regexp.MustCompile("one of `ip,prefix` must be specified"),
			},
		},
	})
}

func TestAccDataSourceNetwork_MissingIP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNetworkMissingIP,
				ExpectError: regexp.MustCompile("one of `ip,prefix` must be specified"),
			},
		},
	})
}

func TestAccDataSourceNetwork_MissingMask(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceNetworkMissingMask,
				ExpectError: regexp.MustCompile("all of `ip,mask` must be specified"),
			},
		},
	})
}

func testAccDataSourceNetworkPrefix(prefix string) string {
	return fmt.Sprintf(`
data "cidr_network" "test" {
  prefix = "%s"
}
`, prefix)
}

func testAccDataSourceNetworkIPMask(ip, mask string) string {
	return fmt.Sprintf(`
data "cidr_network" "test" {
  ip   = "%s"
  mask = "%s"
}
`, ip, mask)
}

const testAccDataSourceNetworkMissingArgs = `
data "cidr_network" "test" {
}
`

const testAccDataSourceNetworkMissingIP = `
data "cidr_network" "test" {
  mask = "255.255.255.0"
}
`

const testAccDataSourceNetworkMissingMask = `
data "cidr_network" "test" {
  ip = "192.168.2.56"
}
`
