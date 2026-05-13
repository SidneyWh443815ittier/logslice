// Package transform provides post-filter line transformation for logslice.
//
// Transformations are applied after a line has matched all filter criteria
// and before it is handed to the output formatter. This allows users to:
//
//   - Rename noisy or inconsistent field names (rename:msg:message)
//   - Drop sensitive or irrelevant fields before display (drop:password)
//   - Inject synthetic constant fields for downstream tooling (set:env:prod)
//
// JSON log lines are parsed and re-serialised so that field-level operations
// work correctly. Plain-text lines are supported for "set" operations only,
// which append key=value pairs to the end of the line.
//
// Operations are applied in the order they are declared, enabling composed
// pipelines such as renaming a field and then dropping the original.
package transform
