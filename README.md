# CSI secret store plugin for vcluster

The goal of this plugin is to support secrets coming from CSI secret store in vcluster. 

## Requirements

- [Install secret store CSI driver](https://azure.github.io/secrets-store-csi-driver-provider-azure/docs/getting-started/installation/#deployment-using-helm) on host machine: 
```shell script
helm repo add csi-secrets-store-provider-azure https://azure.github.io/secrets-store-csi-driver-provider-azure/charts
helm install csi csi-secrets-store-provider-azure/csi-secrets-store-provider-azure
```
- [Install jsPolicy](https://www.jspolicy.com/docs/quickstart) on host machine
```shell script
helm install jspolicy jspolicy -n jspolicy --create-namespace --repo https://charts.loft.sh
```

## Description

The plugin will create all required `SecretProviderClass` on the host machine. However, currently the plugin cannot change the pod's `VolumeAttributes`. Therefore, jsPolicy is required.  This might change once [https://github.com/loft-sh/vcluster/issues/465](https://github.com/loft-sh/vcluster/issues/465) will be added to vcluster.
