#------------------------------------------------------------------------------------------------------------------------
#
# Generic variables
#
#------------------------------------------------------------------------------------------------------------------------
variable "owner" {
  description = "Owner of the resource. This variable is used to set the 'owner' label."
  type        = string
}

variable "project" {
  description = "Company project name."
  type        = string
}

variable "environment" {
  description = "Company environment for which the resources are created (e.g. dev, tst, acc, prd, all)."
  type        = string
}

#------------------------------------------------------------------------------------------------------------------------
#
# crypto_key variables
#
#------------------------------------------------------------------------------------------------------------------------
variable "name" {
  description = "The name of the crypto key. This variable is appended to the crypto key name and used to set the 'name' label."
  type        = string
}

variable "key_ring" {
  description = "The key ring (self link) to add the crypto key to."
  type        = string
}

variable "purpose" {
  description = "Crypto key purpose, see https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys#CryptoKeyPurpose"
  type        = string
  default     = "ENCRYPT_DECRYPT"
}

variable "roles" {
  description = "Map of role name's as `key` and members list as `value` to bind permissions"
  type        = map(map(string))
  default     = {}
}
