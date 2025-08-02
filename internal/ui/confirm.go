package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the confirm component
var (
	confirmTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	confirmMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	confirmYesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	confirmNoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	confirmHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// ConfirmModel represents the confirm component
type ConfirmModel struct {
	title     string
	message   string
	confirmed bool
	answered  bool
	cancelled bool
}

// NewConfirm creates a new confirm component
func NewConfirm(title, message string) ConfirmModel {
	return ConfirmModel{
		title:   title,
		message: message,
	}
}

// Init initializes the confirm component
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Update handles the confirm component updates
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "y", "yes":
			m.confirmed = true
			m.answered = true
			return m, nil
		case "n", "no":
			m.confirmed = false
			m.answered = true
			return m, nil
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, nil
		}
	}
	return m, nil
}

// View renders the confirm component
func (m ConfirmModel) View() string {
	var b strings.Builder

	// If answered, show the result
	if m.answered {
		if m.confirmed {
			b.WriteString(confirmYesStyle.Render("✓ Confirmed"))
		} else {
			b.WriteString(confirmNoStyle.Render("✗ Cancelled"))
		}
		return b.String()
	}

	// Message
	if m.message != "" {
		b.WriteString(confirmMessageStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	// Options
	b.WriteString(confirmYesStyle.Render("[Y]es"))
	b.WriteString(" / ")
	b.WriteString(confirmNoStyle.Render("[N]o"))
	b.WriteString("\n\n")

	// Help text
	b.WriteString(confirmHelpStyle.Render("Press Y for Yes, N for No, or Esc to cancel"))

	return b.String()
}

// IsAnswered returns true if the user has answered (y/n)
func (m ConfirmModel) IsAnswered() bool {
	return m.answered
}

// IsConfirmed returns true if the user confirmed (pressed y/yes)
func (m ConfirmModel) IsConfirmed() bool {
	return m.answered && m.confirmed
}

// IsCancelled returns true if the user cancelled (pressed esc)
func (m ConfirmModel) IsCancelled() bool {
	return m.cancelled
}

// Reset resets the confirm component to its initial state
func (m *ConfirmModel) Reset() {
	m.confirmed = false
	m.answered = false
	m.cancelled = false
}

// Confirm is a convenience function to get a confirmation from the user
func Confirm(title, message string) (bool, error) {
	model := NewConfirm(title, message)
	
	program := tea.NewProgram(confirmWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return false, err
	}

	if wrapper, ok := finalModel.(confirmWrapper); ok {
		if wrapper.model.IsCancelled() {
			return false, nil // User cancelled
		}
		return wrapper.model.IsConfirmed(), nil
	}
	
	return false, nil
}

// confirmWrapper wraps the ConfirmModel to handle quit behavior
type confirmWrapper struct {
	model ConfirmModel
}

func (w confirmWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w confirmWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			newModel, _ := w.model.Update(msg)
			w.model = newModel
			return w, tea.Quit
		default:
			newModel, cmd := w.model.Update(msg)
			w.model = newModel
			if w.model.IsAnswered() {
				return w, tea.Quit
			}
			return w, cmd
		}
	}

	newModel, cmd := w.model.Update(msg)
	w.model = newModel
	return w, cmd
}

func (w confirmWrapper) View() string {
	return w.model.View()
}
