package plugin

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"google.golang.org/api/iam/v1"
)

const (
	serviceAccountRegex = "(?m)[a-z]([-a-z0-9]*[a-z0-9])"
)

type gcpSA struct {
	Name             string
	Project          string
	ServiceAccountID *gcputil.ServiceAccountId
}

func (sa gcpSA) getServiceAccount(iamService *iam.Service) (*iam.ServiceAccount, error) {
	if sa.ServiceAccountID == nil {
		return nil, fmt.Errorf("gcpSA has no service account ID")
	}

	serviceAccount, err := iamService.Projects.ServiceAccounts.Get(sa.ServiceAccountID.ResourceName()).Do()
	if err != nil {
		return nil, nil
	}
	return serviceAccount, nil
}

func (sa gcpSA) setServiceAccountID(name string) error {
	serviceAccountID, err := parseGoogleServiceAccountEmail(name)
	if err != nil {
		return err
	}
	sa.ServiceAccountID = serviceAccountID
	return nil
}

func (sa gcpSA) createOrUpdateIamPolicies() {

}

func parseGoogleServiceAccountEmail(name string) (*gcputil.ServiceAccountId, error) {
	saRegex := regexp.MustCompile(serviceAccountRegex)
	matches := saRegex.FindAllString(name, -1)
	if len(matches) != 5 {
		return nil, errors.New("invalid google service account format")
	}
	return &gcputil.ServiceAccountId{
		Project:   matches[1],
		EmailOrId: name,
	}, nil
}

// "projects/serviceAccounts": {
// 	"iam": {
// 		"v1": IamRestResource{
// 			Name:                      "serviceAccounts",
// 			TypeKey:                   "projects/serviceAccounts",
// 			Service:                   "iam",
// 			IsPreferredVersion:        true,
// 			Parameters:                []string{"resource"},
// 			CollectionReplacementKeys: map[string]string{},
// 			GetMethod: RestMethod{
// 				HttpMethod: "POST",
// 				BaseURL:    "https://iam.googleapis.com/",
// 				Path:       "v1/{+resource}:getIamPolicy",
// 			},
// 			SetMethod: RestMethod{
// 				HttpMethod:    "POST",
// 				BaseURL:       "https://iam.googleapis.com/",
// 				Path:          "v1/{+resource}:setIamPolicy",
// 				RequestFormat: `{"policy": %s}`,
// 			},
// 		},
// 	},
// },
