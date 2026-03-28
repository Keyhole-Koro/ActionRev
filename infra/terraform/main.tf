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

  # Terraform state を GCS に保存する
  # 初回 apply 前に以下のバケットを手動で作成すること:
  #   gsutil mb -p <project_id> -l asia-northeast1 gs://<project_id>-tfstate
  backend "gcs" {
    bucket = "REPLACE_WITH_YOUR_PROJECT_ID-tfstate"
    prefix = "actionrev"
  }
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
  env    = var.env
  labels = {
    env     = var.env
    app     = "actionrev"
    managed = "terraform"
  }
}
