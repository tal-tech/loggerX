package logutils

import (
	"net"
	"strconv"
	"strings"

	"github.com/petermattis/goid"
)

var prefix int64 = 10000000000000000

func GenLoggerId() int64 {
	id := goid.Get()
	return id + prefix
}

func init() {
	getprefix()
}

func getprefix() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	var firstIP net.IP
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				firstIP = ipnet.IP
				break
			}
		}
	}
	addrslice := strings.Split(firstIP.String(), ".")
	if len(addrslice) > 2 {
		cAddr, _ := strconv.ParseInt(addrslice[len(addrslice)-2], 10, 64)
		dAddr, _ := strconv.ParseInt(addrslice[len(addrslice)-1], 10, 64)
		prefix = (cAddr*1000 + dAddr) * 10000000000
	}
}
