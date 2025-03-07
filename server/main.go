package main

import "net_game/server/snet"

func main() {

	server := snet.NewServer("127.0.0.1", 8080, "", "")
	server.Serve()

}
