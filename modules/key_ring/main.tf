locals {
  project     = var.project
  environment = var.environment
  location    = var.location

  # We like to have the project and env in the key ring name, so it's obvious what the key ring belongs to
  # without checking the GCP project it's part of
  # This also allows for central key ring management since the naming is more unique by default.
  name = replace(lower(format("%s-%s-%s", local.project, local.environment, var.purpose)), " ", "-")

  # Map role members to GCP expected format
  roles = { for role, members in try(var.roles, {}) : role => [for member, type in members : format("%s:%s", type, member)] }
}

resource "google_kms_key_ring" "key_ring" {
  # Hacked out while https://github.com/hashicorp/terraform/issues/3116 is open, this needs to be variable for testing etc....
  #lifecycle {
  #  prevent_destroy = true
  #}
  name     = local.name
  location = local.location
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

resource "google_kms_key_ring_iam_policy" "map" {
  for_each = data.google_iam_policy.map

  key_ring_id = google_kms_key_ring.key_ring.self_link
  policy_data = each.value.policy_data
}
