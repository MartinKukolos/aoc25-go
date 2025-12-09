# Copilot Instructions

## Project Overview
- Advent of Code 2025 solutions written in Go 1.22.
- Each day lives in `DayN/` with `main.go`, `main_test.go`, `input.txt`, `sample.txt`, `Readme.md`, and now a `blog.md` write-up.
- Root `go.mod` defines module `aoc25` with packages per day (`aoc25/DayN`).

## Coding Guidelines
- Favor clear, iterative algorithms; avoid premature concurrency.
- Keep files ASCII-only unless puzzle input forces otherwise.
- Each solver should expose `Solve(io.Reader)` and use a small `main` for CLI handling.
- Prefer pure functions and explicit error handling over panics.
- Reuse helper utilities within the same file; avoid cross-day imports.

## Testing & Tooling
- For new days, add `main_test.go` verifying the sample from the puzzle statement.
- Run `go test ./...` before committing.
- Run `gofmt` (or `go fmt ./...`) on all touched Go files.

## Documentation
- Every day folder should include `blog.md` summarizing approach and complexity (match existing style in `Day8/blog.md`).
- When adding new inputs or samples, keep original puzzle wording untouched in `Readme.md`.

## Commit Hygiene
- One logical change per commit (e.g., “Add Day10 solver and tests”).
- If regenerating inputs or large files, mention it explicitly in the commit message.

## Copilot Prompts
- When asking Copilot for help, mention the specific day folder, file, and the puzzle behavior you’re targeting.
- Provide sample input/output whenever possible so Copilot can infer edge cases.
