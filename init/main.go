package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		startCmds()
		return
	}

	if os.Args[1] == "gen-config" {
		genConfig()
		reload()
	}
}
