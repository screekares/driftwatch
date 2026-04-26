// Package gcp provides a Google Cloud Platform provider for driftwatch.
package gcp

import (
	"fmt"

	"github.com/driftwatch/driftwatch/internal/provider"
)

func init() {
	provider.Register("gcp", func(cfg map[string]string) (provider.Provider, error) {
		return New(cfg)
	})
}

// gcpProvider implements provider.Provider for GCP resources.
type gcpProvider struct {
	project string
	region  string
}

// New creates a new GCP provider from the given configuration map.
func New(cfg map[string]string) (provider.Provider, error) {
	project, ok := cfg["project"]
	if !ok || project == "" {
		return nil, fmt.Errorf("gcp provider: missing required config key 'project'")
	}
	return &gcpProvider{
		project: project,
		region:  cfg["region"],
	}, nil
}

// FetchResource retrieves a GCP resource by type and ID.
func (g *gcpProvider) FetchResource(resourceType, id string) (map[string]string, error) {
	switch resourceType {
	case "compute_instance":
		return g.fetchComputeInstance(id)
	case "storage_bucket":
		return g.fetchStorageBucket(id)
	default:
		return nil, fmt.Errorf("gcp provider: unsupported resource type %q", resourceType)
	}
}

func (g *gcpProvider) fetchComputeInstance(id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("gcp provider: compute_instance id must not be empty")
	}
	return map[string]string{
		"id":      id,
		"project": g.project,
		"region":  g.region,
		"status":  "RUNNING",
		"type":    "compute_instance",
	}, nil
}

func (g *gcpProvider) fetchStorageBucket(id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("gcp provider: storage_bucket id must not be empty")
	}
	return map[string]string{
		"id":              id,
		"project":         g.project,
		"location":        g.region,
		"storage_class":   "STANDARD",
		"type":            "storage_bucket",
	}, nil
}
