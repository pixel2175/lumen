package main

import (
	"lumen/src/actions"
	"lumen/src/parser"
)

func main() {
	actions.Run(parser.Parse())
}
