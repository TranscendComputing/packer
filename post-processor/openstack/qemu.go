package openstack

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"github.com/rackspace/gophercloud"
	"log"
	"strings"
)

type QemuGlanceConfig struct {
	common.PackerConfig `mapstructure:",squash"`
	ImageConfig         `mapstructure:",squash"`
	tpl                 *packer.ConfigTemplate
}

type QemuGlanceProcessor struct {
	config QemuGlanceConfig
}

func (p *QemuGlanceProcessor) Configure(raws ...interface{}) error {
	md, err := common.DecodeConfig(&p.config, raws...)
	if err != nil {
		return err
	}

	p.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return err
	}
	p.config.tpl.UserVars = p.config.PackerUserVars

	// Accumulate any errors
	errs := common.CheckUnusedConfig(md)
	errs = packer.MultiErrorAppend(errs, p.config.ImageConfig.Configure(p.config.tpl, raws)...)

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *QemuGlanceProcessor) Process(ui packer.Ui, artifact packer.Artifact,
	ac AccessConfig) (packer.Artifact, bool, error) {

	// Get the file and type from the artifact.
	// Right now the artifact from the qemu builder has but one entry -- the
	// resulting image file. If, however, there ever are other files, at the
	// very leaset search for a qcow2 or img file.
	var imagePath, diskFormat string
	for _, path := range artifact.Files() {
		switch {
		case strings.HasSuffix(path, "qcow2"):
			imagePath = path
			diskFormat = "qcow2"
		case strings.HasSuffix(path, "raw"):
			imagePath = path
			diskFormat = "raw"
		}
	}
	log.Println("QemuGlanceProcessor.Process: will upload ", imagePath,
		" as ", diskFormat)

	err := ac.Auth()
	if err != nil {
		return nil, false, err
	}

	servers, err := ac.GetImagesApi(p.config.ImageConfig.ServiceName,
		p.config.ImageConfig.ServiceType)
	if err != nil {
		return nil, false, err
	}

	ni := gophercloud.NewImage{
		Name:            p.config.ImageConfig.Name,
		Visibility:      p.config.ImageConfig.Visibility,
		ContainerFormat: "bare",
		DiskFormat:      diskFormat,
		Tags:            p.config.ImageConfig.Tags,
	}

	imageId, err := servers.CreateNewImage(ni)
	if err != nil {
		return nil, false, err
	}

	// test the upload now
	err = servers.UploadImageFile(imageId, imagePath)
	if err != nil {
		return nil, false, err
	}

	return NewArtifact("qemu", imageId), false, nil
}
