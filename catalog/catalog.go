// SPDX-FileCopyrightText: 2026 Playground Logic LLC
// SPDX-License-Identifier: Apache-2.0

// Package catalog defines the shared SRE-type catalog schema for the Provabl
// suite. An SRE *type* (e.g. "nih-genomics", "cui-l2") is a named compliance
// profile that maps to the compliance frameworks it must satisfy, where its
// account lands (OU), the tags every account of that type must carry, and the
// baseline stacks to apply.
//
// It lives in provabl/schemas — a single source of truth — because two tools
// consume it: **vendor** (`vendor provision --type <name>` vends an account of a
// type; `vendor catalog` lists/inspects them) and **attest** (attest#98's SRE
// type catalog). Sharing the schema, not the data path, keeps the two from
// drifting — the same reasoning as the qualify↔attest tag-schema decision
// (qualify#32). Neither tool imports the other; both import this.
//
// This package is schema + validation only. It reads/writes catalog JSON and
// checks it is well-formed; it performs no AWS calls and makes no compliance
// claims (attest does that, after a scan).
package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
)

// SchemaVersion is the catalog schema version. It is the cross-repo drift signal:
// a breaking change to the SRE-type shape bumps this, and each consumer pins the
// version it understands (mirrors the attest:* tag schema's SchemaVersion).
const SchemaVersion = 1

// typeKeyRe constrains an SRE-type key to a stable, path/tag-safe token.
var typeKeyRe = regexp.MustCompile(`^[a-z][a-z0-9-]{1,62}$`)

// SREType is one entry in the catalog: a named compliance profile.
type SREType struct {
	// Key is the stable catalog identifier used as `vendor --type <key>`, e.g.
	// "nih-genomics". Lowercase kebab-case; the primary key of the catalog.
	Key string `json:"key"`
	// Name is a human-readable label, e.g. "NIH controlled-access genomics SRE".
	Name string `json:"name"`
	// Description explains what the type is for (one or two sentences).
	Description string `json:"description,omitempty"`
	// Frameworks are the attest framework ids this type must satisfy — passed to
	// `attest compile --frameworks …` as the vending pre-flight (e.g.
	// ["nih-gds", "nist-800-53-moderate"]). At least one is required.
	Frameworks []string `json:"frameworks"`
	// OU is the target Organizational Unit key an account of this type lands in.
	// vendor resolves it against ground-meta's OU ids; "" means the operator must
	// supply --parent explicitly.
	OU string `json:"ou,omitempty"`
	// Tags are the tags every account of this type must carry (data class,
	// cost-center, etc.). Applied at vend time; keys must be non-empty.
	Tags map[string]string `json:"tags,omitempty"`
	// BaselineStacks names the ground CloudFormation baseline stacks to apply to
	// an account of this type (a subset/superset of ground's defaults).
	BaselineStacks []string `json:"baseline_stacks,omitempty"`
}

// Catalog is the full set of SRE types, plus the schema version it conforms to.
type Catalog struct {
	SchemaVersion int       `json:"schema_version"`
	Types         []SREType `json:"types"`
}

// Load parses a catalog from JSON and validates it. A malformed or invalid
// catalog is an error — a consumer should never vend from an unvalidated catalog.
func Load(r io.Reader) (*Catalog, error) {
	var c Catalog
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields() // an unexpected field is a schema drift, not a silent ignore
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("parse catalog: %w", err)
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate checks the catalog is well-formed: a known schema version, unique
// well-formed type keys, and each type naming at least one framework.
func (c *Catalog) Validate() error {
	if c.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported catalog schema_version %d (this build understands %d)", c.SchemaVersion, SchemaVersion)
	}
	seen := make(map[string]bool, len(c.Types))
	for i, t := range c.Types {
		if !typeKeyRe.MatchString(t.Key) {
			return fmt.Errorf("type[%d]: key %q must match %s", i, t.Key, typeKeyRe.String())
		}
		if seen[t.Key] {
			return fmt.Errorf("duplicate type key %q", t.Key)
		}
		seen[t.Key] = true
		if t.Name == "" {
			return fmt.Errorf("type %q: name is required", t.Key)
		}
		if len(t.Frameworks) == 0 {
			return fmt.Errorf("type %q: at least one framework is required (it drives the attest compile pre-flight)", t.Key)
		}
		for k := range t.Tags {
			if k == "" {
				return fmt.Errorf("type %q: a tag key is empty", t.Key)
			}
		}
	}
	return nil
}

// Get returns the SRE type with the given key, or (nil, false) if absent.
func (c *Catalog) Get(key string) (*SREType, bool) {
	for i := range c.Types {
		if c.Types[i].Key == key {
			return &c.Types[i], true
		}
	}
	return nil, false
}

// Keys returns the catalog's type keys, sorted — for stable listing output.
func (c *Catalog) Keys() []string {
	keys := make([]string, 0, len(c.Types))
	for _, t := range c.Types {
		keys = append(keys, t.Key)
	}
	sort.Strings(keys)
	return keys
}
