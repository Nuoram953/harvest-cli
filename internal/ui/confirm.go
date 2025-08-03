package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type ConfirmModel struct {
	title     string
	message   string
	confirmed bool
	answered  bool
	cancelled bool
}

func NewConfirm(title, message string) ConfirmModel {
	return ConfirmModel{
		title:   title,
		message: message,
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

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

func (m ConfirmModel) View() string {
	var b strings.Builder

	if m.answered {
		if m.confirmed {
			b.WriteString(confirmYesStyle.Render("✓ Confirmed"))
		} else {
			b.WriteString(confirmNoStyle.Render("✗ Cancelled"))
		}
		return b.String()
	}

	if m.message != "" {
		b.WriteString("\n\n")
		b.WriteString(confirmMessageStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	b.WriteString(confirmYesStyle.Render("[Y]es"))
	b.WriteString(" / ")
	b.WriteString(confirmNoStyle.Render("[N]o"))
	b.WriteString("\n\n")

	return b.String()
}

func (m ConfirmModel) IsAnswered() bool {
	return m.answered
}

func (m ConfirmModel) IsConfirmed() bool {
	return m.answered && m.confirmed
}

func (m ConfirmModel) IsCancelled() bool {
	return m.cancelled
}

func (m *ConfirmModel) Reset() {
	m.confirmed = false
	m.answered = false
	m.cancelled = false
}

func Confirm(title, message string) (bool, error) {
	model := NewConfirm(title, message)
	
	program := tea.NewProgram(confirmWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return false, err
	}

	if wrapper, ok := finalModel.(confirmWrapper); ok {
		if wrapper.model.IsCancelled() {
		}
		return wrapper.model.IsConfirmed(), nil
	}
	
	return false, nil
}

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
