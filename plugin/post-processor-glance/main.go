package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/mitchellh/packer/post-processor/glance"
)

func main() {
	plugin.ServePostProcessor(new(glance.PostProcessor))
}
