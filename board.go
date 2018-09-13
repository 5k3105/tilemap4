package main

/// ref: github.com/mandykoh/scrubble
type Board struct {
	OffsetX, OffsetY int
	Rows             int
	Columns          int
	Positions        []Position
}

type Position struct {
	PTile *Tile
}

type Coord struct {
	Row    int
	Column int
}

type Tile struct {
	Type     string
	Row, Col int
	X, Y     int
}

func (b *Board) LocateTile(x, y int) (*Position, int, int, int, int) {
	tx, ty := fit_grid(x, y)
	iscale := int(scale)
	r, c := ty/iscale, tx/iscale
	ox, oy := b.OffsetX, b.OffsetY
	r, c = r-oy, c-ox
	pos := b.Position(Coord{r, c})
	return pos, r, c, tx, ty
}

func (b *Board) LocateTile2(x, y int) *Position {
	tx, ty := fit_grid(x, y)
	iscale := int(scale)
	r, c := ty/iscale, tx/iscale
	ox, oy := b.OffsetX, b.OffsetY
	r, c = r-oy, c-ox
	pos := b.Position(Coord{r, c})
	return pos
}

func (b *Board) AddTile(x, y int) {
	pos, r, c, tx, ty := b.LocateTile(x, y)
	if pos != nil {
		pos.PTile = &Tile{selected_tile, r, c, tx, ty}
	}
}

func (b *Board) DeleteTile(x, y int) {
	pos := b.LocateTile2(x, y)
	if pos != nil {
		pos.PTile = nil
	}
}

func (b *Board) AddSelectors() {
	iscale := int(scale)
	start_x, start_y := b.OffsetX*iscale, b.OffsetY*iscale
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Columns; j++ {
			b.Positions[i*b.Columns+j].PTile = &Tile{tilefiles[j+i], j, i, start_x + (j * iscale), start_y + (i * iscale)}
		}
	}
}

func (b *Board) SelectTile(x, y int) {
	pos := b.LocateTile2(x, y)
	if pos != nil {
		selected_tile = pos.PTile.Type
	}
}

func fit_grid(x, y int) (tx, ty int) {
	nxt := offset * 2
	tx, ty = (x/nxt)*nxt, (y/nxt)*nxt
	return
}

func (b *Board) Position(c Coord) *Position {
	if c.Row < 0 || c.Row >= b.Rows || c.Column < 0 || c.Column >= b.Columns {
		return nil
	}
	return &b.Positions[c.Row*b.Columns+c.Column]
}
