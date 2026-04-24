// Package provider defines the Provider interface and a global registry for
// pluggable infrastructure backends used by driftwatch.
//
// # Adding a new provider
//
// 1. Create a sub-package (e.g. internal/provider/aws).
// 2. Implement the Provider interface.
// 3. Call provider.Register("aws", factory) inside an init() function.
// 4. Blank-import the sub-package from cmd/root.go or your entry-point so
//    the init() runs and the provider becomes available.
//
// Example:
//
//	import _ "github.com/yourorg/driftwatch/internal/provider/aws"
package provider
