locals {
  owner       = "myself"
  project     = "testapp"
  environment = var.environment
}

module "key_ring" {
  source = "../../modules/key_ring"

  environment = local.environment
  project     = local.project

  purpose = "terratest"
}

module "crypto_key" {
  source = "../../modules/crypto_key"

  owner       = local.owner
  environment = local.environment
  project     = local.project

  name     = "terratest"
  key_ring = module.key_ring.key_ring
}

output "key_ring_name" {
  value = module.key_ring.key_ring_name
}

output "key_ring" {
  value = module.key_ring.key_ring
}

output "crypto_key_name" {
  value = module.crypto_key.crypto_key_name
}
