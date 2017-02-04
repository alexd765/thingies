package main

import "C"

import (
	"encoding/json"
	"errors"
	"fmt"

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

	result := Result{
		Version:     1,
		OutputToken: task.InputToken,
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
