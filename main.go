package main

import (
	"os"

	"github.com/anthonyho007/vault-gcp-token/plugin"
	log "github.com/golang/glog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: plugin.Factory, // put the backend here
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Error("failed to initialize vault-gcp-token plugin")
		os.Exit(1)
	}
}
