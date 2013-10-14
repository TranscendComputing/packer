package openstack

import (
	"fmt"
)

const BuilderId = "transcend.post-processor.openstack"

type Artifact struct {
	ImageId  string
	Provider string
}

func NewArtifact(provider, imageId string) *Artifact {
	return &Artifact{
		ImageId:  imageId,
		Provider: provider,
	}
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Files() []string {
	return []string{}
}

func (a *Artifact) Id() string {
	return a.ImageId
}

func (a *Artifact) String() string {
	return fmt.Sprintf("'%s' post-processor openstack: %s", a.Provider, a.ImageId)
}

func (a *Artifact) Destroy() error {
	return nil
}
