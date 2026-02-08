<!-- Short, actionable guidance for AI coding agents working on this repo -->
# EpicFreeGamesList — Copilot Instructions

Summary
- Purpose: keeps the Epic Store "free games" list updated and produces `epic_free_games.json` used by the static UI (`index.html` + `renderTable.js`). See [README.md](README.md#L1) for usage examples.

Architecture & dataflow (quick)
- Entry point: [main.go](main.go#L1). The program is a simple CLI with operations: `search`, `rate`, `free`, `free_mobile`, `fix_ratings`, `mobile_discover_page`.
- Fetch layer: files named `graphql_get_*.go` (e.g. [graphql_get_free_games.go](graphql_get_free_games.go#L1)) query Epic endpoints to discover product slugs and ratings.
 - CLI layer: `cli_handler_*.go` — files providing CLI entrypoints. Example: [cli_handler_free.go](cli_handler_free.go#L1).
- Parsing/pagination: `paginated_discover_modules_*.go` contains logic for paginated endpoints and extraction.
- Models: `json_model.go` and the local types in `graphql_get_free_games.go` shape the JSON written to `epic_free_games.json`.
- Output: `epic_free_games.json` is the canonical JSON produced/consumed; the static UI reads it. Dockerfile builds a container image of the updater.
 - Output: `epic_free_games.json` is the canonical JSON produced/consumed; the static UI reads it.

Project-specific conventions & gotchas
 - Note: older commits used the misspelling `cli_hander_*`. Current files use `cli_handler_*`. If you rename files, update `main.go` and any docs.
- Network calls: code uses `net/http` and also includes `github.com/bogdanfinn/tls-client` (see `go.mod`) to handle Epic's protections. Keep HTTP client behavior in mind when editing request logic.
- Country/locale: `GetFreeGames()` hardcodes query params (example: `country=CA`, `locale=en-US`) — change there when altering region behavior ([graphql_get_free_games.go](graphql_get_free_games.go#L1)).
- Rating resolution: `CliHandlerFree()` will call `RateGame()` where available; missing sandboxId is expected and logged rather than fatal.
- Files referenced from README: build and run examples assume Windows CLI (`.\\main` and PowerShell). Cross-platform use: prefer `go run main.go`/`go build` on other OSes.

How to build / run / debug (concrete)
- Recommended build (cross-platform): run from the repository root so Go discovers all packages automatically.
  - Build a binary: `go build -o epic-updater .`
  - Run without producing a binary: `go run . free --inputFile epic_free_games.json --outputFile out.json`
  - Quick local test (Windows PowerShell):
    - Build then run: `go build -o epic-updater.exe .` then `.\epic-updater.exe free`.
    - Or: `go run . free` (PowerShell automatically resolves `go` command).
- Notes: you do not need to list every `.go` file in the `go build`/`go run` command. Using `.` at the repo root (or module root) is preferred.
- Debug tips: run `go run . free` to print fetch logs (functions already print responses). Add scoped prints inside `GetFreeGames()` or `RateGame()` for HTTP response bodies when investigating parsing issues.

Integration tests (quick smoke runs)
- To quickly validate `free` and `free_mobile` after code changes, run them against a minimal empty JSON file (`{}`) so the CLI will append or produce output without needing a full dataset.
- Create an empty JSON file:
  - Unix/macOS / WSL:
    - `echo '{}' > empty.json`
  - Windows PowerShell:
    - `Set-Content -Path .\empty.json -Value '{}'`
- Run the CLI (built binary):
  - Cross-platform (after `go build -o epic-updater .`):
    - `./epic-updater free --inputFile empty.json --outputFile out.json`
    - `./epic-updater free_mobile --inputFile empty.json --outputFile out_mobile.json`
  - Windows PowerShell (built `.exe`):
    - `.\epic-updater.exe free --inputFile .\empty.json --outputFile out.json`
    - `.\epic-updater.exe free_mobile --inputFile .\empty.json --outputFile out_mobile.json`
- Run without building (useful during development):
  - `go run . free --inputFile empty.json --outputFile out.json`
  - `go run . free_mobile --inputFile empty.json --outputFile out_mobile.json`
- Inspect the output files (`out.json`, `out_mobile.json`) and console logs for parsing errors or missing fields. Use these smoke runs in CI or locally after changes to `GetFreeGames()`, `FreeMobileGames()`, or GraphQL parsing logic.

Files to inspect for common edits
- Data fetch / parsing: [graphql_get_free_games.go](graphql_get_free_games.go#L1)
 - CLI handling and JSON wiring: [cli_handler_free.go](cli_handler_free.go#L1)
- Program entry and available commands: [main.go](main.go#L1)
- Frontend integration: [index.html](index.html#L1) and [renderTable.js](renderTable.js#L1)
 - (Removed) Dockerfile and old Docker helpers were previously present and have been deleted from the repository.

Primary components and workflow
- Two main components:
  1. CLI updater: the CLI (`main.go` + `cli_handler_free.go` / `cli_handler_free_mobile.go`) runs `free` and `free_mobile` to update `epic_free_games.json`. This is the core weekly job and is scheduled as a GitHub Action in this project.
  2. GitHub Pages UI: static site (`index.html` + `renderTable.js`) reads `epic_free_games.json` and renders the free games list.
- All other CLI commands (`search`, `rate`, `fix_ratings`, `mobile_discover_page`) are for debugging, enrichment, and manual investigation — they are not part of the scheduled updater.

External integrations
- Epic Store GraphQL endpoints:
  - `https://graphql.epicgames.com/graphql` for search/rating queries
  - `https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions` used by `GetFreeGames()`
- The project depends on `github.com/bogdanfinn/tls-client` and `github.com/bogdanfinn/fhttp` (see [go.mod](go.mod#L1)). Changes to HTTP clients should preserve impersonation/UA behavior.

Editing guidance for AI agents
 - Preserve runtime-visible strings and CLI names unless you update all call sites (e.g., `main.go` and build scripts).
- When modifying GraphQL shapes, update local structs used for JSON unmarshalling in the same file (e.g., `FreeGameResponse` in [graphql_get_free_games.go](graphql_get_free_games.go#L1)).
- Keep region/locale flags centralized: prefer changing them in `GetFreeGames()` rather than scattering edits.
- Do not add heavy async concurrency without running the program locally — the code prints and relies on synchronous ordering for constructing output JSON.

When in doubt / where to ask
- Start with reproducing the behaviour locally using `go run main.go free` and inspect `out.json`/console logs.
- If network issues arise, check usage of `tls-client` in `go.mod` and the approach in README research notes.

Frontend dev testing (index.html / renderTable.js)
- Always serve the static UI locally before verifying changes. Run the platform-specific helper when available:
  - Windows PowerShell: `./serve.ps1`
  - Unix / WSL / macOS: `./serve.sh`
- These scripts run a small HTTP server and handle the GitHub Pages prefix used by this repo. Do not rely on generic fallbacks (`python -m http.server`) — they do not replicate the repo prefix handling and produce false positives.
- Verify UI changes using the Playwright MCP viewer (do not install Playwright locally):
  1. Start the local server using the appropriate helper (above) so the site is reachable at `http://localhost:8000/` (or the repo prefix URL the script prints).
  2. Open the Playwright MCP/viewer and point it at the local URL to interactively inspect the page and confirm visual/DOM changes.
  3. Use Playwright MCP's recording or step-through features to capture the verification steps (e.g., navigate to the page, locate `#free-games-table`, confirm rows render).
  4. For CI or automated checks, export a Playwright MCP test or use the project's preferred CI runner — use the MCP viewer during development to create reliable, reproducible checks.
- Use the serve script + Playwright MCP sequence for every change to `index.html` or `renderTable.js` to avoid false positives caused by file/path handling differences between local file serving and GitHub Pages.

If any section is unclear or you need more examples (specific functions or sample log lines), tell me which area and I will expand or add inline examples.
