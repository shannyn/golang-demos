package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"syscall"
)

type TCPHeader struct {
	Source      uint16
	Destination uint16
	SeqNum      uint32
	AckNum      uint32
	DataOffset  uint8 // 4 bits
	Reserved    uint8 // 3 bits
	ECN         uint8 // 3 bits
	Ctrl        uint8 // 6 bits
	Window      uint16
	Checksum    uint16 // Kernel will set this if it's 0
	Urgent      uint16
	Options     []TCPOption
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

// Parse packet into TCPHeader structure
func NewTCPHeader(data []byte) *TCPHeader {
	var tcp TCPHeader
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &tcp.Source)
	binary.Read(r, binary.BigEndian, &tcp.Destination)
	binary.Read(r, binary.BigEndian, &tcp.SeqNum)
	binary.Read(r, binary.BigEndian, &tcp.AckNum)

	var mix uint16
	binary.Read(r, binary.BigEndian, &mix)
	tcp.DataOffset = byte(mix >> 12)  // top 4 bits
	tcp.Reserved = byte(mix >> 9 & 7) // 3 bits
	tcp.ECN = byte(mix >> 6 & 7)      // 3 bits
	tcp.Ctrl = byte(mix & 0x3f)       // bottom 6 bits

	binary.Read(r, binary.BigEndian, &tcp.Window)
	binary.Read(r, binary.BigEndian, &tcp.Checksum)
	binary.Read(r, binary.BigEndian, &tcp.Urgent)

	return &tcp
}

func main_socket() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Printf("socket err: %s", err)
		return
	}
	defer syscall.Close(fd)

	fmt.Println("start socket fd", fd)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if numRead > 20 {
			ip4Header, _ := ipv4.ParseHeader(buf[:20])
			fmt.Println("ip4Header: ", ip4Header)
		} else {
			fmt.Printf("% X\n", buf[:numRead])
		}

	}
}

func main_icmp() {
	protocol := "icmp"
	netaddr, _ := net.ResolveIPAddr("ip4", "127.0.0.1")
	conn, _ := net.ListenIP("ip4:"+protocol, netaddr)
	defer conn.Close()

	buf := make([]byte, 1024)
	numRead, _, _ := conn.ReadFrom(buf)
	fmt.Printf("% X\n", buf[:numRead])
}

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket err: %s \n", err)
		return
	}
	defer syscall.Close(fd)

	sa := syscall.SockaddrInet4{Port: 8888, Addr: [4]byte{127, 0, 0, 1}}
	err = syscall.Bind(fd, &sa)
	if err != nil {
		fmt.Printf("bind err: %s \n", err)
		return
	}

	err = syscall.Listen(fd, 256)
	if err != nil {
		fmt.Printf("listen err: %s \n", err)
		return
	}

	nfd, csa, err := syscall.Accept(fd)
	if err != nil {
		fmt.Printf("accept err: %s \n", err)
		return
	}
	defer syscall.Close(nfd)
	fmt.Printf("client addr: %s \n", csa)

	_, err = syscall.Write(nfd, []byte("hello, world!"))
	if err != nil {
		fmt.Printf("write err: %s \n", err)
		return
	}
}
