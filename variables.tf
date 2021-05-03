variable "name" {
  type = string 
  default = "leaderboard-cluster"
}

variable "project" {
  type = string
  # Change this with your own project name
  default = "leaderboard-312410"
}

variable "location" {
  type = string 
  default = "europe-west3-a"
}

variable "initial_node_count" {
  type = number
  default = 2
}

variable "machine_type" {
  type = string 
  # 4vCPU - 16GB Memory
  default = "e2-standard-8"
}
