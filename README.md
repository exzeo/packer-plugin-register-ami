# packer-plugin-register-ami

## Description

* Packer post-processor plugin for saving AMI ID to AWS Systems Manager Parameter Store.


## GPG for Github Action

Create Passphrase, this can generated, but keep it on hand, used to create GPG key and export

```
# macOS
gpg --gen-key
gpg --armor --export-secret-key <email> | pbcopy

# Ubuntu (assuming GNU base64)
gpg --gen-key
gpg --armor --export-secret-key <email> | xclip

```

In Github, Save Passphase into GPG_PASSPHRASE and save private key from export into GPG_PRIVATE_KEY


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
