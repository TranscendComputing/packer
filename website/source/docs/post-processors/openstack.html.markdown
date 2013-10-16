---
layout: "docs"
---

# OpenStack Post-Processor

Type: `openstack`

The OpenStack builder is able to upload virtual machine images to the OpenStack Images service. The result of running this post-processor, if successful, is a new image ID.

The post-processor includes multiple sub-postprocessors that allow uploading from various other builders. Currently, only the [Qemu builder](/docs/builders/qemu.html) is supported, but plans currently call for completing support for [VMware](/docs/builders/vmware.html) and [VirtualBox](/docs/builders/virtualbox.html) as well.

If you've never used a post-processor before, please read the documentation on [using post-processors](/docs/templates/post-processors.html) in templates. This knowledge will be expected for the remainder of this document.

## Basic Example

Here is a basic example whithin a larger template. This example is functional so long as you fixup paths to files, URLS for ISOs, checksums and provide a reasonable kickstart file.

<pre class="prettyprint">
{
  "variables": {
	"qemu_image_name": "my-centos-test",
	"openstack_image_visibility": "public"
  },
  "builders":
  [
    {
      "type": "qemu",
      "name": "mycentos",
      "iso_url": "http://mirror.raystedman.net/centos/6/isos/x86_64/CentOS-6.4-x86_64-minimal.iso",
      "iso_checksum": "4a5fa01c81cc300f4729136e28ebe600",
      "iso_checksum_type": "md5",
      "output_directory": "output_centos_mytest",
      "ssh_wait_timeout": "30s",
      "shutdown_command": "shutdown -P now",
      "disk_size": 5000,
      "format": "qcow2",
      "headless": false,
      "accelerator": "kvm",
      "http_directory": "/home/myhomedirectory/packer/httpfiles",
      "http_port_min": 10082,
      "http_port_max": 10089,
      "ssh_host_port_min": 2222,
      "ssh_host_port_max": 2229,
      "ssh_username": "root",
      "ssh_password": "s0m3password",
      "ssh_port": 22,
      "ssh_wait_timeout": "90m",
      "vm_name": "mycentosvm",
      "net_device": "virtio-net",
      "disk_interface": "virtio",
      "qemuargs": [ ["-m", "1024m"], ["-cpu", "core2duo"] ],
      "boot_command":
      [
        "<tab><wait>",
        " ks=http://10.0.2.2:{{ .HTTPPort }}/centos6-ks.cfg<enter>"
      ]
    }
  ],
  "provisioners": [
    {
      "type": "shell",
      "inline": [
        "rm -fv /etc/udev/rules.d/70*",
        "find /etc/sysconfig/net* -name \"*.bak\" -exec rm -fv {} \\;",
        "rm -fv /etc/sysconfig/networking/profiles/default/*",
        "rm -fv /etc/sysconfig/networking/devices/*",
        "awk '!/HWADDR/' /etc/sysconfig/network-scripts/ifcfg-eth0 >/tmp/ifcfg-eth0",
        "awk '!/NM_MANAGED/' /tmp/ifcfg-eth0 >/etc/sysconfig/network-scripts/ifcfg-eth0"
      ]
    }
  ],
  "post-processors": [
    {
	  "type": "openstack",
      "username": "myopenstackusername",
	  "password": "myopenstackpassword",
	  "project": "MyOpenStackTenantName",
	  "provider": "http://my.openstack.com:35357/v2.0/tokens",
      "region": "RegionOne",
      "qemu": {
		"service_name": "glance",
		"service_type": "image",
	    "image_name": "{{user `qemu_image_name`}}",
	    "visibility": "{{user `openstack_image_visibility`}}",
	    "tags": [ "tag1", "tag2", "tag3", "this is just another tag" ]
      }
	}
  ]
}
</pre>

## Configuration Reference

There are many configuration options available for the OpenStack post-processor.  They are organized below into two categories: required and optional. Within each category, the available options are alphabetized and described.

Required:

* `ssh_username` (string) - The username to use to SSH into the machine once the OS is installed.
* `type` (string) - The name of the post processor to use -- always \"openstack\" for this postprocessor.
* `username` (string) - The (keystone) username to use for authentication with OpenStack.
* `password` (string) - The password to use for authentication with OpenStack.
* `project` (string) - The tenant (project) to use for registring private images.
* `provider` (string) - The URL to use for authentication with OpenStack, e.g., http://my.openstack.com:35357/v2.0/tokens)
* `region` (string) - The OpenStack Region in which to register the new image and upload the image file.

Optional:

* `qemu` - (object) Specifies to use the qemu sub-postprocessor.
    * `service_name` (string) - the OpenStack endpoint service name for the Images API.
    * `service_type` (string) - the OpenStack name for the service -- this will always be \"image\" for the time being.
    * `image_name` (string) - the name to set as the image name in OpenStack. This name will appear in the OpenStack user interface when viewing images.
    * `visibility` (string) - specifies whether the new image uploaded will have \"public\" or \"private\" visibility. If private, only those users in the project (tenant) will see the image in image lists.
    * `tags` (array of strings) - tags to apply to the image. Tags are arbitrary, short strings that are attached as metadata to the image record in the OpenStack Image service.

