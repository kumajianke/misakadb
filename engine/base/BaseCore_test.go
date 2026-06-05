package engine_base

import "testing"

func TestEngineLockerSupportGeneratesUniqueDefaultNamespacePerInstance(t *testing.T) {
	first := &EngineLockerSupport{}
	second := &EngineLockerSupport{}

	firstKey := first.lockKey("row:demo")
	secondKey := second.lockKey("row:demo")

	if firstKey == secondKey {
		t.Fatalf("expected different default lock keys, got same key %q", firstKey)
	}
}

func TestEngineLockerSupportDefaultNamespaceIsStableForSameInstance(t *testing.T) {
	locker := &EngineLockerSupport{}

	firstKey := locker.lockKey("engine")
	secondKey := locker.lockKey("engine")

	if firstKey != secondKey {
		t.Fatalf("expected stable lock key for same instance, got %q and %q", firstKey, secondKey)
	}
}

func TestEngineLockerSupportUsesExplicitNamespace(t *testing.T) {
	locker := &EngineLockerSupport{LockNamespace: "tinydb:test"}

	key := locker.lockKey("row:demo")
	if key != "tinydb:test:row:demo" {
		t.Fatalf("expected explicit namespace to be used, got %q", key)
	}
}
