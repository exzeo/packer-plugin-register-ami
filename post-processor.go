package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/hcl/v2/hcldec"

	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type PostProcessor struct {
	testMode bool
	config   Config
	store    StoreInterface
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.plugin.register-ami",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if p.config.Key == "" {
		return errors.New("empty `key` is not allowed. Please make sure that it is set correctly")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	log.Println("Running the post-processor")

	amiId := p.GetImageId(artifact)

	if !p.testMode {
		sess, err := p.config.AccessConfig.Session()
		if err != nil {
			return nil, true, false, err
		}

		p.store, err = NewStore(
			sess.Copy(&aws.Config{Region: sess.Config.Region}),
			p.config,
		)

		if err != nil {
			return nil, true, false, err
		}
	}

	log.Printf("Saving AMI: %s to Parameter Store: %s\n", amiId, p.config.Key)
	err := p.store.SaveParameter(p.config.Key, amiId)
	if err != nil {
		return nil, true, false, *err
	}

	return artifact, true, false, nil
}

func (p *PostProcessor) GetImageId(artifact packer.Artifact) string {
	splitedString := strings.Split(artifact.Id(), ":")
	amiId := splitedString[len(splitedString)-1]

	return amiId
}
