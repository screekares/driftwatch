package azure

import (
	"testing"
)

func TestNew_MissingSubscriptionID(t *testing.T) {
	_, err := New(map[string]string{"resource_group": "rg-test"})
	if err == nil {
		t.Fatal("expected error for missing subscription_id, got nil")
	}
}

func TestNew_MissingResourceGroup(t *testing.T) {
	_, err := New(map[string]string{"subscription_id": "sub-123"})
	if err == nil {
		t.Fatal("expected error for missing resource_group, got nil")
	}
}

func TestNew_ValidConfig(t *testing.T) {
	p, err := New(map[string]string{
		"subscription_id": "sub-123",
		"resource_group":  "rg-test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != providerName {
		t.Errorf("expected name %q, got %q", providerName, p.Name())
	}
}

func TestFetchResource_VirtualMachine(t *testing.T) {
	p, _ := New(map[string]string{"subscription_id": "sub-123", "resource_group": "rg-test"})
	res, err := p.FetchResource("VirtualMachine", "vm-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["type"] != "VirtualMachine" {
		t.Errorf("expected type VirtualMachine, got %q", res["type"])
	}
	if res["id"] != "vm-001" {
		t.Errorf("expected id vm-001, got %q", res["id"])
	}
}

func TestFetchResource_StorageAccount(t *testing.T) {
	p, _ := New(map[string]string{"subscription_id": "sub-123", "resource_group": "rg-test"})
	res, err := p.FetchResource("StorageAccount", "sa-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["sku"] != "Standard_LRS" {
		t.Errorf("expected sku Standard_LRS, got %q", res["sku"])
	}
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	p, _ := New(map[string]string{"subscription_id": "sub-123", "resource_group": "rg-test"})
	_, err := p.FetchResource("CosmosDB", "db-001")
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
}

func TestFetchResource_SubscriptionPropagated(t *testing.T) {
	p, _ := New(map[string]string{"subscription_id": "sub-xyz", "resource_group": "rg-prod"})
	res, err := p.FetchResource("VirtualMachine", "vm-002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["subscription"] != "sub-xyz" {
		t.Errorf("expected subscription sub-xyz, got %q", res["subscription"])
	}
	if res["resource_group"] != "rg-prod" {
		t.Errorf("expected resource_group rg-prod, got %q", res["resource_group"])
	}
}
