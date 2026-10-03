# 0055 - Normalize formatted numeric chart cells

Date: 2026-10-02

## Status

Accepted.

## Context

CSV and Markdown tables commonly format numeric cells with thousands separators
and currency symbols. Strict floating-point parsing rejected these values, so
otherwise valid tables could not be graphed.

## Decision

- Before parsing graph value cells, strip comma separators and Unicode currency
  symbols (Unicode `Sc`, including `$`, `€`, `£`, `¥`, and `₹`).
- Continue rejecting non-numeric or non-finite values after normalization, with
  source row and column context.
- Continue treating empty cells as missing chart values.

## Consequences

- Commonly formatted values such as `$1,234.50` are graphed as `1234.50`.
- Other punctuation and alphabetic text are not silently discarded.

## References

- `utils/graphs/graphs.go` - numeric cell parsing
- `utils/graphs/graphs_test.go` - normalization regression cases
