# 0024 — Dismiss a menu before dispatching its selection

## Status

Accepted.

## Context

`Go to workspace symbol` opened no menu at all. Nothing was logged, so there was
no exception — the code ran to completion and the menu still never appeared.

`selectMenuItem` dispatched the selection and *then* dismissed the menu:

```js
function selectMenuItem(item) {
    _menuOperationInProgress = true;
    if (activeListMenuId !== null) {
        menuSelect(activeListMenuId, item.index);   // runs the handler
    }
    hideListMenu(false);                            // then tears down
}
```

`hideListMenu` is not scoped to a particular menu. It resets the shared listbox:
`activeListMenuId = null`, clears `listItems` / `filteredItems`, unsets
`showSearch` and `dynamicQuery`, and animates the single `listRoot` to
`display: none`.

So any handler that opens another menu **synchronously** has that new menu
destroyed the moment it returns.

This was latent for as long as the LSP submenus existed, but hidden: every one of
them awaited a language-server round trip before calling `showLocalMenu`, so the
new menu opened a microtask *after* the teardown and survived. Moving workspace
symbol search server-side (ADR 0022) removed its only pre-menu `await` in the
common case — the document is already attached, so nothing suspends — and the
latent bug became reproducible.

`goToCurrentLspSymbol` still awaits `NotesLspDocumentSymbols`, which is why one
symbol command worked and its sibling did not.

## Decision

Dismiss first, dispatch second. `activeListMenuId` is captured before the hide,
because `hideListMenu` nulls it:

```js
const menuId = activeListMenuId;
hideListMenu(false);
if (menuId !== null) {
    menuSelect(menuId, item.index);
}
```

Applied to both dispatch paths — mouse selection and the `Enter` key handler,
which had the same ordering.

`menuSelect` also restored focus to whatever was focused before the menu opened,
which would steal focus from a newly opened menu's search field. `showLocalMenu`
re-assigns `_localMenuReturnFocus`, so a non-null value after the handler returns
is a reliable signal that a new menu took ownership:

```js
if (returnTo && !_localMenuReturnFocus) returnTo.focus();
```

## Consequences

- Handlers may open a submenu synchronously or asynchronously; both work.
- Selection callbacks now run against a already-dismissed menu. They must not
  assume menu state is still live — none did, since `hideListMenu` ran
  immediately afterwards anyway.

## Lessons

- **An `await` can be load-bearing.** The original code was only correct because
  every handler happened to suspend first. Removing an unrelated network call
  broke a menu — timing was doing the work that ordering should have done.
- **"Nothing happened, nothing logged" points at teardown, not failure.** No
  exception meant the open path ran fine; something undid it afterwards.
- **Verify a regression test actually regresses.** Restoring the old ordering
  turned the new test red with `expected 'none' to be 'block'`, confirming it
  pins the real defect rather than passing incidentally.

## Notes Printing (2026-10-23)

View/Edit printing prepares a separate `#notes-print-root` from the complete
editor buffer rather than printing a scrolled preview or Monaco's virtualized
DOM. View mode renders the Markdown body through the shared processing pipeline;
Edit mode prints source text, including unsaved Monaco changes. Native details
sections are expanded in the snapshot only. Print layout excludes every other
body child, including the file sidebar, Tools panel, menus, and tooltips, and
removes viewport clipping without changing the live document layout or theme.

The Print action closes shared tooltips and waits for menu dismissal animations,
image decoding, font readiness, and layout before invoking Wails `WindowPrint`.
Repeated requests during preparation share one operation. The prepared document
is retained until `afterprint`, or replaced by the next request if the native
print implementation does not emit that event; the Wails promise alone is not
a reliable indication that the native dialog has closed. `afterprint` events
received while the print binding is still pending are ignored: WebKit can finish
a preview render before Save as PDF needs to render the document again. Removing
the snapshot at that point leaves the later PDF render without its content.
Once the native operation resolves, a subsequent event may clean up normally;
otherwise the next request replaces the retained snapshot. Failed print requests
discard the snapshot. Other document modes clear retained snapshots before
using their existing window-print path. Printing uses the user-supplied iTerm
Display-P3 palette, scoped to print media so the live theme is unchanged.
Foreground and bold text are black; ANSI and link colours retain their P3
components. Page layout requests 15 mm margins and a 5 mm inner gutter, with
cloned box decoration for fragmented pages. Mermaid SVG colours remain those
chosen by the diagram renderer rather than being recoloured through CSS.

Print preparation captures each element's original computed font size before
switching to print media. A print-only rule subtracts `2pt` from each captured
size (with a `1pt` minimum), including plain source output from Monaco. Capturing
all sizes before applying the reduction preserves heading sizes and avoids
compounding the reduction through nested elements. Screen typography is unchanged.

### Native Print Backend

The AI panel's Print action snapshots top-level assistant `text` blocks as
Markdown and uses the same print pipeline. Thinking, tool calls/results/errors/
summaries, user-prompt headings, and nested subagent text are excluded. Offscreen
historical output is fetched by block ID without loading excluded blocks or
changing the displayed history. Legacy transcripts print only classified
`Final Answer` sections; unclassified transcripts are not printed. With no
assistant output, the action shows a notification rather than an empty dialog.

Markdown modals expose Print through both body and table context menus. The
modal retains its original Markdown and passes it to the Notes print hook,
which uses the same isolated print pipeline as Markdown View without changing
the displayed modal. Jupyter Run context menus expose Print but render the
document through the Markdown View pipeline, excluding runtime output and
controls without changing the active tab. Edit output begins with a separate
filename line before the complete source text: `$PROJECT` files use their
project-relative path, while other categories use the basename. Filenames are
inserted as text, not HTML.
Current modal callers are Git worktree information, Git branch/tag information,
and the terminal's View command meta action.

`WApp.WindowPrint` delegates to `utils/printing.Print(ctx)`. On macOS, an
app-owned Objective-C bridge finds the main window's `WKWebView` through Cocoa's
view hierarchy and runs the system print operation on the main thread. It copies
`NSPrintInfo` instead of mutating the shared settings, explicitly selects portrait,
fits horizontal pagination, and leaves vertical pagination automatic. Native
margins remain zero so the existing print CSS controls page spacing. The native
call waits for the modal print operation to finish; cancelling the dialog is not
an application error. Setup errors reject the existing Wails binding promise.

The print view must receive `webView.bounds`: WebKit initially returns a view
with a zero-sized frame. Start the operation with
`runOperationModalForWindow:delegate:didRunSelector:contextInfo:`, not synchronous
`runOperation` inside a main-queue dispatch block. Completion is signalled back
to the waiting Go caller while the main queue remains available for WebKit.
The synchronous path stalled the native PDF regression fixture even after its
frame was corrected.

This bypasses Wails' hardcoded macOS landscape default without modifying Wails
or relying on its private window context. macOS requires cgo and macOS 11 or
later. Non-macOS builds retain `runtime.WindowPrint(ctx)` as their fallback.
Package tests cover cancelled contexts and a missing native application, while
Linux and Windows builds are cross-compiled. An opt-in native test saves a PDF
from a hidden WebKit window using the production operation helper and inspects
its text with PDFKit. It verifies all 120 paragraphs, start/end markers, multiple
portrait pages, and exclusion of application controls. Run it in a macOS
WindowServer session with `TTYPHOON_TEST_NATIVE_PRINT=1 go test ./utils/printing
-run TestNativePrintPDFContainsPreparedDocument -v`. This does not open a print
dialog; physical printer and interactive dialog behavior still require manual
verification.

## Reference

- `frontend/src/popup_menu.js` — `selectMenuItem`, `menuSelect`, `hideListMenu`,
  window `keydown` Enter handler
- `frontend/src/popup_menu.test.js` — `keeps a menu opened synchronously from
  another menu selection`
- ADR 0022 — removing the pre-menu `await` that had been masking this
