package udp

import "net"

type UDP struct {
	SrcIP   net.IP
	DstIP   net.IP
	SrcPort int
	DstPort int
	Content []byte
}

func (u *UDP) Send() error {
	laddr := &net.UDPAddr{IP: u.SrcIP, Port: u.SrcPort}
	raddr := &net.UDPAddr{IP: u.DstIP, Port: u.DstPort}

	if conn, err := net.DialUDP("udp", laddr, raddr); err != nil {
		return err
	} else {
		defer conn.Close()
		if _, err := conn.Write(u.Content); err != nil {
			return err
		} else {
			return nil
		}
	}
}
