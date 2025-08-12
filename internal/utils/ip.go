package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func GetPrivateIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	findIPByName := func(name string) string {
		for _, iface := range ifaces {
			if iface.Name == name && iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
				if ip := findPrivateIPForInterface(&iface); ip != "" {
					return ip
				}
			}
		}
		return ""
	}

	if ip := findIPByName("en0"); ip != "" {
		return ip, nil
	}

	if ip := findIPByName("en1"); ip != "" {
		return ip, nil
	}

	for _, iface := range ifaces {
		if iface.Name == "en0" || iface.Name == "en1" {
			continue
		}
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if ip := findPrivateIPForInterface(&iface); ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("no private ip found")
}

func findPrivateIPForInterface(iface *net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue
		}
		if isPrivateIP(ip) {
			return ip.String()
		}
	}
	return ""
}

// isPrivateIP checks if IP is in RFC1918 private ranges
func isPrivateIP(ip net.IP) bool {
	privateBlocks := []net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}
	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// getPublicIP fetches the external IP from an online service
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(ipBytes))
	return ip, nil
}

func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() != nil || parsedIP.To16() != nil
}