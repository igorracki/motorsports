# Project Guidelines & Standards

This document outlines the core principles, architecture decisions, and coding standards for the F1 Data platform. All agents and developers must strictly adhere to these guidelines to ensure the project remains clean, maintainable, and extensible.

## 1. Project Scope
*   **Active Projects:** 
    *   `backend` (Go): The main business logic and API gateway.
    *   `fastf1_wrapper` (Python): A specialized adapter for the FastF1 library.
*   **Strictly Ignore:** `frontend_remote`. Do not read, modify, or reference this directory.

## 2. Architecture Decisions

### Service Roles
*   **The Brain (`backend`):** Responsible for orchestration, business logic, validation, and final data presentation (formatting).
*   **The Adapter (`fastf1_wrapper`):** Exists solely to bridge the Python-exclusive FastF1 library. Its job is to fetch raw data, normalize types, and pass them to the backend. It should **not** perform complex domain calculations.

### Data Strategy
*   **Internal Transport:** All durations and times must be transported between services as **milliseconds** (`int64`). This maintains precision and simplifies cross-language communication.
*   **Formatting:** The **Go Backend** is the single source of truth for data formatting. It converts milliseconds into human-readable strings (e.g., `"1:12.345"` or `"+5.230"`).
*   **Gap Handling ("Pass-Through"):** Trust the data source. If FastF1 provides a gap to the winner, pass it through to the backend as a dedicated field (`gap_ms`). Avoid fragile reconstruction logic unless the source data is missing.
*   **Composition over Inheritance:** Use composition for data models.
    *   Example: A `DriverResult` contains core fields and optional pointers to `RaceDetails` or `QualifyingDetails`.
    *   In Go, use pointers with `omitempty` tags to ensure clean JSON responses.
*   **Schema Synchronization:**
    *   Any change to the `fastf1_wrapper` output (e.g., nullable fields) **MUST** be immediately reflected in the `backend` models.
    *   Failing to do so will cause JSON unmarshaling errors in the backend.

## 3. Coding Standards

### General Principles
*   **Readability First:** Code must be clean and neat. If a design looks "weird" or cluttered, refactor it.
*   **Small Functions:** Prefer small, clear functions with a single responsibility. Avoid "mega-functions" with many operations.
*   **Explicit Naming:** 
    *   Variables and functions must be descriptive. Use `RaceWeekend` instead of `Event`. Use `elapsed_duration` instead of `time`.
    *   **Forbidden:** Single-letter variables (`s`, `x`), cryptic abbreviations (`td`, `s_type`), and vague terms (`data`, `obj`). **Exceptions:** `i` (loops), `ctx` (context), `ms` (milliseconds).
*   **Extensibility:** Design features (like session types) so they can be easily extended without breaking existing logic.
*   **Observability (Logging):**
    *   **Traceability:** Add INFO logs at key entry and exit points (e.g., "Fetching session results...", "Found X drivers").
    *   **Error Handling:** Always log full error details in the wrapper before returning safe fallbacks (like `None`). Do not swallow errors silently.

### Python Guidelines (`fastf1_wrapper`)
*   **Models:** Always use `dataclasses`.
*   **Isolation:** Keep the `FastF1Provider` clean.
    *   Move complex extraction logic to `src/core/utils/extractors.py`.
    *   Keep simple type conversions in `src/core/utils/converters.py`.
*   **Accuracy:** Always load sessions with `laps=True` to ensure aggregated statistics (Fastest Lap, Lap Counts) are correctly populated by the library.

### Go Guidelines (`backend`)
*   **Layout:** Follow standard Go project structure (`cmd/` for entry points, `internal/` for private logic).
*   **Pointers:** Use pointers for optional nested data structures to allow for `nil` values and clean JSON output.

## 4. Testing Strategy
*   **Behavioral Testing:** Focus on **Use-Cases** and system behavior rather than strictly unit testing every individual function.
*   **Structure:** Adopt a **Given/When/Then** approach for test clarity.
*   **Mocking Policy:** 
    *   Mock only at the boundaries (e.g., the HTTP call to the Python wrapper).
    *   Prefer real implementations where feasible (e.g., using an in-memory SQLite database instead of mocking a repository) to ensure logic is truly verified.

## 5. Agent Instructions & Safety
*   **Graceful Fallbacks:** Handle "Future Events" (sessions that haven't happened yet) by returning a `200 OK` with an empty results list, allowing the frontend to handle it gracefully.
*   **Defensive Programming:** Always check for `NaN`, `None`, or `NaT` when dealing with external data sources like FastF1.
*   **Verification Workflow:** 
    1.  Perform changes.
    2.  Run `go test ./...` in the `backend`.
    3.  Verify the project builds (`go build`).
    4.  Check Python syntax/imports.
*   **Dependency Management:** NEVER install new libraries or modules (e.g., via `pip`, `go get`, `npm`) without first:
    1.  Explaining **what** the library does.
    2.  Explaining **why** it is strictly necessary.
    3.  Receiving explicit **confirmation** from the user.
