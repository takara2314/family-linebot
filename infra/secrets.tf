resource "google_secret_manager_secret" "line_channel_secret" {
  project             = var.project_id
  secret_id           = "linebot-channel-secret"
  deletion_protection = true

  replication {
    auto {}
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.required["secretmanager.googleapis.com"]]
}

resource "google_secret_manager_secret" "line_channel_token" {
  project             = var.project_id
  secret_id           = "linebot-channel-token"
  deletion_protection = true

  replication {
    auto {}
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.required["secretmanager.googleapis.com"]]
}

resource "google_secret_manager_secret_iam_member" "runtime_accessor_channel_secret" {
  project   = var.project_id
  secret_id = google_secret_manager_secret.line_channel_secret.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.runtime.member
}

resource "google_secret_manager_secret_iam_member" "runtime_accessor_channel_token" {
  project   = var.project_id
  secret_id = google_secret_manager_secret.line_channel_token.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.runtime.member
}
