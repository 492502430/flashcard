# FlashCard Mini Program — Agent Rules

## Pre-Commit Review (MANDATORY)

**Every commit MUST pass all checks below. Do NOT skip any file type.**

### 1. JS Files (.js)
- [ ] Read entire file after edit
- [ ] No functions placed after `});` (closing of Page/App)
- [ ] Every `{` has matching `}`, no orphaned braces
- [ ] Comma between every method in Page({...})
- [ ] No duplicate method names
- [ ] `bindtap`/`bindgetuserinfo` handlers exist in JS

### 2. WXML Files (.wxml)
- [ ] All tags properly closed (`<view>` → `</view>`)
- [ ] All `wx:if`/`wx:for`/`wx:else` properly nested
- [ ] `data-*` attributes match handler expectations

### 3. WXSS Files (.wxss) ← DO NOT SKIP
- [ ] Read last 20 lines of file after each edit
- [ ] No orphaned CSS properties (properties without selector)
- [ ] Every `{` has matching `}`
- [ ] No duplicate `}` at file end
- [ ] After any `sed -i` on WXSS, read full file to verify

### 4. Backend (.go)
- [ ] `go build ./cmd/server/` passes
- [ ] No unused imports or variables

## Never Use sed -i for:
- JS files — use write_file or patch
- Complex WXSS edits — use write_file or patch
- OK for: simple single-line replacements, deleting known line ranges

## Commit Rules
- One logical fix per commit
- Commit message describes WHAT was fixed, not HOW
- Push immediately after commit
