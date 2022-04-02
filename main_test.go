package main

import (
	"testing"
)

type Case struct {
	name string
	mode int
	want string
}

func TestMatchFile(t *testing.T) {
	cases := []Case{
		{"Alex.Rider.S02E07.Assassin.2160p.AMZN.WEB-DL.DDP5.1.HDR.H.265-playWEB.mkv", 1, "S02E07"},
		{"The.Fall.2013.S01.E05.Complete.BluRay.720p.x264.AC3-CMCT.mkv", 2, "S01.E05"},
		{"Gold.Panning.2022.E08.1080p.WEB-DL.H264.AAC-OurTV.mp4", 3, "E08"},
		{"[梦蓝字幕组]Crayonshinchan 蜡笔小新[1119][2022.02.19][AVC][1080P][GB_JP][MP4].mp4", 5, "2022.02.19"},
	}

	for _, item := range cases {
		find := MatchFile(item.mode, item.name)
		give := ""
		if len(find) > 0 {
			give = find[0]
		}
		if give != item.want {
			t.Errorf("MatchFile(%s) give: %s, want: %s\n", item.name, give, item.want)
		}
	}
}
