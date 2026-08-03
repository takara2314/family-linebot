resource "google_app_engine_application" "default" {
  project     = var.project_id
  location_id = "asia-northeast2"

  depends_on = [google_project_service.required["appengine.googleapis.com"]]

  lifecycle {
    prevent_destroy = true
  }
}
