package plugin

import (
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/framework"
)

func pathGSA(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Printf("gsa/%s", framework.GenericNameRegex("name")),
			"name": {
				Type: framework.TypeString,
				Description: "Required. Google Service Account Name",
			},
			ExistenceCheck: nil, // existence func
			Callbacks: map[logical.Operation]framework.OperationFunc {
				logical.CreateOperation: nil,
				logical.UpdateOperation: nil,
				logical.DeleteOperation: nil,
			},
		},
		{
			Pattern: fmt.Printf("gsa/%s/token/accessToken", framework.GenericNameRegex("name")),
			"name": {
				Type: framework.TypeString,
				Description: "Required. Google Service Account Name",
			},
			"delegates": {
				Type: framework.TypeCommaStringSlice,
				Description: "Optional. Specify delegation chain if a delegated request flow is used.",
			},
			"scope": {
				Type: framework.TypeCommaStringSlice,
				Description: "List of Oauth scopes to assign to credentials",
			},
			"lifetime": {
				Type: framework.TypeSignedDurationSecond,
				Description: "The duration of the access token in seconds",
			},
			ExistenceCheck: nil, // existence func
			Callbacks: map[logical.Operation]framework.OperationFunc {
				logical.CreateOperation: nil,
			},
		}
	}
}
