package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gohouse/t"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/tools"
)

const PACS_AET = "DCM4CHEE"
const PACS_SOP_IMAGE = "1.2.840.10008.5.1.4.1.1.6.1"
const PACS_SOP_MULTI_FRAME = "1.2.840.10008.5.1.4.1.1.3.1"
const PACS_SOP_SECONDARY_SCREEN = "1.2.840.10008.5.1.4.1.1.7"

type PacsAnswerCount struct {
	Count int
}

type PacsKeyValuePair struct {
	Vr    string
	Value []interface{}
}

type PacsRespContents []map[string]PacsKeyValuePair

type PathCache struct {
	LocalCache string
	LocalThumb string
	GenDate    time.Time
}

func PacsSync(addr string) {
	var err error
	var respCount PacsAnswerCount
	var bs []byte
	var studiesContents PacsRespContents

	log("w", "sync ->", addr)
	u := url.URL{
		Scheme:   "http",
		Host:     addr,
		Path:     fmt.Sprintf("dcm4chee-arc/aets/%s/rs/studies/count", PACS_AET),
		RawQuery: "includefield=all",
	}

	if _, bs, err = tools.HttpGet(u.String()); err != nil {
		log("e", "E1", err.Error())
		return
	} else if err = json.Unmarshal(bs, &respCount); err != nil {
		log("e", "E3", err.Error())
	}

	fmt.Println(respCount.Count)

	u.Path = fmt.Sprintf("dcm4chee-arc/aets/%s/rs/studies", PACS_AET)
	u.RawQuery = "includefield=all&offset=0"

	if _, bs, err = tools.HttpGet(u.String()); err != nil {
		log("e", "E1", err.Error())
		return
	} else if err = json.Unmarshal(bs, &studiesContents); err != nil {
		log("e", "E3", err.Error())
	}

	for i, content := range studiesContents {
		log("t", "Parsing index", i)
		data := content["0020000D"].Value[0]
		studiesID := data.(string)
		for _, series := range PacsGetDataStudies(addr, studiesID) {
			data := series["0020000E"].Value[0]
			seriesID := data.(string)
			respSeries := PacsGetDataSeries(addr, studiesID, seriesID)
			for _, instanceData := range respSeries {
				detail, err := PacsAnalysisInstanceData(instanceData)
				if err != nil {
					log("E", "import E3", err.Error())
					continue
				}

				detail.Server = addr
				detail.AET = PACS_AET

				database.PacsInstanceCreate(detail)
			}
		}
	}
}

func PacsGetDataStudies(addr, studiesID string) (resp PacsRespContents) {
	var err error
	var bs []byte

	u := url.URL{
		Scheme:   "http",
		Opaque:   "",
		User:     nil,
		Host:     addr,
		Path:     fmt.Sprintf("dcm4chee-arc/aets/%s/rs/studies/%s/series", PACS_AET, studiesID),
		RawQuery: "includefield=all&offset=0&orderby=SeriesNumber",
	}

	if _, bs, err = tools.HttpGet(u.String()); err != nil {
		log("e", "E1", err.Error())
		return
	} else if err = json.Unmarshal(bs, &resp); err != nil {
		log("e", "E2", err.Error())
		return
	}

	return resp
}

func PacsGetDataSeries(addr, studiesID, seriesID string) (resp PacsRespContents) {
	var err error
	var bs []byte

	u := url.URL{
		Scheme:   "http",
		Opaque:   "",
		User:     nil,
		Host:     addr,
		Path:     fmt.Sprintf("dcm4chee-arc/aets/%s/rs/studies/%s/series/%s/instances", PACS_AET, studiesID, seriesID),
		RawQuery: "includefield=all&offset=0&orderby=InstanceNumber",
	}

	if _, bs, err = tools.HttpGet(u.String()); err != nil {
		log("e", "E1", err.Error())
		return
	} else if err = json.Unmarshal(bs, &resp); err != nil {
		log("e", "E2", err.Error())
		return
	}

	return resp
}

func PacsAnalysisInstanceData(data map[string]PacsKeyValuePair) (resp database.PacsInstanceDetail, err error) {
	var width, height int
	var studiesId, seriesId, instanceId, instanceType string

	if k, err := pacsGetKeyValue(data, "00280011", 0); err != nil {
		return resp, err
	} else {
		width = t.ParseInt(k.(float64))
	}

	if k, err := pacsGetKeyValue(data, "00280010", 0); err != nil {
		return resp, err
	} else {
		height = t.ParseInt(k.(float64))
	}

	if k, err := pacsGetKeyValue(data, "0020000D", 0); err != nil {
		return resp, err
	} else {
		studiesId = k.(string)
	}

	if k, err := pacsGetKeyValue(data, "0020000E", 0); err != nil {
		return resp, err
	} else {
		seriesId = k.(string)
	}

	if k, err := pacsGetKeyValue(data, "00080018", 0); err != nil {
		return resp, err
	} else {
		instanceId = k.(string)
	}

	if k, err := pacsGetKeyValue(data, "00080016", 0); err != nil {
		return resp, err
	} else {
		instanceType = k.(string)
	}

	resp = database.PacsInstanceDetail{
		Server:       "",
		AET:          "",
		StudiesID:    studiesId,
		SeriesID:     seriesId,
		InstanceID:   instanceId,
		Frames:       0,
		Width:        width,
		Height:       height,
		InstanceType: instanceType,
	}

	switch resp.InstanceType {
	case PACS_SOP_MULTI_FRAME:
		if k, err := pacsGetKeyValue(data, "00280008", 0); err != nil {
			return resp, err
		} else if f, err := strconv.Atoi(k.(string)); err != nil {
			return resp, err
		} else {
			resp.Frames = f
		}

	case PACS_SOP_IMAGE, PACS_SOP_SECONDARY_SCREEN:

	default:
		log("E", "Unknown type")
	}
	return resp, nil
}

func PacsGetInstanceRenderedFrame(data database.InfoInstance, frame int, tp string) (bs []byte, err error) {
	var u string

	if data.Frames == 0 {
		u = fmt.Sprintf("http://%s/dcm4chee-arc/aets/%s/rs/studies/%s/series/%s/instances/%s/rendered", data.Server, data.AET, data.StudiesID, data.SeriesID, data.InstanceID)
	} else if frame < 0 {
		u = fmt.Sprintf("http://%s/dcm4chee-arc/aets/%s/rs/studies/%s/series/%s/instances/%s/rendered", data.Server, data.AET, data.StudiesID, data.SeriesID, data.InstanceID)
	} else if frame <= data.Frames {
		u = fmt.Sprintf("http://%s/dcm4chee-arc/aets/%s/rs/studies/%s/series/%s/instances/%s/frames/%d/rendered", data.Server, data.AET, data.StudiesID, data.SeriesID, data.InstanceID, frame+1)
	} else {
		return nil, errors.New("over frame")
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	switch tp {
	case "png", "PNG":
		req.Header.Set("Accept", "image/png")
	case "jpg", "jpeg":
		req.Header.Set("Accept", "image/jpeg")
	case "mp4":
		log("e", "dcm4chee not support")
		//req.Header.Set("Accept", "video/mp4")
	case "gif":
		req.Header.Set("Accept", "image/gif")

	}

	resp, err := client.Do(req)
	if err != nil {
		log("e", "client", err.Error())
		return nil, err
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return bs, err
	}
}

func PacsImportStudiesAll(gid int) (n int, err error) {
	studiesIds, err := database.PacsGetAllStudiesIds()
	if err != nil {
		return 0, err
	}

	n = 0
	for _, studiesId := range studiesIds {
		if err = database.GroupAddPacsStudiesId(gid, studiesId); err != nil {
			continue
		}
		n++
	}
	return n, nil
}

func PacsSplitStudiesToGroup(srcGid int, destGids []int) (err error) {
	studiesIds, err := database.GroupGetPacsStudiesIds(srcGid)
	if err != nil {
		return err
	}

	gCurrent := 0
	for _, studiesId := range studiesIds {
		if gCurrent >= len(destGids) {
			gCurrent = 0
		}

		if err := database.GroupAddPacsStudiesId(destGids[gCurrent], studiesId); err != nil {
			log("e", err.Error())
			continue
		}
		gCurrent++
	}
	return nil
}

func PacsGetInstanceThumb(instanceId string) (data []byte, err error) {
	info, err := database.PacsInstanceGetInfo(instanceId)
	if err != nil {
		return nil, err
	}

	var pathCache PathCache
	if info.PathCache != "" && json.Unmarshal([]byte(info.PathCache), &pathCache) == nil && pathCache.LocalThumb != "" {
		if bs, err := ioutil.ReadFile(pathCache.LocalThumb); err == nil {
			return bs, nil
		}
	}

	_, data, err = PacsGetInstanceImage(instanceId, true)
	return data, err
}

func PacsGetInstanceImage(instanceId string, forceRegen bool) (bsRaw, bsThumb []byte, err error) {
	thumbSize := 300
	var imgRaw, imgThumb image.Image
	var pathCache PathCache
	var info database.InfoInstance

	if info, err = database.PacsInstanceGetInfo(instanceId); err != nil {
		return nil, nil, err
	}

	if !forceRegen && info.PathCache != "" && json.Unmarshal([]byte(info.PathCache), &pathCache) == nil && pathCache.LocalCache != "" {
		if bsRaw, err = ioutil.ReadFile(pathCache.LocalCache); err == nil {
			return bsRaw, nil, err
		}
	}

	if bsRaw, err = PacsGetInstanceRenderedFrame(info, 0, "png"); err != nil {
		return nil, nil, err
	}

	cacheRoot := path.Join(global.GetAppSettings().SystemMediaPath, "pacs_cache", info.StudiesID, info.SeriesID)
	if _, err = os.Stat(cacheRoot); err != nil {
		os.MkdirAll(cacheRoot, os.ModePerm)
	}

	// 生成缩略图
	r := bytes.NewReader(bsRaw)
	buf := new(bytes.Buffer)
	if imgRaw, _, err = image.Decode(r); err != nil {
		return nil, nil, err
	}
	s := imgRaw.Bounds().Size()
	width, height := tools.CalcResize(s.X, s.Y, thumbSize, thumbSize)
	imgThumb = imaging.Thumbnail(imgRaw, width, height, imaging.Linear)
	err = jpeg.Encode(buf, imgThumb, nil)
	bsThumb = buf.Bytes()

	// 缓存原图/视频
	if info.Frames > 0 {
		pathCache = PathCache{
			LocalCache: path.Join(cacheRoot, fmt.Sprintf("%s.gif", info.InstanceID)),
			LocalThumb: path.Join(cacheRoot, fmt.Sprintf("%s_thumb_%d.jpg", info.InstanceID, thumbSize)),
			GenDate:    time.Now(),
		}

		go func() {
			bsRaw, err = PacsGetInstanceRenderedFrame(info, -1, "gif")
			if err != nil {
				log("e", err.Error())
			} else {
				SaveMedia(pathCache.LocalCache, bsRaw)
			}
		}()

	} else {
		pathCache = PathCache{
			LocalCache: path.Join(cacheRoot, fmt.Sprintf("%s.png", info.InstanceID)),
			LocalThumb: path.Join(cacheRoot, fmt.Sprintf("%s_thumb_%d.jpg", info.InstanceID, thumbSize)),
			GenDate:    time.Now(),
		}
		go SaveMedia(pathCache.LocalCache, bsRaw)

	}
	go SaveMedia(pathCache.LocalThumb, bsThumb)

	jbs, _ := json.Marshal(pathCache)
	info.PathCache = string(jbs)

	database.PacsInstanceUpdate(info)

	return bsRaw, bsThumb, nil
}

func PacsSeriesCheckAuthor(seriesId string, uid int) error {
	info, err := database.PacsSeriesGetInfo(seriesId)
	if err != nil {
		return err
	}

	if info.LabelAuthorUid == 0 || info.LabelAuthorUid == uid {
		return nil
	} else {
		return errors.New(fmt.Sprintf("数据锁定为其他标注用户：%d", info.LabelAuthorUid))
	}
}

func PacsSeriesCheckReview(seriesId string, uid int) error {
	info, err := database.PacsSeriesGetInfo(seriesId)
	if err != nil {
		return err
	}

	if info.LabelReviewUid == 0 || info.LabelReviewUid == uid {
		return nil
	} else {
		return errors.New(fmt.Sprintf("数据锁定为其他审核用户：%d", info.LabelReviewUid))
	}
}

func PacsSetSeriesAuthorLabel(seriesId string, uid int, progress int) error {
	if err := PacsSeriesCheckAuthor(seriesId, uid); err != nil {
		return err
	} else {
		return database.PacsSeriesUpdateAuthor(seriesId, uid, progress)
	}
}

func PacsInstanceUpdateLabel(seriesId, instanceId string, uid int, view, diagnose, interfere string) error {
	err := PacsSeriesCheckAuthor(seriesId, uid)
	if err != nil {
		return err
	} else {
		return database.PacsInstanceUpdateLabel(instanceId, view, diagnose, interfere)
	}
}

func SaveMedia(fn string, data []byte) (err error) {
	fp, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log("e", err.Error())
	}
	_, err = fp.Write(data)
	if err != nil {
		log("e", err.Error())
	}
	fp.Close()
	return nil
}
func pacsGetKeyValue(data map[string]PacsKeyValuePair, key string, index int) (value interface{}, err error) {
	for k, v := range data {
		if k == key {
			vals := v.Value
			if len(vals) > index {
				return vals[index], nil
			} else {
				return nil, errors.New("value not found")
			}
		}
	}
	return nil, errors.New("key not found")
}
