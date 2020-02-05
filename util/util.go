package util

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// ReadFile file content
func ReadFile(path string) (bytes []byte, err error) {
	bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// WriteFile content to file
func WriteFile(path string, bytes []byte) (err error) {
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// FileExists check if file exists
func FileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists check if directory exists
func DirectoryExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetAvailableIp search for an available in cidr against a list of reserved ips
func GetAvailableIp(cidr string, reserved []string) (string, error) {
	addresses, err := GetAllAddressesFromCidr(cidr)
	if err != nil {
		return "", err
	}

	for _, addresse := range addresses {
		ok := true
		for _, r := range reserved {
			if addresse == r {
				ok = false
				break
			}
		}
		if ok {
			return addresse, nil
		}
	}

	return "", errors.New("no more available address from cidr")
}

// GetAllAddressesFromCidr get all ip addresses from cidr
func GetAllAddressesFromCidr(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address (and server ip .1)
	return ips[2 : len(ips)-1], nil
}

// IsIPv6 check if given ip is IPv6
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
