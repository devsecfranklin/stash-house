# GKE Notes

## Setup

```sh
gcloud components install gke-gcloud-auth-plugin || sudo apt-get install google-cloud-sdk-gke-gcloud-auth-plugin
sudo apt install -y kubectl
gcloud components update # use this, or the following giant blob
sudo apt-get update && sudo apt-get --only-upgrade install google-cloud-cli-app-engine-go google-cloud-cli-pubsub-emulator google-cloud-cli-datastore-emulator google-cloud-cli google-cloud-cli-managed-flink-client google-cloud-cli-app-engine-java google-cloud-cli-cloud-run-proxy google-cloud-cli-minikube google-cloud-cli-cbt google-cloud-cli-terraform-tools kubectl google-cloud-cli-app-engine-python-extras google-cloud-cli-log-streaming google-cloud-cli-firestore-emulator google-cloud-cli-kpt google-cloud-cli-anthos-auth google-cloud-cli-package-go-module google-cloud-cli-nomos google-cloud-cli-spanner-migration-tool google-cloud-cli-local-extract google-cloud-cli-app-engine-python google-cloud-cli-istioctl google-cloud-cli-app-engine-grpc google-cloud-cli-config-connector google-cloud-cli-cloud-build-local google-cloud-cli-kubectl-oidc google-cloud-cli-docker-credential-gcr google-cloud-cli-enterprise-certificate-proxy google-cloud-cli-gke-gcloud-auth-plugin google-cloud-cli-spanner-emulator google-cloud-cli-skaffold google-cloud-cli-bigtable-emulator google-cloud-cli-anthoscli
gke-gcloud-auth-plugin --version
gcloud auth list
gcloud auth activate-service-account --key-file={$GOOGLE_APPLICATION_CREDENTIALS}
gcloud container clusters create --binauthz-evaluation-mode=PROJECT_SINGLETON_POLICY_ENFORCE --zone us-central1-a evening_tide
```

## Label Nodes

```sh
k cluster-info
k get nodes
kubectl label node gke-ps-devsecops-gke-build-pool-0ccb2515-yshf node-role.kubernetes.io/build="build"
kubectl label node gke-ps-devsecops-gke-dev-test-6959817f-3oo1  node-role.kubernetes.io/dev-test="dev-test"
```

- Select node like so:

```yaml
      nodeSelector:
        "node-role.kubernetes.io/dev-test": dev-test
```

## Access

- Add your IP as a /32 around line 132 in networks.tf
  - ONLY /32 are allowed, otherwise IT will contact us about a violation.
- Add your IP as a /32 cidr_block in main.tf
  - example:

```hcl
   cidr_blocks {
      cidr_block   = "84.207.227.14/32"
      display_name = "kuba-remote"
    }
```

## Operation

The dashboard is already installed, but here is how you could re-install it:

```sh
kubectl apply --filename https://github.com/tektoncd/dashboard/releases/latest/download/tekton-dashboard-release.yaml
```

Start the dashboard from your local machine like so:

```sh
kubectl --namespace tekton-pipelines port-forward svc/tekton-dashboard 9097:9097
```

Now you can [Navigate to the Dashboard](http://127.0.0.1:9097/#/pipelineruns) in your local browser.

![bugs](https://github.com/devsecfranklin/gke-cluster/blob/main/docs/images/draft2.jpg)

## Troubleshooting

- Stuck in Pending

```sh
kubectl get pods --field-selector=status.phase=Pending
kubectl describe pod ctfd-5cdb8d88f5-j6cmv
kubectl get pv
kubectl describe pv pvc-0df2fc1b-37e0-4258-85e7-67393065dde0
```
