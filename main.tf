resource "google_container_cluster" "default" {
  name        = var.name
  project     = var.project
  description = "Leaderboard GKE Cluster"
  location    = var.location

  remove_default_node_pool = true
  initial_node_count       = var.initial_node_count

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

resource "google_container_node_pool" "default" {
  name       = "${var.name}-node-pool"
  project    = var.project
  location   = var.location
  cluster    = google_container_cluster.default.name
  node_count = 1


  autoscaling {
    max_node_count = 4
    min_node_count = 2
  }

  node_config {
    preemptible  = true
    machine_type = var.machine_type

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}

provider "kubernetes" {
  host  = "https://${data.google_container_cluster.leaderboard_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.leaderboard_cluster.master_auth[0].cluster_ca_certificate,
  )
}
