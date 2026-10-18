package models

import (
	"database/sql"
	"testing"
)

func TestDownloadsAllowed(t *testing.T) {
	cases := []struct {
		share, server, want bool
	}{
		{share: true, server: true, want: true},
		{share: false, server: true, want: false}, // the sharer switched them off
		{share: true, server: false, want: false}, // the server switched them off for every link
		{share: false, server: false, want: false},
	}
	for _, c := range cases {
		s := &Share{AllowDownload: c.share}
		if got := s.DownloadsAllowed(c.server); got != c.want {
			t.Errorf("share=%v server=%v: got %v, want %v", c.share, c.server, got, c.want)
		}
	}
}

func TestPlayLimitReached(t *testing.T) {
	limit := func(n int64) sql.NullInt64 { return sql.NullInt64{Int64: n, Valid: true} }
	cases := []struct {
		name  string
		share Share
		want  bool
	}{
		{"no limit", Share{TotalPlays: 100}, false},
		{"plays left", Share{TotalPlays: 2, MaxTotalPlays: limit(3)}, false},
		{"last play used", Share{TotalPlays: 3, MaxTotalPlays: limit(3)}, true},
		{"over the limit", Share{TotalPlays: 4, MaxTotalPlays: limit(3)}, true},
	}
	for _, c := range cases {
		if got := c.share.PlayLimitReached(); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestWithholdDetails(t *testing.T) {
	five := int64(5)
	info := SharePublicInfo{
		Title: "Film", PosterURL: "/p", BackdropURL: "/b", LogoURL: "/l", ItemType: "Movie",
		RequiresPassword: true, MaxTotalPlays: &five, TotalPlays: 2, AllowDownload: true,
		Overview: "plot", Tagline: "line", RuntimeSeconds: 60, Year: 2020, OfficialRating: "PG",
		CommunityRating: 7, CriticRating: 80, Genres: []string{"Drama"}, Studios: []string{"S"},
		Directors: []string{"D"}, Actors: []ActorInfo{{Name: "A"}}, VideoQuality: &VideoQualityInfo{},
		AudioTracks: []AudioTrack{{Index: 1}}, SubtitleTracks: []SubtitleTrack{{Index: 2}},
	}
	info.WithholdDetails()

	if !info.Locked {
		t.Error("withheld info must say it is locked")
	}
	// What the link points at, and its terms, stay visible before the password.
	if info.Title != "Film" || info.PosterURL == "" || info.BackdropURL == "" || info.LogoURL == "" ||
		info.MaxTotalPlays == nil || info.TotalPlays != 2 || !info.RequiresPassword || !info.AllowDownload {
		t.Errorf("title, artwork and terms must survive: %+v", info)
	}
	// Everything else waits for the password.
	if info.Overview != "" || info.Tagline != "" || info.RuntimeSeconds != 0 || info.Year != 0 ||
		info.OfficialRating != "" || info.CommunityRating != 0 || info.CriticRating != 0 ||
		info.Genres != nil || info.Studios != nil || info.Directors != nil || info.Actors != nil ||
		info.VideoQuality != nil || info.AudioTracks != nil || info.SubtitleTracks != nil {
		t.Errorf("details leaked before the password: %+v", info)
	}
}

func TestHasEpisodes(t *testing.T) {
	for typ, want := range map[string]bool{"Season": true, "Series": true, "Movie": false, "Episode": false} {
		if got := (&Share{ItemType: typ}).HasEpisodes(); got != want {
			t.Errorf("%s: got %v, want %v", typ, got, want)
		}
	}
}
