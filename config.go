//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package main

import (
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessConfig        `mapstructure:",squash"`

	Name         string `mapstructure:"name"`
	SecureString bool   `mapstructure:"secure_string"`
	AmiDataType  bool   `mapstructure:"ami_data_type"`

	ctx interpolate.Context
}
