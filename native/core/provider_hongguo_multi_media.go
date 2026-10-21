package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func hongguoModelBizParam() map[string]any {
	return map[string]any{
		"detail_page_version":       0,
		"device_level":              3,
		"disable_digg_stat":         false,
		"need_all_video_definition": true,
		"need_mp4_align":            false,
		"use_os_player":             false,
		"use_server_dns":            false,
		"video_platform":            1024,
	}
}

func hongguoModelExtraQuery() url.Values {
	return url.Values{
		"host_abi":                {"arm64-v8a"},
		"dragon_device_type":      {"phone"},
		"pv_player":               {"71332"},
		"compliance_status":       {"0"},
		"need_personal_recommend": {"1"},
		"player_so_load":          {"1"},
		"is_android_pad_screen":   {"0"},
	}
}

const hongguoMultiVideoModelPath = "/novel/player/multi_video_model/preload/v1/"

func (downloader *Downloader) resolveHongguoMultiVideoMedia(ctx context.Context, videoID string) (providerMedia, error) {
	if !hongguoNumericID.MatchString(videoID) {
		return providerMedia{}, errors.New("红果视频 ID 无效")
	}
	payload := map[string]any{
		"biz_param":          hongguoModelBizParam(),
		"dr_scene":           "preload",
		"mixed_video_id_map": map[string]any{"1004": []string{videoID}},
	}
	result, err := downloader.hongguoAppRequest(ctx, http.MethodPost, hongguoMultiVideoModelPath, hongguoModelExtraQuery(), payload)
	if err != nil {
		return providerMedia{}, err
	}
	model, err := hongguoModelFromMultiResponse(result, videoID)
	if err != nil {
		return providerMedia{}, err
	}
	fallback, err := hongguoFallbackAPI(model)
	if err != nil {
		return providerMedia{}, err
	}
	videoData, err := downloader.fetchHongguoFallbackVideoData(ctx, fallback)
	if err != nil {
		return providerMedia{}, err
	}
	return hongguoMediaFromFallbackData(model, videoData)
}

func hongguoModelFromMultiResponse(result map[string]any, videoID string) (map[string]any, error) {
	data := nestedMap(result, "data")
	entry, _ := data[videoID].(map[string]any)
	if entry == nil {
		return nil, errors.New("红果取流响应缺少该集的条目")
	}
	raw := entry["video_model"]
	switch value := raw.(type) {
	case map[string]any:
		return value, nil
	case string:
		decoder := json.NewDecoder(strings.NewReader(value))
		decoder.UseNumber()
		var model map[string]any
		if err := decoder.Decode(&model); err != nil || model == nil {
			return nil, errors.New("红果 App 播放信息格式异常")
		}
		return model, nil
	default:
		return nil, errors.New("红果取流响应缺少 video_model")
	}
}

func hongguoFallbackAPI(model map[string]any) (string, error) {
	raw := model["fallback_api"]
	switch value := raw.(type) {
	case string:
		if strings.TrimSpace(value) == "" {
			return "", errors.New("红果取流响应缺少 fallback_api")
		}
		return value, nil
	case map[string]any:
		address := mapString(value, "fallback_api")
		if strings.TrimSpace(address) == "" {
			return "", errors.New("红果取流响应缺少 fallback_api")
		}
		return address, nil
	default:
		return "", errors.New("红果取流响应缺少 fallback_api")
	}
}

func (downloader *Downloader) fetchHongguoFallbackVideoData(ctx context.Context, address string) (map[string]any, error) {
	parsed, err := url.Parse(address)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" {
		return nil, errors.New("红果备用播放地址无效")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", hongguoAppUserAgent)
	request.Header.Set("Accept", "application/json")
	response, err := downloader.doCatalogRequestWithTimeout(request, 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, readErr := io.ReadAll(io.LimitReader(response.Body, providerMaxBodyBytes+1))
	if readErr != nil {
		return nil, readErr
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("红果备用播放接口 HTTP %d", response.StatusCode)
	}
	if len(content) == 0 || len(content) > providerMaxBodyBytes {
		return nil, errors.New("红果备用播放接口未返回有效数据")
	}
	var result map[string]any
	decoder := json.NewDecoder(strings.NewReader(string(content)))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil || result == nil {
		return nil, errors.New("红果备用播放接口返回格式异常")
	}
	videoData := nestedMap(nestedMap(result, "video_info"), "data")
	if len(videoData) == 0 {
		return nil, errors.New("红果备用播放接口缺少 video_info")
	}
	return videoData, nil
}

func hongguoMediaFromFallbackData(model, videoData map[string]any) (providerMedia, error) {
	keySeed, err := hongguoSpadeKeySeed(videoData)
	if err != nil {
		return providerMedia{}, err
	}
	rows := nestedMap(videoData, "video_list")
	if len(rows) == 0 {
		return providerMedia{}, errors.New("红果备用播放接口无可播档位")
	}
	duration, _ := strconv.ParseFloat(mapString(model, "video_duration", "duration"), 64)

	type candidate struct {
		media providerMedia
		score int
		order int
	}
	var candidates []candidate
	var keyErr error
	order := 0
	for key, raw := range rows {
		item, _ := raw.(map[string]any)
		if item == nil {
			continue
		}
		order++
		encrypted := mapString(item, "main_url")
		if encrypted == "" {
			continue
		}
		address, decErr := decryptHongguoSpadeURL(encrypted, keySeed)
		if decErr != nil || !isProviderHTTPMediaURL(address) {
			continue
		}
		meta := nestedMap(item, "video_meta")
		codec := strings.ToLower(mapString(meta, "codec_type"))
		if codec == "" {
			codec = strings.ToLower(mapString(item, "codec_type"))
		}
		if codec == "bytevc2" {
			continue
		}
		height, _ := strconv.Atoi(mapString(item, "vheight"))
		if height == 0 {
			height, _ = strconv.Atoi(mapString(meta, "vheight"))
		}
		media := providerMedia{
			URL:      address,
			Referer:  "",
			Duration: time.Duration(duration * float64(time.Second)),
			Quality:  height,
		}
		if spade := mapString(item, "spade_a"); spade != "" {
			key, keyDecodeErr := hongguoContentKey(spade)
			if keyDecodeErr != nil {
				keyErr = keyDecodeErr
				continue
			}
			media.CENCKey = key
		}
		score := height * 10
		if codec == "h264" || codec == "avc" || codec == "avc1" {
			score++
		}
		candidates = append(candidates, candidate{media: media, score: score, order: order})
		_ = key
	}
	if len(candidates) == 0 {
		if keyErr != nil {
			return providerMedia{}, fmt.Errorf("红果 App 媒体密钥不可用: %w", keyErr)
		}
		return providerMedia{}, errors.New("红果 App 未返回兼容的媒体，已跳过不支持的编码")
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].order < candidates[j].order
	})
	selected := candidates[0].media
	for _, item := range candidates {
		selected.Variants = append(selected.Variants, item.media)
	}
	return selected, nil
}
