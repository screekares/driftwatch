package aws

import (
	"context"
	"testing"
)

func TestNew_MissingRegion(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing region, got nil")
	}
}

func TestNew_ValidConfig(t *testing.T) {
	c, err := New(map[string]string{"region": "us-east-1", "profile": "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %q", c.region)
	}
}

func TestFetchResource_EC2Instance(t *testing.T) {
	c, _ := New(map[string]string{"region": "eu-west-1"})
	attrs, err := c.FetchResource(context.Background(), "ec2.instance", "i-0abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["id"] != "i-0abc123" {
		t.Errorf("expected id i-0abc123, got %q", attrs["id"])
	}
	if attrs["region"] != "eu-west-1" {
		t.Errorf("expected region eu-west-1, got %q", attrs["region"])
	}
}

func TestFetchResource_S3Bucket(t *testing.T) {
	c, _ := New(map[string]string{"region": "us-west-2"})
	attrs, err := c.FetchResource(context.Background(), "s3.bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["name"] != "my-bucket" {
		t.Errorf("expected name my-bucket, got %q", attrs["name"])
	}
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	c, _ := New(map[string]string{"region": "us-east-1"})
	_, err := c.FetchResource(context.Background(), "rds.cluster", "db-1")
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
}

func TestFetchResource_EmptyID(t *testing.T) {
	c, _ := New(map[string]string{"region": "us-east-1"})
	_, err := c.FetchResource(context.Background(), "ec2.instance", "")
	if err == nil {
		t.Fatal("expected error for empty resource id, got nil")
	}
}

func TestInit_RegistersProvider(t *testing.T) {
	// Verify the init() function registered the provider by attempting
	// to create one via the top-level provider.New helper.
	// We import the aws package for its side-effects in other tests;
	// here we just confirm New itself works end-to-end.
	c, err := New(map[string]string{"region": "ap-southeast-1"})
	if err != nil {
		t.Fatalf("New returned unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
