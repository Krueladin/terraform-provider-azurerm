package azure

import (
	"fmt"
	//	"log"
	"regexp"
	//	"github.com/hashicorp/terraform/helper/schema"
	//	"github.com/hashicorp/terraform/helper/validation"
)

// This should later be moced to the validators in the provider
func ValidateFrontDoorName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[a-zA-Z0-9]+([-a-zA-Z0-9]?[a-zA-Z0-9])*$").MatchString(input) {
		errors = append(errors, fmt.Errorf("The name can contain only letters, numbers, and hyphens and must start with a letter or number"))
	}
	if len(input) < 5 || len(input) > 64 {
		errors = append(errors, fmt.Errorf("The name must be between 5 and 64 characters long"))
	}

	return warnings, errors
}
func ValidateBackendPoolName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[a-zA-Z0-9]+([-a-zA-Z0-9]?[a-zA-Z0-9])*$").MatchString(input) {
		errors = append(errors, fmt.Errorf("The name can contain only letters, numbers, and hyphens and must start with a letter or number"))
	}
	if len(input) < 1 || len(input) > 90 {
		errors = append(errors, fmt.Errorf("The name must be between 5 and 64 characters long"))
	}

	return warnings, errors
}
func FrontDoorsNameCheckType() string {
	return "MicrosoftNetworkFrontDoors"
}
