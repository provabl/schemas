# Changelog

All notable changes to schemas will be documented in this file.

The format is based on [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Security

- **Bump Go 1.26.5 → 1.26.6** to clear five Go standard-library advisories
  (**GO-2026-6218**, **GO-2026-6090**, **GO-2026-6088**, **GO-2026-5972**,
  **GO-2026-5026** — `net/url`, `crypto/tls`, `encoding/xml`, `encoding/asn1`, and
  `x/net/idna` via `net/http`), all fixed in go1.26.6. Toolchain bump only — no code
  changes.

## [0.1.0] - 2026-07-21

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

[Unreleased]: https://github.com/provabl/schemas/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/provabl/schemas/releases/tag/v0.1.0
