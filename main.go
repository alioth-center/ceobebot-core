package main

import (
	_ "github.com/ceobebot/qqchannel/plugin/imports"
	"github.com/ceobebot/qqchannel/processor"
)

func main() {
	processor.Serve()
}
