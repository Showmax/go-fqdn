package fqdn

import (
	"net"
	"os"
	"runtime"
	"strings"
)

// Get Fully Qualified Domain Name
// returns empty string or hostname if FQDN is unobtainable
func Get() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		if dnsDomain := os.Getenv("USERDNSDOMAIN"); dnsDomain != "" {
			return strings.Join([]string{hostname, dnsDomain}, "."), nil
		}
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return hostname, err
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip, err := ipv4.MarshalText()
			if err != nil {
				return hostname, err
			}
			hosts, err := net.LookupAddr(string(ip))
			if err != nil || len(hosts) == 0 {
				return hostname, err
			}
			fqdn := hosts[0]
			return strings.TrimSuffix(fqdn, "."), nil // return fqdn without trailing dot
		}
	}
	return hostname, nil
}
