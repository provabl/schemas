// SPDX-FileCopyrightText: 2026 Playground Logic LLC
// SPDX-License-Identifier: Apache-2.0

package catalog

import (
	"strings"
	"testing"
)

const validCatalog = `{
  "schema_version": 1,
  "types": [
    {
      "key": "nih-genomics",
      "name": "NIH controlled-access genomics SRE",
      "description": "dbGaP / GDS controlled-access genomic data.",
      "frameworks": ["nih-gds", "nist-800-53-moderate"],
      "ou": "SensitiveResearch",
      "tags": {"data-class": "GENOMIC", "cost-center": "research"},
      "baseline_stacks": ["logging", "security", "network"]
    },
    {
      "key": "cui-l2",
      "name": "CUI CMMC Level 2 SRE",
      "frameworks": ["cmmc-level-2"]
    }
  ]
}`

func TestLoad_Valid(t *testing.T) {
	c, err := Load(strings.NewReader(validCatalog))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(c.Types) != 2 {
		t.Fatalf("got %d types, want 2", len(c.Types))
	}
	nih, ok := c.Get("nih-genomics")
	if !ok {
		t.Fatal("nih-genomics not found")
	}
	if nih.OU != "SensitiveResearch" || nih.Tags["data-class"] != "GENOMIC" {
		t.Errorf("nih-genomics fields wrong: %+v", nih)
	}
	if len(nih.Frameworks) != 2 {
		t.Errorf("frameworks = %v", nih.Frameworks)
	}
}

func TestKeys_Sorted(t *testing.T) {
	c, _ := Load(strings.NewReader(validCatalog))
	got := c.Keys()
	if len(got) != 2 || got[0] != "cui-l2" || got[1] != "nih-genomics" {
		t.Errorf("Keys() = %v, want [cui-l2 nih-genomics]", got)
	}
}

func TestGet_Missing(t *testing.T) {
	c, _ := Load(strings.NewReader(validCatalog))
	if _, ok := c.Get("nope"); ok {
		t.Error("expected miss for an unknown key")
	}
}

func TestLoad_RejectsBadVersion(t *testing.T) {
	if _, err := Load(strings.NewReader(`{"schema_version": 99, "types": []}`)); err == nil {
		t.Error("unsupported schema_version must error")
	}
}

func TestLoad_RejectsUnknownField(t *testing.T) {
	// DisallowUnknownFields: a stray field signals drift, not silent tolerance.
	j := `{"schema_version": 1, "types": [], "extra": true}`
	if _, err := Load(strings.NewReader(j)); err == nil {
		t.Error("unknown top-level field must error")
	}
}

func TestValidate_TypeRules(t *testing.T) {
	cases := map[string]string{
		"bad key (uppercase)": `{"schema_version":1,"types":[{"key":"Bad","name":"x","frameworks":["f"]}]}`,
		"empty name":          `{"schema_version":1,"types":[{"key":"ok","name":"","frameworks":["f"]}]}`,
		"no frameworks":       `{"schema_version":1,"types":[{"key":"ok","name":"x","frameworks":[]}]}`,
		"empty tag key":       `{"schema_version":1,"types":[{"key":"ok","name":"x","frameworks":["f"],"tags":{"":"v"}}]}`,
	}
	for name, j := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Load(strings.NewReader(j)); err == nil {
				t.Errorf("%s: expected a validation error", name)
			}
		})
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	j := `{"schema_version":1,"types":[
	  {"key":"dup","name":"a","frameworks":["f"]},
	  {"key":"dup","name":"b","frameworks":["g"]}]}`
	if _, err := Load(strings.NewReader(j)); err == nil {
		t.Error("duplicate type key must error")
	}
}

func TestValidate_KeyLengthBounds(t *testing.T) {
	// single char too short (regex requires >=2), and a too-long key.
	short := `{"schema_version":1,"types":[{"key":"a","name":"x","frameworks":["f"]}]}`
	if _, err := Load(strings.NewReader(short)); err == nil {
		t.Error("1-char key must error")
	}
	long := `{"schema_version":1,"types":[{"key":"` + strings.Repeat("a", 80) + `","name":"x","frameworks":["f"]}]}`
	if _, err := Load(strings.NewReader(long)); err == nil {
		t.Error("over-long key must error")
	}
}
