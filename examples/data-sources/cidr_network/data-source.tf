data "cidr_network" "example1" {
  prefix = "192.168.2.56/29"
}

data "cidr_network" "example2" {
  ip   = "192.168.2.57"
  mask = "255.255.255.248"
}
