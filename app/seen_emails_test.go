package app

import (
	"testing"

	"github.com/floatpane/matcha/fetcher"
	"github.com/floatpane/matcha/internal/emailstore"
)

func uids(emails []fetcher.Email) []uint32 {
	out := make([]uint32, len(emails))
	for i, e := range emails {
		out[i] = e.UID
	}
	return out
}

func TestNewlyReceived(t *testing.T) {
	m := &Model{store: &emailstore.Store{FolderEmails: map[string][]fetcher.Email{}}}
	a := fetcher.Email{UID: 1, AccountID: "x"}
	b := fetcher.Email{UID: 2, AccountID: "x"}
	c := fetcher.Email{UID: 2, AccountID: "y"}

	if got := m.newlyReceived("INBOX", []fetcher.Email{a, b}); len(got) != 0 {
		t.Fatalf("initial sync should report nothing new, got %v", uids(got))
	}
	if got := m.newlyReceived("INBOX", []fetcher.Email{a, b}); len(got) != 0 {
		t.Fatalf("unchanged refetch should report nothing new, got %v", uids(got))
	}
	got := m.newlyReceived("INBOX", []fetcher.Email{c, a, b})
	if len(got) != 1 || got[0].AccountID != "y" {
		t.Fatalf("expected only the new email from account y, got %+v", got)
	}

	if got := m.newlyReceived("Empty", nil); len(got) != 0 {
		t.Fatalf("empty initial sync, got %v", uids(got))
	}
	if got := m.newlyReceived("Empty", []fetcher.Email{a}); len(got) != 1 {
		t.Fatalf("first email in a previously empty folder should be new, got %v", uids(got))
	}
}

func TestNewlyReceivedUsesCacheAsBaseline(t *testing.T) {
	a := fetcher.Email{UID: 1, AccountID: "x"}
	b := fetcher.Email{UID: 2, AccountID: "x"}
	m := &Model{store: &emailstore.Store{FolderEmails: map[string][]fetcher.Email{"INBOX": {a}}}}

	got := m.newlyReceived("INBOX", []fetcher.Email{b, a})
	if len(got) != 1 || got[0].UID != 2 {
		t.Fatalf("expected UID 2 to be new relative to the cache, got %v", uids(got))
	}
}
