terraform {
  required_version = ">= 1.7"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }

  # terraform init -backend-config=backends/<env>.hcl
  backend "gcs" {}
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

locals {
  labels = {
    env     = var.env
    app     = "actionrev"
    managed = "terraform"
  }
}

# ---------------------------------------------------------------------------
# GCP APIs
# ---------------------------------------------------------------------------

resource "google_project_service" "apis" {
  for_each = toset([
    "run.googleapis.com",
    "storage.googleapis.com",
    "bigquery.googleapis.com",
    "aiplatform.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudtasks.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "firebase.googleapis.com",
    "identitytoolkit.googleapis.com",
  ])

  project                    = var.project_id
  service                    = each.value
  disable_on_destroy         = false
  disable_dependent_services = false
}

# ---------------------------------------------------------------------------
# Modules
# 依存関係: platform → pipeline → api
#           datastore は独立
# ---------------------------------------------------------------------------

module "platform" {
  source     = "./modules/platform"
  project_id = var.project_id
  region     = var.region
  labels     = local.labels
  depends_on = [google_project_service.apis]
}

module "datastore" {
  source     = "./modules/datastore"
  project_id = var.project_id
  region     = var.region
  labels     = local.labels
  depends_on = [google_project_service.apis]
}

module "pipeline" {
  source            = "./modules/pipeline"
  project_id        = var.project_id
  region            = var.region
  env               = var.env
  labels            = local.labels
  sandbox_sa_email  = module.platform.sandbox_sa_email
  depends_on        = [module.platform]
}

module "api" {
  source               = "./modules/api"
  project_id           = var.project_id
  region               = var.region
  env                  = var.env
  labels               = local.labels
  backend_image        = var.backend_image
  backend_sa_email     = module.platform.backend_sa_email
  uploads_bucket_name  = module.datastore.uploads_bucket_name
  bigquery_dataset_id  = module.datastore.dataset_id
  cloud_tasks_queue    = module.pipeline.queue_name
  sandbox_job_name     = module.pipeline.sandbox_job_name
  sandbox_bucket_name  = module.pipeline.sandbox_bucket_name
  min_instances        = var.backend_min_instances
  max_instances        = var.backend_max_instances
  gemini_model         = var.gemini_model

  stripe_secret_key_id          = module.platform.stripe_secret_key_id
  stripe_webhook_secret_id      = module.platform.stripe_webhook_secret_id
  discord_webhook_url_id        = module.platform.discord_webhook_url_id
  firebase_admin_credentials_id = module.platform.firebase_admin_credentials_id

  depends_on = [module.platform, module.pipeline, module.datastore]
}
