// Package gcp provides a Google Cloud Platform provider for driftwatch.
package gcp

import (
	"context"
	"fmt"

	"github.com/driftwatch/driftwatch/internal/provider"
)

const providerName = "gcp"

func init() {
	provider.Register(providerName, New)
}

// gcpProvider implements provider.Provider for GCP resources.
type gcpProvider struct {
	project string
	region  string
}

// New creates a new GCP provider from the given config map.
// Required keys: "project". Optional: "region".
func New(cfg map[string]string) (provider.Provider, error) {
	project, ok := cfg["project"]
	if !ok || project == "" {
		return nil, fmt.Errorf("gcp provider: missing required config key \"project\"")
	}
	region := cfg["region"]
	return &gcpProvider{project: project, region: region}, nil
}

// FetchResource retrieves a GCP resource by type and ID.
func (g *gcpProvider) FetchResource(ctx context.Context, resourceType, id string) (map[string]string, error) {
	switch resourceType {
	case "compute_instance":
		return g.fetchComputeInstance(ctx, id)
	case "storage_bucket":
		return g.fetchStorageBucket(ctx, id)
	default:
		return nil, fmt.Errorf("gcp provider: unsupported resource type %q", resourceType)
	}
}

func (g *gcpProvider) fetchComputeInstance(_ context.Context, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("gcp provider: compute_instance id must not be empty")
	}
	// Stub: replace with real GCP Compute API call.
	return map[string]string{
		"id":      id,
		"type":    "compute_instance",
		"project": g.project,
		"region":  g.region,
		"status":  "RUNNING",
	}, nil
}

func (g *gcpProvider) fetchStorageBucket(_ context.Context, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("gcp provider: storage_bucket id must not be empty")
	}
	// Stub: replace with real GCP Storage API call.
	return map[string]string{
		"id":      id,
		"type":    "storage_bucket",
		"project": g.project,
		"location": g.region,
	}, nil
}
