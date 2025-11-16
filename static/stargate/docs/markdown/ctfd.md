# Set up the CTFd Instance in GKE

- Do not use node selector, it will need a match on ALL labels.
- You must set a public IP address with`gcloud compute addresses`

```sh
kubectl create ns ctfd
kubectl config set-context --current --namespace=ctfd
kubectl apply -f ctfd-mysql-deployment.yaml
kubectl apply -f ctfd-redis-deployment.yaml
kubectl apply -f ctfd-deployment.yaml
kubectl apply -f ctfd-backend.yaml
kubectl apply -f ctfd-nginx-deployment.yaml
```

- Checkout the setup so far:

```sh
kubectl get all
kubectl describe deployments -n ctfd
validate-ctfd.sh
```

## Test connection

This check simply shows that everything is working correctly.

```sh
kubectl port-forward service/ctfd-nginx -n ctfd
kubectl port-forward service/ctfd-nginx 8080:8080 -n ctfd # this will make it so you can browse to the site
```

- Navigate to `http://127.0.0.1:8080/setup`

## Expose Service

- Also apply the `ingress.yaml`, verify it like so: `kubectl get ingress -n ctfd`
- Run the Terraform to provision the backend service.

```sh
kubectl get svcneg -n ctfd
```
