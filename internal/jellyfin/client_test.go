package jellyfin

import "testing"

func TestContainsEpisode(t *testing.T) {
	episodes := []EpisodeInfo{{ID: "a"}, {ID: "b"}}
	for id, want := range map[string]bool{"a": true, "b": true, "c": false, "": false} {
		if got := ContainsEpisode(episodes, id); got != want {
			t.Errorf("%q: got %v, want %v", id, got, want)
		}
	}
	if ContainsEpisode(nil, "a") {
		t.Error("an empty list contains nothing")
	}
}
