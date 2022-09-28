# CSI secret store plugin for vcluster

The goal of this plugin is to support secrets coming from CSI secret store in vcluster.

The plugin will create all required `SecretProviderClass` on the host machine. 

## Requirements

- [Install secret store CSI driver](https://azure.github.io/secrets-store-csi-driver-provider-azure/docs/getting-started/installation/#deployment-using-helm) on host machine: 
```shell script
helm repo add csi-secrets-store-provider-azure https://azure.github.io/secrets-store-csi-driver-provider-azure/charts
helm install csi csi-secrets-store-provider-azure/csi-secrets-store-provider-azure
```

## Usage

Create a `plugin.yaml` with the following content:

```yaml
vcluster-secrets-store-csi-plugin:
    image: shakedos/vcluster-secrets-store-csi-plugin
    imagePullPolicy: IfNotPresent
    rbac:
      role:
        extraRules:
          - apiGroups: ["secrets-store.csi.x-k8s.io"]
            resources: ["secretproviderclasses"]
            verbs: ["create", "delete", "patch", "update", "get", "list", "watch"]
      clusterRole:
        extraRules:
          - apiGroups: ["apiextensions.k8s.io"]
            resources: ["customresourcedefinitions"]
            verbs: ["get", "list", "watch"]
```

Once your vcluster is up and running, make sure you have the `vcluster-secrets-store-csi-plugin` running:

```shell script
kdp vcluster-0 | grep secrets-store-csi-plugin 
```

## Example

```shell
export TENANT_ID=xxx
export CLIENT_ID=xxx
export CLIENT_SECRET=xxx
```

```yaml
vcluster connect -n csi vcluster -- kubectl apply -f - <<EOF
apiVersion: secrets-store.csi.x-k8s.io/v1alpha1
kind: SecretProviderClass
metadata:
  name: test-spc
spec:
  provider: azure
  parameters:
    keyvaultName: "stage-vaut"
    objects: |
      array:
        - |
          objectAlias: ca.pem
          objectName: myca-cert
          objectType: secret
    tenantId: "$TENANT_ID"
---
apiVersion: v1
kind: Pod
metadata:
  name: secret-store-test-pod
spec:
  containers:
    - name: secret-store-test-pod
      imagePullPolicy: Always
      image: ubuntu:20.04
      command: ["sleep"]
      args: ['infinity']
      volumeMounts:
        - name: secrets-store-inline
          mountPath: /secrets
  volumes:
    - name: secrets-store-inline
      csi:
        driver: secrets-store.csi.k8s.io
        readOnly: true
        volumeAttributes:
          secretProviderClass: "test-spc"
        nodePublishSecretRef:
          name: kvcreds
---
apiVersion: v1
data:
  clientid: $CLIENT_ID
  clientsecret: $CLIENT_SECRET
kind: Secret
metadata:
  name: kvcreds
EOF
```

You should see it up and running in the host cluster:

```shell script
k get pods | grep secret-store
NAME                                                         READY   STATUS    RESTARTS   AGE
secret-store-test-pod-x-default-x-vcluster                   1/1     Running   0          5m48s   10.42.0.183     lima-rancher-desktop   <none>           <none>
```

And also in the virtual cluster:

```shell script
vcluster connect -n csi vcluster -- kubectl get pods
NAME           READY   STATUS    RESTARTS   AGE
test-shaked1   1/1     Running   0          6m41s
```
