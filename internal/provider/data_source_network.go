package provider

import (
	"context"
	"fmt"
	"math/big"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to manipulate network addresses, such as IPv4 and IPv6 addresses, subnets, masks and prefixes.",

		ReadContext: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "The IPv4 or IPv6 address and network mask in CIDR notation. " +
					"The address may denote a specific address (using a host identifier, such as `10.0.0.1/8`), " +
					"or the beginning address of an entire network (using a host identifier of 0, such as `10.0.0.0/8`).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"prefix", "ip"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"ip": {
				Description:      "The IPv4 or IPv6 address (such as `10.0.0.1` or `2001:db8::68`).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				RequiredWith:     []string{"ip", "mask"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"mask": {
				Description:      "The IPv4 or IPv6 network mask (such as `255.0.0.0` or `ffff:ffff:ffff::`).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"mask_bits": {
				Description: "The number of bits in the IPv4 or IPv6 network mask.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"network": {
				Description: "The IPv4 or IPv6 network in CIDR notation (such as `10.0.0.0/8` or `2001:db8::/48`).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"first_ip": {
				Description: "The first IPv4 or IPv6 address in the CIDR range (such as `10.0.0.1` or `2001:db8::`).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_ip": {
				Description: "The last IPv4 or IPv6 address in the CIDR range (such as `10.255.255.254` or `2001:db8:0:ffff:ffff:ffff:ffff:ffff`).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"broadcast": {
				Description: "The IPv4 broadcast address (such as `10.255.255.255`).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"hostnum": {
				Description: "The hostnum (such as `1`).  Can be used in cidrhost as hostnum",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	var broadcast net.IP
	var ip net.IP
	var mask net.IPMask
	var network *net.IPNet
	var hostnum string

	if v, ok := d.GetOk("prefix"); ok {
		ip, network, err = net.ParseCIDR(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("ip"); ok {
		ip = net.ParseIP(v.(string))
	}

	if v, ok := d.GetOk("mask"); ok {
		m := net.ParseIP(v.(string))
		switch {
		case len(m.To4()) == net.IPv4len:
			mask = net.IPv4Mask(m[12], m[13], m[14], m[15])
		case len(m) == net.IPv6len:
			mask = net.IPMask(m)
		}

		network = &net.IPNet{
			IP:   ip.Mask(mask),
			Mask: mask,
		}
	}

	masked := net.IP(network.Mask)

	hostip := make(net.IP, len(ip))
	copy(hostip, ip)

	if len(masked) == net.IPv4len && len(hostip) == net.IPv6len {
		hostip = hostip[12:]
	}

	for i, b := range masked {
		hostip[i] = hostip[i] &^ b
	}

	ipint := big.NewInt(0)
	ipint.SetBytes(hostip)
	hostnum = ipint.Text(10)

	firstIP, lastIP := cidr.AddressRange(network)

	ones, _ := network.Mask.Size()
	if network.IP.To4() != nil {
		if ones < 31 {
			firstIP = cidr.Inc(firstIP)
			lastIP = cidr.Dec(lastIP)
		}

		broadcast = make(net.IP, len(network.IP.To4()))
		for i := 0; i < len(network.IP.To4()); i++ {
			broadcast[i] = network.IP[i] | ^network.Mask[i]
		}
	}

	d.Set("prefix", fmt.Sprintf("%s/%d", ip.String(), ones))
	d.Set("ip", ip.String())
	d.Set("mask", net.IP(network.Mask).String())
	d.Set("mask_bits", ones)
	d.Set("network", network.String())
	d.Set("first_ip", firstIP.String())
	d.Set("last_ip", lastIP.String())
	d.Set("broadcast", broadcast.String())
	d.Set("hostnum", hostnum)

	d.SetId(network.String())

	return nil
}
