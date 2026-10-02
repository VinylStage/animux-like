package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/pet"
	"github.com/VinylStage/animux-like/ui"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Continuously show your pet (Ctrl+C or 'q' to quit)",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
		}
	},
}

type model struct {
	state *pet.State
	err   error
	msg   string
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchState() (*pet.State, error) {
	return daemon.SendRequest("/status", nil)
}

func initialModel() model {
	state, err := fetchState()
	return model{state: state, err: err, msg: "Watching..."}
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "f":
			state, err := daemon.SendRequest("/feed", nil)
			if err == nil {
				m.state = state
				m.msg = "Yum!"
			}
		case "p":
			state, err := daemon.SendRequest("/play", nil)
			if err == nil {
				m.state = state
				m.msg = "Yay!"
			}
		case "c":
			state, err := daemon.SendRequest("/clean", nil)
			if err == nil {
				m.state = state
				m.msg = "Cleaned!"
			}
		}

	case tickMsg:
		state, err := fetchState()
		m.state = state
		m.err = err
		// Every tick, reset the message to normal if it's currently a reaction
		if m.state != nil && m.msg != "Watching..." {
			// Actually let's just keep the reaction for a bit or just revert to watching
			// Simple approach: don't reset, but we could check state to say "hungry"
			if m.state.IsSick {
				m.msg = "I feel sick..."
			} else if m.state.Hunger < 30 {
				m.msg = "So hungry..."
			} else if m.state.Happiness < 30 {
				m.msg = "So bored..."
			} else if m.state.Cleanliness < 30 {
				m.msg = "It smells..."
			} else {
				m.msg = "Watching..."
			}
		}
		return m, doTick()
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return ui.RenderError(m.err) + "\n\nPress 'q' to quit."
	}
	if m.state == nil {
		return "Loading...\n\nPress 'q' to quit."
	}
	
	rendered := ui.RenderResponse(m.state, m.msg)
	help := "\n\n[q] Quit | [f] Feed | [p] Play | [c] Clean"
	return rendered + help
}
