# packer-plugin-register-ami

## Description

* Packer post-processor plugin for saving AMI ID to AWS Systems Manager Parameter Store.

## Installation

```
packer {
  required_plugins {
    register-ami = {
      version = ">= <latest>"
      source = "github.com/exzeo/register-ami"
    }
  }
}
```

## Usage

The following example is a template for registering an AMI ID with a given Parameter Store.

```hcl
build {
  # ... build image
  post-processor "register-ami" {
    name = "/aws/test/ami"
  }
}
```
