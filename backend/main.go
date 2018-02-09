package main

import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/harrydb/go/img/grayscale"
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
}

func handle(task Task) (Result, error) {
	if err := validate(task); err != nil {
		return Result{}, err
	}

	creds := credentials.NewEnvCredentials()
	if _, err := creds.Get(); err != nil {
		log.Fatalf("failed to get creds: %s", err)
	}
	sess := session.New(aws.NewConfig().WithRegion("eu-west-1").WithCredentials(creds))
	s3Client := s3.New(sess)

	result := Result{
		Version:     3,
		OutputToken: strings.TrimSuffix(task.InputToken, filepath.Ext(task.InputToken)) + ".png",
	}

	obj, err := s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String("thingies-input"), Key: aws.String(task.InputToken)})
	if err != nil {
		return Result{}, err
	}
	defer obj.Body.Close()
	img, _, err := image.Decode(obj.Body)
	if err != nil {
		return Result{}, err
	}
	grayImg := grayscale.Convert(img, grayscale.ToGrayLuminance)
	var buf bytes.Buffer
	if err := png.Encode(&buf, grayImg); err != nil {
		return Result{}, err
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("thingies-output"),
		Key:    aws.String(result.OutputToken),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func validate(task Task) error {
	if task.Version != 1 {
		return fmt.Errorf("wrong version: want 1, got %d", task.Version)
	}
	if task.InputToken == "" {
		return errors.New("input token is empty")
	}
	return nil
}

func main() {
	lambda.Start(handle)
}
