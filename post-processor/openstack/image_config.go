package openstack

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
)

// ImageConfig common configuration related to Image Uploads
type ImageConfig struct {
	Name        string   `mapstructure:"image_name"`
	Visibility  string   `mapstructure:"visibility"`
	ServiceName string   `mapstructure:"service_name"`
	ServiceType string   `mapstructure:"service_type"`
	Tags        []string `mapstructure:"tags"`
}

func (ic *ImageConfig) Configure(t *packer.ConfigTemplate, raws ...interface{}) []error {
	errs := make([]error, 0)

	if t == nil {
		var err error
		t, err = packer.NewConfigTemplate()
		if err != nil {
			return []error{err}
		}
	}

	templates := map[string]*string{
		"image_name":   &ic.Name,
		"visibility":   &ic.Visibility,
		"service_name": &ic.ServiceName,
		"service_type": &ic.ServiceType,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = t.Process(*ptr, nil)
		if err != nil {
			errs = append(
				errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	// TODO: validate the tags

	if len(errs) > 0 {
		return errs
	}

	if ic.Visibility == "" {
		ic.Visibility = "public"
	}

	return nil
}
