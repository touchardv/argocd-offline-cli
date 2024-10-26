# argocd-offline-cli

A small [Argo CD](https://argo-cd.readthedocs.io/en/stable/) CLI utility, built on top of Argo CD Go packages, that can be used to preview the Kubernetes resources generated from an [ApplicationSet](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/applicationset-specification/), without the need to connect to an actual Argo CD server.

## Requirements

* A recent version of [Helm v3](https://helm.sh/).

## Limitations

Only a few [generators](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/Generators/) and Helm source repositories are supported.

## Usage

### Configuration

The `HELM_REPO_USERNAME` and `HELM_REPO_PASSWORD` environment variables can be specified in order to provide the default credentials that should be used to authenticate to Helm repositories. If not specified, the local `helm` command settings may be used to authenticate (if present).

### Preview Application(s) from an ApplicationSet

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest
```

#### Example: filter by application name

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest -n app-name
```

#### Example: filter by application name, display in yaml format

```shell
argocd-offline-cli appset preview-apps /path/to/application-set-manifest -n app-name -o yaml
```

### Preview Resource manifest(s) from an ApplicationSet

```shell
argocd-offline-cli appset preview-resources /path/to/application-set-manifest
```
