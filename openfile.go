package main

import (
	"fmt"
	

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


var cursorStyle = lipgloss.NewStyle().
	Background(secondaryColor)

type fileModel struct{
	path string
	fileContent []string
	cp cursorPos
}
type cursorPos struct{
	x int
	y int 
}

func initFileModel(path string)tea.Model{
	fc, _ := readFile(path)
	return fileModel{
		fileContent: fc,
		cp : cursorPos{
			x: 0,
			y: 0,
		},
	}
}

func (fm fileModel) Init() tea.Cmd{
	return nil
}


func (fm fileModel) View() string{
	s := ""
	for i, l := range fm.fileContent{
		if fm.cp.y == i{
			r := []rune(l)
			before := r[:fm.cp.x]
			after := r[fm.cp.x + 1:]
			c := r[fm.cp.x]		
			cursorStyle.Render(string(c))
			s += fmt.Sprintf("%s%s%s",before,c,after)
		}	
		s += fmt.Sprintf("%s\n",l)
	}
	return s

}

func (fm fileModel) Update(msg tea.Msg)(tea.Model, tea.Cmd){
			switch  msg := msg.(type){
				case tea.KeyMsg:
					switch msg.String(){
					case "h":
					case "j":
						fm.MoveDown()
					case "k":
						fm.MoveUp()
					case "l":
					}

			}
			return fm,nil
}


func (fm fileModel) MoveDown() {

	if fm.cp.y < len(fm.fileContent){
		fm.cp.y += 1	
	}
}


func (fm fileModel) MoveUp(){
	if fm.cp.y > 0 {
		fm.cp.y -= 1
	}

}