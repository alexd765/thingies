package main

import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

// Task is the input of our lambda function.
type Task struct {
	Version    int    `json:"version"`
	InputToken string `json:"input-token"`
}

// Result is the output of our lambda function
type Result struct {
	Version     int    `json:"version"`
	OutputToken string `json:"output-token"`
	Size        int    `json:"size"`
}

var s3Client *s3.S3

func init() {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("failed to get creds: %s", err)
	}
	sess := session.New(aws.NewConfig().WithRegion("eu-west-1").WithCredentials(creds))
	s3Client = s3.New(sess)
}

// Handle is our lambda function.
func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	var task Task
	if err := json.Unmarshal(evt, &task); err != nil {
		return nil, err
	}
	if err := validate(task); err != nil {
		return nil, err
	}

	obj, err := s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String("thingies-input"), Key: aws.String(task.InputToken)})
	if err != nil {
		return nil, err
	}
	var size int
	if obj.ContentLength != nil {
		size = int(*obj.ContentLength)
	}
	result := Result{
		Version:     2,
		OutputToken: task.InputToken,
		Size:        size,
	}
	return result, nil
}

func validate(task Task) error {
	if task.Version != 1 {
		return fmt.Errorf("wrong version: want 1, got %d", task.Version)
	}
	if task.InputToken == "0" {
		return errors.New("input token is empty")
	}
	return nil
}
