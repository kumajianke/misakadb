package dataset

import "testing"

func TestGetWithPageRangeQueries(t *testing.T) {
	skipList := NewSkipList[*PassStuct]()
	for _, key := range []int{10, 20, 30, 40, 50} {
		skipList.Insert(key, &PassStuct{Name: "v"})
	}

	rows, err := skipList.GetWithPage("=", 30, 0, 10)
	if err != nil {
		t.Fatalf("equal query returned error: %v", err)
	}
	if len(rows) != 1 || rows[0].Key != 30 {
		t.Fatalf("unexpected equal rows: %#v", rows)
	}

	rows, err = skipList.GetWithPage(">", 20, 1, 2)
	if err != nil {
		t.Fatalf("greater query returned error: %v", err)
	}
	if len(rows) != 2 || rows[0].Key != 40 || rows[1].Key != 50 {
		t.Fatalf("unexpected greater rows: %#v", rows)
	}

	rows, err = skipList.GetWithPage("<", 40, 1, 2)
	if err != nil {
		t.Fatalf("less query returned error: %v", err)
	}
	if len(rows) != 2 || rows[0].Key != 20 || rows[1].Key != 30 {
		t.Fatalf("unexpected less rows: %#v", rows)
	}

	rows, err = skipList.GetWithPage("<=", 40, 3, 10)
	if err != nil {
		t.Fatalf("less-or-equal query returned error: %v", err)
	}
	if len(rows) != 1 || rows[0].Key != 40 {
		t.Fatalf("unexpected less-or-equal rows: %#v", rows)
	}
}

func TestDeleteWithRangeQueries(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		skipList := NewSkipList[*PassStuct]()
		for _, key := range []int{10, 20, 30, 40, 50} {
			skipList.Insert(key, &PassStuct{Name: "v"})
		}

		deleted, err := skipList.DeleteWith("=", 30, 10)
		if err != nil {
			t.Fatalf("equal delete returned error: %v", err)
		}
		if deleted != 1 {
			t.Fatalf("expected 1 deleted key, got %d", deleted)
		}
		if skipList.Contains(30) {
			t.Fatal("expected key 30 to be deleted")
		}
		if !skipList.Contains(40) {
			t.Fatal("expected key 40 to remain")
		}
	})

	t.Run("greater", func(t *testing.T) {
		skipList := NewSkipList[*PassStuct]()
		for _, key := range []int{10, 20, 30, 40, 50} {
			skipList.Insert(key, &PassStuct{Name: "v"})
		}

		deleted, err := skipList.DeleteWith(">", 20, 2)
		if err != nil {
			t.Fatalf("greater delete returned error: %v", err)
		}
		if deleted != 2 {
			t.Fatalf("expected 2 deleted keys, got %d", deleted)
		}
		if skipList.Contains(30) || skipList.Contains(40) {
			t.Fatal("expected keys 30 and 40 to be deleted")
		}
		if !skipList.Contains(50) {
			t.Fatal("expected key 50 to remain because of limit")
		}
	})

	t.Run("greater or equal unlimited", func(t *testing.T) {
		skipList := NewSkipList[*PassStuct]()
		for _, key := range []int{10, 20, 30, 40, 50} {
			skipList.Insert(key, &PassStuct{Name: "v"})
		}

		deleted, err := skipList.DeleteWith(">=", 40, 0)
		if err != nil {
			t.Fatalf("greater-or-equal delete returned error: %v", err)
		}
		if deleted != 2 {
			t.Fatalf("expected 2 deleted keys, got %d", deleted)
		}
		if skipList.Contains(40) || skipList.Contains(50) {
			t.Fatal("expected keys 40 and 50 to be deleted")
		}
		if !skipList.Contains(30) {
			t.Fatal("expected key 30 to remain")
		}
	})
}
