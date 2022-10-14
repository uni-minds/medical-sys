package module

import (
	"gitee.com/uni-minds/bridge-pacs/global"
	"gitee.com/uni-minds/bridge-pacs/pacs_server"
	"gitee.com/uni-minds/medical-sys/database"
)

var pi *pacs_server.PacsServer

func PacsGetInstanceThumb(id string) (bs []byte, err error) {
	pi = database.BridgeGetPacsServerHandler()

	return pi.InstanceMediaGetThumb(id)
}

func PacsGetInstanceMedia(id string, targetType string) (bs []byte, mediaType string, err error) {
	pi = database.BridgeGetPacsServerHandler()
	switch targetType {
	case "image":
		return pi.InstanceMediaGet(id, "image")
	case "mp4":
		return pi.InstanceMediaGet(id, "mp4")
	default:
		return pi.InstanceMediaGet(id, "ogv")
	}
}

func PacsGetStudiesInfo(id string) (info global.StudiesInfo, err error) {
	pi = database.BridgeGetPacsServerHandler()
	return pi.FindStudiesById(id)
}

func PacsGetSeriesInfo(id string) (info global.SeriesInfo, err error) {
	pi = database.BridgeGetPacsServerHandler()
	return pi.FindSeriesByIdLocal(id)
}

func PacsSetSeriesMemo(id, memo string) error {
	pi = database.BridgeGetPacsServerHandler()
	return pi.SeriesUpdateLabelMemo(id, memo)
}

func PacsGetSeriesMemo(id string) (memo string, err error) {
	pi = database.BridgeGetPacsServerHandler()
	info, err := pi.FindSeriesByIdLocal(id)
	if err != nil {
		return "", err
	} else {
		return info.LabelMemo, nil
	}
}

func PacsGetStudiesInfoAuthor(id string) (uid int, err error) {
	info, err := PacsGetStudiesInfo(id)
	if err != nil {
		return 0, err
	}
	return info.LabelUidAuthor, nil
}

func PacsGetStudiesInfoReview(id string) (uid int, err error) {
	info, err := PacsGetStudiesInfo(id)
	if err != nil {
		return 0, err
	}
	return info.LabelUidReview, nil
}

func PacsGetStudiesIdFromSeriesId(seriesId string) (studiesId string, err error) {
	info, err := PacsGetSeriesInfo(seriesId)
	if err != nil {
		return "", err
	}

	return info.StudiesId, nil
}
