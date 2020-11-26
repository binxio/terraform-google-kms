
# Module `key_ring`

Provider Requirements:
* **google:** (any version)

## Input Variables
* `app` (required): app name (e.g. mbs, data or gitlab)
* `environment` (required): Environment for which the resources are created (e.g. dev, test, prod or tool)
* `gcp_project` (default `"global-iam"`): Project in which the key_ring is created.
* `location` (default `"europe-west4"`): The GCP region of the key_ring.
* `name` (required): The name of the key ring. This variable is appended to the key ring name and used to set the 'name' label.
* `roles` (required): Map of role name's as key and members list as value to bind permissions

## Output Values
* `key_ring`: Self link to the key ring
* `key_ring_name`: Name of the key ring

## Managed Resources
* `google_kms_key_ring.key_ring` from `google`
* `google_kms_key_ring_iam_policy.map` from `google`

## Data Resources
* `data.google_iam_policy.map` from `google`

