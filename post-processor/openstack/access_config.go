package openstack

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
	"github.com/rackspace/gophercloud"
	"os"
)

// AccessConfig is for common configuration related to openstack access
type AccessConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Project  string `mapstructure:"project"`
	Provider string `mapstructure:"provider"`
	Region   string `mapstructure:"region"`
	access   *gophercloud.Access
	context  *gophercloud.Context
}

// Auth returns a valid Auth object for access to openstack services, or
// an error if the authentication couldn't be resolved.
func (ac *AccessConfig) Auth() error {
	username := ac.Username
	password := ac.Password
	project := ac.Project
	provider := ac.Provider

	if username == "" {
		username = os.Getenv("SDK_USERNAME")
	}
	if password == "" {
		password = os.Getenv("SDK_PASSWORD")
	}
	if project == "" {
		project = os.Getenv("SDK_PROJECT")
	}
	if provider == "" {
		provider = os.Getenv("SDK_PROVIDER")
	}

	authoptions := gophercloud.AuthOptions{
		Username:    username,
		Password:    password,
		AllowReauth: true,
	}

	if project != "" {
		authoptions.TenantName = project
	}

	ac.context = gophercloud.TestContext()
	ac.context.RegisterProvider(ac.Provider, gophercloud.Provider{ac.Provider})
	access, err := ac.context.Authenticate(ac.Provider, authoptions)
	if err != nil {
		return err
	}
	ac.access = access

	return nil
}

func (ac *AccessConfig) Configure(t *packer.ConfigTemplate, raws ...interface{}) []error {
	errs := make([]error, 0)

	if t == nil {
		var err error
		t, err = packer.NewConfigTemplate()
		if err != nil {
			return []error{err}
		}
	}

	templates := map[string]*string{
		"username": &ac.Username,
		"password": &ac.Password,
		"provider": &ac.Provider,
		"region":   &ac.Region,
	}

	for n, ptr := range templates {
		var err error
		*ptr, err = t.Process(*ptr, nil)
		if err != nil {
			errs = append(
				errs, fmt.Errorf("Error processing %s: %s", n, err))
		}
	}

	if ac.Region == "" {
		errs = append(errs, fmt.Errorf("region must be specified"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (ac *AccessConfig) GetImagesApi(service_name, service_type string) (gophercloud.CloudServersProvider, error) {
	api := gophercloud.ApiCriteria{
		Name:      service_name,
		Type:      service_type,
		Region:    ac.Region,
		VersionId: "",
		UrlChoice: gophercloud.PublicURL,
	}

	servers, err := ac.context.ServersApi(ac.access, api)
	if err != nil {
		return nil, err
	}

	return servers, err
}
