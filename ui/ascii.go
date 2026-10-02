package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/VinylStage/animux-like/pet"
)

var (
	colorHappy  = lipgloss.Color("205")
	colorSad    = lipgloss.Color("240")
	colorSick   = lipgloss.Color("196")
	colorHungry = lipgloss.Color("214")
	colorInfo   = lipgloss.Color("86")
)

var catAscii = `
 /\___/\
( o   o )
(  =^=  )
(        )
(         )
(          )))))))))))
`

var dogAscii = `
  __      _
o'')}____//
 'c-___,-'
`

var penguinAscii = `
   _
 ('v')
//-=-\\
(\_=_/)
 ^^ ^^
`

func getAscii(species pet.SpeciesType) string {
	switch species {
	case pet.Cat:
		return catAscii
	case pet.Dog:
		return dogAscii
	case pet.Penguin:
		return penguinAscii
	default:
		return catAscii
	}
}

func RenderResponse(state *pet.State, message string) string {
	ascii := getAscii(state.Species)

	style := lipgloss.NewStyle().Foreground(colorHappy)
	if state.IsSick {
		style = style.Foreground(colorSick)
	} else if state.Hunger < 30 {
		style = style.Foreground(colorHungry)
	} else if state.Happiness < 30 {
		style = style.Foreground(colorSad)
	}

	styledAscii := style.Render(strings.TrimSpace(ascii))

	infoStyle := lipgloss.NewStyle().Foreground(colorInfo).Bold(true)
	info := infoStyle.Render(fmt.Sprintf("%s the %s [Hunger: %d | Happy: %d | Clean: %d]", 
		state.Name, state.Species, state.Hunger, state.Happiness, state.Cleanliness))
	
	msgStyle := lipgloss.NewStyle().Italic(true).MarginTop(1).MarginBottom(1)
	styledMsg := msgStyle.Render(message)

	return lipgloss.JoinVertical(lipgloss.Center, info, styledAscii, styledMsg)
}

func RenderError(err error) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("Error: " + err.Error())
}
