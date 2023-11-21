package main

import "flag"

func parseFlags(addr flag.Value) {
	flag.Var(addr, "a", "net address host:port")
	flag.Parse()
}
