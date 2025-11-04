# Tekton

## Install `pass`

Store credentials in "pass" manager.

```sh
gpg --list-keys # get your public key id
pass init C25565E4701F4ED36A0711AA114F3606EFD923BB # Replace the string with the ID of your public GPG key
pass insert gh-username
pass ls
pass insert gh-password
pass ls
pass show gh-username
pass show gh-password
bin/create_secret.sh
```

* You probably want to do another Terraform apply at this point since the secret must exist for Terraform
  to apply the YAML that creates the service user.

## Install Tekton

```sh
kubectl config set-context --current --namespace=tekton-pipelines
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml
kubectl get namespace tekton-pipelines
# you can wait to see it RUNNING
watch kubectl get pods -n tekton-pipelines
```

* Delete the webhook configurations if you see an error like so:
`Error from server (InternalError): error when creating "./pipeline-test.yml": Internal error occurred: failed calling webhook "webhook.pipeline.tekton.dev":`

```sh
kubectl get mutatingwebhookconfigurations | grep tekton
kubectl get validatingwebhookconfigurations | grep tekton
kubectl delete mutatingwebhookconfigurations webhook.pipeline.tekton.dev
kubectl delete validatingwebhookconfigurations config.webhook.pipeline.tekton.dev
kubectl delete validatingwebhookconfigurations validation.webhook.pipeline.tekton.dev
```

## GCP Firewall

[Need to add firewall rule in GCP/GKE](https://stackoverflow.com/questions/59461747/tekton-on-private-kubernetes-cluster-on-gcp-gke)

* Name the rule something like `allow-tekton-webhook`

```sh
# https://github.com/kubernetes/kubernetes/issues/79739
resource "google_compute_firewall" "allow-tekton-webhook" {
  name        = "${var.name}-allow-tekton-webhook" # here is where you configure the name
  project     = var.project_id
  network     = google_compute_network.vpc.name
  description = "Allow Tekton webhook ingress traffic"

  allow {
    protocol = "tcp"
    ports = ["443",
      "8008",
      "8080",
      "8443",
      "9090",
      "9443" # GitHub webhook
    ]
  }
  allow {
    protocol = "udp"
  }

  direction = "INGRESS"
  # kubectl get endpoints -n tekton-pipelines
  # kubectl describe endpoints tekton-pipelines-webhook
  source_ranges = ["10.254.0.0/16"] # match the `master_ipv4_cidr_block` in your `private_cluster_config` in [terraform/main.tf](terraform/main.tf)
}
```
