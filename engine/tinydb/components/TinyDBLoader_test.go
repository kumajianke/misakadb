package components

import (
	"errors"
	mson "misakadb/engine/Mson"
	"os"
	"strings"
	"testing"
)

func TestInitLoaderReturnsStatErrorWhenPathCheckFails(t *testing.T) {
	tmpDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("chdir to temp dir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	err = os.WriteFile("db-datas", []byte("not-a-directory"), 0600)
	if err != nil {
		t.Fatalf("create fake db-datas file failed: %v", err)
	}

	loader := &TinyDBLoaderImp{DBName: "demo"}
	err = loader.InitLoader(mson.MsonParse{Name: "demo"})
	if err == nil {
		t.Fatal("expected InitLoader to return stat error")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected non-IsNotExist stat error, got: %v", err)
	}
	if err.Error() == "database is exist!" {
		t.Fatalf("expected real stat error instead of exist message, got: %v", err)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "directory") {
		t.Fatalf("expected directory-related stat error, got: %v", err)
	}
}
