package core

import (
	"encoding/json"
	"strings"
	"testing"
)

const hongguoPlayerFixturePage = `<script type="application/ld+json">{
  "@context": "https://schema.org",
  "@type": "VideoObject",
  "url": "https://hongguoduanju.com/player/7693823855940357144",
  "contentUrl": "https://hongguoduanju.com/video/public/7693825728432852030/360p.mp4"
}</script>
<script type="application/ld+json">{
  "@context": "https://schema.org",
  "@type": "VideoObject",
  "url": "https://hongguoduanju.com/player/7693823855940357144",
  "contentUrl": "https://hongguoduanju.com/video/public/7693825728432852030/480p.mp4"
}</script>
<script type="application/ld+json">{
  "@context": "https://schema.org",
  "@type": "VideoObject",
  "url": "https://hongguoduanju.com/player/7693823855940357144",
  "contentUrl": "https://hongguoduanju.com/video/public/9999999999999999999/360p.mp4"
}</script>
<script type="application/ld+json">{"@type":"Product","name":"not a video"}</script>`

func TestHongguoWebPublicAddresses(t *testing.T) {
	got := hongguoWebPublicAddresses(hongguoPlayerFixturePage, "7693825728432852030")
	if len(got) != 2 {
		t.Fatalf("want 2 addresses bound to this vid, got %d: %v", len(got), got)
	}
	for _, address := range got {
		if !strings.Contains(address, "/video/public/7693825728432852030/") {
			t.Fatalf("address not bound to requested vid: %s", address)
		}
		if strings.Contains(address, "9999999999999999999") {
			t.Fatalf("leaked a foreign video address: %s", address)
		}
	}
}

func TestHongguoNormalizeWebURL(t *testing.T) {
	raw := `https://cdn.example/v.mp4?x=1\u0026y=2&amp;z=3&Range=0-1024`
	out := hongguoNormalizeWebURL(raw)
	if out != "https://cdn.example/v.mp4?x=1&y=2&z=3" {
		t.Fatalf("normalize wrong: %s", out)
	}
}

func TestHongguoPlayerPage(t *testing.T) {
	loader := map[string]any{
		"player_(series_id)/(vid)/layout": map[string]any{"vid": "800001"},
		"player_(series_id)/(vid)/page": map[string]any{
			"vid":        "800001",
			"series_id":  "700001",
			"isSuccess":  true,
			"publicVideoSeo": map[string]any{"available": true, "height": 360, "durationSeconds": 174},
			"seriesDetail":   map[string]any{"accessible_episode_cnt": 10, "episode_cnt": 68,
				"vid_list": []any{"800001", "800002"}},
		},
	}
	data := map[string]any{"loaderData": loader}
	page := hongguoPlayerPage(data)
	if page == nil {
		t.Fatal("player page not selected")
	}
	if mapString(page, "vid") != "800001" {
		t.Fatalf("picked wrong page entry: %v", page["vid"])
	}
	detail, _ := page["seriesDetail"].(map[string]any)
	if mapInt(detail, "accessible_episode_cnt") != 10 {
		t.Fatalf("accessible gate not read: %v", detail)
	}
	if idx := hongguoEpisodeIndex(anyList(detail["vid_list"]), "800002"); idx != 2 {
		t.Fatalf("episode index wrong: %d", idx)
	}
}

func TestHongguoEpisodeIndexNotFound(t *testing.T) {
	if got := hongguoEpisodeIndex([]any{"800001", "800002"}, "999999"); got != 0 {
		t.Fatalf("unknown vid must yield 0, got %d", got)
	}
}

func TestParseHongguoAppModel(t *testing.T) {
	encoded, _ := json.Marshal(map[string]any{
		"video_duration": 174,
		"video_list": []any{
			map[string]any{"main_url": "https://cdn.example/x.mp4", "video_meta": map[string]any{"codec_type": "bytevc2", "definition": "720p"},
				"encrypt_info": map[string]any{"spade_a": "00112233", "encryption_method": "cenc-aes-ctr"}},
		},
	})
	entry := map[string]any{"video_model": string(encoded)}
	model := parseHongguoAppModel(entry)
	if len(model) == 0 {
		t.Fatal("model not parsed from encoded string")
	}
	obj := map[string]any{"video_model": map[string]any{"video_list": []any{}}}
	if m := parseHongguoAppModel(obj); len(m) == 0 {
		t.Fatal("object-form video_model not parsed")
	}
	if m := parseHongguoAppModel(map[string]any{"other": 1}); m != nil {
		t.Fatalf("expected nil model, got %v", m)
	}
}
