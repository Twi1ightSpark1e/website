package util

import (
	"net"
	"net/http"
)

func GetRemoteAddr(req *http.Request) net.IP {
	if val, ok := req.Header["X-Real-Ip"]; ok {
		return net.ParseIP(val[0])
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return net.ParseIP(ip)
}
