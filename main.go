package main

import (
	"flag"

	_ "vk.com/m/docs"
	"vk.com/m/routes"
	"vk.com/m/utils"
)

var (
	addr = flag.String("addr", ":8000", "TCP address to listen to")
)

func main() {
	flag.Parse()
	utils.InitLogger()
	routes.Routes(addr)

}
