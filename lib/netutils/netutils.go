package netutils

import (
	"math/big"
	"net"
)

func BroadcastAddress(subnet *net.IPNet) net.IP {
	n := len(subnet.IP)
	out := make(net.IP, n)
	var m byte
	for i := 0; i < n; i++ {
		m = subnet.Mask[i] ^ 0xff
		out[i] = subnet.IP[i] | m
	}
	return out
}

func NextSubnet(subnet *net.IPNet) *net.IPNet {
	n := len(subnet.IP)
	out := BroadcastAddress(subnet)
	var c byte = 1
	for i := n - 1; i >= 0; i-- {
		out[i] = out[i] + c
		if out[i] == 0 && c > 0 {
			c = 1
		} else {
			c = 0
		}

	}
	if c == 1 {
		return nil
	}
	return &net.IPNet{IP: out.Mask(subnet.Mask), Mask: subnet.Mask}
}

func IncrementIP(ip net.IP, increment int) net.IP {
	if ip == nil {
		return nil
	}

	if increment == 0 {
		return ip.To16()
	}

	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}

	ipInt := new(big.Int)
	ipInt.SetBytes(ip)

	incrementInt := big.NewInt(int64(increment))

	result := new(big.Int).Add(ipInt, incrementInt)

	resultBytes := result.Bytes()

	if ip.To4() != nil {
		if len(resultBytes) > 4 {
			return nil
		}
		resultIP := make(net.IP, 4)
		copy(resultIP[4-len(resultBytes):], resultBytes)
		return resultIP
	} else {
		if len(resultBytes) > 16 {
			return nil
		}
		resultIP := make(net.IP, 16)
		copy(resultIP[16-len(resultBytes):], resultBytes)
		return resultIP
	}
}
