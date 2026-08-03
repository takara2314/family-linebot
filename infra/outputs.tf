output "workload_identity_provider" {
  description = "Provider resource name used by google-github-actions/auth."
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "deployer_service_account" {
  value = google_service_account.deployer.email
}

output "runtime_service_account" {
  value = google_service_account.runtime.email
}
