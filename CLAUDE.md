## Purpose

* **Goal**: Build and maintain `filebundle`, a Go-based CLI tool that crawls local directories, filters files by extension/path, and flattens them into a single, TOC-indexed text file for AI-friendly consumption.

* **Behavior**: Prefer small, focused changes. Mirror the "interactive prompt" pattern for any missing flags.

## Development environment

* **Dev container**:
* Use the official Go image (Go 1.25+).
* **Do**: Install tools (linters, etc.) via Go modules or the dev container config.


* **Dependencies**:
* **Respect** `go.mod`/`go.sum`.
* **Primary Libraries**: `github.com/spf13/cobra` for CLI and `github.com/schollz/progressbar/v3` for crawl feedback.



## Architecture and code organization

* **Overall structure**:
* **`main.go`**: Single entry point. Defines Cobra flags and the interactive `promptUser` flow.
* **`crawler`**: Handles recursive directory walking, ignore-pattern logic, and file discovery.
* **`bundler`**: Responsible for reading file contents and applying the specific text formatting (dividers, paths).
* **`writer`**: Handles TOC generation and final disk IO for the output bundle.


* **Separation of concerns**:
* `main`: Orchestrates (Prompt → Crawl → Bundle → Write).
* `crawler`: Filters based on `--include` (subdirs), `--extensions`, and `--exclude`.
* `bundler`: Transforms raw file bytes into the delimited string format.



## CLI design and flag handling

* **Use Cobra**:
* Attach a single `Run` function to `rootCmd` that calls `executeBundle`.


* **Flag conventions**:
* `--input, -i`: Root directory (Default: `.`).
* `--include, -d`: Subdirectories to include (Default: `*`).
* `--extensions, -e`: Comma-separated list (Default: `md,txt`).
* `--output, -o`: Resulting filename (Default: `bundle.txt`).
* `--exclude, -x`: Patterns to skip (Default: `.git,node_modules,bin`).


* **Validation**:
* Check if the input path exists and is a directory.
* Validate that the output filename is not the same as an input file being bundled.



## Interactive, stepped input flow

* **Pattern**: Support non-interactive (flags only) and interactive (missing flags requested via `promptUser`).
* **Order of prompts**:
1. Input directory (`--input`).
2. Target subdirectories (`--include`).
3. Extensions to bundle (`--extensions`).
4. Output filename (`--output`).


* **Confirmation**:
* Print a summary: "Crawling [Input] for [Extensions] in [Subdirs]... Output to [Output]."
* Require `y/n` confirmation before reading files.



## Bundle Format Standards

The tool must generate a text file with these specific sections:

1. **Header**: "TABLE OF CONTENTS" section listing all relative paths found.
2. **Dividers**: Each file content block must be wrapped in:
```text
-----------------------------------------------
relative/path/to/file.ext
-----------------------------------------------

```



## Error handling and user feedback

* **Centralized error handling**: Use `handleError(step string, err error)` to exit with clear context.
* **Progress feedback**:
* Use a progress bar when reading and bundling files to show activity in large repos (like HAProxy docs).


* **Graceful skips**: If a file exists but is unreadable (permissions), log the warning but continue bundling other files.

## How Claude should work in this repo

* **Mirroring**: Look at `main.go` for the `promptUser` and `handleError` implementations and use them exactly for new features.
* **Focus**: When adding a new filter (like regex exclusion), place the logic in the `crawler` package, not in `main.go`.
* **LLM Optimization**: Always ensure the final output format remains "flat" and clearly delimited, as this is the primary value of the tool for AI users.