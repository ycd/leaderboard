
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

1. Create Deployment, ConfigMap and Service for Redis

```bash
kubectl apply -f redis-deployment-service.yaml
```

2. Create Deployment for API service.

```bash
kubectl apply -f api-deployment.yaml
```

3. Apply and expose the API service through a load balancer.

```bash
kubectl apply -f api-service.yaml
```

That's it.

### Creating mock data

You can create mock data for testing with the following command.

```bash
make insert-mock-data BASE_URL=http://HOST:PORT
``` 

# API Reference

### Get the global leaderboard

```http
  GET /leaderboard
```


#### Response 

```json
{
    "data": [
        {
            "rank": 1,
            "points": 99,
            "display_name": "test_451",
            "country": "ug"
        },
        {
            "rank": 2,
            "points": 98,
            "display_name": "test_376",
            "country": "fr"
        }
    ],
    "error": null,
    "success": true
}
```


### Get leaderboard of the specific country

```http
  GET /leaderboard/:country_code
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `country_code`      | `string` | **Required**. ISO 639-1 Code of the Country |

#### Response

```json
{
    "data": [
        {
            "rank": 1,
            "points": 99,
            "display_name": "test_443",
            "country": "fr"
        },
        {
            "rank": 2,
            "points": 98,
            "display_name": "test_376",
            "country": "fr"
        }
    ],
    "error": null,
    "success": true
}
```

### Submit Score

```http
  GET /score/submit
```

#### Request 

```json
{
	"score_worth": 1234.6,
	"user_id": "f1607032-aaf3-41df-9d61-84d06b97a322"
}
```

#### Response

```json
{
    "data": {
        "score_worth": 1234.6,
        "user_id": "f1607032-aaf3-41df-9d61-84d06b97a322",
        "timestamp": 1619972382
    },
    "error": "",
    "success": true
}
```

### Create new user

```http
  GET /user/create
```


#### Request

```json
{
	"display_name": "ycd_123",
	"country": "us"
}
```

#### Response

```json
{
    "data": {
        "user_id": "0db88bec-6b99-4972-b72e-c35c9ee38cf9",
        "display_name": "ycd_123",
        "points": 0,
        "rank": 100
    },
    "error": "",
    "success": true
}
```


### Get user profile

```http
  GET /user/profile/:guid
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `guid`      | `string` | **Required**. GUID(UUID) |

#### Response

```json
{
    "data": {
        "user_id": "64f6faea-db27-415f-a7bc-6b8cd8e893d1",
        "display_name": "yaguasdf",
        "points": 0,
        "rank": 100
    },
    "error": "",
    "success": true
}
```
