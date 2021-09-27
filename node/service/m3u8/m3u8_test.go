package m3u8

import "testing"

func TestNewTransit(t *testing.T) {
	transit, err := NewTransit("https://vod1.bdzybf1.com//20201015/bwcEwYbn/index.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	println(transit)
}
