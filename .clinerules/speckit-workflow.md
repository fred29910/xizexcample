## Brief overview
This Cline rule file outlines the development workflow using the `speckit` command suite. It provides guidelines on how to use commands like `/speckit.specify`, `/speckit.plan`, and `/speckit.implement` to manage the software development lifecycle, from specification to implementation.

## Development Workflow
The typical development workflow follows these steps:
1.  **Specify**: Use `/speckit.specify` to create a feature specification from a natural language description.
2.  **Clarify**: Use `/speckit.clarify` to resolve any ambiguities in the specification.
3.  **Plan**: Use `/speckit.plan` to create a technical plan, including architecture and data models.
4.  **Tasks**: Use `/speckit.tasks` to generate a detailed, actionable task list.
5.  **Checklist**: Use `/speckit.checklist` to generate quality checklists for requirements.
6.  **Implement**: Use `/speckit.implement` to execute the implementation plan.
7.  **Analyze**: Use `/speckit.analyze` to perform a consistency and quality analysis.

## Command Usage Guidelines
-   **`/speckit.specify`**: Creates the initial feature specification (`spec.md`). Focus on *what* the feature should do, not *how*.
-   **`/speckit.clarify`**: Asks targeted questions to remove ambiguity from `spec.md`.
-   **`/speckit.plan`**: Generates `plan.md`, `data-model.md`, and other design artifacts.
-   **`/speckit.tasks`**: Creates `tasks.md` with a detailed, dependency-ordered task list.
-   **`/speckit.checklist`**: Generates checklists to validate the quality of requirements.
-   **`/speckit.implement`**: Executes the tasks in `tasks.md`.
-   **`/speckit.analyze`**: A read-only command to check for inconsistencies across artifacts.
-   **`/speckit.constitution`**: Manages the project's guiding principles.

## Technology Stack and Code Style
-   **Active Technologies**: Go (latest stable version) + Zinx, Protobuf.
-   **Project Structure**: Source code should be in `src/` and tests in `tests/`.
-   **Code Style**:
    -   All Go code must be formatted with `gofmt`.
    -   All Go code must adhere to `golint` rules.
    -   Mandatory pre-commit `lint` checks must be performed.
