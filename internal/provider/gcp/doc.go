// Package gcp implements the driftwatch provider interface for Google Cloud
// Platform resources.
//
// # Supported resource types
//
//   - compute_instance – GCP Compute Engine VM instances
//   - storage_bucket   – GCP Cloud Storage buckets
//
// # Configuration
//
// The GCP provider is configured via the providers section of the driftwatch
// configuration file:
//
//	providers:
//	  gcp:
//	    project: my-gcp-project
//	    region:  us-central1   # optional
//
// The "project" key is required. "region" is optional and used as a hint when
// fetching region-scoped resources.
package gcp
