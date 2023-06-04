package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"net"
	"regexp"
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
