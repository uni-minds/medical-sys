/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: class_media.go
 * Date: 2022/09/25 17:43:25
 */

package database

import (
	"encoding/json"
	"fmt"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/utils/tools"
	"os"
	"runtime"
	"strings"
)

var dbMedia *DbMedia

type DbMedia struct {
	TableName  string
	ModuleName string
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
type MediaMetadata struct {
	Folder string
}

func GetMedia() *DbMedia {
	if dbMedia == nil {
		dbMedia = CreateMediaDB(os.Stdout)
	}
	return dbMedia
}

func CreateMediaDB(fp *os.File) *DbMedia {
	db := new(DbMedia)
	db.TableName = global.DefaultDatabaseMediaTable
	db.ModuleName = "DB_M"
	return db
}

func InitMediaDB() {
	dbSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
    id             INTEGER                   not null
        primary key autoincrement,
    media_uuid     TEXT    default ''        not null,
    display_name   TEXT    default ''        not null,
    path           TEXT    default ''        not null,
    width          INTEGER default 0         not null,
    height         INTEGER default 0         not null,
    duration       REAL    default 0         not null,
    frames         INTEGER default 0         not null,
    fps            INTEGER default 0         not null,
    media_type     TEXT    default 'default' not null,
    upload_uid     INTEGER default 0         not null,
    upload_time    TEXT    default ''        not null,
    memo           TEXT    default ''        not null,
    patient_id     TEXT    default ''        not null,
    machine_id     TEXT    default ''        not null,
    metadata       TEXT    default ''        not null,
    crf_define     TEXT    default ''        not null,
    keywords       TEXT    default ''        not null,
    media_data     TEXT    default '{}'      not null,
    label_author   TEXT    default ''        not null,
    label_reviewer TEXT    default ''        not null,
    label_progress INTEGER default 0         not null,
    cowork         TEXT    default 'single'  not null,
    media_hash     TEXT    default ''        not null)`, global.DefaultDatabaseMediaTable)

	_, err := DB().Execute(dbSql)
	if err != nil {
		log.Error(err.Error())
	}
}

func (p *DbMedia) Trace(msg ...any) {
	p.Write("t", tools.InterfaceExpand(msg))
}
func (p *DbMedia) Warn(msg ...any) {
	p.Write("w", tools.InterfaceExpand(msg))
}
func (p *DbMedia) Error(msg ...any) {
	p.Write("e", tools.InterfaceExpand(msg))
}
func (p *DbMedia) Write(level string, msg string) {
	_, file, line, ok := runtime.Caller(2)
	if ok && global.FlagGetDebug() {
		logger.Write(p.ModuleName, level, fmt.Sprintf("%s:%d %s", file, line, msg))
	} else {
		logger.Write(p.ModuleName, level, msg)
	}
}

func (p *DbMedia) Get(i interface{}) (info DbStructMedia, err error) {
	switch i.(type) {
	case int:
		if err = DB().Table(&info).Where("id", i).Select(); err != nil {
			return info, err
		}

	case string:
		if err = DB().Table(&info).Where("media_uuid", i).Select(); err != nil {
			return info, err
		}
	}

	if info.Id < 1 {
		return info, fmt.Errorf(global.EMediaNotExist)
	} else {
		return info, nil
	}
}
func (p *DbMedia) GetAll() (infos []DbStructMedia, err error) {
	err = DB().Table(&infos).OrderBy("id").Select()
	return infos, err
}

func (p *DbMedia) Create(info DbStructMedia) (mediaUUID string, err error) {
	mediaUUID = info.MediaUUID

	if mediaUUID == "" {
		mediaUUID = tools.RandString0f(32)
	}

	for {
		_, err = p.Get(mediaUUID)
		if err != nil {
			break
		}
		mediaUUID = tools.RandString0f(32)
	}

	info.Id = 0
	info.MediaUUID = mediaUUID

	_, err = DB().Table(p.TableName).Insert(info)

	info, err = p.Get(info.MediaUUID)
	return info.MediaUUID, err
}

func (p *DbMedia) UpdateManual(i interface{}, data map[string]interface{}) error {
	info, err := p.Get(i)
	if err != nil {
		return err
	}

	_, err = DB().Table(p.TableName).Data(data).Where("id", info.Id).Update()
	return err
}

func (p *DbMedia) UpdateHash(i interface{}, hash string) error {
	data := map[string]interface{}{"media_hash": hash}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) UpdateDetail(i interface{}, mediaData interface{}) error {
	if _, err := p.Get(i); err != nil {
		return err
	}

	mediaType := ""

	switch mediaData.(type) {
	case MediaInfoUltrasonicImage:
		mediaType = global.MediaTypeUltrasonicImage

	case MediaInfoUltrasonicVideo:
		mediaType = global.MediaTypeUltrasonicVideo

	default:
		return fmt.Errorf(global.EMediaUnknownType)
	}

	jbs, _ := json.Marshal(mediaData)

	data := map[string]interface{}{"media_type": mediaType, "media_data": string(jbs)}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) UpdateMetadata(i interface{}, metadata MediaMetadata) error {
	jbs, _ := json.Marshal(metadata)
	data := map[string]interface{}{"metadata": string(jbs)}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) UpdateAuthorReview(i interface{}, authorUid, reviewUid int) error {
	data := map[string]interface{}{"label_author": authorUid, "label_reviewer": reviewUid}
	return p.UpdateManual(i, data)
}

func (p *DbMedia) GetDetail(i interface{}) (interface{}, error) {
	info, err := p.Get(i)
	if err != nil {
		return nil, err
	}

	switch info.MediaType {
	case global.MediaTypeUltrasonicVideo:
		var data MediaInfoUltrasonicVideo
		err = json.Unmarshal([]byte(info.MediaData), &data)
		return data, err

	case global.MediaTypeUltrasonicImage:
		var data MediaInfoUltrasonicImage
		err = json.Unmarshal([]byte(info.MediaData), &data)
		return data, err

	default:
		return nil, fmt.Errorf(global.EMediaUnknownType)
	}
}

// Frames and duration
func (p *DbMedia) UpdateFramesAndDuration(i interface{}, fps float64, duration float64, frames int) error {
	data := map[string]interface{}{"frames": frames, "fps": fps, "duration": duration}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) UpdateWidthAndHeight(i interface{}, width, height int) error {
	data := map[string]interface{}{"width": width, "height": height}
	return p.UpdateManual(i, data)
}

// Views
func (p *DbMedia) SetCrfViews(i interface{}, views []string) error {
	return p.UpdateViews(i, strings.Join(views, ","))
}
func (p *DbMedia) GetCrfView(i interface{}) (views []string, err error) {
	info, err := p.Get(i)
	if err != nil {
		return
	}

	switch info.CrfDefine {
	case "", "[]", "null":
		return make([]string, 0), nil
	default:
		return strings.Split(info.CrfDefine, ","), nil
	}
}
func (p *DbMedia) AddCrfView(i interface{}, view string) error {
	views, err := p.GetCrfView(i)
	if err != nil {
		return err
	}
	for _, v := range views {
		if v == view {
			return nil
		}
	}
	views = append(views, view)
	return p.SetCrfViews(i, views)
}
func (p *DbMedia) UpdateViews(i interface{}, viewsStr string) error {
	data := map[string]interface{}{"crf_define": viewsStr}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) RemoveView(i interface{}, view string) error {
	views, err := p.GetCrfView(i)
	if err != nil {
		return err
	}

	for i, v := range views {
		if v == view {
			views = append(views[:i], views[i+1:]...)
			return p.SetCrfViews(i, views)
		}
	}
	return nil
}

// Keywords
func (p *DbMedia) SetKeywords(i interface{}, keywords []string) error {
	if keywords == nil {
		keywords = make([]string, 0)
	}

	jbs, _ := json.Marshal(keywords)
	return p.UpdateKeywords(i, string(jbs))
}
func (p *DbMedia) GetKeywords(i interface{}) (keywords []string, err error) {
	info, err := p.Get(i)
	if err != nil {
		return
	}

	switch info.Keywords {
	case "", "[]", "null":
		return make([]string, 0), nil
	default:
		err = json.Unmarshal([]byte(info.Keywords), &keywords)
		return keywords, err
	}
}
func (p *DbMedia) AddKeyword(i interface{}, keyword string) error {
	keywords, err := p.GetKeywords(i)
	if err != nil {
		return err
	}

	for _, v := range keywords {
		if v == keyword {
			return nil
		}
	}

	keywords = append(keywords, keyword)
	return p.SetKeywords(i, keywords)
}
func (p *DbMedia) UpdateKeywords(i interface{}, keywordStr string) error {
	data := map[string]interface{}{"keywords": keywordStr}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) RemoveKeyword(i interface{}, keyword string) error {
	keywords, err := p.GetKeywords(i)
	if err != nil {
		return err
	}

	for i, v := range keywords {
		if v == keyword {
			keywords = append(keywords[:i], keywords[i+1:]...)
			return p.SetKeywords(i, keywords)
		}
	}
	return nil
}

// Memo
func (p *DbMedia) UpdateMemo(i interface{}, memo string) error {
	data := map[string]interface{}{"memo": memo}
	return p.UpdateManual(i, data)
}

// Label
func (p *DbMedia) LabelUpdateAuthor(i interface{}, authorData string) error {
	data := map[string]interface{}{"label_author": authorData}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) LabelUpdateReview(i interface{}, reviewData string) error {
	data := map[string]interface{}{"label_reviewer": reviewData}
	return p.UpdateManual(i, data)
}
func (p *DbMedia) LabelUpdateProgress(i interface{}, progress int) error {
	data := map[string]interface{}{"label_progress": progress}
	return p.UpdateManual(i, data)
}

func (p *DbMedia) Selector(mediaUUIDs []string, field string, asc bool) (infos []DbStructMedia, err error) {
	uuids := make([]interface{}, 0)
	for _, id := range mediaUUIDs {
		uuids = append(uuids, id)
	}

	db := DB().Table(&infos).WhereIn("media_uuid", uuids)
	if field != "" {
		if asc {
			field = fmt.Sprintf("%s asc", field)
		} else {
			field = fmt.Sprintf("%s desc", field)
		}
		db.Order(field)
	}

	err = db.Select()
	return infos, err
}

func (p *DbMedia) Delete(i interface{}) error {
	info, err := p.Get(i)
	if err != nil {
		return err
	}
	_, err = DB().Table(p.TableName).Where("media_uuid", info.MediaUUID).Delete()
	return err
}
