package kanban

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/rpcarvs/faz/internal/model"
)

// TestLightTerminalRendersReadableCardAndModal guards the light-theme fix:
// adaptive colors must resolve to a light card background with dark text, and
// modals must carry an explicit foreground so default dark terminal text is
// never drawn on the dark modal background.
func TestLightTerminalRendersReadableCardAndModal(t *testing.T) {
	oldProfile := lipgloss.ColorProfile()
	oldDark := lipgloss.HasDarkBackground()
	lipgloss.SetColorProfile(termenv.ANSI256)
	lipgloss.SetHasDarkBackground(false)
	defer func() {
		lipgloss.SetColorProfile(oldProfile)
		lipgloss.SetHasDarkBackground(oldDark)
	}()

	now := time.Now()
	issue := model.Issue{ID: "p-1", Title: "Short title", Type: "task", Priority: 1, Status: "open", CreatedAt: now, UpdatedAt: now}

	card := NewModel(stubService{}).renderCard(issue, false, 32)
	assertMarkerStyleContains(t, card, "Short title", "48;5;254") // light card background
	assertMarkerStyleContains(t, card, "Short title", "38;5;235") // dark, readable title text

	m := NewModel(stubService{})
	m.applyWindowSize(80, 24)
	help := m.renderHelp()
	if !strings.Contains(help, "48;5;254") || !strings.Contains(help, "38;5;235") {
		t.Fatalf("light help modal missing readable bg/fg codes: %q", help)
	}
}
