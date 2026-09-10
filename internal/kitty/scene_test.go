package kitty

import (
	"strings"
	"testing"
)

func TestSceneStableHeight(t *testing.T) {
	first := Compose(Idle, 0)
	hSky := strings.Count(first.Sky, "\n")
	hCat := strings.Count(first.Cat, "\n")
	hGround := strings.Count(first.Ground, "\n")
	for mood := Idle; mood <= Done; mood++ {
		for i := 0; i < 24; i++ {
			s := Compose(mood, i)
			if strings.Count(s.Sky, "\n") != hSky {
				t.Fatalf("sky height mood %d tick %d", mood, i)
			}
			if strings.Count(s.Cat, "\n") != hCat {
				t.Fatalf("cat height mood %d tick %d", mood, i)
			}
			if strings.Count(s.Ground, "\n") != hGround {
				t.Fatalf("ground height mood %d tick %d", mood, i)
			}
			if !strings.Contains(s.Cat, `/\_/\`) && !strings.Contains(s.Cat, `/‾_/\`) {
				t.Fatalf("missing ears mood %d tick %d\n%s", mood, i, s.Cat)
			}
		}
	}
}
