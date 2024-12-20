package cmd

type Cursor struct {
	scrollOffset  int
	row           int
	col           int
	cursorVisible bool
}
