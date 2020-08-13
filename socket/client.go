package main

import (
	"fmt"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket err: %s \n", err)
		return
	}
	defer syscall.Close(fd)

	sa := syscall.SockaddrInet4{Port: 8888, Addr: [4]byte{127, 0, 0, 1}}
	err = syscall.Connect(fd, &sa)
	if err != nil {
		fmt.Printf("connect err: %s \n", err)
		return
	}

	buf := make([]byte, 256)
	_, err = syscall.Read(fd, buf)
	if err != nil {
		fmt.Printf("read err: %s \n", err)
		return
	}

	fmt.Printf("read message from server: %s \n", buf)
}
