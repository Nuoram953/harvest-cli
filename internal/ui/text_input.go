package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the text input component
var (
	textInputTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	textInputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	textInputErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))

	textInputSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46"))

	textInputHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))
)

// TextInputOptions configures the text input behavior
type TextInputOptions struct {
	Title        string
	Prompt       string
	Placeholder  string
	Required     bool
	Password     bool
	CharLimit    int
	Width        int
	ValidateFunc func(string) error
	DefaultValue string
}

// TextInputModel represents the text input component
type TextInputModel struct {
	textInput textinput.Model
	options   TextInputOptions
	err       error
	submitted bool
	cancelled bool
}

// NewTextInput creates a new text input component
func NewTextInput(options TextInputOptions) TextInputModel {
	ti := textinput.New()

	// Set defaults if not provided
	if options.Prompt == "" {
		options.Prompt = "Enter text:"
	}
	if options.Width == 0 {
		options.Width = 50
	}
	if options.CharLimit == 0 {
		options.CharLimit = 200
	}

	// Configure the text input
	ti.Placeholder = options.Placeholder
	ti.Width = options.Width
	ti.CharLimit = options.CharLimit
	ti.Focus()

	// Set password mode if requested
	if options.Password {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	// Set default value if provided
	if options.DefaultValue != "" {
		ti.SetValue(options.DefaultValue)
	}

	return TextInputModel{
		textInput: ti,
		options:   options,
	}
}

// Init initializes the text input component
func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles the text input updates
func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Get the input value
			value := strings.TrimSpace(m.textInput.Value())

			// Check if required and empty
			if m.options.Required && value == "" {
				m.err = fmt.Errorf("this field is required")
				return m, nil
			}

			// Run custom validation if provided
			if m.options.ValidateFunc != nil {
				if err := m.options.ValidateFunc(value); err != nil {
					m.err = err
					return m, nil
				}
			}

			// Success!
			m.err = nil
			m.submitted = true
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil
	}

	// Update the text input
	m.textInput, cmd = m.textInput.Update(msg)

	// Clear error when user starts typing again
	if m.err != nil && msg != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyRunes {
			m.err = nil
		}
	}

	return m, cmd
}

// View renders the text input component
func (m TextInputModel) View() string {
	var b strings.Builder

	// Title
	if m.options.Title != "" {
		b.WriteString(textInputTitleStyle.Render(m.options.Title))
		b.WriteString("\n\n")
	}

	// If submitted successfully, show the result
	if m.submitted {
		b.WriteString(textInputSuccessStyle.Render("✓ Input submitted successfully!"))
		if !m.options.Password {
			b.WriteString(fmt.Sprintf("\nYou entered: %s", m.textInput.Value()))
		}
		return b.String()
	}

	// If cancelled, show cancellation message
	if m.cancelled {
		b.WriteString(textInputErrorStyle.Render("✗ Input cancelled"))
		return b.String()
	}

	// Prompt
	prompt := m.options.Prompt
	if m.options.Required {
		prompt += " *"
	}
	b.WriteString(prompt)
	b.WriteString("\n\n")

	// Text input
	b.WriteString(textInputStyle.Render(m.textInput.View()))
	b.WriteString("\n\n")

	// Error message if any
	if m.err != nil {
		b.WriteString(textInputErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Help text
	helpLines := []string{
		"Press Enter to submit",
		"Press Esc to cancel",
	}

	if m.options.CharLimit > 0 {
		remaining := m.options.CharLimit - len(m.textInput.Value())
		helpLines = append(helpLines, fmt.Sprintf("Characters remaining: %d", remaining))
	}

	for _, line := range helpLines {
		b.WriteString(textInputHelpStyle.Render(line))
		b.WriteString("\n")
	}

	return b.String()
}

// IsSubmitted returns true if input has been submitted successfully
func (m TextInputModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if input was cancelled
func (m TextInputModel) IsCancelled() bool {
	return m.cancelled
}

// GetValue returns the submitted value, or empty string if not submitted
func (m TextInputModel) GetValue() string {
	if m.submitted {
		return strings.TrimSpace(m.textInput.Value())
	}
	return ""
}

// Reset resets the input component to its initial state
func (m *TextInputModel) Reset() {
	m.textInput.SetValue(m.options.DefaultValue)
	m.err = nil
	m.submitted = false
	m.cancelled = false
	m.textInput.Focus()
}

// SetValue sets the current input value
func (m *TextInputModel) SetValue(value string) {
	m.textInput.SetValue(value)
}

// Focus focuses the text input
func (m *TextInputModel) Focus() {
	m.textInput.Focus()
}

// Blur removes focus from the text input
func (m *TextInputModel) Blur() {
	m.textInput.Blur()
}

// TextInput is a convenience function to get text input from the user
func TextInput(options TextInputOptions) (string, error) {
	model := NewTextInput(options)

	program := tea.NewProgram(textInputWrapper{model})
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if wrapper, ok := finalModel.(textInputWrapper); ok {
		if wrapper.model.IsCancelled() {
			return "", fmt.Errorf("input cancelled")
		}
		return wrapper.model.GetValue(), nil
	}

	return "", fmt.Errorf("failed to get text input")
}

// SimpleTextInput prompts for basic text input
func SimpleTextInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:  title,
		Prompt: prompt,
	})
}

// RequiredTextInput prompts for required text input
func RequiredTextInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:    title,
		Prompt:   prompt,
		Required: true,
	})
}

// PasswordInput prompts for password input
func PasswordInput(title, prompt string) (string, error) {
	return TextInput(TextInputOptions{
		Title:    title,
		Prompt:   prompt,
		Password: true,
		Required: true,
	})
}

// ValidatedTextInput prompts for text input with custom validation
func ValidatedTextInput(title, prompt string, validateFunc func(string) error) (string, error) {
	return TextInput(TextInputOptions{
		Title:        title,
		Prompt:       prompt,
		Required:     true,
		ValidateFunc: validateFunc,
	})
}

// textInputWrapper wraps the TextInputModel to handle quit behavior
type textInputWrapper struct {
	model TextInputModel
}

func (w textInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w textInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			newModel, _ := w.model.Update(msg)
			w.model = newModel
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

func (w textInputWrapper) View() string {
	return w.model.View()
}
