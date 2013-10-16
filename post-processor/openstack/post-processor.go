// openstack implements the packer.PostProcessor interface and adds a
// post-processor that uploads artifacts to an OpenStack image repository
package openstack

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
)

var builtins = map[string]string{
	"mitchellh.virtualbox": "virtualbox",
	"mitchellh.vmware":     "vmware",
	"transcend.qemu":       "qemu",
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessConfig        `mapstructure:",squash"`
}

type PostProcessor struct {
	config      Config
	premade     map[string]OpenStackProcessor
	extraConfig map[string]interface{}
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	_, err := common.DecodeConfig(&p.config, raws...)
	if err != nil {
		return err
	}

	tpl, err := packer.NewConfigTemplate()
	if err != nil {
		return err
	}
	tpl.UserVars = p.config.PackerUserVars

	// Store extra configuration we'll send to each post-processor type
	p.extraConfig = make(map[string]interface{})
	p.extraConfig["packer_build_name"] = p.config.PackerBuildName
	p.extraConfig["packer_builder_type"] = p.config.PackerBuilderType
	p.extraConfig["packer_debug"] = p.config.PackerDebug
	p.extraConfig["packer_force"] = p.config.PackerForce
	p.extraConfig["packer_user_variables"] = p.config.PackerUserVars

	// check and configure access
	errs := &packer.MultiError{make([]error, 0)}
	errs = packer.MultiErrorAppend(errs, p.config.AccessConfig.Configure(tpl, raws)...)

	// look for and create any subprocessors
	// TODO(mitchellh): Properly handle multiple raw configs. This isn't
	// very pressing at the moment because at the time of this comment
	// only the first member of raws can contain the actual type-overrides.
	var mapConfig map[string]interface{}
	if err := mapstructure.Decode(raws[0], &mapConfig); err != nil {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Failed to decode config: %s", err))
		return errs
	}

	p.premade = make(map[string]OpenStackProcessor)
	for k, raw := range mapConfig {
		pp, err := p.openstackProcessor(k, raw, p.extraConfig)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
			continue
		}

		if pp == nil {
			continue
		}

		p.premade[k] = pp
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	ppName, ok := builtins[artifact.BuilderId()]
	if !ok {
		return nil, false, fmt.Errorf("Unknown artifact type, can't upload image: %s", artifact.BuilderId())
	}

	// Use the premade PostProcessor if we have one. Otherwise, we
	// create it and configure it here.
	pp := p.premade[ppName]
	if pp == nil {
		log.Printf("Premade post-processor for '%s' not found. Creating.", ppName)

		var err error
		pp, err = p.openstackProcessor(ppName, nil, p.extraConfig)
		if err != nil {
			return nil, false, err
		}

		if pp == nil {
			return nil, false, fmt.Errorf("OpenStack upload post-processor not found: %s", ppName)
		}
	}

	ui.Say(fmt.Sprintf("Uploading (to OpenStack) image for '%s' provider", ppName))
	return pp.Process(ui, artifact, p.config.AccessConfig)
}

func (p *PostProcessor) openstackProcessor(key string, specific interface{}, extra map[string]interface{}) (OpenStackProcessor, error) {
	gp := keyToPostProcessor(key)
	if gp == nil {
		return nil, nil
	}

	if err := gp.Configure(extra, specific); err != nil {
		return nil, err
	}

	return gp, nil
}

// keyToPostProcessor maps a configuration key to the actual post-processor
// it will be configuring. This returns a new instance of that post-processor.
func keyToPostProcessor(key string) OpenStackProcessor {
	switch key {
	case "virtualbox":
		return new(VBoxOpenStackProcessor)
	case "vmware":
		return new(VMwareOpenStackProcessor)
	case "qemu":
		return new(QemuOpenStackProcessor)
	default:
		return nil
	}
}
