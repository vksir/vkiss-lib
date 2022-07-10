package udpredirect

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"net-tool/public/packet"
	"net-tool/public/udp"
)

var packetChan chan *packet.Packet

func findDevByIp(ip net.IP) (*pcap.Interface, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		for _, address := range device.Addresses {
			if address.IP.Equal(ip) {
				return &device, nil
			}
		}
	}
	return nil, errors.New("find device failed by ip")
}

func capture(dev *pcap.Interface, filter string) {
	if h, err := pcap.OpenLive(dev.Name, 4096, true, -1); err != nil {
		log.Panicln(err)
	} else if err := h.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	} else {
		defer h.Close()
		for {
			if packetData, _, err := h.ReadPacketData(); err != nil {
				log.Panicln(err)
			} else if p, err := packet.NewPacket(packetData); err != nil {
				log.Panicln(err)
			} else {
				fmt.Println(p.String())
				packetChan <- p
			}
		}
	}
}

func redirect(ctx context.Context, srcIP, dstIP net.IP) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-packetChan:
			u := udp.UDP{
				SrcIP:   srcIP,
				DstIP:   dstIP,
				SrcPort: int(p.UDP.SrcPort),
				DstPort: int(p.UDP.DstPort),
				Content: p.Payload.LayerContents(),
			}
			if err := u.Send(); err != nil {
				log.Printf("send p failed: %+v", u)
			}
		}
	}
}

func Run(ctx context.Context, devIp, filter, srcIp, dstIp string) {
	dev, err := findDevByIp(net.ParseIP(devIp))
	if err != nil {
		log.Panicf("find device by ip failed: ip=%s, err=%s", devIp, err)
	}

	packetChan = make(chan *packet.Packet, 64)
	go capture(dev, filter)
	redirect(ctx, net.ParseIP(srcIp), net.ParseIP(dstIp))
}
