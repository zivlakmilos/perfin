package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to perfin server",
	Run:   runLogin,
}

var (
	username string
	password string
)

func init() {
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "perfin account username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "perfin account password")

	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(initLoginModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "exited with error: %v\n", err)
		os.Exit(1)
	}
}

const (
	txtUsername = iota
	txtPassword
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
)

type loginModel struct {
	inputs  []textinput.Model
	focused int
}

func initLoginModel() loginModel {
	inputs := make([]textinput.Model, 2)
	inputs[txtUsername] = textinput.New()
	inputs[txtUsername].Placeholder = "username"
	inputs[txtUsername].CharLimit = 50
	inputs[txtUsername].Width = 50
	inputs[txtUsername].SetValue(username)

	inputs[txtPassword] = textinput.New()
	inputs[txtPassword].Placeholder = "password"
	inputs[txtPassword].CharLimit = 50
	inputs[txtPassword].Width = 50
	inputs[txtPassword].SetValue(password)
	inputs[txtPassword].EchoMode = textinput.EchoPassword

	focused := txtUsername
	if username == "" {
		inputs[txtUsername].Focus()
	} else {
		inputs[txtPassword].Focus()
		focused = txtPassword
	}

	return loginModel{
		inputs:  inputs,
		focused: focused,
	}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m loginModel) View() string {
	return fmt.Sprintf(
		` Login to perfin:

 %s
 %s

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Username"),
		m.inputs[txtUsername].View(),
		inputStyle.Width(30).Render("Password"),
		m.inputs[txtPassword].View(),
		continueStyle.Render("Login ->"),
	) + "\n"
}

func (m *loginModel) prevInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *loginModel) nextInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
