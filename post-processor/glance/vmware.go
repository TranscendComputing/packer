package glance

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

type VMwareGlanceConfig struct {
	common.PackerConfig `mapstructure:",squash"`
	tpl                 *packer.ConfigTemplate
}

type VMwareGlanceProcessor struct {
	config VMwareGlanceConfig
}

func (p *VMwareGlanceProcessor) Configure(raws ...interface{}) error {
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
	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *VMwareGlanceProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	// generate a UUID for glance (and as the artifact id))
	imageId, err := GenerateUUID()
	if err != nil {
		return nil, false, err
	}

	// TODO: Get the image file from the artifact

	return NewArtifact("vmware", imageId), false, nil
}
