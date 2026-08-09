package store

import (
	"context"
	"strings"
	"testing"
)

func TestInitSchemaRunsFinalSchemaForEmptyDatabase(t *testing.T) {
	record := &tokenSecurityRecord{}
	db := openTokenSecurityTestDB(t, record)
	defer db.Close()

	store := &MySQLStore{db: db}
	if err := store.initSchema(context.Background()); err != nil {
		t.Fatalf("initSchema: %v", err)
	}

	findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS roles")
	findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS tenant_access_paths")
	findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS tenant_remote_access_defaults")
	findTokenSecurityCall(t, record, "INSERT IGNORE INTO field_enum")
	chainHops := findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS chain_hops")
	normalizedChainHops := normalizedQuery(chainHops.Query)
	if !strings.Contains(normalizedChainHops, "candidate_index INT NOT NULL") {
		t.Fatalf("chain_hops missing candidate_index: %s", normalizedChainHops)
	}
	if !strings.Contains(normalizedChainHops, "PRIMARY KEY (chain_id, hop_index, candidate_index)") {
		t.Fatalf("chain_hops primary key does not preserve candidate priority: %s", normalizedChainHops)
	}
	chainProbeResults := findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS chain_probe_results")
	if normalized := normalizedQuery(chainProbeResults.Query); !strings.Contains(normalized, "blocking_group_index INT NOT NULL DEFAULT -1") {
		t.Fatalf("chain_probe_results missing blocking group: %s", normalized)
	}
	accessPaths := normalizedQuery(findTokenSecurityCall(t, record, "CREATE TABLE IF NOT EXISTS node_access_paths").Query)
	if !strings.Contains(accessPaths, "remote_protocol VARCHAR(16) NOT NULL DEFAULT ''") {
		t.Fatalf("node_access_paths missing remote protocol: %s", accessPaths)
	}
	for _, removed := range []string{"target_node_id", "entry_node_id", "relay_node_ids_json"} {
		if strings.Contains(accessPaths, removed) {
			t.Fatalf("node_access_paths persists derived field %s: %s", removed, accessPaths)
		}
	}
}

func TestInitSchemaSkipsNonEmptyDatabase(t *testing.T) {
	record := &tokenSecurityRecord{tableCount: 1}
	db := openTokenSecurityTestDB(t, record)
	defer db.Close()

	store := &MySQLStore{db: db}
	if err := store.initSchema(context.Background()); err != nil {
		t.Fatalf("initSchema: %v", err)
	}

	for _, call := range record.snapshot() {
		if strings.HasPrefix(normalizedQuery(call.Query), "CREATE TABLE") {
			t.Fatalf("unexpected schema statement: %s", normalizedQuery(call.Query))
		}
	}
}
