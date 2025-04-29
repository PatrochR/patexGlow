package main

//TODO: Searching world in textarea

//FIXME: Text Highlight

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var linesCount int
var showSearchTable bool = false

type sessionState uint

const (
	codeView sessionState = iota
	tableView
	searchTableView
	inputView
)

var (
	modelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))
	modelFoucseStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("240")).
				Inherit(modelStyle)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))
	inputFoucseStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("240")).
				Inherit(modelStyle)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Center, lipgloss.Center)
	tableFoucseStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("69")).
				Inherit(tableStyle)

	// highlightStyle = lipgloss.NewStyle().
	// 			Bold(true).
	// 			Foreground(lipgloss.Color("#7D56F4"))

	// stringHighlightStyle = lipgloss.NewStyle().
	// 	Italic(true).
	// 	Bold(false).
	// 	Foreground(lipgloss.Color("#27ae60")).
	// 	Inherit(highlightStyle)

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))
)

type errMsg error

type model struct {
	state       sessionState
	textarea    textarea.Model
	table       table.Model
	searchTable table.Model
	textInput   textinput.Model
	err         error
}

var source_file string
var searchingLine []int
var save_message string

const hint = "(ctrl+c to quit) (ctrl+s to save)"

var file_is_saved bool
var backup_file string

func main() {

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func search(text string, key string) []int {
	lines := strings.Split(text, "\n")
	var result []int
	for i, s := range lines {
		if strings.Contains(s, key) {
			result = append(result, i)
		}
	}
	return result
}

func generate_status_bar(width int, status, format, filename, help string) string {
	width += 6
	w := lipgloss.Width

	statusKey := statusStyle.Render(status)
	encoding := encodingStyle.Render(format)
	fishCake := fishCakeStyle.Render(filename)
	statusVal := statusText.
		Width(width - w(statusKey) - w(encoding) - w(fishCake)).
		Render(help)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	)
	return statusBarStyle.Width(width).Render(bar)
}

func read_file(textarea *textarea.Model, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		source_file = filename
		return
	}
	var value string
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value += line
		value += "\n"
	}
	textarea.SetValue(value + "\n")

}

func get_all_dir() []table.Row {
	var value []table.Row
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for i, e := range entries {
		value = append(value, table.Row{})
		value[i] = append(value[i], e.Name())
		if e.IsDir() {
			value[i] = append(value[i], "Dir")
		} else {
			value[i] = append(value[i], "File")
		}
	}
	return value
}

func initialModel() model {

	input := textinput.New()
	input.Placeholder = "Search"
	input.CharLimit = 156

	searchColums := []table.Column{
		{Title: "Line", Width: 15},
	}

	st := table.New(
		table.WithColumns(searchColums),
		table.WithFocused(false),
	)

	row := get_all_dir()
	colums := []table.Column{
		{Title: "Type", Width: 15},
		{Title: "Name", Width: 15},
	}
	t := table.New(
		table.WithColumns(colums),
		table.WithRows(row),
		table.WithFocused(false),
		//table.WithHeight(),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	st.SetStyles(s)
	ti := textarea.New()
	ti.Focus()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	if len(os.Args) > 1 {
		source_file = os.Args[1]
		read_file(&ti, source_file)
	} else {
		source_file = "out.txt"
		//FIXME:Test For Debug
		// source_file = "alo.go"
		// read_file(&ti, source_file)
	}

	return model{
		textarea:    ti,
		err:         nil,
		table:       t,
		searchTable: st,
		textInput:   input,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlS:
			saveFile(m.textarea.Value())
		case tea.KeyEnter:
			if m.state == tableView {
				if m.table.SelectedRow()[1] != "Dir" {
					m.textarea.Reset()
					read_file(&m.textarea, m.table.SelectedRow()[0])
					source_file = m.table.SelectedRow()[0]
				}
			} else if m.state == inputView {
				searchingLine = search(m.textarea.Value(), m.textInput.Value())
				m.state = searchTableView
				m.textInput.Blur()
				m.searchTable.Focus()
				showSearchTable = true
				if len(searchingLine) != 0 {
					var rows []table.Row
					for i, s := range searchingLine {
						rows = append(rows, table.Row{})
						rows[i] = append(rows[i], strconv.Itoa(s+1))
					}
					m.searchTable.SetRows(rows)
				}
			} else if m.state == searchTableView {
				showSearchTable = false
				m.state = codeView
				m.searchTable.Blur()
				m.textarea.Focus()
				m.table.Blur()
				m.textInput.Blur()
				linesCount := len(strings.Split(m.textarea.Value(), "\n"))
				go_up(&m, linesCount)
				v, _ := strconv.Atoi(m.searchTable.SelectedRow()[0])
				go_down(&m, v-1)
				msg.Type = tea.KeyCtrlA
			}

		case tea.KeyCtrlF:
			m.state = inputView
			m.searchTable.Blur()
			m.textarea.Blur()
			m.table.Blur()
			m.textInput.Focus()

		case tea.KeyTab:
			{
				if showSearchTable {
					if m.state == codeView {
						m.state = searchTableView
						m.textarea.Blur()
						m.searchTable.Focus()
					} else {
						m.state = codeView
						m.searchTable.Blur()
						m.textarea.Focus()
					}
				} else if m.state == codeView {
					m.state = tableView
					m.textarea.Blur()
					m.table.Focus()
				} else {
					m.state = codeView
					m.table.Blur()
					m.textarea.Focus()
				}

			}
		default:
			check_file_saved(m.textarea.Value())
		}
		switch m.state {
		case searchTableView:
			m.searchTable, cmd = m.searchTable.Update(msg)
			cmds = append(cmds, cmd)
		case codeView:
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		case inputView:
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - (msg.Width / 4))
		m.textarea.SetHeight(msg.Height - 8)
		m.table.SetWidth((msg.Width / 4) - 5)
		m.table.SetHeight(msg.Height - 5)
		m.searchTable.SetWidth((msg.Width / 4) - 5)
		m.searchTable.SetHeight(msg.Height - 5)
		m.textInput.Width = (msg.Width - 20)

	case errMsg:
		m.err = msg
		return m, nil
	}

	if file_is_saved {
		save_message = "File is saved"
	} else {
		save_message = ""
	}

	return m, tea.Batch(cmds...)
}

func check_file_saved(value string) {
	if value == backup_file {
		file_is_saved = true
	} else {
		file_is_saved = false
	}
}

func saveFile(value string) {
	if value == "" {
		return
	}
	err := os.WriteFile(source_file, []byte(value), 0644)
	if err != nil {
		panic(err)
	}
	file_is_saved = true
	backup_file = value
}

func (m model) View() string {

	var s string
	// var coloredSyntax []string
	// GoHighlight := [25]string{ "break" , "default" , "func" , "interface" , "select" , "case" ,  "defer" , "go", "map" , "struct","chan" , "else" , "goto" , "package" ,"switch" , "const" ,"fallthrough" , "if" , "range" , "type" , "continue" , "for" ,  "import" , "return" , "var"}
	// highlight := m.textarea.View()
	// for _, s := range GoHighlight{
	// 	coloredSyntax =	append(coloredSyntax, highlightStyle.Render(s))
	// }

	// characters := strings.Split(highlight , "")
	// var isdoubleQuotation bool
	// var isBacktick bool
	// var isBackSlash bool

	// for i ,s := range characters{
	// 	if s == "\\" {
	// 		isBackSlash = true
	// 		isdoubleQuotation = false
	// 		isBacktick = false
	// 	}
	// 	if s == `"` {
	// 		if isdoubleQuotation{
	// 			isdoubleQuotation = false
	// 			isBacktick = false
	// 		}else if isBackSlash{
	// 			isBackSlash = false
	// 		}else{
	// 			isdoubleQuotation = true
	// 		}
	// 		continue
	// 	}else if s == "`" {
	// 		if isBacktick{
	// 			isBacktick = false
	// 			isdoubleQuotation = false
	// 		}else{
	// 			isBacktick = true
	// 		}
	// 		continue
	// 	}
	// 	if isdoubleQuotation {
	// 		characters[i] = stringHighlightStyle.Render(s)
	// 	}else if isBacktick {
	// 		characters[i] = stringHighlightStyle.Render(s)
	// 	}
	// }
	// highlight = strings.Join(characters , "")

	// for i, s := range coloredSyntax{
	// 	highlight = strings.ReplaceAll(highlight , GoHighlight[i] , s)
	// }

	// stringsInView := strings.Split(highlight, "\"")

	// for i,s := range stringsInView  {
	// 	if !IsDigit(i){
	// 		if strings.LastIndex(s,"`") != len(s) - 1 {
	// 			if strings.Index(s ,"`") != 0 {
	// 				if strings.Index(s , "\\") != len(s) {
	// 					stringsInView[i] = stringHighlightStyle.Render(s)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// highlight = strings.Join(stringsInView , `"`)

	//Check it later
	//row := strconv.Itoa(m.textarea.LineInfo().ColumnOffset)

	format := strings.ToUpper(strings.Split(source_file, ".")[1])
	if showSearchTable {
		if m.state == codeView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableStyle.Render(m.searchTable.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputFoucseStyle.Render(m.textInput.View()), s)
		} else if m.state == searchTableView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelFoucseStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableFoucseStyle.Render(m.searchTable.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputFoucseStyle.Render(m.textInput.View()), s)
		} else {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelFoucseStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableStyle.Render(m.searchTable.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputStyle.Render(m.textInput.View()), s)
		}
	} else {
		if m.state == codeView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableStyle.Render(m.table.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputFoucseStyle.Render(m.textInput.View()), s)
		} else if m.state == tableView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelFoucseStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableFoucseStyle.Render(m.table.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputFoucseStyle.Render(m.textInput.View()), s)
		} else {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelFoucseStyle.Render(fmt.Sprintf("%s\n\n%s", m.textarea.View(), "\n"+generate_status_bar(m.textarea.Width(), save_message, format, source_file, hint))), tableStyle.Render(m.table.View()))
			s = lipgloss.JoinVertical(lipgloss.Center, inputStyle.Render(m.textInput.View()), s)
		}
	}

	return s

}

func IsDigit(n int) bool {
	return n%2 == 0
}

func go_up(m *model, count int) {
	for i := 0; i < count; i++ {
		m.textarea.CursorUp()
	}
}

func go_down(m *model, count int) {
	for i := 0; i < count; i++ {
		m.textarea.CursorDown()
	}
}
