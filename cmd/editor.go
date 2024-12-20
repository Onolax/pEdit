package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type Editor struct {
	app         *tview.Application
	txtArea     *tview.TextView
	lineNumbers *tview.TextView
	layout      *tview.Flex
	cursor      *Cursor
	blob        *Blob
}

func InitDisplay(app *tview.Application, blob *Blob) *Editor {
	Area := tview.NewTextView()
	Area.SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(false).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Text Editor")

	lineNum := tview.NewTextView()
	lineNum.SetDynamicColors(true).SetScrollable(true).SetBorder(true)

	cursor := &Cursor{scrollOffset: 0, row: 0, col: 0, cursorVisible: true}
	layout := tview.NewFlex().AddItem(lineNum, 6, 1, false).
		AddItem(Area, 0, 1, true)
	return &Editor{app, Area, lineNum, layout, cursor, blob}
}

func (ed *Editor) Setter(filetxt []string) string {
	var txt strings.Builder
	for i, line := range filetxt {
		if i == ed.cursor.row {
			for j, char := range []rune(line) {
				if j == ed.cursor.col && ed.cursor.cursorVisible {
					txt.WriteString("[black:yellow:]")
					txt.WriteRune(char)
					txt.WriteString("[-:-:-]")
				} else {
					txt.WriteRune(char)
				}
			}
			if ed.cursor.col == len(line) && ed.cursor.cursorVisible {
				txt.WriteString("[black:yellow:] [-:-:-]")
			}
		} else {
			txt.WriteString(line)
		}
		txt.WriteString("\n")
	}
	return txt.String()
}

func (ed *Editor) Render() {
	filetxt := strings.Split(ed.blob.content, "\n")
	ed.SetScreen(filetxt)
	ed.txtArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			ed.app.Stop()
			return nil
		case tcell.KeyEnter:
			line := filetxt[ed.cursor.row]
			leftPart := line[:ed.cursor.col]
			rightPart := line[ed.cursor.col:]
			filetxt[ed.cursor.row] = leftPart
			filetxt = append(filetxt[:ed.cursor.row+1], append([]string{rightPart}, filetxt[ed.cursor.row+1:]...)...)

			ed.cursor.row++
			ed.cursor.col = 0
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)

			return nil
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if ed.cursor.col > 0 {
				line := []rune(filetxt[ed.cursor.row])
				filetxt[ed.cursor.row] = string(append(line[:ed.cursor.col-1], line[ed.cursor.col:]...))
				ed.cursor.col--
			} else if ed.cursor.row > 0 {
				prevLine := filetxt[ed.cursor.row-1]
				filetxt[ed.cursor.row-1] = prevLine + filetxt[ed.cursor.row]
				filetxt = append(filetxt[:ed.cursor.row], filetxt[ed.cursor.row+1:]...)
				ed.cursor.row--
				ed.cursor.col = len(prevLine)
			}
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)

			return nil
		case tcell.KeyCtrlS:
			content := strings.Join(filetxt, "\n")
			ed.blob.Save(content)
			return nil
		case tcell.KeyRight:
			if ed.cursor.col < len(filetxt[ed.cursor.row]) {
				ed.cursor.col++
			} else if ed.cursor.row < len(filetxt)-1 {
				ed.cursor.row++
				ed.cursor.col = 0
			}
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)

			return nil

		case tcell.KeyLeft:
			if ed.cursor.col > 0 {
				ed.cursor.col--
			} else if ed.cursor.row > 0 {
				ed.cursor.row--
				ed.cursor.col = len(filetxt[ed.cursor.row])
			}
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)

			return nil
		case tcell.KeyDown:
			if ed.cursor.row < len(filetxt)-1 {
				ed.cursor.row++
				ed.cursor.col = min(len(filetxt[ed.cursor.row]), ed.cursor.col)
			}
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)

			return nil
		case tcell.KeyUp:
			if ed.cursor.row > 0 {
				ed.cursor.row--
				ed.cursor.col = min(len(filetxt[ed.cursor.row]), ed.cursor.col)
			}
			ed.cursor.cursorVisible = true
			ed.SetScreen(filetxt)
			return nil
		default:
			if event.Rune() != 0 {
				line := []rune(filetxt[ed.cursor.row])
				filetxt[ed.cursor.row] = string(append(line[:ed.cursor.col], append([]rune{event.Rune()}, line[ed.cursor.col:]...)...))
				ed.cursor.col++
				ed.cursor.cursorVisible = true
				ed.SetScreen(filetxt)

				return nil
			}
		}
		return event

	})
	// go function that runs always to blink the cursor
	go func() {
		for {
			time.Sleep(400 * time.Millisecond) // half a second interval
			ed.cursor.cursorVisible = !ed.cursor.cursorVisible
			ed.app.QueueUpdateDraw(func() {
				ed.SetScreen(filetxt)

			})
		}
	}()
}

func (ed *Editor) setLineNumber(filetxt []string) {
	var count = len(filetxt)
	var numbers string
	for i := 0; i < count; i++ {
		numbers += fmt.Sprintln(i)
	}
	ed.lineNumbers.SetText(numbers)
}

func (ed *Editor) SetScreen(filetxt []string) {
	ed.Scroll()
	ed.setLineNumber(filetxt)
	ed.txtArea.SetText(ed.Setter(filetxt))
}

func (ed *Editor) Scroll() {
	if ed.cursor.row < ed.cursor.scrollOffset {
		ed.cursor.scrollOffset = ed.cursor.row
	} else if ed.cursor.row > ed.cursor.scrollOffset+20 {
		ed.cursor.scrollOffset = ed.cursor.row - 19
	}
	ed.txtArea.ScrollTo(ed.cursor.scrollOffset, 0)
	ed.lineNumbers.ScrollTo(ed.cursor.scrollOffset, 0)
}
