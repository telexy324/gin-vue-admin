package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"net"
	"net/netip"
	"regexp"
	"strconv"
)

func GetManageIp() (ip net.IP, err error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		ads, er := iFace.Addrs()
		if er != nil {
			err = er
			continue
		}
		for _, addr := range ads {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			ok, er := regexp.Match(consts.ManageIpPrefix, []byte(ip.String()))
			if er != nil {
				err = er
				continue
			}
			if ok {
				return ip, nil
			} else {
				continue
			}
		}
	}
	return nil, err
}

func hosts(cidr string) (ips []string, err error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return
	}
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr.String())
	}

	if len(ips) < 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func Hosts(ip, mask string) ([]string, error) {
	prefix, _ := net.IPMask(net.ParseIP(mask).To4()).Size()
	return hosts(ip + `/` + strconv.Itoa(prefix))
}

func HostsCIDR(cidr string) ([]string, error) {
	return hosts(cidr)
}
