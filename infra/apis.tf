locals {
  required_services = toset([
    "aiplatform.googleapis.com",
    "appengine.googleapis.com",
    "cloudbuild.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iamcredentials.googleapis.com",
    "secretmanager.googleapis.com",
    "speech.googleapis.com",
    "sts.googleapis.com",
    "translate.googleapis.com",
  ])
}

resource "google_project_service" "required" {
  for_each = local.required_services

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}
