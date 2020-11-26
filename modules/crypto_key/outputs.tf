output "crypto_key" {
  description = "Self link to the crypto key"
  value       = google_kms_crypto_key.crypto_key.self_link
}
output "crypto_key_name" {
  description = "Name of the crypto key"
  value       = google_kms_crypto_key.crypto_key.name
}
