---
page_title: "Provider: CIDR"
subcategory: ""
description: |-
  The CIDR provider is used to manipulate network addresses, such as IPv4 and IPv6 addresses, subnets, masks and prefixes.
---

# CIDR Provider

The CIDR provider is used to manipulate network addresses, such as IPv4 and IPv6 addresses, subnets, masks and prefixes.

## Example Usage

```terraform
data "cidr_network" "example1" {
  prefix = "192.168.2.56/29"
}

data "cidr_network" "example2" {
  ip   = "192.168.2.57"
  mask = "255.255.255.248"
}
```