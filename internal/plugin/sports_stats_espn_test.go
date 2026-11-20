package plugin

import (
	"encoding/json"
	"testing"
)

func TestESPNSummaryKeepsClockForStreamCalibration(t *testing.T) {
	t.Parallel()
	var summary espnStatsSummary
	if err := json.Unmarshal([]byte(`{"header":{"id":"1","competitions":[{"status":{"period":2,"displayClock":"8:40","type":{"state":"in"}}}]}}`), &summary); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(espnGameStats(SportsEvent{}, summary))
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["period"] != float64(2) || payload["clock"] != "8:40" {
		t.Fatalf("game clock lost in stats response: period=%v clock=%v", payload["period"], payload["clock"])
	}
}

func TestESPNFootballDriveMomentsKeepPlayIdentityAndInjuryType(t *testing.T) {
	t.Parallel()
	var summary espnStatsSummary
	if err := json.Unmarshal([]byte(`{"header":{"id":"1","competitions":[{"status":{"type":{"state":"in"}}}]},"drives":{"current":{"plays":[{"id":"injury-1","text":"Injury timeout","type":{"text":"Injury Timeout"}},{"id":"score-1","text":"Touchdown","scoringPlay":true}]}}}`), &summary); err != nil {
		t.Fatal(err)
	}
	plays := espnGameStats(SportsEvent{LeagueID: "nfl"}, summary).Plays
	if len(plays) != 2 {
		t.Fatalf("football drive moments missing: %+v", plays)
	}
	encoded, err := json.Marshal(plays)
	if err != nil {
		t.Fatal(err)
	}
	var payload []map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload[0]["id"] != "score-1" || payload[0]["scoring"] != true || payload[1]["id"] != "injury-1" || payload[1]["injury"] != true {
		t.Fatalf("play metadata lost: %s", encoded)
	}
}

func TestESPNSummaryProbablesAcceptObjectStatistics(t *testing.T) {
	t.Parallel()
	var summary espnStatsSummary
	raw := `{"header":{"id":"1","competitions":[{"status":{"type":{"state":"pre"}},"competitors":[{"homeAway":"home","team":{"id":"10"},"probables":[{"athlete":{"displayName":"Gerrit Cole"},"record":"14-6","statistics":{"splits":{"categories":[]}}}]}]}]}}`
	if err := json.Unmarshal([]byte(raw), &summary); err != nil {
		t.Fatalf("summary probables with object statistics must decode: %v", err)
	}
	result := SportsGameStats{}
	applyESPNCompetitionDetail(&result, summary.Header.Competitions[0], "baseball/mlb")
	if len(result.Probables) != 1 || result.Probables[0].Summary != "14-6" {
		t.Fatalf("unexpected probables %+v", result.Probables)
	}
}

func TestESPNSummaryHighlightsKeepPlayableStreamAndWebPage(t *testing.T) {
	t.Parallel()
	var summary espnStatsSummary
	raw := `{"header":{"id":"1","competitions":[{"status":{"type":{"state":"post"}}}]},"videos":[
		{"headline":"Game Highlights","thumbnail":"https://a.espncdn.com/thumb.jpg","links":{"source":{"href":"https://media.video-cdn.espn.com/clip.mp4"},"web":{"href":"https://www.espn.com/video/clip?id=1"}}},
		{"headline":"Web only","links":{"web":{"href":"https://www.espn.com/video/clip?id=2"}}},
		{"headline":"Insecure","links":{"source":{"href":"http://media.video-cdn.espn.com/clip.mp4"}}}
	]}`
	if err := json.Unmarshal([]byte(raw), &summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	videos := espnGameStats(SportsEvent{}, summary).Videos
	if len(videos) != 2 {
		t.Fatalf("expected two https highlights, got %+v", videos)
	}
	if videos[0].Stream != "https://media.video-cdn.espn.com/clip.mp4" || videos[0].URL != "https://www.espn.com/video/clip?id=1" {
		t.Fatalf("highlight must carry the mp4 stream and the ESPN page, got %+v", videos[0])
	}
	if videos[1].Stream != "" || videos[1].URL != "https://www.espn.com/video/clip?id=2" {
		t.Fatalf("web-only highlight must not claim a stream, got %+v", videos[1])
	}
}
