package component

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

type Comsumer interface {
	Consume(chan gopacket.Packet)
}
type PacketConsumer struct {
	parent NetFlow
	ch     chan gopacket.Packet
	//数据解析曾
	layerType []gopacket.Layer
}

func (receiver PacketConsumer) Consume(flows chan gopacket.Packet) {
	var ip4 layers.IPv4
	var ip6 layers.IPv4
	var tcp layers.TCP
	var udp layers.UDP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ip4, &ip6, &tcp, &udp)
	layerData := []gopacket.LayerType{}
	for {
		packet := <-flows
		layer := packet.LinkLayer()
		flow := layer.LinkFlow()
		metaInfo := receiver.parent.GetFlow(flow)
		meta, ok := metaInfo.(MetaFlow)
		if !ok {
			log.Fatal()
			return
		}
		b_InData := (receiver.parent.ip == meta.src)
		err := parser.DecodeLayers(packet.Data(), &layerData)
		if err != nil {
			log.Fatal(err, " 解析协议层失败，继续解析下一个数据包")
			continue
		}
		for _, ltype := range layerData {
			switch ltype {
			case layers.LayerTypeUDP:
				if b_InData {
					meta.Out_udp += 1
				} else {
					meta.In_Udp += 1
				}
			case layers.LayerTypeTCP:
				meta.Detail()
			}
		}
	}
}
