package main

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type NavigateToCreateFileMsg struct{}
type NavigateToEditFileMsg struct{}
type NavigateToHelpMsg struct{}
type NavigateBack struct{}
type NavigateToOpenFile struct{
    filePath string
}

var (
	primaryColor   = lipgloss.Color("#FF6AD5")  // Hot pink
    secondaryColor = lipgloss.Color("#C774E8")  // Purple
    accentColor    = lipgloss.Color("#94D2F5")  // Cyan blue
    bgColor        = lipgloss.Color("#241734")  // Dark purple bg
    textColor      = lipgloss.Color("#E4E4E7")  // Light gray
)
var (
    // Title style
    titleStyle = lipgloss.NewStyle().
        Foreground(primaryColor).
        Bold(true).
        Padding(1, 0).
        MarginBottom(1)
    
    // Menu item style (unselected)
    itemStyle = lipgloss.NewStyle().
        Foreground(textColor).
        PaddingLeft(4)
    
    // Selected menu item
    selectedItemStyle = lipgloss.NewStyle().
        Foreground(accentColor).
        Bold(true).
        PaddingLeft(2).
        Background(secondaryColor)
    
    // Container for centering
    containerStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(secondaryColor).
        Padding(2, 4).
        Align(lipgloss.Center)
    
    // Footer hint
    footerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#666666")).
        Italic(true).
        MarginTop(1)
)



type mainMenuModel struct{
	cursor int
	choices []choice
	windowWidth  int
    windowHeight int

}

type choice struct{
	name string
	msg tea.Msg
}

func initialMenu()mainMenuModel{
	return mainMenuModel{
		choices: []choice{
			{name: "Edit file", msg: NavigateToEditFileMsg{}},
			{name: "Create file", msg: NavigateToCreateFileMsg{}},
			{name: "Help", msg: NavigateToHelpMsg{}},
		},
	}
}
func (m mainMenuModel)Init()tea.Cmd{
	return nil
}

func (m mainMenuModel)Update(msg tea.Msg)(tea.Model,tea.Cmd){
	switch msg := msg.(type){
		case tea.KeyMsg :
			switch msg.String(){
				case "q", "ctrl+c":
					return m, tea.Quit
				case "up", "k":
					if m.cursor > 0{
						m.cursor --
					}
				case "down", "j":
					if m.cursor < len(m.choices)-1{
						m.cursor ++
					}
				case "enter", " ":
					return m,func() tea.Msg{
						return m.choices[m.cursor].msg
					}
			}	
	}
	return m,nil
}

func (m mainMenuModel)View()string{
	title := titleStyle.Render("🎵 SYNTHWAVE EDITOR 🎵")

	var menuItems []string
    for i, choice := range m.choices {
        if m.cursor == i {
            menuItems = append(menuItems, selectedItemStyle.Render("▶ "+choice.name))
        } else {
            menuItems = append(menuItems, itemStyle.Render("  "+choice.name))
        }
    }
	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	footer := footerStyle.Render("Press q or Ctrl+C to exit")

	content := lipgloss.JoinVertical(
        lipgloss.Center,
        title,
        menu,
        footer,
    )
	styledContent := containerStyle.Render(content)
	return lipgloss.Place(
        80, 24,  // Width, Height (you'll get these from tea.WindowSizeMsg)
        lipgloss.Center, lipgloss.Center,
        styledContent,
    )


}

