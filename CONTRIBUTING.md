# Contributing to ubgo/jsonslice

Thanks for your interest in `ubgo/jsonslice`. This repository is licensed under the **Apache License 2.0**. Pull requests are welcome.

## Workflow

1. Open an issue first for anything beyond a tiny fix. Discussing the design upfront avoids wasted work.
2. Fork + branch named after the issue: `fix/123-null-scan`, `feat/456-typed-element`.
3. Run local checks: `task ci`.
4. Use Conventional Commits for the PR title.

## Code conventions

- **Zero third-party deps.** `go.mod` for the core must stay stdlib-only — the gqlgen integration works via duck typing precisely so consumers never inherit a gqlgen import.
- **Race detector mandatory.** Every test must pass under `-race`.
- **Coverage target.** ≥ 90% line coverage.
- **Public API stability.** Once the module reaches v1.0.0, breaking changes require a major version bump and a strong rationale.
- **No comments explaining what the code does.** Names should make that clear. Reserve comments for the *why* — non-obvious invariants, hidden constraints, surprising tradeoffs.

## Testing locally

```sh
task test           # standard tests
task test:race      # with race detector
task test:coverage  # with coverage report
task lint           # golangci-lint
task ci             # everything
```

## License of contributions

By submitting a pull request, you agree that your contribution is provided under the same Apache License 2.0 as the rest of the repository.
