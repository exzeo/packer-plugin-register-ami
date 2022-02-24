package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type StoreInterface interface {
	SaveParameter(string, string) *error
}

// Store is a wrapper of aws-sdk client
type Store struct {
	ssmconn ssmiface.SSMAPI
	config  Config
	now     time.Time
}

func NewStore(sess *session.Session, config Config) (*Store, error) {
	store := &Store{
		ssmconn: ssm.New(sess),
		config:  config,
		now:     time.Now().UTC(),
	}

	return store, nil
}

func (s *Store) SaveParameter(name string, ami string) *error {

	input := &ssm.PutParameterInput{
		Name:      aws.String(name),
		DataType:  aws.String("aws:ec2:image"),
		Overwrite: aws.Bool(true),
		Value:     aws.String(ami),
	}

	_, err := s.ssmconn.PutParameter(input)

	if err != nil {
		return &err
	}

	return nil
}
