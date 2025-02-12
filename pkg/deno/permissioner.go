package deno

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type ErrorNameUnmatched struct {
	Expected PermissionName
	Actual   PermissionName
}

func (e *ErrorNameUnmatched) Error() string {
	return fmt.Sprintf("name unmatched, expected `%v`, actual `%v`", e.Expected, e.Actual)
}

type ErrorInvalidIP struct {
	Value string
}

func (e *ErrorInvalidIP) Error() string {
	return fmt.Sprintf("invalid ip: %v", e.Value)
}

var ErrAllHost = errors.New("network permission without host is disallowed")

type ErrorNoIP struct {
	Host string
}

func (e *ErrorNoIP) Error() string {
	return fmt.Sprintf("resolve to no ip: %v", e.Host)
}

type ErrorGlobalUnicast struct {
	IP net.IP
}

func (e *ErrorGlobalUnicast) Error() string {
	return fmt.Sprintf("global unicast: %v", e.IP)
}

type ErrorInterfaceLocalMulticast struct {
	IP net.IP
}

func (e *ErrorInterfaceLocalMulticast) Error() string {
	return fmt.Sprintf("interface local multicast: %v", e.IP)
}

type ErrorLinkLocalUnicast struct {
	IP net.IP
}

func (e *ErrorLinkLocalUnicast) Error() string {
	return fmt.Sprintf("link local unicast: %v", e.IP)
}

type ErrorLinkLocalMulticast struct {
	IP net.IP
}

func (e *ErrorLinkLocalMulticast) Error() string {
	return fmt.Sprintf("link local multicast: %v", e.IP)
}

type ErrorLoopback struct {
	IP net.IP
}

func (e *ErrorLoopback) Error() string {
	return fmt.Sprintf("loopback: %v", e.IP)
}

type ErrorMulticast struct {
	IP net.IP
}

func (e *ErrorMulticast) Error() string {
	return fmt.Sprintf("multicast: %v", e.IP)
}

type ErrorPrivate struct {
	IP net.IP
}

func (e *ErrorPrivate) Error() string {
	return fmt.Sprintf("private: %v", e.IP)
}

type ErrorUnspecified struct {
	IP net.IP
}

func (e *ErrorUnspecified) Error() string {
	return fmt.Sprintf("unspecified: %v", e.IP)
}

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
		return false, &ErrorNameUnmatched{
			Expected: PermissionNameNet,
			Actual:   pd.Name,
		}
	}

	if pd.Host == nil {
		return false, ErrAllHost
	}

	var ips []net.IP
	switch {
	case pd.Host.IPv4 != nil:
		ips = append(ips, pd.Host.IPv4)
	case pd.Host.IPv6 != nil:
		ips = append(ips, pd.Host.IPv6)
	default:
		addrs, err := p.resolver.LookupHost(ctx, pd.Host.Host)
		if err != nil {
			return false, err
		}
		for _, addr := range addrs {
			ip := net.ParseIP(addr)
			if ip == nil {
				return false, &ErrorInvalidIP{
					Value: addr,
				}
			}
			ips = append(ips, ip)
		}
	}

	if len(ips) <= 0 {
		return false, &ErrorNoIP{
			Host: pd.Host.Host,
		}
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
		return false, &ErrorGlobalUnicast{ip}
	}
	return true, nil
}

func DisallowInterfaceLocalMulticast(ip net.IP) (bool, error) {
	if ip.IsInterfaceLocalMulticast() {
		return false, &ErrorInterfaceLocalMulticast{ip}
	}
	return true, nil
}

func DisallowLinkLocalUnicast(ip net.IP) (bool, error) {
	if ip.IsLinkLocalUnicast() {
		return false, &ErrorLinkLocalUnicast{ip}
	}
	return true, nil
}

func DisallowLinkLocalMulticast(ip net.IP) (bool, error) {
	if ip.IsLinkLocalMulticast() {
		return false, &ErrorLinkLocalMulticast{ip}
	}
	return true, nil
}

func DisallowLoopback(ip net.IP) (bool, error) {
	if ip.IsLoopback() {
		return false, &ErrorLoopback{ip}
	}
	return true, nil
}

func DisallowMulticast(ip net.IP) (bool, error) {
	if ip.IsMulticast() {
		return false, &ErrorMulticast{ip}
	}
	return true, nil
}

func DisallowPrivate(ip net.IP) (bool, error) {
	if ip.IsPrivate() {
		return false, &ErrorPrivate{ip}
	}
	return true, nil
}

func DisallowUnspecified(ip net.IP) (bool, error) {
	if ip.IsUnspecified() {
		return false, &ErrorUnspecified{ip}
	}
	return true, nil
}
