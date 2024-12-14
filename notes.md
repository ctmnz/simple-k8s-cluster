## Simple

# Create CSR

```
openssl genrsa -out daniel.key 2048
openssl req -new -key daniel.key -out daniel.csr -subj "/CN=daniel/O=app1/O=app2"
```

# Approve CSR

```bash
openssl x509 -req -in daniel.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial -out daniel.crt -days 365

cat daniel.csr | base64 | tr -d "\n" | pbcopy
```

### replace the bas64-encoded-csr with the key created and copied into the clipboard above.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: daniel-csr
spec:
  request: $(cat daniel.csr | base64 | tr -d '\n')
  signerName: kubernetes.io/kube-apiserver-client
  expirationSeconds: 86400  # one day
  usages:
  - client auth
EOF
```


### Apply, check and approve the csr and extract the crt
```bash
kubectl get csr
kubectl certificate approve daniel-csr
kubectl get csr daniel-csr -o jsonpath='{.status.certificate}'| base64 -d > daniel.crt
kubectl config set-credentials daniel --client-key=daniel.key --client-certificate=daniel.crt --embed-certs=true
kubectl config set-context local-cluster-lab --cluster=kind-simple-cluster --user=daniel
kubectl config use-context local-cluster-lab
kubectl config view --raw --minify --flatten > daniel-kubeconfig.yaml
```

## Create and bind a role
### Docs: https://kubernetes.io/docs/reference/kubectl/generated/kubectl_create/kubectl_create_role/
## Global Admin bind the builtin role cluster-admin to user daniel
```bash
kubectl create clusterrolebinding daniel-admin-binding --clusterrole=cluster-admin --user=daniel
```

## More granular role

```bash
kubectl create ns role-lab01-ns
kubectl create role pod-reader --verb=get --verb=list --verb=watch --resource=pods -n role-lab01-ns
kubectl create rolebinding daniel-admin-binding --role=pod-reader --user=daniel -n role-lab01-ns
```

## Same but with yaml

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: role-lab01-ns
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - list
  - watch
```

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: daniel-admin-binding
  namespace: role-lab01-ns
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: pod-reader
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: daniel
```



# Test the role

# positive
```bash
kubectl get po -n role-lab01-ns --kubeconfig daniel-kube.yaml
```

# false positive
```bash
kubectl run po -n role-lab01-ns --image=ubuntu  --kubeconfig daniel-kube.yaml
```


