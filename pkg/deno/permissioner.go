package deno

import (
	"context"
	"fmt"
	"net"
)

type Permissioner interface {
	RequestPermission(ctx context.Context, pd PermissionDescriptor) (bool, error)
}

type IPPolicyPermissioner struct {
	resolver *net.Resolver
	disallow []IPPolicy
}

func DisallowIPPolicy(policies ...IPPolicy) IPPolicyPermissioner {
	return IPPolicyPermissioner{
		disallow: policies,
	}
}

func (p IPPolicyPermissioner) RequestPermission(ctx context.Context, pd PermissionDescriptor) (bool, error) {
	if pd.Name != PermissionNameNet {
		return false, fmt.Errorf("name != net")
	}

	if pd.Host == nil {
		return false, fmt.Errorf("any host")
	}

	var ips []net.IP
	if pd.Host.IPv4 != nil {
		ips = append(ips, pd.Host.IPv4)
	} else if pd.Host.IPv6 != nil {
		ips = append(ips, pd.Host.IPv6)
	} else {
		addrs, err := p.resolver.LookupHost(ctx, pd.Host.Host)
		if err != nil {
			return false, err
		}
		for _, addr := range addrs {
			ip := net.ParseIP(addr)
			if ip == nil {
				return false, fmt.Errorf("invalid ip: %v", addr)
			}
			ips = append(ips, ip)
		}
	}

	if len(ips) <= 0 {
		return false, fmt.Errorf("resolved no ip")
	}

	for _, ip := range ips {
		for _, policy := range p.disallow {
			_, err := policy(ip)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

type IPPolicy func(ip net.IP) (bool, error)

func DisallowGlobalUnicast(ip net.IP) (bool, error) {
	if ip.IsGlobalUnicast() {
		return false, fmt.Errorf("ip is global unicast")
	}
	return true, nil
}

func DisallowInterfaceLocalMulticast(ip net.IP) (bool, error) {
	if ip.IsInterfaceLocalMulticast() {
		return false, fmt.Errorf("ip is interface local multicast")
	}
	return true, nil
}

func DisallowLinkLocalUnicast(ip net.IP) (bool, error) {
	if ip.IsLinkLocalUnicast() {
		return false, fmt.Errorf("ip is link local unicast")
	}
	return true, nil
}

func DisallowLinkLocalMulticast(ip net.IP) (bool, error) {
	if ip.IsLinkLocalMulticast() {
		return false, fmt.Errorf("ip is link local multicast")
	}
	return true, nil
}

func DisallowLoopback(ip net.IP) (bool, error) {
	if ip.IsLoopback() {
		return false, fmt.Errorf("ip is loopback")
	}
	return true, nil
}

func DisallowMulticast(ip net.IP) (bool, error) {
	if ip.IsMulticast() {
		return false, fmt.Errorf("ip is multicast")
	}
	return true, nil
}

func DisallowPrivate(ip net.IP) (bool, error) {
	if ip.IsPrivate() {
		return false, fmt.Errorf("ip is private")
	}
	return true, nil
}

func DisallowUnspecified(ip net.IP) (bool, error) {
	if ip.IsUnspecified() {
		return false, fmt.Errorf("ip is unspecified")
	}
	return true, nil
}
