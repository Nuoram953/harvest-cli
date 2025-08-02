package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the date input
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// DateInputModel represents the date input component
type DateInputModel struct {
	textInput   textinput.Model
	err         error
	validDate   *time.Time
	submitted   bool
	title       string
	placeholder string
}

// NewDateInput creates a new date input component
func NewDateInput(title string) DateInputModel {
	ti := textinput.New()
	ti.Placeholder = "YYYY-MM-DD"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return DateInputModel{
		textInput:   ti,
		title:       title,
		placeholder: "YYYY-MM-DD",
	}
}

// SetPlaceholder sets a custom placeholder for the input
func (m *DateInputModel) SetPlaceholder(placeholder string) {
	m.placeholder = placeholder
	m.textInput.Placeholder = placeholder
}

// Init initializes the date input component
func (m DateInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles the date input updates
func (m DateInputModel) Update(msg tea.Msg) (DateInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Validate the date when Enter is pressed
			dateStr := strings.TrimSpace(m.textInput.Value())
			
			var parsedDate time.Time
			var err error
			
			if dateStr == "" {
				// Use today's date if no input provided
				parsedDate = time.Now()
			} else {
				// Try to parse the provided date
				parsedDate, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					m.err = fmt.Errorf("invalid date format. Use YYYY-MM-DD")
					return m, nil
				}
			}

			// Success!
			m.validDate = &parsedDate
			m.err = nil
			m.submitted = true
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			// Don't quit here, let the parent handle it
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the date input component
func (m DateInputModel) View() string {
	var b strings.Builder

	// Title
	if m.title != "" {
		b.WriteString(titleStyle.Render(m.title))
		b.WriteString("\n\n")
	}

	// If submitted successfully, show the result
	if m.submitted && m.validDate != nil {
		return b.String()
	}

	// Input prompt
	b.WriteString("Enter a date:")
	b.WriteString("\n\n")

	// Text input
	b.WriteString(inputStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	// Error message if any
	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Help text
	b.WriteString(helpStyle.Render(fmt.Sprintf("Format: %s (e.g., 2024-12-25)", m.placeholder)))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press Enter to submit (empty = today's date)"))

	return b.String()
}

// IsSubmitted returns true if a valid date has been submitted
func (m DateInputModel) IsSubmitted() bool {
	return m.submitted && m.validDate != nil
}

// GetDate returns the submitted date, or nil if not submitted
func (m DateInputModel) GetDate() *time.Time {
	if m.submitted && m.validDate != nil {
		return m.validDate
	}
	return nil
}

// GetDateISO8601 returns the submitted date in ISO 8601 format (YYYY-MM-DD)
func (m DateInputModel) GetDateISO8601() string {
	if m.submitted && m.validDate != nil {
		return m.validDate.Format("2006-01-02")
	}
	return ""
}

// Reset resets the input component to its initial state
func (m *DateInputModel) Reset() {
	m.textInput.SetValue("")
	m.err = nil
	m.validDate = nil
	m.submitted = false
	m.textInput.Focus()
}

// TextInputDate is a convenience function to get a date input from the user in ISO 8601 format
func TextInputDate(title string) (string, error) {
	model := NewDateInput(title)
	
	program := tea.NewProgram(dateInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if wrapper, ok := finalModel.(dateInputWrapper); ok {
		return wrapper.model.GetDateISO8601(), nil
	}
	
	return "", fmt.Errorf("failed to get date input")
}

// TextInputDateAsTime is a convenience function to get a date input as *time.Time
func TextInputDateAsTime(title string) (*time.Time, error) {
	model := NewDateInput(title)
	
	program := tea.NewProgram(dateInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	if wrapper, ok := finalModel.(dateInputWrapper); ok {
		return wrapper.model.GetDate(), nil
	}
	
	return nil, fmt.Errorf("failed to get date input")
}

// dateInputWrapper wraps the DateInputModel to handle quit behavior
type dateInputWrapper struct {
	model DateInputModel
}

func (w dateInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w dateInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return w, tea.Quit
		case tea.KeyEnter:
			newModel, cmd := w.model.Update(msg)
			w.model = newModel
			if w.model.IsSubmitted() {
				return w, tea.Quit
			}
			return w, cmd
		}
	}

	newModel, cmd := w.model.Update(msg)
	w.model = newModel
	return w, cmd
}

func (w dateInputWrapper) View() string {
	view := w.model.View()
	if !w.model.IsSubmitted() {
		view += "\n" + helpStyle.Render("Press Esc to cancel")
	}
	return view
}
