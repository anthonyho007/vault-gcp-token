package plugin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

func (b *backend) IamClient(req *logical.Request) (*iam.Service, error) {
	b.Logger().Debug("trying to get http client")
	client, err := b.httpClient(req)
	if err != nil {
		return nil, err
	}
	b.Logger().Debug("get http client")
	iamClient, err := iam.New(client)
	if err != nil {
		return nil, err
	}

	iamClient.UserAgent = useragent.String()

	return iamClient, nil
}

func (b *backend) httpClient(req *logical.Request) (*http.Client, error) {
	ctx := context.Background()
	gcpCreds, err := b.getGcpCredentials(ctx, req)
	if err != nil {
		return nil, err
	}
	b.Logger().Debug("get gcp credentials")
	ctx = context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())
	return oauth2.NewClient(ctx, gcpCreds.TokenSource), nil
}

type gcpConfig struct {
	GcpCredentials *gcputil.GcpCredentials `gcp_credentials`
}

func (c gcpConfig) setCredentials(data *framework.FieldData) error {
	if json, ok := data.GetOk("gcp_credentials"); ok {
		credentials, err := gcputil.Credentials(json.(string))
		if err != nil {
			return err
		}

		c.GcpCredentials = credentials
	}
	return nil
}

func (c gcpConfig) getMarshalCredentials() ([]byte, error) {
	if c.GcpCredentials == nil {
		return []byte{}, nil
	}

	return json.Marshal(originalCredentials{
		GcpCredentials: c.GcpCredentials,
		Type:           "service_account",
	})
}

type originalCredentials struct {
	*gcputil.GcpCredentials
	Type string `json: "type"`
}

func (b *backend) getGcpCredentials(ctx context.Context, req *logical.Request) (*google.Credentials, error) {
	config, err := b.readConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	marshalGcpCreds, err := config.getMarshalCredentials()
	if err != nil {
		return nil, err
	}

	var gcpCreds *google.Credentials
	if len(marshalGcpCreds) > 0 {
		gcpCreds, err = google.CredentialsFromJSON(ctx, marshalGcpCreds, iam.CloudPlatformScope)
		if err != nil {
			return nil, err
		}
	} else {
		gcpCreds, err = google.FindDefaultCredentials(ctx, iam.CloudPlatformScope)
		if err != nil {
			return nil, err
		}
	}

	return gcpCreds, err
}
