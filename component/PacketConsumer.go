package component

import (
	"Awesome/component/model"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
)

type Comsumer interface {
	Consume(chan gopacket.Packet)
}
type PacketConsumer struct {
	parent *NetFlow
	ch     chan gopacket.Packet
	////数据解析曾
	//all_layer []gopacket.DecodingLayer
}

func (receiver PacketConsumer) Consume() {
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ip4, &ip6, &tcp, &udp)
	layerData := []gopacket.LayerType{}
	for {
		packet := <-receiver.ch
		layer := packet.LinkLayer()
		flow := layer.LinkFlow()
		metaInfo := receiver.parent.GetFlow(flow)
		meta := metaInfo.(*model.MetaFlow)
		//b_InData := (receiver.parent.device.Addresses == meta.Dst)
		err := parser.DecodeLayers(packet.Data(), &layerData)
		if err != nil {
			log.Fatal(err, " 解析协议层失败，继续解析下一个数据包")
			continue
		}
		for _, ltype := range layerData {
			switch ltype {
			case layers.LayerTypeUDP:
				meta.In_total++
			case layers.LayerTypeTCP:
				meta.In_total++
			}
		}
	}
}
