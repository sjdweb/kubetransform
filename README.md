# Kubetransform

Kubernetes utilities to transform YAML files.

## Install
```
go get -u github.com/sjdweb/kubetransform
```

## Commands

### Deployment
Upgrade v1.5-1.7 `initContainers` annotations to the new `initContainers` spec.

Example
```bash
cat deployment.yml | kubetransform deployment
```

### Secret
Decode secrets.
```bash
echo "$(kubectx get secret thing -o yaml)" | kubetransform secret
```