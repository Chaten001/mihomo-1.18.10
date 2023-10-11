package inbound

import (
	"net"
	"net/netip"
	"strconv"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/transport/socks5"
)

// NewSocket receive TCP inbound and return ConnContext
func NewSocket(target socks5.Addr, conn net.Conn, source C.Type, additions ...Addition) (net.Conn, *C.Metadata) {
	metadata := parseSocksAddr(target)
	metadata.NetWork = C.TCP
	metadata.Type = source
	additions = append(additions, WithSrcAddr(conn.RemoteAddr()), WithInAddr(conn.LocalAddr()))
	for _, addition := range additions {
		addition.Apply(metadata)
	}

	return conn, metadata
}

func NewInner(conn net.Conn, address string) (net.Conn, *C.Metadata) {
	metadata := &C.Metadata{}
	metadata.NetWork = C.TCP
	metadata.Type = C.INNER
	metadata.DNSMode = C.DNSNormal
	metadata.Process = C.ClashName
	if h, port, err := net.SplitHostPort(address); err == nil {
		if port, err := strconv.ParseUint(port, 10, 16); err == nil {
			metadata.DstPort = uint16(port)
		}
		if ip, err := netip.ParseAddr(h); err == nil {
			metadata.DstIP = ip
		} else {
			metadata.Host = h
		}
	}

	return conn, metadata
}
