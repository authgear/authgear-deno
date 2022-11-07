package deno

import (
	"net"
	"net/url"
)

type HostPort struct {
	Host string
	IPv4 net.IP
	IPv6 net.IP
	Port string
}

func ParseHostPort(s string) (*HostPort, error) {
	if s == "" {
		return nil, nil
	}

	u, err := url.Parse("http://" + s)
	if err != nil {
		return nil, err
	}

	host := u.Hostname()
	port := u.Port()

	var ipv4 net.IP
	var ipv6 net.IP

	ip := net.ParseIP(host)
	if ip != nil {
		v4 := ip.To4()
		if v4 != nil {
			ipv4 = v4
		} else {
			ipv6 = ip.To16()
		}
	}

	return &HostPort{
		Host: host,
		IPv4: ipv4,
		IPv6: ipv6,
		Port: port,
	}, nil
}

func (p *HostPort) String() string {
	if p == nil {
		return ""
	}

	host := p.Host
	if p.IPv4 != nil {
		host = p.IPv4.String()
	}
	if p.IPv6 != nil {
		host = "[" + p.IPv6.String() + "]"
	}
	if p.Port == "" {
		return host
	}
	return host + ":" + p.Port
}

func (p HostPort) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *HostPort) UnmarshalText(text []byte) error {
	parsed, err := ParseHostPort(string(text))
	if err != nil {
		return err
	}
	*p = *parsed
	return nil
}
