// Package azure implements the driftwatch provider interface for
// Microsoft Azure cloud resources.
//
// # Configuration
//
// The azure provider requires the following keys in the provider
// config block of driftwatch.yaml:
//
//	provider:
//	  name: azure
//	  config:
//	    subscription_id: "<azure-subscription-id>"
//	    resource_group:  "<resource-group-name>"
//
// # Supported Resource Types
//
//   - VirtualMachine
//   - StorageAccount
//
// The provider is automatically registered via its init function.
package azure
