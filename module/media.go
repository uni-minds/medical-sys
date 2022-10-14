package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

type MediaSummaryInfo struct {
	DisplayName string
	Memo        string
	Duration    float64
	Frames      int
	Views       string
	Keywords    string
	Width       int
	Height      int
	Hash        string
}

//type MediaSummaryAuthorInfo struct {
//	Realname   string
//	Frames     int
//	Counts     int
//	UpdateTime string
//	Progress   string
//	Hash       string
//	Memo       string
//	IsReviewed bool
//	IsModified bool
//}
//type MediaSummaryReviewInfo struct {
//	Realname   string
//	UpdateTime string
//	Progress   string
//	Tips       string
//	Hash       string
//	Memo       string
//	Author     string
//	AuthorTime string
//}
type MediaImportJson struct {
	Source    string `json:"source"`
	Filename  string `json:"target"`
	Descript  string `json:"descript"`
	View      string `json:"view"`
	Keywords  string `json:"keywords"`
	GroupName string `json:"groupname"`
	Fcode     string `json:"fcode"`
	PatientID string `json:"patientid"`
	MachineID string `json:"machineid"`
}

//生成mid的标注作者和审核摘要信息
func MediaGetSummary(mid int) (summary MediaSummaryInfo, err error) {
	mi, err := database.MediaGet(mid)
	if err != nil {
		return
	}

	summary.Hash = mi.Hash
	summary.DisplayName = mi.DisplayName
	summary.Memo = mi.Memo

	// HASH再识别,存在问题？
	//if len(mi.Hash) != 32 {
	//	log.Println("Find old import hash", mi.Hash, mi.Path)
	//	hash := tools.GetFileMD5(mi.Path)
	//	log.Println("Update hash", hash)
	//	database.MediaUpdateHash(mi.Mid, hash)
	//	li,err :=database.LabelGet(mi.Hash)
	//	if err == nil {
	//		err = database.LabelUpdateMediaHash(li.Lid,hash)
	//		if err!= nil {
	//			log.Println(err.Error())
	//		}
	//		log.Println("Label UPD",li.Lid,mi.Hash,hash)
	//	}
	//	summary.Hash = hash
	//}

	// 切面再识别
	switch mi.IncludeViews {
	case "", "null", "[]":
		views := MediaAutoGenViews(mi.DisplayName)
		if len(views) > 0 {
			_ = database.MediaSetViews(mid, views)
			jb, _ := json.Marshal(views)
			summary.Views = string(jb)
		} else {
			summary.Views = "[]"
		}

	default:
		summary.Views = mi.IncludeViews
	}
	// 关键字再识别
	switch mi.Keywords {
	case "", "null", "[]":
		keywords := MediaAutoGenKeywords(mi.DisplayName)
		if len(keywords) > 0 {
			_ = database.MediaSetKeywords(mid, keywords)
			jb, _ := json.Marshal(keywords)
			summary.Keywords = string(jb)
		} else {
			summary.Keywords = "[]"
		}

	default:
		summary.Keywords = mi.Keywords
	}
	// 媒体信息再识别
	summary.Duration = mi.Duration
	summary.Frames = mi.Frames
	summary.Width = mi.Width
	summary.Height = mi.Height
	if mi.Duration == 0 || mi.Frames == 0 || mi.Width == 0 || mi.Height == 0 {
		switch strings.ToLower(filepath.Ext(mi.Path)) {
		case ".ogv":
			w, h, f, d, _, err := MediaInfo(mi.Path)
			if err == nil {
				summary.Frames = f
				summary.Duration = d
				summary.Width = w
				summary.Height = h
				_ = database.MediaUpdateFramesAndDuration(mid, f, d)
				_ = database.MediaUpdateWidthAndHeight(mid, w, h)
			}
		}
	}
	// 标注信息再识别
	//switch mi.LabelAuthorUid {
	//case "", "[]", "null":
	//	labelAuthorsUid := make([]int, 0)
	//	labelAuthorsLid := make([]int, 0)
	//	lis, err := database.LabelGetAll(mid, 0, global.LabelTypeAuthor)
	//	if err != nil {
	//		log.Println(err.Error())
	//		break
	//	}
	//	for _, li := range lis {
	//		labelAuthorsUid = append(labelAuthorsUid, li.Uid)
	//		labelAuthorsLid = append(labelAuthorsLid, li.Lid)
	//	}
	//	MediaUpdateLabelAuthorsUidLid(mid, labelAuthorsUid, labelAuthorsLid)
	//	summary.AuthorUids = labelAuthorsUid
	//	summary.AuthorLids = labelAuthorsLid
	//
	//default:
	//	json.Unmarshal([]byte(mi.LabelAuthorUid), &summary.AuthorUids)
	//	json.Unmarshal([]byte(mi.LabelAuthorsLid), &summary.AuthorLids)
	//}
	//// 审阅信息再识别
	//switch mi.LabelReviewUid {
	//case "", "[]", "null":
	//	labelReviewsUid := make([]int, 0)
	//	labelReviewsLid := make([]int, 0)
	//	lis, err := database.LabelGetAll(mid, 0, global.LabelTypeReview)
	//	if err != nil {
	//		log.Println(err.Error())
	//		break
	//	}
	//	for _, li := range lis {
	//		labelReviewsUid = append(labelReviewsUid, li.Uid)
	//		labelReviewsLid = append(labelReviewsLid, li.Lid)
	//	}
	//	MediaUpdateLabelReviewsUidLid(mid, labelReviewsUid, labelReviewsLid)
	//	summary.ReviewUids = labelReviewsUid
	//	summary.ReviewLids = labelReviewsLid
	//
	//default:
	//	json.Unmarshal([]byte(mi.LabelReviewUid), &summary.ReviewUids)
	//	json.Unmarshal([]byte(mi.LabelReviewsLid), &summary.ReviewLids)
	//}

	return
}
func MediaGetMid(hash string) int {
	mi, err := database.MediaGet(hash)
	if err != nil {
		log.Println("MediaGet E", err.Error())
		return -1
	}
	return mi.Mid
}

/*
//用于JsGrid生成标注作者按钮
func MediaGetLabelAuthorsSummary(uids, lids []int) []MediaSummaryAuthorInfo {
	authors := make([]MediaSummaryAuthorInfo, 0)
	if len(uids) > 0 {
		for i, uid := range uids {
			li, _ := database.LabelGet(lids[i])
			updateTime := ""
			if li.ModifyTime != "" {
				updateTime = li.ModifyTime
			} else {
				updateTime = li.CreateTime
			}

			var labelInfoDataAuthor database.LabelInfoAuthorData
			if li.Type != global.LabelTypeAuthor {
				log.Println("error label_author wrong type of lid", li.Lid)
				continue
			}

			if len(li.Data) > 10 {
				json.Unmarshal([]byte(li.Data), &labelInfoDataAuthor)
				if labelInfoDataAuthor.Json == "" {
					labelInfoDataAuthor.Json = li.Data
					database.LabelUpdateLabelData(li.Lid, li.Uid, li.Frames, li.Counts, labelInfoDataAuthor, li.Progress)
					li, _ = database.LabelGet(li.Lid)
				}
			}
			authors = append(authors, MediaSummaryAuthorInfo{
				Realname:   UserGetRealname(uid),
				Frames:     li.Frames,
				Counts:     li.Counts,
				UpdateTime: updateTime,
				Progress:   li.Progress,
				Hash:       li.Hash,
				Memo:       li.Memo,
			})
		}
	}
	return authors
}

//用于JsGrid生成标注审核按钮
func MediaGetLabelReviewersSummary(uids, lids []int) []MediaSummaryReviewInfo {
	reviews := make([]MediaSummaryReviewInfo, 0)
	if len(uids) > 0 {
		for i, uid := range uids {
			li, err := database.LabelGet(lids[i])
			if err != nil {
				log.Println(err.Error())
			}
			updateTime := ""
			if li.ModifyTime != "" {
				updateTime = li.ModifyTime
			} else {
				updateTime = li.CreateTime
			}

			var labelInfoDataReview database.LabelInfoReviewerData
			if li.Type != global.LabelTypeReview {
				log.Println("error label_review wrong type for lid", li.Lid)
				continue
			}

			var tips = "Pure"
			if len(li.Data) > 10 {
				json.Unmarshal([]byte(li.Data), &labelInfoDataReview)
				if labelInfoDataReview.Json == "" {
					labelInfoDataReview.Json = li.Data
					database.LabelUpdateLabelData(li.Lid, li.Uid, li.Frames, li.Counts, labelInfoDataReview, li.Progress)
				} else {
					if labelInfoDataReview.BasedAuthor > 0 {
						tips = fmt.Sprintf("审阅基于[%s]@%s",
							UserGetRealname(labelInfoDataReview.BasedAuthor), labelInfoDataReview.BasedTime)
					}
				}
			}

			reviews = append(reviews, MediaSummaryReviewInfo{
				Realname:   UserGetRealname(uid),
				UpdateTime: updateTime,
				Progress:   li.Progress,
				Tips:       tips,
				Hash:       li.Hash,
				Memo:       li.Memo,
			})
		}
	}
	return reviews
}
*/
func MediaImport(input, displayName, memo string, ownerUid int) (mi database.MediaInfo, err error) {
	filefull := path.Base(input)
	fileext := path.Ext(input)
	filename := filefull[:len(filefull)-len(fileext)]
	datefolder := time.Now().Format("20060102")

	md5 := tools.GetFileMD5(input)

	//_,err = MediaSearch(md5)
	//if err != errors.New(EMediaNotExist) {
	//	err = errors.New(EMediaAlreadyExisted)
	//	return
	//}

	if displayName == "" {
		displayName = filename
	}

	basefolder := "."
	//global.GetMediaRoot()

	rawroot := path.Join(basefolder, "raw", datefolder)
	ogvroot := path.Join(basefolder, "ogv", datefolder)
	gifroot := path.Join(basefolder, "gif", datefolder)

	perm := os.ModePerm
	if err = os.MkdirAll(rawroot, perm); err != nil {
		return
	}
	if err = os.MkdirAll(ogvroot, perm); err != nil {
		return
	}
	if err = os.MkdirAll(gifroot, perm); err != nil {
		return
	}

	rawfile := path.Join(rawroot, filefull)
	if err = tools.CopyFile(input, rawfile); err != nil {
		return
	}

	ogvfile := path.Join(ogvroot, filename+".ogv")
	if err = tools.FFmpegToOGV(input, ogvfile); err != nil {
		return
	}
	//
	//ogvhash, _ := tools.GetFileMD5(ogvfile)
	//w, h, f, d, _, err := MediaInfo(ogvfile)
	//if err != nil {
	//	return
	//}

	giffile := path.Join(gifroot, filename+".gif")
	if err = tools.FFmpegToGIF(input, giffile); err != nil {
		return
	}

	mi = database.MediaInfo{
		DisplayName: displayName,
		Path:        rawfile,
		Hash:        md5,
		UploadTime:  time.Now().Format(time.RFC3339),
		Memo:        memo,
	}
	log.Println(mi)
	return
}
func MediaInfo(mediafile string) (width, height, frames int, duration float64, codec string, err error) {
	info, err := tools.FFprobe(mediafile)
	if err != nil {
		return
	}
	var w, h, d, f string

	cfg, _ := goconfig.LoadFromData([]byte(info))

	codec, _ = cfg.GetValue("STREAM", "codec_name")

	w, _ = cfg.GetValue("STREAM", "width")
	width, _ = strconv.Atoi(w)

	h, _ = cfg.GetValue("STREAM", "height")
	height, _ = strconv.Atoi(h)

	d, _ = cfg.GetValue("STREAM", "duration")
	duration, _ = strconv.ParseFloat(d, 32)

	switch codec {
	case "h264":
		f, _ = cfg.GetValue("STREAM", "nb_frames")
	case "theora", "vp8":
		f, _ = cfg.GetValue("STREAM", "duration_ts")
	default:
		f = "0"
	}
	frames, _ = strconv.Atoi(f)

	return
}
func MediaImportDir(source, store string, uid, gid int) error {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	destFolder := filepath.Join(store, time.Now().Format("20060102-15H"))
	_, err = os.Stat(destFolder)
	if err != nil {
		_ = os.MkdirAll(destFolder, 0777)
	}

	for _, v := range files {
		if v.IsDir() {
			log.Println("Subdir, forbidden")
			continue
		}

		srcFile := filepath.Join(source, v.Name())

		disps := strings.Split(v.Name(), ".")

		switch strings.ToLower(filepath.Ext(v.Name())) {
		case ".ogv":
			mid, err := MediaImportUsVideoOgv(srcFile, destFolder, disps[0], []string{"4AP"}, uid)
			if err != nil {
				continue
			}
			database.GroupAddMedia(gid, mid)
		}
	}
	return nil
}
func MediaImportUsVideoOgv(srcFile, destFolder, dispname string, views []string, uid int) (mid int, err error) {
	checksum := tools.GetFileMD5(srcFile)
	if checksum == "" {
		return 0, errors.New("HASH校验值为空:" + srcFile)
	}

	log.Println("checksum", checksum)
	mi, err := database.MediaGet(checksum)
	if err == nil {
		return mi.Mid, errors.New(global.EMediaAlreadyExisted)
	}

	destFilename := filepath.Base(srcFile)
	destFile := filepath.Join(destFolder, destFilename)
	for {
		_, err = os.Stat(destFile)
		if err == nil {
			data := strings.Split(destFilename, ".")
			destFile = filepath.Join(destFolder, data[0]+"_"+tools.GenSaltString(5)+".ogv")
		} else {
			break
		}
	}
	_ = tools.CopyFile(srcFile, destFile)
	bv, _ := json.Marshal(views)

	width, height, frames, duration, encoder, err := MediaInfo(srcFile)
	mi = database.MediaInfo{
		Mid:            0,
		DisplayName:    dispname,
		Path:           destFile,
		Hash:           checksum,
		Duration:       duration,
		Frames:         frames,
		Width:          width,
		Height:         height,
		Status:         0,
		UploadTime:     "",
		UploadUid:      uid,
		PatientID:      "",
		MachineID:      "",
		FolderName:     filepath.Base(filepath.Dir(srcFile)),
		Fcode:          "",
		IncludeViews:   string(bv),
		Keywords:       "",
		Memo:           "",
		MediaType:      "",
		MediaData:      "",
		LabelAuthorUid: 0,
		LabelReviewUid: 0,
	}
	mid, err = database.MediaCreate(mi)

	detail := database.MediaInfoUltrasonicVideo{
		PathRaw:  srcFile,
		HashRaw:  checksum,
		PathJpgs: "[]",
		Encoder:  encoder,
	}
	_ = database.MediaUpdateDetail(mid, detail)
	return mid, nil
}
func MediaImportFromJson(uid int, srcFolder, destFolder string, data []MediaImportJson) error {
	totalLen := len(data)
	for i, v := range data {
		prog := i * 100 / totalLen
		if prog%10 == 0 {
			fmt.Println("Import progress:")
		}
		srcFile := filepath.Join(srcFolder, v.Filename)
		dispname := strings.Split(filepath.Base(v.Source), ".")[0]
		var views []string
		_ = json.Unmarshal([]byte(v.View), &views)
		mid, err := MediaImportUsVideoOgv(srcFile, destFolder, dispname, views, uid)
		if err != nil {
			log.Println("导入过程中错误：", err.Error())
			continue
		}
		_ = database.MediaUpdateFolderName(mid, filepath.Dir(v.Source))

		var gid int
		if v.GroupName == "" {
			gid, _ = database.GroupCreate(database.GroupInfo{
				GroupName:   global.DefGroupUngrouped,
				DisplayName: global.DefGroupUngroupedName,
			})
		} else {
			gi, err := database.GroupGet(v.GroupName)
			if err != nil {
				fmt.Println("创建分组：", v.GroupName)
				gid, _ = database.GroupCreate(database.GroupInfo{
					GroupName:   v.GroupName,
					DisplayName: v.GroupName,
				})
			} else {
				gid = gi.Gid
			}
		}
		database.GroupAddMedia(gid, mid)
	}
	return nil
}
func MediaGetRealpath(hash string, uid int) string {
	mi, err := userGetMediaInfo(uid, hash)
	log.Println("Find", hash, uid)
	log.Println(mi.Path)
	if err != nil {
		return ""
	} else {
		return mi.Path
	}
}
func MediaUpdateLabel(mid, uid, lid int, labeltype string) error {
	mi, err := database.MediaGet(mid)
	if err != nil {
		return err
	}

	switch labeltype {
	case global.LabelTypeAuthor:
		//db := parseUidLidStringToDbMap(mi.LabelAuthorUid, mi.LabelAuthorsLid)
		//db[uid] = lid
		//uidstr, lidstr := parseUidLidMapMapToString(db)
		return database.MediaUpdateLabelAuthorUidLid(mi.Mid, uid, lid)

	case global.LabelTypeReview:
		//db := parseUidLidStringToDbMap(mi.LabelReviewUid, mi.LabelReviewsLid)
		//db[uid] = lid
		//uidstr, lidstr := parseUidLidMapMapToString(db)
		return database.MediaUpdateLabelReviewUidLid(mi.Mid, uid, lid)

	default:
		return errors.New(global.EMediaUnknownType)
	}
}

func MediaAutoGenViews(dispname string) []string {
	views := make([]string, 0)

	if strings.Contains(dispname, "3VT") || strings.Contains(dispname, "三血管气管") {
		views = append(views, "3VT")
	} else if strings.Contains(dispname, "3V") || strings.Contains(dispname, "三血管") {
		views = append(views, "3V")
	}

	if strings.Contains(dispname, "L") || strings.Contains(dispname, "左室") {
		views = append(views, "L")
	}

	if strings.Contains(dispname, "R") || strings.Contains(dispname, "右室") {
		views = append(views, "R")
	}

	if strings.Contains(dispname, "AP") || strings.Contains(dispname, "四腔") {
		views = append(views, "4AP")
	} else if strings.Contains(dispname, "AA") || strings.Contains(dispname, "主动脉弓") {
		views = append(views, "AA")
	} else if strings.Contains(dispname, "AC") || strings.Contains(dispname, "腹横切") {
		views = append(views, "AC")
	} else if strings.Contains(dispname, "A") || strings.Contains(dispname, "腹") {
		views = append(views, "A")
	}

	if strings.Contains(dispname, "VC") {
		views = append(views, "VC")
	}

	return views
}
func MediaAutoGenKeywords(dispname string) []string {
	keywords := make([]string, 0)
	if strings.Contains(dispname, "异常") {
		keywords = append(keywords, "异常")
	}
	if strings.Contains(dispname, "正常") {
		keywords = append(keywords, "正常")
	}
	return keywords
}
func MediaAutoGenFcode(dispname string) string {
	for _, v := range strings.Split(dispname, " ") {
		for _, w := range strings.Split(v, "-") {
			if w == "" {
				continue
			}

			if strings.ToLower(w[0:1]) == "f" {
				return w
			}
		}
	}
	return ""
}
func MediaDeleteLabelAll(mid int) error {
	mi, err := database.MediaGet(mid)
	if err != nil {
		return err
	}
	if err = database.MediaRemoveLabel(mid); err != nil {
		return err
	}
	return database.LabelDelete(mi.Hash)
}

func MediaSetLabelAuthorJson(mid int, jsonstr string, authorUid int) error {
	mi, err := database.MediaGet(mid)
	if err != nil {
		return err
	}

	var lid int
	li, err := database.LabelGet(mi.Hash)
	if err != nil {
		return err
	}
	if li.Data == jsonstr {
		return errors.New("same data")
	}
	err = database.LabelUpdateJsonDataOnly(li.Lid, jsonstr)
	if err != nil {
		return err
	}

	return database.MediaUpdateLabelAuthorUidLid(mid, authorUid, lid)
}
func MediaGetAll() map[int][]string {
	mis, _ := database.MediaGetAll()

	data := make(map[int][]string, 0)
	for _, v := range mis {
		data[v.Mid] = []string{v.Hash, v.DisplayName, v.Path}
	}
	return data
}
