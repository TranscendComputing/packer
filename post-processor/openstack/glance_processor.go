package openstack

import (
	"github.com/mitchellh/packer/packer"
)

type OpenStackProcessor interface {
	// Configuration
	Configure(...interface{}) error

	// Processor (generally upload)
	Process(packer.Ui, packer.Artifact, AccessConfig) (packer.Artifact, bool, error)
}
