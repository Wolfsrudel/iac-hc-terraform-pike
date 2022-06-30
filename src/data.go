package pike

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
)

// GetAPI determines which APIs a resource targets
func GetAPI(resource string) string {

	m := map[string]string{
		"aws_instance":                           "ec2",
		"aws_s3_bucket":                          "s3",
		"aws_s3_bucket_accelerate_configuration": "s3",
		"aws_s3_bucket_policy":                   "s3",
		"aws_s3_bucket_logging":                  "s3",
		"aws_vpc":                                "vpc",
	}

	api := m[resource]
	if api == "" {
		log.Printf("Resource:%s does not map", resource)
		return ""
	}

	return api
}

// GetResources retrieves all the resources in a tf file
func GetResources(file fs.FileInfo, dirname string) ([]string, string, string) {
	var result []string

	src, _ := ioutil.ReadFile(dirname + file.Name())
	resources := strings.Split(string(src), "}")

	for _, resource := range resources {
		if len(resource) > 1 {
			s := strings.Split(resource, "\"")
			verb, template := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
			if verb == "resource" {
				result = append(result, template)
			}
		}
	}

	return result, dirname + file.Name(), string(src)
}

// GetProvider retrieves the provider from the resource
func GetProvider(resource string) string {
	return strings.Split(resource, "_")[0]
}
