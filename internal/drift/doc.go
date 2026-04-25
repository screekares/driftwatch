// Package drift provides the core drift-detection engine for driftwatch.
//
// It compares live infrastructure resources (fetched via a provider) against
// their declared state and produces a [Report] describing any discrepancies.
//
// # Workflow
//
//  1. Create a [Detector] with [New], passing a configured provider.
//  2. Call [Detector.Check] with the declared resources from the config file.
//  3. Inspect the returned [Report] with [Report.HasDrift] or render it using
//     a [Formatter] obtained from [NewFormatter].
//
// # Formatters
//
// Two output formats are supported:
//   - [FormatText] — a human-readable tabular layout written to any io.Writer.
//   - [FormatJSON] — indented JSON suitable for machine consumption or CI pipelines.
package drift
