locals {
  project     = var.project
  environment = var.environment
  owner       = var.owner

  key_ring = var.key_ring
  purpose  = var.purpose

  labels = {
    "project" = substr(replace(lower(local.project), "/[^\\p{Ll}\\p{Lo}\\p{N}_-]+/", "_"), 0, 63)
    "env"     = substr(replace(lower(local.environment), "/[^\\p{Ll}\\p{Lo}\\p{N}_-]+/", "_"), 0, 63)
    "owner"   = substr(replace(local.owner, "/[^\\p{Ll}\\p{Lo}\\p{N}_-]+/", "_"), 0, 63)
    "creator" = "terraform"
  }

  rotation_period = "2629800s" # 1 month

  # We like to have the project and env in the key ring name, so it's obvious what the key ring belongs to
  # without checking the GCP project it's part of
  # This also allows for central key ring management since the naming is more unique by default.
  name = replace(lower(format("%s-%s-%s", local.project, local.environment, var.name)), " ", "-")

  # Map role members to GCP expected format
  roles = { for role, members in try(var.roles, {}) : role => [for member, type in members : format("%s:%s", type, member)] }
}


resource "google_kms_crypto_key" "crypto_key" {
  name            = local.name
  key_ring        = local.key_ring
  rotation_period = local.rotation_period
  purpose         = local.purpose
  labels          = local.labels
}

data "google_iam_policy" "map" {
  for_each = (length(keys(local.roles)) > 0 ? { policy = {} } : {})

  dynamic "binding" {
    for_each = local.roles

    content {
      role    = binding.key
      members = binding.value
    }
  }
}

resource "google_kms_crypto_key_iam_policy" "map" {
  for_each = data.google_iam_policy.map

  crypto_key_id = google_kms_crypto_key.crypto_key.self_link
  policy_data   = each.value.policy_data
}
