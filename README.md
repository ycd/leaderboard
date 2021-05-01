
# Getting Started


# Deployment

## Deploy on Kubernetes with Terraform 

1. Download [Terraform](https://www.terraform.io/downloads.html)
1. Initialize terraform with `terraform init`
1. If you haven't got your credentials yet, get them with `gcloud auth application-default login`
1. Provision the cluster setup with `terraform apply -auto-approve`
1. Get the kubeconfig file by running the following command 

```bash
gcloud container clusters get-credentials $( echo var.name | terraform console ) --zone $( echo var.location | terraform console ) --project $( echo var.project | terraform console )
``` 

# API Reference
