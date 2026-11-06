package tui

import (
	"unicode/utf8"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/floatpane/matcha/config"
)

// QuitKeyAware is implemented by screens that know whether the user is
// currently typing. Screens that don't implement it are treated as typing,
// so a printable quit key (e.g. "q") never quits from inside a text field.
type QuitKeyAware interface {
	AcceptsQuitKey() bool
}

// IsQuitKey reports whether key should quit the app while screen is active.
// ctrl+c always quits. The configured quit key always quits when it is a
// modifier/special key (e.g. "ctrl+q"); when it is a printable character it
// only quits on screens that report they are not accepting text input.
func IsQuitKey(screen tea.Model, key string) bool {
	if key == keyCtrlC {
		return true
	}
	quit := config.Keybinds.Global.Quit
	if quit == "" || key != quit {
		return false
	}
	if !isPrintableKey(quit) {
		return true
	}
	if a, ok := screen.(QuitKeyAware); ok {
		return a.AcceptsQuitKey()
	}
	return false
}

func isPrintableKey(key string) bool {
	return key == keySpace || utf8.RuneCountInString(key) == 1
}

func (m Choice) AcceptsQuitKey() bool { return true }

func (m *EmailView) AcceptsQuitKey() bool { return true }

func (m *Drafts) AcceptsQuitKey() bool {
	return m.list.FilterState() != list.Filtering
}

func (m *Inbox) AcceptsQuitKey() bool {
	return m != nil && m.searchOverlay == nil && m.list.FilterState() != list.Filtering
}

func (m *FolderInbox) AcceptsQuitKey() bool {
	if m.jumpingToFolder {
		return false
	}
	return m.inbox.AcceptsQuitKey()
}

func (m *Settings) AcceptsQuitKey() bool {
	if m.activePane != PaneContent {
		return true
	}
	switch m.activeCategory {
	case CategoryAccounts:
		return !m.isCryptoConfig
	case CategoryEncryption:
		return false
	case CategoryPlugins:
		return !m.pluginEditing
	case CategoryGeneral, CategoryTheme, CategoryMailingLists, CategoryContacts:
		return true
	}
	return true
}
