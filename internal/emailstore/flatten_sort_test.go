package emailstore

import (
	"testing"
	"time"

	"github.com/floatpane/matcha/fetcher"
)

func BenchmarkFlattenAndSortLarge(b *testing.B) {
	emails := make([]fetcher.Email, 27000)
	for i := range emails {
		emails[i] = fetcher.Email{UID: uint32(i), Date: time.Unix(int64((i*7919)%100000), 0)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := map[string][]fetcher.Email{"a": append([]fetcher.Email(nil), emails...)}
		_ = FlattenAndSort(in)
	}
}

func TestFlattenAndSortNewestFirst(t *testing.T) {
	now := time.Now()
	got := FlattenAndSort(map[string][]fetcher.Email{
		"a": {{UID: 1, Date: now.Add(-2 * time.Hour)}, {UID: 2, Date: now}},
		"b": {{UID: 3, Date: now.Add(-time.Hour)}},
	})
	want := []uint32{2, 3, 1}
	for i, e := range got {
		if e.UID != want[i] {
			t.Fatalf("position %d: got UID %d, want %d", i, e.UID, want[i])
		}
	}
}
