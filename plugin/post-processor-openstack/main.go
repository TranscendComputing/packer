package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/mitchellh/packer/post-processor/openstack"
)

func main() {
	plugin.ServePostProcessor(new(openstack.PostProcessor))
}
