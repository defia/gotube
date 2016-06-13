package gotube

import (
	"errors"
	"log"
)

type VideoRawInfo struct {
	Title          string
	Adaptive_fmts  []KV
	Fmt_stream_map []KV
}

type VideoInfo struct {
	Title                 string
	AvailableDownloadInfo []*DownloadInfo
}

type DownloadInfo struct {
	FileExtension string
	Url           string
	Type          string
	itag          string
}

type KV map[string]string

func (kv KV) Get(k string) string {
	return kv[k]
}

//Thanks to https://github.com/gantt/downloadyoutube
var (
	FORMAT_LABEL = map[string]string{"5": "FLV 240p", "18": "MP4 360p", "22": "MP4 720p", "34": "FLV 360p", "35": "FLV 480p", "37": "MP4 1080p", "38": "MP4 2160p", "43": "WebM 360p", "44": "WebM 480p", "45": "WebM 720p", "46": "WebM 1080p", "135": "MP4 480p - no audio", "137": "MP4 1080p - no audio", "138": "MP4 2160p - no audio", "139": "M4A 48kbps - audio", "140": "M4A 128kbps - audio", "141": "M4A 256kbps - audio", "264": "MP4 1440p - no audio", "266": "MP4 2160p - no audio", "298": "MP4 720p60 - no audio", "299": "MP4 1080p60 - no audio"}

	FORMAT_TYPE  = map[string]string{"5": "flv", "18": "mp4", "22": "mp4", "34": "flv", "35": "flv", "37": "mp4", "38": "mp4", "43": "webm", "44": "webm", "45": "webm", "46": "webm", "135": "mp4", "137": "mp4", "138": "mp4", "139": "m4a", "140": "m4a", "141": "m4a", "264": "mp4", "266": "mp4", "298": "mp4", "299": "mp4"}
	FORMAT_ORDER = []string{"5", "18", "34", "43", "35", "135", "44", "22", "298", "45", "37", "299", "46", "264", "38", "266", "139", "140", "141"}
)

func GenerateVideoInfo(b []byte) (*VideoRawInfo, error) {
	title, err := parseTitle(b)
	if err != nil {
		log.Println(err)
	}
	fmt1, err := parseAdaptivefmts(b)
	if err != nil {
		log.Println(err)

	}

	fmt2, err := parsefmtstreamap(b)
	if err != nil {
		log.Println(err)
	}
	return &VideoRawInfo{
		Title:          title,
		Adaptive_fmts:  toMaps(fmt1),
		Fmt_stream_map: toMaps(fmt2),
	}, nil

}

func (rawinfo *VideoRawInfo) toVideoInfo() *VideoInfo {
	info := new(VideoInfo)
	info.Title = rawinfo.Title

	dinfos := make([]*DownloadInfo, 0)
	for _, v := range rawinfo.Adaptive_fmts {
		d, err := generateDownloadInfo(v)
		if err != nil {
			continue
		}
		dinfos = append(dinfos, d)

	}
	for _, v := range rawinfo.Fmt_stream_map {
		d, err := generateDownloadInfo(v)
		if err != nil {
			continue
		}
		dinfos = append(dinfos, d)

	}
	info.AvailableDownloadInfo = dinfos

	return info

}

func generateDownloadInfo(fmts map[string]string) (*DownloadInfo, error) {
	url := fmts["url"]
	itag := fmts["itag"]
	var err error
	if url == "" {
		return nil, errors.New("contains no url")
	}
	label := FORMAT_LABEL[itag]
	if label == "" {
		err = errors.New("not a common format")
	}
	return &DownloadInfo{
		FileExtension: FORMAT_TYPE[itag],
		Url:           fmts["url"],
		Type:          FORMAT_LABEL[itag],
		itag:          itag,
	}, err
}
