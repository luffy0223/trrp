package main

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	filter := "tcp and port "
	portList := strings.Join(getPortList(), " or ")
	filter = filter + portList
	fmt.Println("filter:" + filter)

	if handle, err := pcap.OpenLive("any", 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(filter); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet) // Do something with a packet here.
		}
	}
}

func handlePacket(packet gopacket.Packet) {
	/*fmt.Println(packet.NetworkLayer())
	fmt.Println("-------------------------------")
	fmt.Println(packet.Metadata().CaptureInfo)*/
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			// Get actual TCP data from this layer
			tcp, _ := tcpLayer.(*layers.TCP)
			fmt.Printf("%s:%s -> %s:%s\n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)
		}
	}
}

func getHostByPid(pid string) (result string, err error) {
	//cmd:=exec.Command("bash","-c","netstat -anvp tcp |grep "+pid+" ｜awk -F \" \" -v OFS=\",\" '{print $4}' ")
	cmd := exec.Command("bash", "-c", "netstat -anvp tcp |grep "+pid+" |awk -F \" \" -v OFS=\"+++\" '{print $4,$5}'")
	if runtime.GOOS == "windows" {
		return "", nil
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	return out.String(), err1

}

func getPortList() (portList []string) {
	output, err1 := getHostByPid("49423")
	if err1 != nil {
		fmt.Println("err")
	}
	//fmt.Println(output)
	hostList := strings.Split(output, "\n")
	//fmt.Println(hostList)
	//var portList []string
	for _, value := range hostList {
		if len(value) < 1 {
			continue
		}
		sourceHost := strings.Split(value, "+++")[0]
		sourcePort := strings.Split(sourceHost, ".")
		result := sourcePort[len(sourcePort)-1]
		//fmt.Println(result)
		portList = append(portList, result)
	}
	fmt.Println("--------------")
	fmt.Println(portList)
	return
}
