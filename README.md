# loadtest-controller

## Install Tools

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) - Kubernetes CLI
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) - Local Kubernetes clusters (See NOTE below)
- [jq](https://stedolan.github.io/jq/download/) - JSON processor
- [vegeta](https://github.com/tsenart/vegeta#install) - Load testing tool

NOTE: The `kind` tool relies on docker, if you can not use docker, there is experimental support for [kind with podman](https://kind.sigs.k8s.io/docs/user/rootless/#creating-a-kind-cluster-with-rootless-podman).

## Setup

Create a local cluster.

```sh
kind create cluster --name loadtesting
```
