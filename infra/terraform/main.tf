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

  # 初期化時に -backend-config=backends/<env>.hcl を指定する
  # 例: terraform init -backend-config=backends/stage.hcl
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
# Artifact Registry
# ---------------------------------------------------------------------------

resource "google_artifact_registry_repository" "actionrev" {
  repository_id = "actionrev"
  location      = var.region
  format        = "DOCKER"
  description   = "ActionRev Docker images"
  labels        = local.labels

  cleanup_policies {
    id     = "keep-latest-10"
    action = "KEEP"
    most_recent_versions { keep_count = 10 }
  }

  cleanup_policies {
    id     = "delete-untagged-after-30d"
    action = "DELETE"
    condition {
      tag_state  = "UNTAGGED"
      older_than = "2592000s"
    }
  }

  depends_on = [google_project_service.apis]
}

# ---------------------------------------------------------------------------
# Modules
# ---------------------------------------------------------------------------

module "storage" {
  source              = "./modules/storage"
  project_id          = var.project_id
  region              = var.region
  env                 = var.env
  labels              = local.labels
  uploads_bucket_name = var.uploads_bucket_name
  depends_on          = [google_project_service.apis]
}

module "bigquery" {
  source     = "./modules/bigquery"
  project_id = var.project_id
  region     = var.region
  dataset_id = var.bigquery_dataset_id
  labels     = local.labels
  depends_on = [google_project_service.apis]
}

module "iam" {
  source              = "./modules/iam"
  project_id          = var.project_id
  sandbox_bucket_name = module.storage.sandbox_bucket_name
  depends_on          = [google_project_service.apis]
}

module "secrets" {
  source     = "./modules/secrets"
  project_id = var.project_id
  labels     = local.labels
  depends_on = [google_project_service.apis]
}

module "tasks" {
  source     = "./modules/tasks"
  project_id = var.project_id
  region     = var.region
  env        = var.env
  depends_on = [google_project_service.apis]
}

module "cloudrun" {
  source               = "./modules/cloudrun"
  project_id           = var.project_id
  region               = var.region
  env                  = var.env
  labels               = local.labels
  backend_image        = var.backend_image
  sandbox_image        = var.sandbox_image
  backend_sa_email     = module.iam.backend_sa_email
  sandbox_sa_email     = module.iam.sandbox_sa_email
  uploads_bucket_name  = module.storage.uploads_bucket_name
  sandbox_bucket_name  = module.storage.sandbox_bucket_name
  bigquery_dataset_id  = module.bigquery.dataset_id
  cloud_tasks_queue    = module.tasks.queue_name
  min_instances        = var.backend_min_instances
  max_instances        = var.backend_max_instances
  gemini_model         = var.gemini_model

  stripe_secret_key_id          = module.secrets.stripe_secret_key_id
  stripe_webhook_secret_id      = module.secrets.stripe_webhook_secret_id
  discord_webhook_url_id        = module.secrets.discord_webhook_url_id
  firebase_admin_credentials_id = module.secrets.firebase_admin_credentials_id

  depends_on = [module.iam, module.secrets]
}
