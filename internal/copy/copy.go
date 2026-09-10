package copy

import (
	"math/rand/v2"

	"github.com/euphoricair7/zenkitty/internal/kitty"
)

type Picker struct {
	last map[kitty.Mood]string
}

func New() *Picker {
	return &Picker{last: map[kitty.Mood]string{}}
}

func (p *Picker) Line(mood kitty.Mood) string {
	pool := pools[mood]
	if len(pool) == 0 {
		return ""
	}
	choice := pool[rand.IntN(len(pool))]
	if len(pool) > 1 {
		for i := 0; i < 4 && choice == p.last[mood]; i++ {
			choice = pool[rand.IntN(len(pool))]
		}
	}
	p.last[mood] = choice
	return choice
}

func Leave() string {
	return leave[rand.IntN(len(leave))]
}

var pools = map[kitty.Mood][]string{
	kitty.Idle: {
		"whenever you're ready",
		"no rush",
		"the wind is soft today",
		"i'll be here",
	},
	kitty.Settling: {
		"ok. let's sit",
		"getting comfy",
		"here we go",
	},
	kitty.Present: {
		"",
		"",
		"",
		"",
		"",
		"",
		"mm",
		"still here",
		"the air is quiet",
	},
	kitty.Paused: {
		"taking a breath",
		"i'll keep the spot warm",
		"whenever you're back",
	},
	kitty.Napping: {
		"still here",
		"i'll nap too",
		"spot's warm",
	},
	kitty.Winding: {
		"almost there",
		"mm. winding down",
		"a little more",
	},
	kitty.Done: {
		"that was a good sit",
		"nice stretch",
		"that was a nice one",
	},
}

var leave = []string{
	"see you",
	"spot'll be here",
	"until next time",
}
