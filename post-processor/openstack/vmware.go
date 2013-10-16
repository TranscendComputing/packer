package openstack

import (
	"fmt"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

type VMwareOpenStackConfig struct {
	common.PackerConfig `mapstructure:",squash"`
	tpl                 *packer.ConfigTemplate
}

type VMwareOpenStackProcessor struct {
	config VMwareOpenStackConfig
}

func (p *VMwareOpenStackProcessor) Configure(raws ...interface{}) error {
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

func (p *VMwareOpenStackProcessor) Process(ui packer.Ui, artifact packer.Artifact,
	access AccessConfig) (packer.Artifact, bool, error) {

	// TODO: implement

	var ErrNotImplemented = fmt.Errorf("Not implemented")

	return nil, false, ErrNotImplemented
}
