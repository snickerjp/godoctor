# Go Development Guidelines
All code contributed to this project must adhere to the [Google Go Style Guide](https://google.github.io/styleguide/go/) and the following principles.

## 1. Formatting
All Go code **must** be formatted with `gofmt` before being submitted.

## 2. Naming Conventions
- **Packages:** Use short, concise, all-lowercase names.
- **Variables, Functions, and Methods:** Use `camelCase` for unexported identifiers and `PascalCase` for exported identifiers.
- **Interfaces:** Name interfaces for what they do (e.g., `io.Reader`), not with a prefix like `I`.

## 3. Error Handling
- Errors are values. Do not discard them.
- Handle errors explicitly using the `if err != nil` pattern.
- Provide context to errors using `fmt.Errorf("context: %w", err)`.

## 4. Simplicity and Clarity
- "Clear is better than clever." Write code that is easy to understand.
- Avoid unnecessary complexity and abstractions.
- Prefer returning concrete types, not interfaces.

## 5. Documentation
- All exported identifiers (`PascalCase`) **must** have a doc comment.
- Comments should explain the *why*, not the *what*.

## 6. Project structure
- cmd/ contains source code for target binaries (e.g. server, client)
- internal/ contains source code for packages not meant to be exported (e.g. internal/tools/hello)
- bin/ contains the compiled binaries
- At the root place README.md, go.mod and go.sum

## 7. Execution Environment
- All Go commands (build, run, test, fmt) must be executed within the Docker container using `docker compose`.
- Example: `sudo docker compose exec app go <command>`

## 8. Documentation Retrieval
- **Always** use the `read_docs` tool to retrieve documentation for Go packages or symbols.
- This is mandatory when:
    - Seeing an import for the first time in a session.
    - After a new dependency is installed (e.g., via `go get`).

## 9. Binary Distribution
- Build artifacts in `bin/` must be **standalone binaries** that can run on the host environment (linux) without the Go toolchain.
- Use `CGO_ENABLED=0` to ensure static linking.

## 10. Testing and Validation
- When running or testing server-like binaries that may block (e.g., `server`), **always use a timeout** (e.g., `timeout 2s ./bin/server`) to prevent the execution from hanging indefinitely.

