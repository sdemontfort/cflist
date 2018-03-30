package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	svc := cloudformation.New(sess)

	args := os.Args[1:]
	cmd := args[0]

	if cmd == "list" {
		filter := args[1]
		listStacks(svc, filter)
	} else if cmd == "diff" {
		diffStack(svc)
	}
}

func listStacks(svc *cloudformation.CloudFormation, filter string) {
	resp, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		return
	}

	var counter = 1

	fmt.Println(fmt.Sprintf("Cloudformation stacks matching \"%s\":", filter))
	for _, stack := range resp.Stacks {
		if fuzzy.Match(filter, *stack.StackName) {
			fmt.Println(fmt.Sprintf("%d: %s", counter, *stack.StackName))
			counter++
		}
	}
}

// Diffs a change set with the existing stack template and parameters
func diffStack(svc *cloudformation.CloudFormation) error {
	resp, err := svc.GetTemplate(&cloudformation.GetTemplateInput{
		StackName: aws.String("stack-name-here"),
	})
	if err != nil {
		return err
	}

	template := string(*resp.TemplateBody)
	oldTemplate := ``
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(template, oldTemplate, false)
	fmt.Println(dmp.DiffPrettyText(diffs))

	return nil
}
