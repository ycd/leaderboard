
# Getting Started


# Deployment

## Deploy on Kubernetes with Terraform 

1. Download [Terraform](https://www.terraform.io/downloads.html)
1. Initialize terraform with `terraform init`
1. If you haven't got your credentials yet, get them with `gcloud auth application-default login`
1. Provision the cluster setup with `terraform apply -auto-approve`
1. Get the kubeconfig file by running the following command 

```bash
gcloud container clusters get-credentials $( echo var.name | terraform console ) \
    --zone $( echo var.location | terraform console )   \
    --project $( echo var.project | terraform console )
``` 

6. After a few minutes of provisioning your cluster must be ready, test it with `kubectl cluster-info`.

## Deploy PostgreSQL on Kubernetes

1. Create Persistent volume. 

```bash
kubectl apply -f psql-persistent-vol.yaml
```

2. Apply The PostgreSQL deployment.

```bash
kubectl apply -f psql-service.yaml
```

3. Apply the PostgreSQL service. 

```bash
kubectl apply -f psql-service.yaml
```

## Deploying the services
TODO

# API Reference
