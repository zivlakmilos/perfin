package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	inputStyle    = lipgloss.NewStyle().Foreground(primaryColor)
	continueStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	errorStyle    = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	successStyle  = lipgloss.NewStyle().Foreground(successColor).Bold(true)
)

type loginStatus struct {
	code    int
	message string
}

type loginModel struct {
	inputs      []textinput.Model
	focused     int
	loginStatus loginStatus
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
				username = m.inputs[txtUsername].Value()
				password = m.inputs[txtPassword].Value()
				return m, handleLogin()
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
	case loginStatus:
		{
			m.loginStatus = msg
			if msg.code == 0 {
				return m, tea.Quit
			}
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m loginModel) View() string {
	var statusMsg string
	if m.loginStatus.message != "" {
		if m.loginStatus.code == 0 {
			statusMsg = successStyle.Render(fmt.Sprintf("login success: %s", m.loginStatus.message))
		} else {
			statusMsg = errorStyle.Render(fmt.Sprintf("error: %s", m.loginStatus.message))
		}
	}

	return fmt.Sprintf(
		` Login to perfin:

 %s
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
		statusMsg,
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

func handleLogin() tea.Cmd {
	return func() tea.Msg {
		body, err := json.Marshal(map[string]any{
			"username": username,
			"password": password,
		})
		if err != nil {
			return loginStatus{
				code:    -1,
				message: err.Error(),
			}
		}
		res, err := http.Post(getUrl(config.ApiBaseUrl, "/auth/login"), "application/json", bytes.NewReader(body))
		if err != nil {
			return loginStatus{
				code:    -1,
				message: err.Error(),
			}
		}
		defer func() { _ = res.Body.Close() }()

		if res.StatusCode != http.StatusOK {
			data := make(map[string]any)
			resJson, err2 := io.ReadAll(res.Body)
			if err2 != nil {
				return loginStatus{
					code:    -1,
					message: res.Status,
				}
			}
			err2 = json.Unmarshal(resJson, &data)
			if err2 != nil {
				return loginStatus{
					code:    -1,
					message: res.Status,
				}
			}

			return loginStatus{
				code:    -1,
				message: fmt.Sprintf("%s", data["error"]),
			}
		}

		data := make(map[string]any)
		resJson, err := io.ReadAll(res.Body)
		if err != nil {
			return loginStatus{
				code:    -1,
				message: err.Error(),
			}
		}
		err = json.Unmarshal(resJson, &data)
		if err != nil {
			return loginStatus{
				code:    -1,
				message: err.Error(),
			}
		}

		switch token := data["token"].(type) {
		case string:
			config.Token = token
			saveConfig()
			return loginStatus{
				code:    0,
				message: fmt.Sprintf("%s", data["token"]),
			}
		}

		return loginStatus{
			code:    -1,
			message: "unknown",
		}
	}
}
