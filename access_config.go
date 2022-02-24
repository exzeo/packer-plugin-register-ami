package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/go-cleanhttp"
)

// AccessConfig is for common configuration related to AWS access
type AccessConfig struct {
	AccessKey   string `mapstructure:"access_key"`
	ProfileName string `mapstructure:"profile"`
	SecretKey   string `mapstructure:"secret_key"`

	session *session.Session
}

// Session returns a valid session.Session object for access to AWS services, or
// an error if the authentication and region couldn't be resolved
func (c *AccessConfig) Session() (*session.Session, error) {
	if c.session != nil {
		return c.session, nil
	}

	// Create new AWS config
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)

	config = config.WithHTTPClient(cleanhttp.DefaultClient())
	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.Proxy = http.ProxyFromEnvironment

	// Figure out which possible credential providers are valid; test that we
	// can get credentials via the selected providers, and set the providers in
	// the config.
	creds, err := c.GetCredentials(config)
	if err != nil {
		return nil, err
	}
	config.WithCredentials(creds)

	// Create session options based on our AWS config
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}

	if c.ProfileName != "" {
		opts.Profile = c.ProfileName
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}

	log.Printf("Found region %s", *sess.Config.Region)
	c.session = sess

	cp, err := c.session.Config.Credentials.Get()

	if IsAWSErr(err, "NoCredentialProviders", "") {
		return nil, c.NewNoValidCredentialSourcesError(err)
	}

	if err != nil {
		return nil, fmt.Errorf("error loading credentials for AWS Provider: %s", err)
	}

	log.Printf("[INFO] AWS Auth provider used: %q", cp.ProviderName)

	return c.session, nil
}

// GetCredentials gets credentials from the environment, shared credentials,
// the session (which may include a credential process), or ECS/EC2 metadata
// endpoints. GetCredentials also validates the credentials and the ability to
// assume a role or will return an error if unsuccessful.
func (c *AccessConfig) GetCredentials(config *aws.Config) (*awsCredentials.Credentials, error) {
	// Reload values into the config used by the Packer-Terraform shared SDK
	awsbaseConfig := &awsbase.Config{
		AccessKey:    c.AccessKey,
		DebugLogging: false,
		Profile:      c.ProfileName,
		SecretKey:    c.SecretKey,
	}

	return awsbase.GetCredentials(awsbaseConfig)
}

// IsAWSErr returns true if the error matches all these conditions:
//  * err is of type awserr.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
func IsAWSErr(err error, code string, message string) bool {
	if err, ok := err.(awserr.Error); ok {
		return err.Code() == code && strings.Contains(err.Message(), message)
	}
	return false
}

// NewNoValidCredentialSourcesError returns user-friendly errors for authentication failed.
func (c *AccessConfig) NewNoValidCredentialSourcesError(err error) error {
	return fmt.Errorf("no valid credential sources found for register-ami post processor. "+
		"Error: %w", err)
}
