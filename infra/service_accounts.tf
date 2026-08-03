resource "google_service_account" "runtime" {
  project      = var.project_id
  account_id   = "family-linebot-runtime"
  display_name = "family-linebot App Engine runtime"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_service_account" "deployer" {
  project      = var.project_id
  account_id   = "github-appengine-deployer"
  display_name = "GitHub App Engine deployer"

  lifecycle {
    prevent_destroy = true
  }
}

locals {
  runtime_project_roles = toset([
    "roles/aiplatform.user",
    "roles/cloudtranslate.user",
    "roles/speech.client",
  ])

  deployer_project_roles = toset([
    "roles/appengine.deployer",
    "roles/appengine.serviceAdmin",
    "roles/cloudbuild.builds.editor",
    "roles/serviceusage.serviceUsageConsumer",
    "roles/storage.objectAdmin",
  ])
}

resource "google_project_iam_member" "runtime" {
  for_each = local.runtime_project_roles

  project = var.project_id
  role    = each.value
  member  = google_service_account.runtime.member
}

resource "google_project_iam_member" "deployer" {
  for_each = local.deployer_project_roles

  project = var.project_id
  role    = each.value
  member  = google_service_account.deployer.member
}

resource "google_service_account_iam_member" "deployer_can_use_runtime" {
  service_account_id = google_service_account.runtime.name
  role               = "roles/iam.serviceAccountUser"
  member             = google_service_account.deployer.member
}
