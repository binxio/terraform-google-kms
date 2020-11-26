
# Module `crypto_key`

Provider Requirements:
* **google:** (any version)

## Input Variables
* `app` (required): app name (e.g. mbs, data or gitlab)
* `environment` (required): Environment for which the resources are created (e.g. dev, test, prod or tool)
* `key_ring` (required): The key ring (self link) to add the crypto key to.
* `name` (required): The name of the crypto key. This variable is appended to the crypto key name and used to set the 'name' label.
* `owner_email` (required): Owner of the resources. This variable is used to set the 'owner' label. (Due to GCP label restrictions, only provide the part before '@example.com')
* `purpose` (default `"ENCRYPT_DECRYPT"`): Crypto key purpose, see https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys#CryptoKeyPurpose

## Output Values
* `crypto_key`: Self link to the crypto key
* `crypto_key_name`: Name of the crypto key

## Managed Resources
* `google_kms_crypto_key.crypto_key` from `google`

