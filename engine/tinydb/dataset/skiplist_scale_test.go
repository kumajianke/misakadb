package dataset

import "testing"

const millionDatasetSize = 1_000_000

func buildMillionSkipList() *SkipList[*PassStuct] {
	skipList := NewSkipList[*PassStuct]()
	for key := range millionDatasetSize {
		skipList.Insert(key, &PassStuct{Name: "bulk"})
	}
	return skipList
}

func TestGetWithPageMillionData(t *testing.T) {
	skipList := buildMillionSkipList()

	rows, err := skipList.GetWithPage(">=", 500_000, 100, 5)
	if err != nil {
		t.Fatalf("GetWithPage returned error: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(rows))
	}

	expectedKeys := []int{500_100, 500_101, 500_102, 500_103, 500_104}
	for index, expectedKey := range expectedKeys {
		if rows[index].Key != expectedKey {
			t.Fatalf("unexpected key at index %d: want %d, got %d", index, expectedKey, rows[index].Key)
		}
		if len(rows[index].Data) != 1 {
			t.Fatalf("expected one value for key %d, got %d", rows[index].Key, len(rows[index].Data))
		}
	}
}

func TestDeleteWithMillionData(t *testing.T) {
	skipList := buildMillionSkipList()

	deleted, err := skipList.DeleteWith("like", "12%", 1_000)
	if err != nil {
		t.Fatalf("DeleteWith returned error: %v", err)
	}
	if deleted != 1_000 {
		t.Fatalf("expected 1000 deleted keys, got %d", deleted)
	}

	if skipList.Contains(12) {
		t.Fatal("expected key 12 to be deleted")
	}
	if !skipList.Contains(13) {
		t.Fatal("expected key 13 to remain")
	}
	if !skipList.Contains(999_999) {
		t.Fatal("expected high key to remain")
	}
}
