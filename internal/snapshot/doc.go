// Package snapshot provides types and functions for capturing
// point-in-time representations of live infrastructure resources.
//
// A snapshot records the attributes of each resource at the moment
// of capture and can be persisted to disk as JSON. Saved snapshots
// can be reloaded and fed into the drift detector to compare against
// the current infrastructure-as-code declarations.
//
// Typical usage:
//
//	snap, err := snapshot.Capture(provider, resourceIDs)
//	if err != nil { ... }
//	if err := snapshot.Save(snap, "./snapshots/latest.json"); err != nil { ... }
//
// Later:
//
//	snap, err := snapshot.Load("./snapshots/latest.json")
package snapshot
