package kitty

import "strings"

type Scene struct {
	Sky    string
	Cat    string
	Ground string
}

type dot struct {
	x, y int
	ch   rune
}

func Compose(mood Mood, tick int) Scene {
	cat := Frame(mood, tick)
	w := lineWidth(cat)
	if w < 18 {
		w = 18
		cat = padBlock(cat, w)
	}
	return Scene{
		Sky:    paint(w, 2, skyDots(mood, tick, w)),
		Cat:    cat,
		Ground: ground(mood, tick, w),
	}
}

func skyDots(mood Mood, tick, w int) []dot {
	var dots []dot
	switch mood {
	case Napping:
		dots = append(dots, risingZ(tick, w)...)
	case Done:
		dots = append(dots, sparkle(tick, w)...)
		dots = append(dots, cloud(tick, w)...)
	case Winding:
		dots = append(dots, fireflies(tick, w, 2)...)
		dots = append(dots, risingZ(tick/2, w)...)
	case Settling:
		dots = append(dots, dust(tick, w)...)
	case Paused:
		dots = append(dots, fireflies(tick/2, w, 1)...)
	default:
		dots = append(dots, cloud(tick, w)...)
		dots = append(dots, fireflies(tick, w, 2)...)
		if mood == Idle && tick%18 < 7 {
			dots = append(dots, driftingLeaf(tick, w))
		}
	}
	return dots
}

func ground(mood Mood, tick, w int) string {
	soot := paint(w, 1, sootDots(mood, tick, w))
	grass := sway(tick, w)
	return soot + "\n" + grass
}

func cloud(tick, w int) []dot {
	x := mod(tick/3, w)
	return []dot{
		{x: x, y: 0, ch: '\''},
		{x: mod(x+w/2, w), y: 1, ch: '\''},
	}
}

func fireflies(tick, w, n int) []dot {
	chars := []rune{'.', '*', '.', '.'}
	var dots []dot
	for i := 0; i < n; i++ {
		if (tick+i*3)%4 == 0 {
			continue
		}
		x := mod(tick+i*9, w)
		y := (tick/3 + i) % 2
		dots = append(dots, dot{x: x, y: y, ch: chars[(tick+i)%len(chars)]})
	}
	return dots
}

func driftingLeaf(tick, w int) dot {
	x := mod(tick, w)
	y := (tick / 3) % 2
	return dot{x: x, y: y, ch: '*'}
}

func risingZ(tick, w int) []dot {
	x := w/2 + 4
	if x >= w {
		x = w - 3
	}
	phase := tick % 6
	var dots []dot
	switch phase {
	case 0, 1:
		dots = append(dots, dot{x: x, y: 1, ch: 'z'})
	case 2, 3:
		dots = append(dots, dot{x: x + 1, y: 0, ch: 'z'})
		dots = append(dots, dot{x: x, y: 1, ch: 'z'})
	default:
		dots = append(dots, dot{x: x + 2, y: 0, ch: 'Z'})
		dots = append(dots, dot{x: x, y: 1, ch: '.'})
	}
	return dots
}

func sparkle(tick, w int) []dot {
	return []dot{
		{x: mod(tick*2, w), y: 0, ch: '*'},
		{x: mod(tick*2+w/2, w), y: 1, ch: '·'},
		{x: mod(3+tick, w), y: 0, ch: '.'},
		{x: mod(w/3+tick/2, w), y: 1, ch: '*'},
	}
}

func dust(tick, w int) []dot {
	return []dot{
		{x: mod(2+tick, w), y: 1, ch: '.'},
		{x: mod(5+tick*2, w), y: 0, ch: '·'},
		{x: mod(w-4-tick, w), y: 1, ch: '.'},
	}
}

func sootDots(mood Mood, tick, w int) []dot {
	if mood == Done {
		return nil
	}
	path := []int{0, 1, 2, 2, 1, 0, 0, 1}
	bounce := path[tick%len(path)]
	left := 1 + bounce
	right := w - 4 - bounce
	if right < 0 {
		right = 0
	}
	ch := '.'
	if mood == Napping || mood == Winding {
		ch = '.'
	} else if tick%8 < 2 {
		ch = '*'
	}
	return []dot{
		{x: left, y: 0, ch: ch},
		{x: right, y: 0, ch: '.'},
	}
}

func sway(tick, w int) string {
	waves := []string{
		"~  ~~   ~  ~~  ~",
		" ~~  ~  ~~   ~ ~",
		"~ ~~   ~~  ~  ~",
		"  ~  ~~  ~ ~~  ~",
	}
	src := waves[tick%len(waves)]
	var b strings.Builder
	for utf8Len(b.String()) < w {
		b.WriteString(src)
		b.WriteByte(' ')
	}
	s := b.String()
	return string([]rune(s)[:w])
}

func paint(w, h int, dots []dot) string {
	cells := make([][]rune, h)
	for y := range cells {
		cells[y] = []rune(strings.Repeat(" ", w))
	}
	for _, d := range dots {
		if d.y >= 0 && d.y < h && d.x >= 0 && d.x < w {
			cells[d.y][d.x] = d.ch
		}
	}
	lines := make([]string, h)
	for y := range cells {
		lines[y] = string(cells[y])
	}
	return strings.Join(lines, "\n")
}

func padBlock(s string, w int) string {
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		n := utf8Len(ln)
		if n < w {
			lines[i] = ln + strings.Repeat(" ", w-n)
		}
	}
	return strings.Join(lines, "\n")
}

func utf8Len(s string) int {
	return len([]rune(s))
}

func mod(n, m int) int {
	if m <= 0 {
		return 0
	}
	n %= m
	if n < 0 {
		n += m
	}
	return n
}
