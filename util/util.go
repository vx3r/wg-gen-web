package util

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func ReadFile(path string) (bytes []byte, err error) {
	bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func WriteFile(path string, bytes []byte) (err error) {
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func FileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirectoryExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

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

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}