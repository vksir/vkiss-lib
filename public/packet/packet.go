package packet

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strings"
)

type Packet struct {
	Ethernet layers.Ethernet
	IP4      layers.IPv4
	TCP      layers.TCP
	UDP      layers.UDP
	Payload  gopacket.Payload

	Decoded []gopacket.LayerType
}

func NewPacket(packetData []byte) (*Packet, error) {
	var p Packet
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &p.Ethernet, &p.IP4, &p.TCP, &p.UDP, &p.Payload)
	if err := parser.DecodeLayers(packetData, &p.Decoded); err != nil {
		return nil, err
	} else {
		return &p, nil
	}
}

func (p *Packet) String() string {
	var info []string
	for _, layerType := range p.Decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			info = append(info,
				p.Ethernet.SrcMAC.String()+" > "+p.Ethernet.DstMAC.String(),
				"Ethernet Type: "+p.Ethernet.EthernetType.String(),
			)
		case layers.LayerTypeIPv4:
			info = append(info,
				p.IP4.SrcIP.String()+" > "+p.IP4.DstIP.String(),
				"Protocol: "+p.IP4.Protocol.String(),
			)
		case layers.LayerTypeTCP:
			info = append(info,
				p.TCP.SrcPort.String()+" > "+p.TCP.DstPort.String(),
			)
		case layers.LayerTypeUDP:
			info = append(info,
				p.UDP.SrcPort.String()+" > "+p.UDP.DstPort.String(),
			)
		case gopacket.LayerTypePayload:
			info = append(info,
				"Content: "+string(p.Payload.LayerContents()),
			)
		}
	}
	return strings.Join(info, " | ")
}
