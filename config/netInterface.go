package config

import (
	"fmt"
	"net"
)

// NetInterface 本机网络
type NetInterface struct {
	Name    string
	Address []string
}

// GetNetInterface 获得网卡地址
// 返回ipv4, ipv6地址
func GetNetInterface() (ipv4NetInterfaces []NetInterface, ipv6NetInterfaces []NetInterface, err error) {
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return ipv4NetInterfaces, ipv6NetInterfaces, err
	}

	// https://en.wikipedia.org/wiki/IPv6_address#General_allocation
	_, ipv6Unicast, _ := net.ParseCIDR("2000::/3")

	for _, netInterface := range allNetInterfaces {
		if (netInterface.Flags & net.FlagUp) != 0 {
			addrs, _ := netInterface.Addrs()
			ipv4 := 0
			ipv6 := 0

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
					ones, bits := ipnet.Mask.Size()
					// 需匹配全局单播地址
					if bits == 128 && ones <= bits && ipv6Unicast.Contains(ipnet.IP) {
						ipv6NetInterfaces = append(
							ipv6NetInterfaces,
							NetInterface{
								Name:    fmt.Sprintf(" %v / %v (%v) ", netInterface.Name, ones, ipv4),
								Address: []string{ipnet.IP.String()},
							},
						)
						ipv4++
					}
					if bits == 32 {
						ipv4NetInterfaces = append(
							ipv4NetInterfaces,
							NetInterface{
								Name:    fmt.Sprintf(" %v / %v (%v) ", netInterface.Name, ones, ipv6),
								Address: []string{ipnet.IP.String()},
							},
						)
						ipv6++
					}
				}
			}

		}
	}

	return ipv4NetInterfaces, ipv6NetInterfaces, nil
}
