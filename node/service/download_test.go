package service

import "testing"

func TestD(t *testing.T) {
	download, err := Download("https://ts1.lslkkyj.com/20200617/aAUCQ5Hf/1000kb/hls/BaaIxt2d.ts")
	println(err, len(download))
}
