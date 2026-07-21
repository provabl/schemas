# Changelog

All notable changes to schemas will be documented in this file.

The format is based on [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Initial module** — `github.com/provabl/schemas`, the suite's shared-schema home
  (imported inward like the evidence kernel; imports no tool). Go 1.26.5,
  Apache-2.0 / Playground Logic LLC, CI (Check + Lint) + weekly Security Scan.
- **`catalog`** — the **SRE-type catalog** schema (`SchemaVersion = 1`): an `SREType`
  maps a stable catalog key (e.g. `nih-genomics`) to its compliance `Frameworks`, target
  `OU`, required `Tags`, and `BaselineStacks`. `Load` parses + validates catalog JSON
  (`DisallowUnknownFields`, well-formed unique keys, ≥1 framework per type); `Get`/`Keys`
  read it. Shared by **vendor** (`vendor provision --type …`) and **attest** (attest#98) —
  one schema, not two (mirrors the qualify#32 tag-schema decision). Fully unit-tested.
