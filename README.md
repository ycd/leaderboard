
# Getting Started


# Development

## Running locally

### Prerequisites

1. GNU Make >=3.5
1. Docker Engine >= 20
1. Docker Compose

#### Build image

```bash
make build
```

#### Run locally 

```bash
make run-dev
```

#### Testing

1. Create a temporary test database

```bash
make test-db
```

2. Run tests

```bash
go test ./...
```

# Deployment

## Provisioning a new Kubernetes Cluster with Terraform 

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

1. Create Deployment for API service.

```bash
kubectl apply -f api-deployment.yaml
```

2. Apply and expose the API service through a load balancer.

```bash
kubectl apply -f api-service.yaml
```


### Creating mock data

You can create mock data for testing with the following command.

```bash
make insert-mock-data BASE_URL=http://HOST:PORT
``` 

# API Reference

