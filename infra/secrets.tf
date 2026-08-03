locals {
  line_secrets = toset([
    "linebot-channel-secret",
    "linebot-channel-token",
  ])
}

resource "google_secret_manager_secret" "line" {
  for_each = local.line_secrets

  project             = var.project_id
  secret_id           = each.value
  deletion_protection = true

  replication {
    auto {}
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.required["secretmanager.googleapis.com"]]
}

resource "google_secret_manager_secret_iam_member" "runtime_accessor" {
  for_each = google_secret_manager_secret.line

  project   = var.project_id
  secret_id = each.value.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.runtime.member
}
