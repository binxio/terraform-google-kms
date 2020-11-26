output "key_ring_name" {
  description = "Name of the key ring"
  value       = google_kms_key_ring.key_ring.name
}

output "key_ring" {
  description = "Self link to the key ring"
  value       = google_kms_key_ring.key_ring.self_link
}
