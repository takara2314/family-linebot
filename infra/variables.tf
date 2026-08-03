variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "region" {
  description = "Default Google Cloud region."
  type        = string
}

variable "github_repository_id" {
  description = "Immutable numeric GitHub repository ID."
  type        = string
}

variable "github_repository_owner_id" {
  description = "Immutable numeric GitHub repository owner ID."
  type        = string
}

variable "github_deploy_ref" {
  description = "Git ref allowed to deploy."
  type        = string
}
