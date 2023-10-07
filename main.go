package main

import (
	_ "studio.sunist.work/sunist-c/ceobebot-qqchanel/plugin/imports"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor"
)

func main() {
	processor.Serve()
}
