// Package gcp implements the driftwatch provider interface for
// Google Cloud Platform resources.
//
// Supported resource types:
//
//   - compute_instance: Represents a GCP Compute Engine VM instance.
//   - storage_bucket:   Represents a GCP Cloud Storage bucket.
//
// Configuration keys:
//
//   - project (required): The GCP project ID.
//   - region  (optional): The default region for resource lookups.
//
// The provider is automatically registered under the name "gcp"
// via its init() function when the package is imported.
package gcp
