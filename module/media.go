/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: media.go
 */

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	pacs_dcm4chee "gitee.com/uni-minds/bridge-pacs/dcm4chee"
	pacs_global "gitee.com/uni-minds/bridge-pacs/global"
	pacs_tools "gitee.com/uni-minds/bridge-pacs/tools"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/utils/media"
	"gitee.com/uni-minds/utils/tools"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MediaImportJson struct {
	Source    string   `json:"source"`
	Filename  string   `json:"target"`
	Descript  string   `json:"descript"`
	View      string   `json:"view"`
	Keywords  []string `json:"keywords"`
	GroupName string   `json:"groupname"`
	Fcode     string   `json:"fcode"`
	PatientID string   `json:"patientid"`
	MachineID string   `json:"machineid"`
}

type MediaInfo struct {
	Id             int
	MediaUUID      string
	DisplayName    string
	Path           string
	Width          int
	Height         int
	Duration       float64
	Frames         int
	Fps            float64
	MediaType      string
	UploadUid      int
	UploadTime     time.Time
	Memo           string
	PatientId      string
	MachineId      string
	Metadata       string
	CrfDefine      string
	Keywords       []string
	MediaData      string
	LabelAuthors   map[int]string
	LabelReviewers map[string]int
	LabelProgress  int
	CoworkType     string
	MediaHash      string
}

type MediaSummary struct {
	Id          int
	MediaUUID   string
	DisplayName string
	Memo        string
	Duration    float64
	Frames      int
	Views       string
	Keywords    string
	Width       int
	Height      int
}

func MediaGet(i interface{}) (info MediaInfo, err error) {
	dbi, err := database.GetMedia().Get(i)
	if err == nil {
		// support common media
		return ConvertMediaInfoFromDbMedia(dbi), nil
	}

	switch i.(type) {
	case string:
		str := i.(string)
		// support instance_id
		if global.IsDicomUUID(str) {
			ps := database.BridgeGetPacsServerHandler()
			info, err := ps.FindInstanceByIdLocal(str)
			if err == nil {
				return ConvertMediaInfoFromInstanceInfo(info), nil
			}
		}
	}

	return info, fmt.Errorf("cannot get media: %v", i)
}

func MediaSetMemo(mediaUUID string, memo string) (err error) {
	if _, err = MediaGet(mediaUUID); err != nil {
		return err
	}
	if global.IsDicomUUID(mediaUUID) {
		ps := database.BridgeGetPacsServerHandler()
		return ps.InstanceUpdateLabelMemo(mediaUUID, memo)

	} else {
		return database.GetMedia().UpdateMemo(mediaUUID, memo)

	}
}

func MediaSetLabelProgress(mediaUUID string, progress int) (err error) {
	if _, err := MediaGet(mediaUUID); err != nil {
		return err
	}
	return database.GetMedia().LabelUpdateProgress(mediaUUID, progress)
}

func ConvertSummaryFromMediaInfo(mediaInfo MediaInfo) (summary MediaSummary) {

	if mediaInfo.Duration == 0 || mediaInfo.Frames == 0 || mediaInfo.Width == 0 || mediaInfo.Height == 0 {
		switch strings.ToLower(filepath.Ext(mediaInfo.Path)) {
		case ".ogv", ".mp4", ".m3u8":
			width, height, frames, duration, fps, _, err := MediaGetInfoFromFile(mediaInfo.Path)
			if err == nil {
				summary.Frames = frames
				summary.Duration = duration
				summary.Width = width
				summary.Height = height

				_ = database.GetMedia().UpdateFramesAndDuration(mediaInfo.MediaUUID, fps, duration, frames)
				_ = database.GetMedia().UpdateWidthAndHeight(mediaInfo.MediaUUID, width, height)
			}
		}
	}

	jbs, _ := json.Marshal(mediaInfo.Keywords)

	summary = MediaSummary{
		Id:          mediaInfo.Id,
		DisplayName: mediaInfo.DisplayName,
		Memo:        mediaInfo.Memo,
		Duration:    mediaInfo.Duration,
		Frames:      mediaInfo.Frames,
		Views:       mediaInfo.CrfDefine,
		Keywords:    string(jbs),
		Width:       mediaInfo.Width,
		Height:      mediaInfo.Height,
		MediaUUID:   mediaInfo.MediaUUID,
	}
	return summary
}
func ConvertSummaryFromMediaInfos(infos []MediaInfo) (summaries []MediaSummary) {
	summaries = make([]MediaSummary, 0)
	for _, info := range infos {
		summaries = append(summaries, ConvertSummaryFromMediaInfo(info))
	}
	return summaries
}
func ConvertMediaInfoFromDbMedia(dbi database.DbStructMedia) MediaInfo {
	info := MediaInfo{
		Id:             dbi.Id,
		MediaUUID:      dbi.MediaUUID,
		DisplayName:    dbi.DisplayName,
		Path:           dbi.Path,
		Width:          dbi.Width,
		Height:         dbi.Height,
		Duration:       dbi.Duration,
		Frames:         dbi.Frames,
		Fps:            dbi.Fps,
		MediaType:      dbi.MediaType,
		UploadUid:      dbi.UploadUid,
		Memo:           dbi.Memo,
		PatientId:      dbi.PatientId,
		MachineId:      dbi.MachineId,
		Metadata:       dbi.Metadata,
		CrfDefine:      dbi.CrfDefine,
		MediaData:      dbi.MediaData,
		LabelProgress:  dbi.LabelProgress,
		CoworkType:     dbi.CoworkType,
		MediaHash:      dbi.MediaHash,
		UploadTime:     time.Unix(dbi.UploadTime, 0),
		LabelAuthors:   nil,
		LabelReviewers: nil,
		Keywords:       nil,
	}

	if dbi.LabelAuthor != "" && dbi.LabelAuthor != "null" {
		var authors map[int]string
		json.Unmarshal([]byte(dbi.LabelAuthor), &authors)
		info.LabelAuthors = authors
	} else {
		info.LabelAuthors = make(map[int]string)
	}

	if dbi.LabelReviewer != "" && dbi.LabelReviewer != "null" {
		var reviews map[string]int
		json.Unmarshal([]byte(dbi.LabelReviewer), &reviews)
		info.LabelReviewers = reviews
	} else {
		info.LabelReviewers = make(map[string]int)
	}

	if dbi.Keywords != "" {
		var keyworks []string
		json.Unmarshal([]byte(dbi.Keywords), &keyworks)
		info.Keywords = keyworks
	} else {
		info.Keywords = make([]string, 0)
	}

	return info
}
func ConvertMediaInfoToDbMedia(dbi MediaInfo) database.DbStructMedia {
	authors := ""
	reviews := ""
	keywords := ""

	if len(dbi.LabelAuthors) > 0 {
		tmp, _ := json.Marshal(dbi.LabelAuthors)
		authors = string(tmp)
	}

	if len(dbi.LabelReviewers) > 0 {
		tmp, _ := json.Marshal(dbi.LabelReviewers)
		reviews = string(tmp)
	}

	if len(dbi.Keywords) > 0 {
		tmp, _ := json.Marshal(dbi.Keywords)
		keywords = string(tmp)
	}

	db := database.DbStructMedia{
		Id:            dbi.Id,
		MediaUUID:     dbi.MediaUUID,
		DisplayName:   dbi.DisplayName,
		Path:          dbi.Path,
		Width:         dbi.Width,
		Height:        dbi.Height,
		Duration:      dbi.Duration,
		Frames:        dbi.Frames,
		Fps:           dbi.Fps,
		MediaType:     dbi.MediaType,
		UploadUid:     dbi.UploadUid,
		UploadTime:    dbi.UploadTime.Unix(),
		Memo:          dbi.Memo,
		Metadata:      dbi.Metadata,
		PatientId:     dbi.PatientId,
		MachineId:     dbi.MachineId,
		CrfDefine:     dbi.CrfDefine,
		LabelAuthor:   authors,
		LabelReviewer: reviews,
		LabelProgress: dbi.LabelProgress,
		CoworkType:    dbi.CoworkType,
		MediaData:     dbi.MediaData,
		MediaHash:     dbi.MediaHash,
		Keywords:      keywords,
	}
	return db
}
func ConvertMediaInfoFromInstanceInfo(dbi pacs_global.InstanceInfo) MediaInfo {
	type pacsInstance struct {
		StudiesId string
		SeriesId  string
	}

	bMetadata, _ := json.Marshal(pacsInstance{
		StudiesId: dbi.StudiesId,
		SeriesId:  dbi.SeriesId,
	})

	labelAuthors := make(map[int]string)
	labelReviewer := make(map[string]int)
	labelProgress := 0

	if labelInfos, err := LabelGetByMediaUUID(dbi.InstanceId); err == nil {
		if len(labelInfos) > 0 {
			labelInfo := labelInfos[0]
			labelProgress = labelInfo.Progress
			labelAuthors[labelInfo.AuthorUid] = labelInfo.LabelUUID
			labelReviewer[labelInfo.LabelUUID] = labelInfo.ReviewUid
		}
	}

	info := MediaInfo{
		Id:             int(dbi.RecordDatetime),
		MediaUUID:      dbi.InstanceId,
		DisplayName:    dbi.InstanceId,
		Width:          dbi.MediaWidth,
		Height:         dbi.MediaHeight,
		Duration:       dbi.Duration,
		Frames:         dbi.Frames,
		Fps:            float64(dbi.Frames) / dbi.Duration,
		MediaType:      dbi.MediaType,
		UploadUid:      0,
		UploadTime:     time.Unix(dbi.RecordDatetime, 0),
		Memo:           dbi.LabelMemo,
		PatientId:      dbi.StudiesId,
		MachineId:      "backend",
		Metadata:       string(bMetadata),
		CrfDefine:      dbi.LabelView,
		Keywords:       []string{dbi.LabelDiagnose, dbi.LabelInfoAttend, dbi.LabelInfoAttend},
		MediaData:      "",
		LabelAuthors:   labelAuthors,
		LabelReviewers: labelReviewer,
		LabelProgress:  labelProgress,
		CoworkType:     "single",
		MediaHash:      "",
	}

	switch dbi.MediaType {
	case pacs_dcm4chee.MEDIA_TYPE_IMAGE:
		cache := pacs_tools.DecodeCachePathImage(dbi.CacheLocalPath)
		info.Path = cache.ImageFilename

	case pacs_dcm4chee.MEDIA_TYPE_MULTI_FRAME:
		cache := pacs_tools.DecodeCachePathVideo(dbi.CacheLocalPath)
		info.Path = cache.VideoFilename

	}
	return info
}

func ConvertMediaInfoFromInstanceInfos(dbi []pacs_global.InstanceInfo) []MediaInfo {
	result := make([]MediaInfo, 0)
	for _, info := range dbi {
		result = append(result, ConvertMediaInfoFromInstanceInfo(info))
	}
	return result
}

func MediaRescan(mediaUUID string) error {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return err
	}

	target := info.Path
	// 存在绝对路径则以绝对路径为准
	if len(target) > 0 && target[0] != '/' {
		target = path.Join(global.GetPaths().Media)
	}

	log.Println("t", "rescan media:", target)

	width, height, frames, duration, fps, _, err := MediaGetInfoFromFile(target)
	if frames == 0 {
		frames = int(math.Floor(duration * fps))
	}

	if err != nil {
		return err
	}

	database.GetMedia().UpdateFramesAndDuration(mediaUUID, fps, duration, frames)
	database.GetMedia().UpdateWidthAndHeight(mediaUUID, width, height)
	return nil
}

func MediaImportM3U8(groupId int, uid int, disp, videoFolder, machineId, tagView, tagCustom string) error {
	vfile := path.Join(videoFolder, "video.m3u8")

	width, height, frames, duration, fps, _, err := MediaGetInfoFromFile(vfile)
	if err != nil {
		return err
	}

	checksum, _ := tools.FileGetMD5(vfile)

	info := MediaInfo{
		MediaUUID:      tools.RandString0f(32),
		DisplayName:    disp,
		Path:           vfile,
		Width:          width,
		Height:         height,
		Duration:       duration,
		Frames:         frames,
		Fps:            fps,
		MediaType:      global.MediaTypeM3u8,
		UploadUid:      uid,
		UploadTime:     time.Now(),
		Memo:           "",
		PatientId:      tagCustom,
		MachineId:      machineId,
		Metadata:       "",
		CrfDefine:      tagView,
		Keywords:       nil,
		MediaData:      "",
		LabelAuthors:   nil,
		LabelReviewers: nil,
		LabelProgress:  0,
		CoworkType:     "",
		MediaHash:      checksum,
	}

	dbi := ConvertMediaInfoToDbMedia(info)
	mediaUUID, err := database.GetMedia().Create(dbi)

	if err != nil {
		return err
	}

	if groupId < 0 {
		//导入默认组
		groupId = GroupGetGid("default_stream")
		if groupId > 0 {
			GroupAddMedia(groupId, []string{mediaUUID})
		} else {
			log.Println("e", "cannot get group: default_media")
		}
	}
	return nil
}

//生成mid的标注作者和审核摘要信息

func GetSummaryInstance(instanceId string) (summary MediaSummary, err error) {
	ps := database.BridgeGetPacsServerHandler()
	info, err := ps.FindInstanceByIdLocal(instanceId)
	if err != nil {
		return summary, err
	}

	summary.MediaUUID = info.InstanceId
	summary.DisplayName = info.InstanceId
	summary.Memo = info.LabelMemo
	summary.Views = info.LabelView

	summary.Duration = info.Duration
	summary.Frames = info.Frames
	summary.Width = info.MediaWidth
	summary.Height = info.MediaHeight

	return summary, nil
}
func GetSummaryMedia(mediaUUID string) (summary MediaSummary, err error) {
	info, err := MediaGet(mediaUUID)
	if err != nil {
		return
	}

	return ConvertSummaryFromMediaInfo(info), nil
}

func MediaGetInfoFromFile(mediafile string) (width, height, frames int, duration float64, fps float64, codec string, err error) {
	if _, err = os.Stat(mediafile); err != nil {
		return 0, 0, 0, 0, 0, "", err
	}

	if info, err := media.FfprobeMedia(mediafile); err != nil {
		return 0, 0, 0, 0, 0, "", err

	} else {
		stream := info.Streams[0]
		duration, _ = strconv.ParseFloat(info.Format.Duration, 10)

		if strings.Contains(stream.AvgFrameRate, "/") {
			n := strings.Split(stream.AvgFrameRate, "/")
			n1, _ := strconv.ParseFloat(n[0], 10)
			n2, _ := strconv.ParseFloat(n[1], 10)
			if n1*n2 > 0 {
				fps = n1 / n2
			}
		}

		if fps == 0 && strings.Contains(stream.RFrameRate, "/") {
			n := strings.Split(stream.RFrameRate, "/")
			n1, _ := strconv.ParseFloat(n[0], 10)
			n2, _ := strconv.ParseFloat(n[1], 10)
			if n1*n2 > 0 {
				fps = n1 / n2
			}
		}

		switch stream.CodecName {
		case "h264":
			frames, _ = strconv.Atoi(stream.NbFrames)
		case "theora", "vp8":
			frames = stream.DurationTs
		case "hevc":
			return 0, 0, 0, 0, fps, "", errors.New("unsupport hevc")

		case "vp9":

		default:
		}

		if frames == 0 {
			frames = int(math.Floor(duration * fps))
		}

		return stream.Width, stream.Height, frames, duration, fps, stream.CodecName, nil
	}
}

func InstanceGetVideo(instanceId string, uid int) ([]byte, error) {
	ps := database.BridgeGetPacsServerHandler()
	//pi, err := ps.FindInstanceById(instanceId)
	//if err != nil {
	//	return nil, err
	//}

	bs, _, err := ps.InstanceMediaGet(instanceId, "ogv")
	if err != nil {
		return nil, err
	} else {
		return bs, nil
	}
}
