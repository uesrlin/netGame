package main

import "net"

func main() {
	_, err := net.Listen("", "")
	if err != nil {
		return
	}

}
