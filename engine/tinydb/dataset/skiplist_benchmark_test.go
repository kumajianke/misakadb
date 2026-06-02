package dataset

import (
	"sync"
	"testing"
)

var (
	benchmarkOnce     sync.Once
	benchmarkSkipList *SkipList[*PassStuct]
)

func getBenchmarkSkipList() *SkipList[*PassStuct] {
	benchmarkOnce.Do(func() {
		benchmarkSkipList = buildMillionSkipList()
	})
	return benchmarkSkipList
}

func BenchmarkGetMillionData(b *testing.B) {
	skipList := getBenchmarkSkipList()
	b.ResetTimer()

	index := 0
	for b.Loop() {
		key := (index*7919 + 12345) % millionDatasetSize
		values, ok := skipList.Get(key)
		if !ok || len(values) == 0 {
			b.Fatalf("expected key %d to exist", key)
		}
		index++
	}
}

func BenchmarkGetWithPageMillionData(b *testing.B) {
	skipList := getBenchmarkSkipList()
	b.ResetTimer()

	for b.Loop() {
		rows, err := skipList.GetWithPage(">=", 500_000, 100, 5)
		if err != nil {
			b.Fatalf("GetWithPage returned error: %v", err)
		}
		if len(rows) != 5 {
			b.Fatalf("expected 5 rows, got %d", len(rows))
		}
	}
}
