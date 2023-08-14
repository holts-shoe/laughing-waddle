package main

import (
	"fmt"
	"os"
	"time"

  "time_cli/db"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	stopwatch stopwatch.Model
	keymap    keymap
	help      help.Model
	quitting  bool
  capturingName bool
  capturingDescription bool
  taskName string
  taskDescription string
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m model) View() string {
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.
	if m.capturingName {
    return "Enter task name: " + m.taskName
  }

  if m.capturingDescription {
    return "Enter a description: " + m.taskDescription
  }

  s := m.stopwatch.View() + "\n"
	if !m.quitting {
		s = "Elapsed: " + s
		s += m.helpView()
	}
	return s
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
    if m.capturingName {
      switch msg.String() {
      case "enter" :
        m.capturingName = false
        m.capturingDescription = true
        return m, nil
      case "backspace":
        if len(m.taskName) > 0 {
          m.taskName = m.taskName[:len(m.taskName)-1]
        }
        return m, nil
      default:
        m.taskName += msg.String()
        return m, nil
      }
    }
    
     if m.capturingDescription {
      switch msg.String() {
      case "enter" :
        m.capturingDescription = false
        db.CreateTask(m.taskName, m.taskDescription, time.Now())
        return m, m.stopwatch.Reset()
      case "backspace":
        if len(m.taskDescription) > 0 {
          m.taskDescription = m.taskDescription[:len(m.taskDescription)-1]
        }
        return m, nil
      default:
        m.taskDescription += msg.String()
        return m, nil
      }
    }



		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
      m.capturingName = true
      m.taskName = ""
			return m, m.stopwatch.Reset()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):

      m.keymap.stop.SetEnabled(!m.stopwatch.Running())
			m.keymap.start.SetEnabled(m.stopwatch.Running())
			return m, m.stopwatch.Toggle()
		}
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func main() {
  db.InitDB()
	m := model{
		stopwatch: stopwatch.NewWithInterval(time.Second),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
