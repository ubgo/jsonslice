# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.1.0] - 2026-05-01

### Added

- Initial release.
- `JsonSlice []any` with `encoding/json` round-trip and explicit handling of
  the JSON `null` value.
- `database/sql/driver` integration via `Value()` and `Scan()` for storing
  the slice in JSON / JSONB columns.
- gqlgen-compatible `MarshalGQL` / `UnmarshalGQL` methods detected by duck
  typing — no gqlgen import is required.

[Unreleased]: https://github.com/ubgo/jsonslice/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/ubgo/jsonslice/releases/tag/v0.1.0
