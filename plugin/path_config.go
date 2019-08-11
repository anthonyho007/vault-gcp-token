package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"gcp_credentials": {
				Type:        framework.TypeString,
				Description: "GCP Service Account credentials for vault plugin",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.readConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	gcpCredsRaw, ok := data.GetOk("gcp_credentials")
	if ok {
		gcpCreds, err := gcputil.Credentials(gcpCredsRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid credentials JSON: %v", err)), nil
		}
		config.GcpCredentials = gcpCreds
	}

	storageEntry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, storageEntry); err != nil {
		return nil, err
	}
	b.Logger().Debug("put google service account config")

	config, _ = b.readConfig(ctx, req)
	b.Logger().Debug("%s, %s", config.GcpCredentials.ClientEmail, config.GcpCredentials.ProjectId)
	return nil, nil
}

func (b *backend) readConfig(ctx context.Context, req *logical.Request) (*gcpConfig, error) {
	json, err := req.Storage.Get(ctx, "config")
	config := &gcpConfig{}

	if err != nil {
		return nil, err
	} else if json == nil {
		return config, nil
	}

	if err := json.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}
