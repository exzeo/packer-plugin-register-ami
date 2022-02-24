package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/exzeo/packer-plugin-register-ami/mocks"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func testUI() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessor_Configure_validConfig(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{
		"name": "/aws/server/ami",
	})

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_Configure_missingName(t *testing.T) {
	p := new(PostProcessor)
	err := p.Configure(map[string]interface{}{})

	if err == nil {
		t.Fatal("should cause validation errors")
	}
	if err.Error() != "empty `name` is not allowed. Please make sure that it is set correctly" {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
}

func TestPostProcessor_PostProcess(t *testing.T) {
	mockSSM := &mocks.SSMAPI{}

	mockSSM.On("PutParameter", mock.Anything).Return(nil, nil)

	p := PostProcessor{
		testMode: true,
		store: &Store{
			ssmconn: mockSSM,
		},
		config: Config{
			Name: "/aws/test/ami",
		},
	}

	_, keep, forceOverride, err := p.PostProcess(context.Background(), testUI(), &packer.MockArtifact{IdValue: "us-east-1:ami-12345679abc"})
	if err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if !keep {
		t.Fatal("should keep")
	}
	if forceOverride {
		t.Fatal("should not override")
	}

	mockSSM.AssertExpectations(t)
}
