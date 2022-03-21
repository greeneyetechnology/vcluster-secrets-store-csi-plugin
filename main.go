package main

import (
	"github.com/Greeneye-Technology/vcluster-secrets-store-csi-plugin/syncers"
	"github.com/loft-sh/vcluster-sdk/plugin"
)

func main() {
	ctx := plugin.MustInit("secrets-store-csi-plugin")
	plugin.MustRegister(syncers.NewSecretStoreSyncer(ctx))
	plugin.MustStart()
}
