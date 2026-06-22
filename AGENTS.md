# FlashCard Mini Program — Agent Rules

## Pre-Commit Review (MANDATORY)

Before every `git commit`, you MUST:

1. **Syntax check every .js file changed** — read the full file, scan for:
   - Missing commas between Page() methods
   - Duplicate closing `});` or extra braces
   - Functions placed outside `Page({...})` after `});`
   - Unmatched `{`/`}` counts

2. **Check every .wxml change** — verify:
   - All `wx:if`/`wx:for` attributes are properly closed
   - `bindtap` handlers exist in the corresponding .js file
   - `data-*` attributes match what the handler expects

3. **Check every .wxss change** — verify:
   - No CSS appended after `@keyframes` without proper separator
   - Selectors match the WXML elements

4. **Verify the full file is syntactically valid** by reading from top to bottom.

## Commit Guidelines

- One fix per commit
- Never use `sed -i` for JS edits — always use read_file + write_file
- After any sed operation, re-read the file to verify correctness

## Must-Have for Every UI Change

- No emoji anywhere (WXML, JS strings, CSS content)
- All user-facing text in Chinese
- Buttons use `::after { border: none }` to remove WeChat's default button border
