package pike

import (
	"bytes"
	_ "embed" //required for embed
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed terraform.policy.template
var policyTemplate []byte

//go:embed aws_iam_role.tf
var roleTemplate []byte

// NewAWSPolicy constructor
func NewAWSPolicy(Actions []string) Policy {
	something := Policy{Version: "2012-10-17"}

	sort.Strings(Actions)

	var categories []string
	for _, action := range Actions {
		prefix := strings.Split(action, ":")[0]
		categories = append(categories, prefix)
	}

	sections := unique(categories)
	var statements []Statement

	for count, section := range sections {
		var myActions []string
		for _, action := range Actions {
			mySection := section + ":"
			if strings.Contains(action, mySection) {
				myActions = append(myActions, action)
			}
		}

		state := Statement{Sid: "VisualEditor" + strconv.Itoa(count), Effect: "Allow", Action: myActions, Resource: "*"}

		statements = append(statements, state)
	}

	something.Statements = statements
	return something
}

// GetPolicy creates new iam polices from a list of Permissions
func GetPolicy(actions Sorted) (OutputPolicy, error) {
	var OutPolicy OutputPolicy
	var Empty bool
	Empty = true

	v := reflect.ValueOf(actions)
	typeOfV := v.Type()
	values := make([]interface{}, v.NumField())

	var err error
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		switch typeOfV.Field(i).Name {
		case "AWS":
			if actions.AWS == nil {
				continue
			}

			Empty = false
			//dedupe
			AWSPermissions := unique(actions.AWS)
			OutPolicy.AWS, err = AWSPolicy(AWSPermissions)

			if err != nil {
				log.Print(err)
				continue
			}

		case "GCP":
			if actions.GCP == nil {
				continue
			}

			Empty = false
			//dedupe
			GCPPermissions := unique(actions.GCP)
			OutPolicy.GCP, err = GCPPolicy(GCPPermissions)
			if err != nil {
				log.Print(err)
				continue
			}

		case "AZURE":
			if actions.AZURE == nil {
				continue
			}

			Empty = false
			//dedupe
			AZUREPermissions := unique(actions.AZURE)
			OutPolicy.AZURE, err = AZUREPolicy(AZUREPermissions)
			if err != nil {
				log.Print(err)
				continue
			}
		}

	}
	if Empty {
		return OutPolicy, errors.New("no permissions found")
	}
	return OutPolicy, nil
}

// AWSPolicy create an IAM policy
func AWSPolicy(Permissions []string) (AwsOutput, error) {
	var OutPolicy AwsOutput
	Policy := NewAWSPolicy(Permissions)
	b, err := json.MarshalIndent(Policy, "", "    ")
	if err != nil {
		fmt.Println(err)
		return OutPolicy, err
	}

	type PolicyDetails struct {
		Policy      string
		Name        string
		Path        string
		Description string
	}

	//PolicyName := "terraform" + randSeq(8)
	PolicyName := "terraform_pike"
	theDetails := PolicyDetails{string(b), PolicyName, "/", "Pike Autogenerated policy from IAC"}

	var output bytes.Buffer
	tmpl, err := template.New("test").Parse(string(policyTemplate))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&output, theDetails)

	if err != nil {
		panic(err)
	}
	OutPolicy.Terraform = output.String()
	OutPolicy.JSONOut = string(b) + "\n"

	return OutPolicy, nil
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	sort.Strings(result)
	return result
}
