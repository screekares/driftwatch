// Package aws provides an AWS infrastructure provider for driftwatch.
// It fetches live resource state from AWS using the SDK.
package aws

import (
	"context"
	"fmt"

	"github.com/driftwatch/driftwatch/internal/provider"
)

const providerName = "aws"

func init() {
	provider.Register(providerName, func(cfg map[string]string) (provider.Provider, error) {
		return New(cfg)
	})
}

// Client implements provider.Provider for AWS.
type Client struct {
	region  string
	profile string
}

// New creates a new AWS provider client from the given configuration map.
// Recognised keys: "region", "profile".
func New(cfg map[string]string) (*Client, error) {
	region, ok := cfg["region"]
	if !ok || region == "" {
		return nil, fmt.Errorf("aws provider: missing required config key \"region\"")
	}
	return &Client{
		region:  region,
		profile: cfg["profile"],
	}, nil
}

// FetchResource retrieves the live state of a named resource from AWS.
// resourceType is a dot-separated string, e.g. "ec2.instance".
func (c *Client) FetchResource(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	switch resourceType {
	case "ec2.instance":
		return c.fetchEC2Instance(ctx, resourceID)
	case "s3.bucket":
		return c.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("aws provider: unsupported resource type %q", resourceType)
	}
}

func (c *Client) fetchEC2Instance(_ context.Context, id string) (map[string]string, error) {
	// Stub: real implementation would call ec2.DescribeInstances.
	if id == "" {
		return nil, fmt.Errorf("aws provider: ec2.instance id must not be empty")
	}
	return map[string]string{
		"id":     id,
		"region": c.region,
		"state":  "running",
	}, nil
}

func (c *Client) fetchS3Bucket(_ context.Context, name string) (map[string]string, error) {
	// Stub: real implementation would call s3.GetBucketLocation / HeadBucket.
	if name == "" {
		return nil, fmt.Errorf("aws provider: s3.bucket name must not be empty")
	}
	return map[string]string{
		"name":   name,
		"region": c.region,
	}, nil
}
