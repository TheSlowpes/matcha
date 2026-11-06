package tui

import (
	"testing"

	"github.com/floatpane/matcha/config"
)

func TestIsQuitKey(t *testing.T) {
	orig := config.Keybinds
	t.Cleanup(func() { config.Keybinds = orig })

	inbox := NewInbox(nil, nil)
	composer := NewComposer("", "", "", "", false)

	config.Keybinds.Global.Quit = "q"
	if !IsQuitKey(composer, "ctrl+c") {
		t.Error("ctrl+c should always quit")
	}
	if !IsQuitKey(inbox, "q") {
		t.Error("q should quit from the inbox")
	}
	if IsQuitKey(composer, "q") {
		t.Error("q should not quit while typing in the composer")
	}
	if IsQuitKey(nil, "q") {
		t.Error("q should not quit on unknown screens")
	}

	config.Keybinds.Global.Quit = "ctrl+q"
	if !IsQuitKey(composer, "ctrl+q") {
		t.Error("non-printable quit key should quit from any screen")
	}
	if IsQuitKey(inbox, "q") {
		t.Error("q should not quit when quit is remapped")
	}
}
