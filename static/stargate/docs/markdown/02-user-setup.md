# User Setup

Each user gets:

* A dedicated namespace
* A pre-configured amount of CPU and Memory to use.

## Cluster Credentials

* Add cluster credentials to `~/.config/gcloud`
* Verify you can access the cluster nodes.

```sh
alias k kubectl
gcloud container clusters get-credentials secure-dev-lab --zone us-central1
k config get-clusters
k config set-cluster gke_gcp-gcs-pso_us-central1_secure-dev-lab
k get nodes
```

## Certificates

* Generate user certificates. Move the resultant files to your `~/.kube` folder.
* Add user role and binding.

```sh
go run bin/csr-gen.go client <user-name>;
bin/csr-tool.sh <user-name> client.csr
bin/role.sh <user-name>
```

## Tooling

Toolss you may find useful.

* [Install Tekton CLI](https://github.com/tektoncd/cli#installing-tkn)
* Optional tools: helm, skaffold
* Optional VSCode Plugins
  * [Cloud Code](https://marketplace.visualstudio.com/items?itemName=GoogleCloudTools.cloudcode)
  * [Draw Integration](https://marketplace.visualstudio.com/items?itemName=hediet.vscode-drawio) for diagrams you can save in your repo.
  * [Golang](https://marketplace.visualstudio.com/items?itemName=golang.Go)
  * [Kubernetes](https://marketplace.visualstudio.com/items?itemName=ms-kubernetes-tools.vscode-kubernetes-tools)
  * [Python](https://marketplace.visualstudio.com/items?itemName=ms-python.python)
  * [YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
