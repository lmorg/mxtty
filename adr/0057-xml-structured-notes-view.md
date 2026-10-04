# 0057 - XML uses source-preserving structured Notes View

Date: 2026-10-04

## Status

Accepted.

## Context

Notes already provides a tree-based View for JSON and YAML so users can inspect
and edit individual values without working in raw text. XML documents also
benefit from that interaction, but parsing and serializing the whole document
would discard comments and formatting that are unrelated to the edit.

## Decision

- Route `.xml` Notes files and text files beginning with an XML declaration
  through the existing structured View and tree editor.
- Map XML attributes to `@name`, element text to scalar values, and repeated
  sibling elements to arrays. Preserve mixed text content as `#text`.
- Treat newly added keys prefixed with `@` as attributes, matching the tree's
  representation of existing XML attributes.
- Parse XML into the tree model, but apply edits as targeted changes to the
  original source spans. Do not serialize the whole XML document after an edit.
- Support editing values and attributes, renaming and deleting nodes, and adding
  attributes and child elements while retaining surrounding source text.

## Consequences

- XML joins JSON and YAML in the structured Notes workflow and uses the same
  inline tree controls.
- Comments, declarations, whitespace, and untouched siblings remain as written
  when an individual node is edited.
- The XML editor owns its source scanning and edit behavior; it does not require
  a generic XML serializer that normalizes the document.

## References

- `frontend/src/xml-structured-editor.js` - XML tree mapping and source edits
- `frontend/src/structured-editors.js` - structured editor registration
- `frontend/src/notes.js` - Notes file routing and structured View lifecycle
- `frontend/src/xml-structured-editor.test.js` - XML editor behavior tests
- `frontend/src/notes.test.js` - XML Notes View integration test