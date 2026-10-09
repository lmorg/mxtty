# 0063 - Notes uses the shared tooltip instead of native webview tooltips

## Status

Accepted on 2026-10-09; migration of existing Notes controls is incremental.

## Context

Native HTML `title` tooltips are rendered by the host webview, so their timing,
positioning, styling, and content presentation are outside Notes' control.
Notes has a shared tooltip overlay that follows the pointer, uses application
styling, and can render Markdown.

## Decision

- Use the shared Notes tooltip for visual hover help. Do not rely on native
  webview tooltips from `title` attributes for that help.
- For dynamically rendered Notes content, expose tooltip text through
  `data-tooltip` and use delegated event handling rather than attaching
  listeners to each item.
- Route feature-specific tooltip events through the shared Notes tooltip
  integration. Use the existing show and close functions when a component
  already owns that integration.
- Keep accessible control names and semantic accessibility attributes. A
  visual tooltip supplements them; it does not replace them or make hover the
  only way to understand an essential action.
- Migrate existing Notes `title` attributes when those controls are next
  changed. Retain `title` only where it is not being used as hover help.

## Consequences

- Tooltip appearance and positioning are consistent across Notes and host
  webviews.
- Dynamically added controls can use the same delegated tooltip behavior.
- The shared tooltip is pointer-following and non-interactive; it is not a
  substitute for accessible names or focusable help content.

## Reference

- `frontend/src/shared_tooltip.js` - shared tooltip rendering and pointer tracking
- `frontend/src/notes.js` - Notes tooltip integration
- `frontend/src/notes-git-panel.js` - delegated `data-tooltip` handling