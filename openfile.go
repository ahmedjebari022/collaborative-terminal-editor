package main

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type Mode string
const (
	NormalMode = "Normal"
	InsertMode = "Insert" 
)

var (
	cursorStyle = lipgloss.NewStyle().
		Background(secondaryColor).
		Foreground(textColor)
	
	editorBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2)
	
	statusBarStyle = lipgloss.NewStyle().
		Background(secondaryColor).
		Foreground(textColor).
		Padding(0, 1).
		Bold(true)
	
	modeNormalStyle = lipgloss.NewStyle().
		Background(accentColor).
		Foreground(bgColor).
		Padding(0, 1).
		Bold(true)
	
	modeInsertStyle = lipgloss.NewStyle().
		Background(primaryColor).
		Foreground(bgColor).
		Padding(0, 1).
		Bold(true)
	
	lineNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Width(4).
		Align(lipgloss.Right).
		MarginRight(1)
	
	currentLineNumberStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Width(4).
		Align(lipgloss.Right).
		MarginRight(1).
		Bold(true)
)

type fileModel struct{
	path string
	fileContent []string
	cp cursorPos
	vp viewPort
	m Mode
}
type cursorPos struct{
	x int
	y int 
}

type viewPort struct{
	ch int
	h int
}

func initFileModel(path string)tea.Model{
	fc, _ := readFile(path)
	if len(fc)== 0{
		fc = append(fc, "")
	}
	for i,l := range fc{
		fc[i] = replaceTabulation(l)
	}
	vp := viewPort{
		ch: 0,
		h: 15,
	}
	if len(fc)< 15{
		vp.h = len(fc)
	}
	return fileModel{
		path: path,
		fileContent: fc,
		cp : cursorPos{
			x: 0,
			y: 0,
		},
		vp : vp,
		m : NormalMode,
	}
}

func (fm fileModel) Init() tea.Cmd{
	return nil
}

func Min(a, b int)int{
	if a < b {
		return a
	}
	return b	
}


func (fm fileModel) View() string{
	// Status bar
	modeStyle := modeNormalStyle
	if fm.m == InsertMode {
		modeStyle = modeInsertStyle
	}
	
	modeIndicator := modeStyle.Render(string(fm.m))
	posInfo := statusBarStyle.Render(fmt.Sprintf("Ln %d, Col %d", fm.cp.y+1, fm.cp.x+1))
	fileInfo := statusBarStyle.Render(fmt.Sprintf("📄 %s", fm.path))
	
	statusBar := lipgloss.JoinHorizontal(
		lipgloss.Left,
		modeIndicator,
		" ",
		posInfo,
		" ",
		fileInfo,
	)
	
	// Editor content
	rl := len(fm.fileContent)-fm.vp.ch
	min := Min(rl, fm.vp.h)
	var editorLines []string
	
	for i := fm.vp.ch ; i < fm.vp.ch+min ; i++{
		lineNum := ""
		if fm.cp.y == i {
			lineNum = currentLineNumberStyle.Render(fmt.Sprintf("%d", i+1))
		} else {
			lineNum = lineNumberStyle.Render(fmt.Sprintf("%d", i+1))
		}
		
		if fm.cp.y == i{
			line := fm.fileContent[i]
			
			if len(line) == 0 {
				c := cursorStyle.Render(" ")
				editorLines = append(editorLines, lineNum + c)
			} else if fm.cp.x >= len(line) {
				editorLines = append(editorLines, lineNum + line + cursorStyle.Render(" "))
			} else {
				l := []rune(line)
				bc := string(l[:fm.cp.x])
				c := cursorStyle.Render(string(l[fm.cp.x]))
				ac := ""
				if fm.cp.x+1 < len(l) {
					ac = string(l[fm.cp.x+1:])
				}
				editorLines = append(editorLines, lineNum + bc + c + ac)
			}
		} else {
			editorLines = append(editorLines, lineNum + fm.fileContent[i])
		}
	}
	
	editorContent := lipgloss.JoinVertical(lipgloss.Left, editorLines...)
	styledEditor := editorBorderStyle.Render(editorContent)
	
	// Help text
	helpText := footerStyle.Render("Normal: hjkl=move i=insert ^w=save | Insert: esc=normal ^c=quit")
	
	// Combine all parts
	fullView := lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		styledEditor,
		helpText,
	)
	
	return fullView
}

func (fm fileModel) Update(msg tea.Msg)(tea.Model, tea.Cmd){
			
	switch  msg := msg.(type){
		case tea.KeyMsg:
			if fm.m == NormalMode{
				switch msg.String(){
					case "ctrl+c":
						return fm,tea.Quit
					case "h":
						fm.MoveLeft()
					case "j":
						fm.MoveDown()
					case "k":
						fm.MoveUp()
					case "l":
						fm.MoveRight()
					case "i" : 
						fm.m = InsertMode
					case "ctrl+w":{
						updateFile(fm.path,fm.fileContent)
					}
				}
			}else{
				switch msg.String(){
					case "esc":
						fm.m = NormalMode
					case "ctrl+c":
						return fm,tea.Quit
					case "backspace":
						fm.Delete()
					case "enter":
						fm.AddLine()
					default :
						fm.Write(msg.String())
					}

			}

			}
			return fm,nil
}


func (fm *fileModel) MoveDown() {
	if fm.cp.y < len(fm.fileContent) - 1 {
		if fm.cp.y == fm.vp.ch + fm.vp.h - 1{
			fm.vp.ch += 1
		}
		if fm.cp.x >= len(fm.fileContent[fm.cp.y + 1]){
			if len(fm.fileContent[fm.cp.y + 1]) == 0{
				fm.cp.x = 0
			}else{
				fm.cp.x = len(fm.fileContent[fm.cp.y + 1]) - 1
			}
		}
		fm.cp.y += 1
		
	}
}

func (fm *fileModel) MoveRight(){
	if fm.cp.x < len(fm.fileContent[fm.cp.y]) - 1{
		fm.cp.x += 1
	}else {
		fm.MoveDown()
		fm.cp.x = 0
	}
}


func (fm *fileModel) MoveLeft(){
		if fm.cp.x > 0 {
			fm.cp.x -= 1
		}else{
			fm.MoveUp()
			fm.cp.x = len(fm.fileContent[fm.cp.y]) - 1
		}
}

func (fm *fileModel) MoveUp(){
	if fm.cp.y > 0 {
		
		if fm.cp.y == fm.vp.ch {
			fm.vp.ch -= 1
		}
		if fm.cp.x >= len(fm.fileContent[fm.cp.y - 1]){
			if len(fm.fileContent[fm.cp.y - 1]) == 0{
				fm.cp.x = 0
			}else{
				fm.cp.x = len(fm.fileContent[fm.cp.y - 1]) - 1
			}
		}
		fm.cp.y -= 1
	}

}

func (fm *fileModel)Write(input string) {
	if fm.fileContent[fm.cp.y] == ""{
		fm.fileContent[fm.cp.y] = input
		fm.cp.x += 1
		return
	}

	l := []rune(fm.fileContent[fm.cp.y])
	bi := l[:fm.cp.x]
	ai := l[fm.cp.x:]
	NewLine := string(bi) + input + string(ai)
	fm.fileContent[fm.cp.y] = NewLine
	fm.cp.x += 1	
}

func (fm *fileModel)Delete(){
	l := []rune(fm.fileContent[fm.cp.y])
	if fm.cp.x == 0 && fm.cp.y == 0{
		return
	}
	if len(l) == 0{
		fm.fileContent = slices.Delete(fm.fileContent,fm.cp.y,fm.cp.y+1)
		fm.MoveLeft()
		return
	}
	if fm.cp.x == 0{
		if fm.fileContent[fm.cp.y - 1] == ""{
			fm.fileContent[fm.cp.y - 1] += fm.fileContent[fm.cp.y]
			fm.fileContent = slices.Delete(fm.fileContent,fm.cp.y,fm.cp.y+1)
			fm.MoveUp()
			fm.cp.x = 0
			return
		} 
		ol := fm.fileContent[fm.cp.y]
		fm.fileContent = slices.Delete(fm.fileContent,fm.cp.y,fm.cp.y+1)
		fm.MoveLeft()
		fm.fileContent[fm.cp.y] += ol
		return
	}
	if fm.cp.x == len(fm.fileContent[fm.cp.y]){
		fm.fileContent = fm.fileContent[:len(fm.fileContent)]
		fm.MoveLeft()
		return
			
	}
	bc := l[:fm.cp.x-1]
	ac := l[fm.cp.x:]
	NewLine := string(bc) + string(ac)
	fm.fileContent[fm.cp.y] = NewLine
	fm.cp.x -= 1
}

func (fm *fileModel)AddLine(){
	newl := []string{""}
	if fm.fileContent[fm.cp.y] != "" && fm.cp.x != len(fm.fileContent[fm.cp.y]) {
		l := []byte(fm.fileContent[fm.cp.y])
		a := l[fm.cp.x:]
		b := l[:fm.cp.x]
		fm.fileContent[fm.cp.y] = string(b)
		newl[0] = string(a)
	}
	fm.fileContent = append(fm.fileContent[:fm.cp.y+1],append(newl,fm.fileContent[fm.cp.y+1:]...)...)
	fm.vp.h = Min(15,len(fm.fileContent))
	fm.MoveDown()
	fm.cp.x = 0
}


func replaceTabulation (line string)string{
	spacePerTab := 4
	replacement := strings.Repeat(" ", spacePerTab)
	return strings.ReplaceAll(line,"\t",replacement)
}