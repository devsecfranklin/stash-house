# Tekton Configuration

* Add a task.

```sh
kubectl config set-context --current --namespace=tekton-pipelines

```

## Dashboard

```sh
kubectl apply --filename https://github.com/tektoncd/dashboard/releases/latest/download/tekton-dashboard-release.yaml
```
