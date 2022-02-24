//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package main

import (
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessConfig        `mapstructure:",squash"`

	Key string `mapstructure:"key" required:"true"`

	ctx interpolate.Context
}
