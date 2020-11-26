#------------------------------------------------------------------------------------------------------------------------
# 
# Generic variables
#
#------------------------------------------------------------------------------------------------------------------------
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
# secret variables
#
#------------------------------------------------------------------------------------------------------------------------
variable "purpose" {
  description = "The purpose of the secret. This variable is appended to the secret_id and used to set the 'purpose' label."
  type        = string
}

variable "location" {
  description = "The GCP location (region) of the key_ring."
  type        = string
  default     = "europe-west4"
}

variable "roles" {
  description = "Map of role name's as `key` and members list as `value` to bind permissions"
  type        = map(map(string))
  default     = {}
}
