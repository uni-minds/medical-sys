package database

import "testing"

func TestMediaCreate(t *testing.T) {
	mid, err := MediaCreate(MediaInfo{
		DisplayName: "Ala1",
		Hash:        "AABB",
		UploadUid:   1,
		Memo:        "m1",
	})
	t.Log(mid, err)
	mid, err = MediaCreate(MediaInfo{
		DisplayName: "Ala2",
		Hash:        "AACC",
		UploadUid:   1,
		Memo:        "m1",
	})
	t.Log(mid, err)
	mid, err = MediaCreate(MediaInfo{
		DisplayName: "Ala3",
		Hash:        "AABB",
		UploadUid:   1,
		Memo:        "m1",
	})
	t.Log(mid, err)
}

func TestMediaGet(t *testing.T) {
	t.Log(MediaGet("AABB"))
	t.Log(MediaGet(1))
}

func TestMediaGetAll(t *testing.T) {
	t.Log(MediaGetAll())
}

func TestMediaUpdateMemo(t *testing.T) {
	MediaUpdateMemo(1, "RB")
	t.Log(MediaGet(1))
}

func TestMediaDelete(t *testing.T) {
	MediaDelete(1)
	t.Log(MediaGetAll())
}

func TestMediaUpdateDetail(t *testing.T) {
	MediaUpdateDetail(2, MediaInfoUltrasonicVideo{
		Width:   2,
		Height:  3,
		HashRaw: "ffaa",
		Frames:  10,
	})
}

func TestMediaGetDetail(t *testing.T) {
	mdi, _ := MediaGetDetail(2)
	c := mdi.(MediaInfoUltrasonicVideo)
	t.Log(c.Width)
}
