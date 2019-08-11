package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGSA(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Sprintf("gsa/%s", framework.GenericNameRegex("name")),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Required. Google Service Account Name",
				},
			},
			ExistenceCheck: nil, // existence func
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathGSACreateOrUpdate,
				logical.UpdateOperation: b.pathGSACreateOrUpdate,
				logical.DeleteOperation: nil,
			},
		},
		{
			Pattern: fmt.Sprintf("gsa/%s/generateAccessToken", framework.GenericNameRegex("name")),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Required. Google Service Account Name",
				},
				"delegates": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Optional. Specify delegation chain if a delegated request flow is used.",
				},
				"scope": {
					Type:        framework.TypeCommaStringSlice,
					Description: "List of Oauth scopes to assign to credentials",
				},
				"lifetime": {
					Type:        framework.TypeSignedDurationSecond,
					Description: "The duration of the oauth access token in seconds",
				},
			},
			ExistenceCheck: nil, // existence func
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: nil,
			},
		},
	}
}

func (b *backend) pathGSACreateOrUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rawName, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing service account name"), nil
	}

	name := rawName.(string)

	sa, err := getGcpSA(name, ctx, req)
	if err != nil {
		return nil, err
	}
	if sa == nil {
		accountID, err := parseGoogleServiceAccountEmail(name)
		if err != nil {
			return nil, err
		}

		sa = &gcpSA{
			Name:             accountID.EmailOrId,
			Project:          accountID.Project,
			ServiceAccountID: accountID,
		}
	}

	iamService, err := b.IamClient(req)
	if err != nil {
		return nil, err
	}

	googleServiceAccount, err := sa.getServiceAccount(iamService)
	if err != nil {
		return nil, err
	}
	if googleServiceAccount == nil {
		return logical.ErrorResponse("Failed to find google service account"), nil
	}
	b.Logger().Debug("found google service account")

	return nil, nil
}

func getGcpSA(name string, ctx context.Context, req *logical.Request) (*gcpSA, error) {
	data, err := req.Storage.Get(ctx, fmt.Sprintf("gcp/%s", name))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	sa := &gcpSA{}
	if err := data.DecodeJSON(sa); err != nil {
		return nil, err
	}

	return sa, nil
}
