# Agent Guidelines for Terraform Provider Azure DevOps

This file provides context, rules, and commands for AI coding agents (such as Cursor, Copilot, or CLI agents) operating in the `terraform-provider-azuredevops` repository. 

## 1. Project Overview

*   **Language:** Go 1.24
*   **Provider SDK:** `github.com/hashicorp/terraform-plugin-sdk/v2`
*   **Testing Framework:** `github.com/hashicorp/terraform-plugin-testing`
*   **Architecture:**
    *   `azuredevops/internal/service/<domain>/`: Contains implementation for specific API domains (e.g., `git`, `build`, `core`).
    *   `azuredevops/internal/acceptancetests/`: Contains all acceptance tests (`TestAcc...`).
    *   `azuredevops/internal/client/`: API client configurations.

## 2. Build, Lint, and Test Commands

We use a `GNUmakefile` for orchestration. When running commands, prefer using the provided `make` targets.

### Build and Format
*   **Build the provider:** `make build` (Runs format, dependency checks, and compiles).
*   **Format Go code:** `make fmt` and `make fumpt` (Uses `gofmt` and `gofumpt`).
*   **Format Terraform blocks in tests/docs:** `make terrafmt`.

### Linting
*   **Run Go Linters:** `make lint` (Runs `golangci-lint` using `.golangci.yml`).
*   **Lint Documentation/Website:** `make website-lint`.
*   **Check Dependencies:** `make depscheck` (Ensures `go.mod`, `go.sum`, and `vendor/` are in sync).

### Testing
*   **Run Unit Tests:** `make test`
*   **Run All Acceptance Tests:** `make testacc` (Requires valid `.env` file with Azure DevOps credentials).
*   **Run a Single Acceptance Test:**
    ```bash
    make testacc TEST=./azuredevops/internal/acceptancetests TESTARGS="-run ^TestAccName$"
    ```
    *Alternatively, using Go directly:*
    ```bash
    TF_ACC=1 go test ./azuredevops/internal/acceptancetests -v -run ^TestAccName$ -timeout 120m
    ```

## 3. Code Style & Guidelines

### Formatting and Imports
*   **Formatting:** All Go code must be formatted using `gofumpt` and `goimports`. This is strictly enforced by `golangci-lint`. Always run `make fmt` and `make fumpt` before committing.
*   **Import Grouping:**
    1. Standard library imports.
    2. Third-party imports (e.g., `github.com/hashicorp/...`).
    3. Local module imports (`github.com/microsoft/terraform-provider-azuredevops/...`).

### Error Handling
*   Check all returned errors. Do not silently ignore them (e.g., `_ = err`), except for specifically allowed methods like `ResourceData.Set` or `io.Close` (as defined in `.golangci.yml`).
*   Many existing errors in the codebase use `fmt.Errorf("...: %+v", err)`. When writing new code, align with the surrounding file's convention, but standard Go error wrapping (`%w`) is generally preferred if writing entirely new modules.

### Naming Conventions
*   **Resources & Data Sources:** The file naming convention is `<type>_<resource_name>.go` (e.g., `resource_git_repository.go`, `data_agent_pool.go`).
*   **Acceptance Tests:** Prefix test functions with `TestAcc...` (e.g., `TestAccGitRepository_basic`).
*   **Variables:** Use standard `camelCase`. Exported functions/structs use `PascalCase`.

### Terraform Plugin SDK v2 Specifics
*   Always define `Schema` clearly with appropriate types (`TypeString`, `TypeList`, `TypeSet`, etc.).
*   For complex lists/sets, utilize `TypeSet` with a custom `Set` hash function where appropriate to avoid unnecessary state diffs.
*   Provide robust error messages when API calls fail, ideally including the correlation ID from the Azure DevOps API if available.

### Documentation
*   Documentation files are located in `docs/` and generated via `tfplugindocs`.
*   When adding a new resource or data source, ensure you include examples in the `examples/` directory so they are picked up during doc generation.

## 4. Workflows for Agents

1.  **Exploration:** Use `grep` and `glob` to find relevant domains in `azuredevops/internal/service/` before writing code.
2.  **Implementation:** Create/update the provider logic in the `service` directory.
3.  **Testing:** Every new resource or data source *must* have a corresponding acceptance test in `azuredevops/internal/acceptancetests/`.
4.  **Verification:** Always run `make fumpt`, `make terrafmt`, and `make lint` after making changes to ensure CI will pass.
5.  **No Arbitrary Dependencies:** Do not add third-party libraries unless absolutely necessary. Rely on the Azure DevOps Go SDK (`github.com/microsoft/azure-devops-go-api`) and Terraform Plugin SDK v2.
