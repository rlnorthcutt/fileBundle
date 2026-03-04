# filebundle

`filebundle` is a Go-based CLI tool that crawls a local directory tree, filters files by subdirectory, extension, and exclude patterns, and flattens everything into a single, TOC-indexed text file that is easy to feed into AI systems.

The primary use case is to bundle documentation, notes, or code into one or more clearly delimited text file so an LLM can ingest it as high-quality context.

## Features

- **Crawl a local directory** starting from a configurable root.
- **Filter by subdirectory** (`--include`) and **extension** (`--extensions`).
- **Exclude noisy paths** (e.g. `.git`, `node_modules`, `bin`) with simple substring patterns.
- **Generate a TABLE OF CONTENTS** section listing every bundled file by relative path.
- **Bundle file contents** into a single flat text file using clear dividers:

  ```text
  -----------------------------------------------
  relative/path/to/file.ext
  -----------------------------------------------

  <file contents>

  ```

- **Interactive or flag-driven usage**:
  - If flags are missing, `filebundle` will prompt for values in a friendly, stepped flow.
  - If you provide all flags, it runs non-interactively.
- **Progress bar** while reading and bundling files, so you can see activity on large trees.
- **Graceful skips** for unreadable files (permissions issues are logged but don’t stop the run).

## Installation

### Easy: Build locally

You’ll need [Go](https://golang.org/doc/install) installed.

1. Clone the repository:

   ```bash
   git clone https://github.com/rlnorthcutt/fileBundle.git
   cd fileBundle
   ```

2. Build the CLI tool:

   ```bash
   go build -o filebundle .
   ```

   For a smaller binary, you can use:

   ```bash
   go build -ldflags="-s -w" -o filebundle .
   ```

This will generate the `filebundle` binary in the project root.

## Usage

Once built, you can run the tool from the command line. The tool supports both interactive prompts and command-line flags.

### Interactive Usage

```bash
./filebundle
```

### Example Interactive Session

```bash
$ ./filebundle
Enter the root directory to crawl (default: '.'):
Enter subdirectories to include (comma-separated, '*' for all):
Enter file extensions to include (e.g., 'md,txt'):
Enter patterns to exclude (e.g., '.git,node_modules'):
Enter the output filename (default: 'bundle.txt'):

Preparing to bundle files with the following settings:
  Input Root:  .
  Include:     *
  Extensions:  md,txt
  Exclude:     .git,node_modules,bin
  Output File: bundle.txt

Do you want to proceed? (y/n): y

Successfully bundled 42 files into bundle.txt
```

### Command-Line Options (Non-Interactive)

If you prefer to pass flags instead of interactive prompts, you can run:

```bash
./filebundle \
  --input="./docs" \
  --include="*,guides" \
  --extensions="md,txt" \
  --exclude=".git,node_modules,bin" \
  --output="docs-bundle.txt"
```

Or, use the short flags:

```bash
./filebundle \
  -i="./docs" \
  -d="*,guides" \
  -e="md,txt" \
  -x=".git,node_modules,bin" \
  -o="docs-bundle.txt"
```

### Flags

- `--input, -i`  
  **Root directory** to crawl.  
  Default: `.` (current directory)

- `--include, -d`  
  **Subdirectories to include**, as a comma-separated list of top-level subdirs under the root.  
  Examples: `"*"`, `"docs,src,notes"`  
  Default: `*` (include all subdirectories)

- `--extensions, -e`  
  **File extensions to bundle**, comma-separated.  
  Accepts values with or without leading dots (e.g. `md,txt` or `.md,.txt`).  
  Default: `md,txt`

- `--exclude, -x`  
  **Patterns to skip**, as comma-separated substrings. Any file or directory whose **relative path contains** one of these substrings will be excluded.  
  Example: `.git,node_modules,bin`  
  Default: `.git,node_modules,bin`

- `--output, -o`  
  **Output filename** (relative to the current working directory).  
  Default: `bundle.txt`

## Output Format

The generated bundle file is a single text file with two main sections:

1. **TABLE OF CONTENTS**  
   A header listing all bundled files by relative path, one per line:

   ```text
   TABLE OF CONTENTS
   docs/intro.md
   docs/setup.md
   notes/todo.txt

   ```

2. **Bundled file contents**  
   For each file, `filebundle` appends a clearly delimited block:

   ```text
   -----------------------------------------------
   docs/intro.md
   -----------------------------------------------

   # Intro
   Welcome to the docs...

   -----------------------------------------------
   notes/todo.txt
   -----------------------------------------------

   - item one
   - item two

   ```

This flat, divider-based format is optimized so LLMs can quickly understand where each file begins and ends.

## Project Structure

```bash
fileBundle/
├── main.go        # CLI entry point, flags, interactive flow
├── crawler/       # Recursive directory walking and filtering
│   └── crawler.go
├── bundler/       # Reads file contents and applies divider formatting
│   └── bundler.go
├── writer/        # Generates TOC and writes the final bundle to disk
│   └── writer.go
├── go.mod         # Go module file with dependencies
├── go.sum         # Go module dependency checksums
├── CLAUDE.md      # Internal design and implementation notes
└── README.md      # Project documentation
```

## Dependencies

`filebundle` uses the following Go packages:

- [`github.com/spf13/cobra`](https://github.com/spf13/cobra) - For CLI command and flag management.
- [`github.com/schollz/progressbar/v3`](https://github.com/schollz/progressbar) - For showing progress bars while bundling files.

## Contributing

Issues and pull requests are welcome for new features, bug fixes, and general improvements.

## License

This project is licensed under the MIT License.

