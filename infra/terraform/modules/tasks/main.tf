resource "google_cloud_tasks_queue" "processing" {
  name     = "actionrev-processing-${var.env}"
  location = var.region

  rate_limits {
    max_concurrent_dispatches = 10
    max_dispatches_per_second = 5
  }

  retry_config {
    max_attempts  = 3
    min_backoff   = "10s"
    max_backoff   = "300s"
    max_doublings = 4
  }

  stackdriver_logging_config {
    sampling_ratio = 1.0
  }
}
