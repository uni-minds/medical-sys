package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"uni-minds.com/medical-sys/global"
)

func (*MediaInfo) TableName() string {
	return global.DefaultDatabaseMediaTable
}

type MediaInfo struct {
	Mid             int     `gorose:"mid"`
	DisplayName     string  `gorose:"displayname"`
	Path            string  `gorose:"path"`
	Hash            string  `gorose:"hash"`
	Duration        float64 `gorose:"duration"`
	Frames          int     `gorose:"frames"`
	Width           int     `gorose:"width"`
	Height          int     `gorose:"height"`
	Status          int     `gorose:"status"`
	UploadTime      string  `gorose:"uploadtime"`
	UploadUid       int     `gorose:"uploaduid"`
	PatientID       string  `gorose:"patientid"`
	MachineID       string  `gorose:"machineid"`
	FolderName      string  `gorose:"foldername"`
	Fcode           string  `gorose:"fcode"`
	IncludeViews    string  `gorose:"includeviews"`
	Keywords        string  `gorose:"keywords"`
	Memo            string  `gorose:"memo"`
	MediaType       string  `gorose:"mediatype"`
	MediaData       string  `gorose:"mediadata"`
	LabelAuthorsUid string  `gorose:"labelauthorsuid"`
	LabelAuthorsLid string  `gorose:"labelauthorslid"`
	LabelReviewsUid string  `gorose:"labelreviewsuid"`
	LabelReviewsLid string  `gorose:"labelreviewslid"`
}
type MediaInfoUltrasonicImage struct {
	Width  int
	Height int
}
type MediaInfoUltrasonicVideo struct {
	PathRaw  string
	HashRaw  string
	PathJpgs string
	Encoder  string
}

func initMediaDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	"mid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"displayname" TEXT NOT NULL default "",
	"path" TEXT NOT NULL default "",
	"hash" TEXT NOT NULL default "",
	"duration" NUMERIC NOT NULL default 0,
	"frames" INTERGER NOT NULL default 0,
	"status" INTERGET NOT NULL default 0,
	"uploadtime" TEXT NOT NULL default "2000-01-01 00:00:00",
	"uploaduid" INTERGER NOT NULL default 0,
	"patientid" TEXT NOT NULL default "",
	"machineid" TEXT NOT NULL default "",
	"foldername" TEXT NOT NULL default "",
	"fcode" TEXT NOT NULL default "",
	"includeviews" TEXT NOT NULL default "[]",
	"keywords" TEXT NOT NULL default "[]",
	"memo" TEXT NOT NULL default "",
	"mediatype" TEXT NOT NULL default "",
	"mediadata" TEXT NOT NULL default "{}",
	"labelauthorsuid" TEXT NOT NULL default "[]",
	"labelauthorslid" TEXT NOT NULL default "[]",  
	"labelreviewsuid" TEXT NOT NULL default "[]",
	"labelreviewslid" TEXT NOT NULL default "[]")`, global.DefaultDatabaseMediaTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Panic(err.Error())
	}
}

func MediaCreate(mi MediaInfo) (mid int, err error) {
	if mi.Hash == "" {
		return 0, errors.New(global.EMediaRawHashNull)
	}
	mt, err := MediaGet(mi.Hash)
	if err == nil {
		return mt.Mid, errors.New(global.EMediaAlreadyExisted)
	}

	mi.Mid = 0
	// Upgrade override
	// mi.UploadTime = time.Now().Format(global.TimeFormat)
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Data(mi).Insert()

	mi, err = MediaGet(mi.Hash)
	return mi.Mid, err
}

func MediaGet(i interface{}) (mi MediaInfo, err error) {
	switch i.(type) {
	case int:
		err = DB().Table(&mi).Where("mid", "=", i).Select()
		if mi.Mid == 0 {
			if err != nil {
				log.Println("E DB", err.Error())
			}
			err = errors.New(global.EMediaNotExist)
		}
		return

	case string:
		key := i.(string)
		err = DB().Table(&mi).Where("hash", "=", key).Select()
		if mi.Mid == 0 {
			if err != nil {
				log.Println("E DB", err.Error())
			}
			err = errors.New(global.EMediaNotExist)
		}
	}
	return
}
func MediaGetAll() (ml []MediaInfo, err error) {
	err = DB().Table(&ml).OrderBy("mid").Select()
	return
}

func mediaUpdate(mid int, data interface{}) (err error) {
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Data(data).Where("mid", "=", mid).Update()
	return
}
func MediaUpdateDetail(mid int, mediaData interface{}) error {
	_, err := MediaGet(mid)
	if err != nil {
		return err
	}

	mediaT := ""
	jb, err := json.Marshal(mediaData)
	if err != nil {
		return err
	}
	mediaD := string(jb)

	switch mediaData.(type) {
	case MediaInfoUltrasonicImage:
		mediaT = global.MediaTypeUltrasonicImage

	case MediaInfoUltrasonicVideo:
		mediaT = global.MediaTypeUltrasonicVideo

	default:
		return errors.New(global.EMediaUnknownType)
	}

	data := map[string]interface{}{"mediatype": mediaT, "mediadata": mediaD}
	return mediaUpdate(mid, data)
}

func MediaUpdatePath(mid int, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	data := map[string]interface{}{"path": path}
	return mediaUpdate(mid, data)
}
func MediaUpdateFolderName(mid int, foldername string) error {
	data := map[string]interface{}{"foldername": foldername}
	return mediaUpdate(mid, data)
}
func MediaGetDetail(mid int) (mediaData interface{}, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}

	switch mi.MediaType {
	case global.MediaTypeUltrasonicVideo:
		var data MediaInfoUltrasonicVideo
		err = json.Unmarshal([]byte(mi.MediaData), &data)
		return data, err

	case global.MediaTypeUltrasonicImage:
		var data MediaInfoUltrasonicImage
		err = json.Unmarshal([]byte(mi.MediaData), &data)
		return data, err

	default:
		return mediaData, errors.New(global.EMediaUnknownType)
	}
}

// Frames and duration
func MediaUpdateFramesAndDuration(mid int, frames int, duration float64) error {
	data := map[string]interface{}{"frames": frames, "duration": duration}
	return mediaUpdate(mid, data)
}
func MediaUpdateWidthAndHeight(mid, width, height int) error {
	data := map[string]interface{}{"width": width, "height": height}
	return mediaUpdate(mid, data)
}
func MediaUpdateHash(mid int, hash string) error {
	data := map[string]interface{}{"hash": hash}
	return mediaUpdate(mid, data)
}

// Views
func MediaGetViews(mid int) (views []string, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}
	switch mi.IncludeViews {
	case "", "[]", "null":
		return make([]string, 0), nil
	default:
		err = json.Unmarshal([]byte(mi.Keywords), &views)
		return
	}
}
func MediaAddView(mid int, keyword string) error {
	views, err := MediaGetViews(mid)
	if err != nil {
		return err
	}
	for _, v := range views {
		if v == keyword {
			return nil
		}
	}
	views = append(views, keyword)
	return MediaSetViews(mid, views)
}
func MediaSetViews(mid int, views []string) error {
	if views == nil {
		views = make([]string, 0)
	}
	jb, err := json.Marshal(views)
	if err != nil {
		return err
	}
	return MediaUpdateViews(mid, string(jb))
}
func MediaUpdateViews(mid int, viewsStr string) error {
	data := map[string]interface{}{"includeviews": viewsStr}
	return mediaUpdate(mid, data)
}
func MediaRemoveView(mid int, view string) error {
	views, err := MediaGetViews(mid)
	if err != nil {
		return err
	}

	for i, v := range views {
		if v == view {
			views = append(views[:i], views[i+1:]...)
			return MediaSetViews(mid, views)
		}
	}
	return nil
}

// Keywords
func MediaGetKeywords(mid int) (keywords []string, err error) {
	mi, err := MediaGet(mid)
	if err != nil {
		return
	}

	if mi.Keywords == "" || mi.Keywords == "[]" {
		return make([]string, 0), nil
	}

	err = json.Unmarshal([]byte(mi.Keywords), &keywords)
	return
}
func MediaAddKeyword(mid int, keyword string) error {
	keywords, err := MediaGetKeywords(mid)
	if err != nil {
		return err
	}

	for _, v := range keywords {
		if v == keyword {
			return nil
		}
	}

	keywords = append(keywords, keyword)
	return MediaSetKeywords(mid, keywords)
}
func MediaSetKeywords(mid int, keywords []string) error {
	if keywords == nil {
		keywords = make([]string, 0)
	}

	jb, err := json.Marshal(keywords)
	if err != nil {
		return err
	}
	return MediaUpdateKeywords(mid, string(jb))
}
func MediaUpdateKeywords(mid int, keywordStr string) error {
	data := map[string]interface{}{"keywords": keywordStr}
	return mediaUpdate(mid, data)
}
func MediaRemoveKeyword(mid int, keyword string) error {
	keywords, err := MediaGetKeywords(mid)
	if err != nil {
		return err
	}

	for i, v := range keywords {
		if v == keyword {
			keywords = append(keywords[:i], keywords[i+1:]...)
			return MediaSetKeywords(mid, keywords)
		}
	}
	return nil
}

// Memo
func MediaUpdateMemo(mid int, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return mediaUpdate(mid, data)
}

// DisplayName
func MediaUpdateDisplayName(mid int, name string) error {
	data := map[string]interface{}{"displayname": name}
	return mediaUpdate(mid, data)
}

func MediaUpdateLabelAuthorUidLid(mid int, uidstr, lidstr string) error {
	data := map[string]interface{}{"labelauthorsuid": uidstr, "labelauthorslid": lidstr}
	return mediaUpdate(mid, data)
}

func MediaUpdateLabelReviewUidLid(mid int, uidstr, lidstr string) error {
	data := map[string]interface{}{"labelreviewsuid": uidstr, "labelreviewslid": lidstr}
	return mediaUpdate(mid, data)
}

func MediaDelete(mid int) (err error) {
	_, err = DB().Table(global.DefaultDatabaseMediaTable).Where("mid", "=", mid).Delete()
	return
}
