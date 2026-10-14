# 69. Notes spellcheck suggestions open on right-click

Date: 2026-10-14

## Status

Accepted.

## Context

Notes rendered spelling suggestions when the user left-clicked a misspelled word.
This interrupted ordinary caret placement and text selection. Notes already
uses right-click for editor context menus, so suggestions should use that
gesture only when the pointer is over a spelling underline.

## Decision

- Open the spelling suggestion menu on `contextmenu` over a detected misspelled
  word, positioned below the word as before.
- Prevent the browser menu and stop propagation only when a misspelling is hit,
  so the suggestion menu replaces the general Notes editor menu at that point.
- Do not show suggestions on left-click. Right-clicks elsewhere continue to
  the existing editor context-menu handler.
- Apply a selected suggestion through the existing text replacement and input
  event path so spellcheck and editor state refresh normally.

## Consequences

Ordinary left-clicks remain available for editing; misspellings expose their
suggestions through the pointer's context-menu gesture. Tests cover both the
spelling hit and pass-through behavior elsewhere.

## Reference

- `frontend/src/spellcheck.js`
- `frontend/src/spellcheck.test.js`
- `frontend/src/notes.js` - general editor context menu