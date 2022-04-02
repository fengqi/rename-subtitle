package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	m1 *regexp.Regexp
	m2 *regexp.Regexp
	m3 *regexp.Regexp
	m4 *regexp.Regexp
	m5 *regexp.Regexp
)

func init() {
	m1, _ = regexp.Compile("(?i)S[0-9]+E[0-9]+")
	m2, _ = regexp.Compile("(?i)S[0-9]+\\.E[0-9]+")
	m3, _ = regexp.Compile("(?i)E[0-9]+")
	m4, _ = regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
	m5, _ = regexp.Compile("[0-9]{4}\\.[0-9]{2}\\.[0-9]{2}")
}

func main() {
	mode := "1: S10E12, 2: S01.E12, 3: E08, 4: 2022-03-12, 5: 2022.02.12"

	subExt := flag.String("sext", "srt", "字幕后缀")
	subMode := flag.Int("smod", 1, "字幕匹配方式 "+mode)

	videoExt := flag.String("vext", "mkv", "视频后缀")
	videoMode := flag.Int("vmod", 1, "视频匹配方式 "+mode)

	try := flag.Bool("try", true, "测试模式, 不实际执行")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	subList := listSub(pwd, *subExt, *subMode)
	videoList := listVideo(pwd, *videoExt, *videoMode)
	rename(*try, *subExt, *videoExt, subList, videoList)
}

func rename(try bool, subExt, videoExt string, subList map[string]string, videoList map[string]string) {
	if len(videoList) == 0 {
		log.Println("video list empty")
		return
	}

	if len(subList) == 0 {
		log.Println("sub list empty")
		return
	}

	if try {
		log.Println("run in try mode")
	}

	//path = "/volume1/down/shows/Alex.Rider.S01.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTG"
	//subExt = "chs.eng.ass"
	//videoExt = "mkv"

	for k, v := range videoList {
		index := strings.LastIndex(v, videoExt)
		if index < 0 {
			continue
		}

		// 视频名称
		// 新文件名称= 视频名称.subExt
		videoName := v[0 : index-1]
		if sub, ok := subList[k]; ok && videoName != "" {
			oldSubFile := "./" + sub
			newSubFile := "./" + videoName + "." + subExt

			if oldSubFile == newSubFile {
				log.Printf("no need rename: %s\n", oldSubFile)
				continue
			}

			log.Printf("try rename %s to %s\n", oldSubFile, newSubFile)
			if try {
				continue
			}

			err := os.Rename(oldSubFile, newSubFile)
			if err != nil {
				panic(err)
			}
			log.Printf("rename %s to %s done\n", oldSubFile, newSubFile)
		}
	}
}

func listVideo(path, ext string, mode int) map[string]string {
	//path = "/volume1/down/shows/Alex.Rider.S01.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTG"
	//ext = "mkv"

	pathInfo, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	list := make(map[string]string, 0)
	for _, file := range pathInfo {
		index := strings.LastIndex(file.Name(), ext)
		if index < 0 {
			continue
		}

		find := MatchFile(mode, file.Name())
		if len(find) == 0 {
			continue
		}

		list[find[0]] = file.Name()
	}

	return list
}

func listSub(path, ext string, mode int) map[string]string {
	//path = "/volume1/down/shows/Alex.Rider.S01.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTG"
	//ext = "chs.eng.ass"

	pathInfo, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	list := make(map[string]string, 0)
	for _, file := range pathInfo {
		index := strings.LastIndex(file.Name(), ext)
		if index < 0 {
			continue
		}

		find := MatchFile(mode, file.Name())
		if len(find) == 0 {
			continue
		}

		list[find[0]] = file.Name()
	}

	return list
}

func MatchFile(mode int, file string) []string {
	switch mode {
	case 1:
		return m1.FindStringSubmatch(file)
	case 2:
		return m2.FindStringSubmatch(file)
	case 3:
		return m3.FindStringSubmatch(file)
	case 4:
		return m4.FindStringSubmatch(file)
	case 5:
		return m5.FindStringSubmatch(file)
	}

	return []string{}
}
