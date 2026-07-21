# schemas

**Shared schemas for the [Provabl](https://provabl.dev) suite.**

A small, dependency-light Go module holding schema definitions that more than one
Provabl tool must agree on — so the contract lives in one place and can't drift
between consumers. It makes no AWS calls and no compliance claims; it is types +
validation only.

> Part of the [Provabl](https://provabl.dev) suite. Like the [evidence](https://github.com/provabl/evidence)
> kernel, this is a **shared dependency imported inward** — tools import `schemas`;
> `schemas` imports no tool.

## Packages

| Package | What | Consumers |
|---|---|---|
| `catalog` | The **SRE-type catalog** schema — a named compliance profile (`nih-genomics`, `cui-l2`, …) mapping to frameworks, OU placement, required tags, and baseline stacks. | **vendor** (`vendor provision --type …` / `vendor catalog`) and **attest** (attest#98). One schema, not two. |

## Why a shared module

When two tools need to agree on the shape of an artifact, embedding the schema in
one of them and importing it into the other couples their release cycles and lets
the definition drift. A dedicated schema module — imported by both, importing
neither — keeps a single source of truth. This mirrors the decision behind the
`attest:*` tag schema (qualify#32): the schema is the shared artifact; the code
paths stay independent.

## Versioning

Each schema carries a `SchemaVersion` constant — the cross-repo drift signal. A
breaking change to a schema's shape bumps it, and each consumer pins the version
it understands. The module follows [semver](https://semver.org/spec/v2.0.0.html);
a schema-version bump is a minor (new, back-compatible) or major (breaking) change.

## License

Apache 2.0. Copyright 2026 Playground Logic LLC.
