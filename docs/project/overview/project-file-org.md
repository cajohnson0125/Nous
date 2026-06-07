# File Organization

```text
Nous
* /

- go.mod
- go.sum
- Makefile
- README.md

** docs/
*** project/
*** coding-guides/
*** obsolete/

** cmd/
*** nous/
- main.go

** internal/
- app.go
*** config/
*** styles/
*** layout/
- layout.go
**** components/
- <layout-component1>.go
*** state-engine/
*** library/
*** archives/
*** status/
*** intuition/
*** experts/
*** tooling/
*** llm/
```

## Organization Principles

- Use kebab-case for file and folder names.
- Packages should represent architectural responsibilities rather than implementation patterns.

The primary package boundaries are derived directly from the system architecture:

* State Engine
* Library
* Archives
* Status
* Intuition
* Experts
* Tooling
* LLM Providers

Additional files and sub-packages should be created only when justified by implementation complexity.

Package organization should emerge from the code rather than being predetermined.
