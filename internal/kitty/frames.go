package kitty

import (
	"strings"
	"unicode/utf8"
)

func init() {
	maxW, maxH := 0, 0
	type art struct {
		lines []string
		width int
	}
	all := map[Mood][]art{}
	for mood, set := range frames {
		parsed := make([]art, len(set))
		for i, fr := range set {
			lines := strings.Split(fr, "\n")
			w := 0
			for _, ln := range lines {
				ln = strings.TrimRight(ln, " ")
				if n := utf8.RuneCountInString(ln); n > w {
					w = n
				}
			}
			if len(lines) > maxH {
				maxH = len(lines)
			}
			if w > maxW {
				maxW = w
			}
			parsed[i] = art{lines: lines, width: w}
		}
		all[mood] = parsed
	}
	for mood, set := range all {
		out := make([]string, len(set))
		for i, fr := range set {
			lines := fr.lines
			for len(lines) < maxH {
				lines = append(lines, "")
			}
			left := (maxW - fr.width) / 2
			for j, ln := range lines {
				ln = strings.TrimRight(ln, " ")
				padRight := maxW - left - utf8.RuneCountInString(ln)
				if padRight < 0 {
					padRight = 0
				}
				lines[j] = strings.Repeat(" ", left) + ln + strings.Repeat(" ", padRight)
			}
			out[i] = strings.Join(lines, "\n")
		}
		frames[mood] = out
	}
}

func Frame(mood Mood, i int) string {
	set, ok := frames[mood]
	if !ok || len(set) == 0 {
		set = frames[Present]
	}
	return set[i%len(set)]
}

func Width() int {
	return lineWidth(Frame(Idle, 0))
}

func lineWidth(s string) int {
	max := 0
	for _, ln := range strings.Split(s, "\n") {
		if n := utf8.RuneCountInString(ln); n > max {
			max = n
		}
	}
	return max
}

func hold(pose string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = pose
	}
	return out
}

func seq(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

var (
	rest = norm(`
         /\_/\
       (> • ω • <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	breathe = norm(`
         /\_/\
       (> • ω • <)
        /       \
       (   ∪ ∪   )
        \_____/
`)

	blink = norm(`
         /\_/\
       (> - ω - <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	halfBlink = norm(`
         /\_/\
       (> = ω = <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	ear = norm(`
         /‾_/\
       (> • ω • <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	lookLeft = norm(`
         /\_/\
       (>• ω • <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	lookRight = norm(`
         /\_/\
       (> • ω •<)
        /     \
       (  ∪ ∪  )
        \___/
`)

	tail = norm(`
         /\_/\
       (> • ω • <)
        /     \~
       (  ∪ ∪  )
        \___/
`)

	leafHat = norm(`
          *
         /\_/\
       (> • ω • <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	tilt = norm(`
           /\_/\
         (> • ω • <)
          /     \
         (  ∪ ∪  )
          \___/
`)

	peek = norm(`
           /\_/\
         (> • ω  •<)
          /     \
         (  ∪ ∪  )
          \___/
`)

	sleepy = norm(`
         /\_/\
       (> ˘ ω ˘ <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	sleepyBlink = norm(`
         /\_/\
       (> - ω - <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	nap = norm(`
         /\_/\
       (> - ω - <)
        /     \
       (  ∪ ∪  )
        \___/
`)

	napBreathe = norm(`
         /\_/\
       (> u ω u <)
        /       \
       (   ∪ ∪   )
        \_____/
`)

	stretch = norm(`
       /\_/\
     (> • ω • <)~~
      /       \
    _(         )_
     \_________/
`)

	stretchOut = norm(`
      /\_/\
    (> • ω • <)~~~~
     /         \
   _(           )_
    \___________/
`)

	yawn = norm(`
         /\_/\
       (> • ω • <)
        \  ▽  /
       /       \
      (   ∪ ∪   )
       \_____/
`)

	happy = norm(`
         /\_/\
       (> ^ ω ^ <)
        /     \
       (  ∪ ∪  )
        \___/
`)
)

var frames = map[Mood][]string{
	Idle: seq(
		hold(rest, 3),
		hold(breathe, 1),
		hold(rest, 2),
		hold(ear, 1),
		hold(rest, 2),
		hold(tail, 1),
		hold(rest, 2),
		hold(halfBlink, 1),
		hold(blink, 1),
		hold(rest, 2),
		hold(leafHat, 2),
		hold(rest, 2),
	),
	Settling: seq(
		hold(stretch, 2),
		hold(stretchOut, 2),
		hold(stretch, 1),
		hold(breathe, 1),
		hold(rest, 2),
		hold(happy, 2),
	),
	Present: seq(
		hold(rest, 3),
		hold(lookLeft, 2),
		hold(rest, 2),
		hold(breathe, 1),
		hold(rest, 2),
		hold(blink, 1),
		hold(rest, 2),
		hold(lookRight, 2),
		hold(rest, 2),
		hold(tail, 1),
		hold(rest, 2),
	),
	Paused: seq(
		hold(tilt, 3),
		hold(peek, 2),
		hold(tilt, 3),
		hold(ear, 1),
		hold(tilt, 2),
	),
	Napping: seq(
		hold(nap, 3),
		hold(napBreathe, 2),
		hold(nap, 3),
		hold(sleepyBlink, 1),
		hold(nap, 2),
	),
	Winding: seq(
		hold(sleepy, 3),
		hold(sleepyBlink, 1),
		hold(sleepy, 2),
		hold(breathe, 1),
		hold(sleepy, 3),
	),
	Done: seq(
		hold(yawn, 2),
		hold(stretchOut, 2),
		hold(happy, 2),
		hold(breathe, 1),
		hold(happy, 2),
		hold(tail, 1),
	),
}

func norm(raw string) string {
	raw = strings.Trim(raw, "\n")
	lines := strings.Split(raw, "\n")
	maxW := 0
	for _, ln := range lines {
		if w := utf8.RuneCountInString(ln); w > maxW {
			maxW = w
		}
	}
	for i, ln := range lines {
		pad := maxW - utf8.RuneCountInString(ln)
		if pad > 0 {
			lines[i] = ln + strings.Repeat(" ", pad)
		}
	}
	return strings.Join(lines, "\n")
}
