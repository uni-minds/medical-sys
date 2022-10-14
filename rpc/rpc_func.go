package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/global"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/medical-sys/tools"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parsePacs(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		str.WriteString("pacs export view 4ap to group GN1\n" +
			"pacs list view 4CV group GN1\n" +
			"pacs load from GN1 to GN2 ...")
		return "", errors.New(str.String())
	}

	switch p[0] {
	// pacs list view 3V
	case "list":
		var srcView, srcGroup string
		for count := 1; count < len(p); count++ {
			switch p[count] {
			case "view":
				srcView = p[count+1]
				count += 1
			case "group":
				srcGroup = p[count+1]
				count += 1
			}
		}

		str.WriteString(fmt.Sprintf("-- pacs list view=%s group=%s\n", srcView, srcGroup))

		ps := database.BridgeGetPacsServerHandler()
		infos, err := ps.FindInstancesByView(srcView)
		if err != nil {
			return "", err
		}

		for i, v := range infos {
			str.WriteString(fmt.Sprintln(i, v.InstanceId, v.LabelView))
		}

		str.WriteString(fmt.Sprintf("%d records found.\n", len(infos)))

		return str.String(), nil
	// pacs export view ac group pacs_screen to group pacs_label_ac
	case "export":
		var srcView, srcGroup, dstGroup string
		var srcGid, dstGid int
		var destSelect bool
		var err error

		for count := 1; count < len(p); count++ {
			para := p[count]
			switch para {
			case "view":
				srcView = p[count+1]
				count += 1
			case "to":
				destSelect = true
				continue
			case "group":
				if destSelect {
					dstGroup = p[count+1]
				} else {
					srcGroup = p[count+1]
				}
				count += 1
			}
		}

		ps := database.BridgeGetPacsServerHandler()
		var gi database.GroupInfo

		if srcGid, err = strconv.Atoi(srcGroup); err != nil {
			gi, err = database.GroupGet(srcGroup)
			if err != nil {
				return "", errors.New(fmt.Sprintf("Unknown group name: %s", err.Error()))
			} else {
				srcGid = gi.Gid
			}
		}

		if dstGid, err = strconv.Atoi(dstGroup); err != nil {
			gi, err = database.GroupGet(dstGroup)
			if err != nil {
				return "", errors.New(fmt.Sprintln("Unknown group name:", err.Error()))
			} else {
				dstGid = gi.Gid
			}
		}

		str.WriteString(fmt.Sprintf("gids %d => %d\npacs %s => %s\n", srcGid, dstGid, srcView, dstGroup))

		infos, err := ps.FindInstancesByView(srcView)
		if err != nil {
			return "", err

		} else {
			for _, info := range infos {
				err = database.GroupAddContain(dstGid, info.InstanceId)
				if err != nil {
					fmt.Println("export e:", err.Error())
				}
			}
		}
		return str.String(), nil
	// pacs load from pacs_label_l to pacs_label_l_1 pacs_label_l_2 pacs_label_l_3 pacs_label_l_4
	case "load":
		if p[1] == "from" && p[3] == "to" {
			gMain := p[2]
			gSubs := p[4:]

			gMainId := module.GroupGetGid(gMain)
			if gMainId <= 0 {
				return "", errors.New(fmt.Sprintf("source group is not exist: %s", gMain))
			}

			gSubsId := make([]int, 0)
			for _, gname := range gSubs {
				gid := module.GroupGetGid(gname)
				if gid <= 0 {
					log("e", "subs group is not exist:", gname)
					continue
				}
				gSubsId = append(gSubsId, gid)
			}

			if len(gSubsId) == 0 {
				return "", errors.New("target group is not existed")
			}

			str.WriteString(fmt.Sprintf("src_group:\t%s [%d]\ndest_groups:\t%v [%v]\n", gMain, gMainId, gSubs, gSubsId))

			gSubsInstanceIds := make([]string, 0)
			groupInstanceIds := make(map[int][]string, 0)
			gImpoInstanceIds := module.GroupGetDicom(gMainId)

			// 检查是否目录组中已经导入对应id
			for _, gSubId := range gSubsId {
				containIds, _, err := database.GroupGetContains(gSubId)
				groupInstanceIds[gSubId] = containIds

				if err != nil {
					log("e", "group get contains:", err.Error())
					continue
				}
				gSubsInstanceIds = append(gSubsInstanceIds, containIds...)
			}

			log("i", "sub groups contain instance:", len(gSubsInstanceIds))
			log("i", "impo group contain instance:", len(gImpoInstanceIds))

			// 排除已经导入的id
			remainIds := tools.StringsExcept(gImpoInstanceIds, gSubsInstanceIds)

			log("i", "remain instances to import:", len(remainIds))

			// 轮询并分组
			gidx := 0
			for _, instanceId := range remainIds {
				if gidx == len(gSubsId) {
					gidx = 0
				}
				gsubId := gSubsId[gidx]
				groupInstanceIds[gsubId] = append(groupInstanceIds[gsubId], instanceId)
				gidx++
			}

			// 导入
			for gid, ids := range groupInstanceIds {
				err := database.GroupAddContains(gid, ids)
				if err != nil {
					log("e", "group import ids:", err.Error())
				} else {
					str.WriteString(fmt.Sprintf("group %d has %d instance(s) imported.\n", gid, len(ids)))
				}
			}
			str.WriteString("ok")
			return str.String(), nil
		}
		return "", errors.New("unrecognized format")

	default:
		return "", errors.New(fmt.Sprintf("unknown params: %s", p))
	}
}

func parseUser(p []string) (string, error) {
	str := strings.Builder{}
	if len(p) == 0 {
		str.WriteString("user create u1 p1 (description)\n" +
			"user remove u1 u2...\n" +
			"user activate u1\n" +
			"user password u1 p1\n" +
			"user merge u_old into u_new\n" +
			"user list")
		return "", errors.New(str.String())
	}

	switch p[0] {
	case "create":
		break

	case "remove":
		confirm := false
		var userlist []string
		trueCommand := strings.Builder{}
		if p[1] == "-y" {
			confirm = true
			userlist = p[2:]
		} else {
			trueCommand.WriteString("user remove -y ")
			userlist = p[1:]
		}

		for _, u := range userlist {
			uid := module.UserGetUid(u)
			str.WriteString(fmt.Sprintf("delete user: %d\n", uid))
			if confirm {
				err := database.UserDelete(uid)
				if err != nil {
					log("e", err.Error())
					return str.String(), err
				}
			} else {
				trueCommand.WriteString(fmt.Sprintf("%d ", uid))
			}
		}

		if !confirm {
			str.WriteString("\n!!! Dry mode.\nExecute database use command:\n" + trueCommand.String())
		}

		return str.String(), nil

	//case "merge":
	//	uo := p[1]
	//	un := p[3]
	//	uido := module.UserGetUid(uo)
	//	uidn := module.UserGetUid(un)
	//	if uido > 0 || uidn > 0 {
	//		var answer string
	//		fmt.Printf("DANGER!! User %s [UID=%d] will controll all label made by User %s [UID=%d]\n", un, uidn, uo, uido)
	//		fmt.Printf("Confirm? [Y/n]:")
	//		fmt.Scanln(&answer)
	//		if strings.ToLower(answer) == "y" {
	//			fmt.Println("User confirmed.")
	//			lis, _ := database.LabelGetAll()
	//			for _, li := range lis {
	//				change := false
	//				aid := li.AuthorUid
	//				rid := li.ReviewUid
	//				if aid == uido {
	//					aid = uidn
	//					change = true
	//				}
	//				if rid == uido {
	//					rid = uidn
	//					change = true
	//				}
	//				if change {
	//					li.AuthorUid = aid
	//					li.ReviewUid = rid
	//					database.LabelUpdate(li)
	//					fmt.Printf("Label owner: A<-%d, R<-%d\n", aid, rid)
	//				}
	//			}
	//			return
	//		}
	//	} else {
	//		fmt.Println("User is not existed or user is admin.")
	//	}
	//	fmt.Println("Canceled.")
	//	break

	case "activate":
		username := p[1]
		uid := module.UserGetUid(username)
		if uid > 0 {
			err := module.UserSetActive(uid)
			if err != nil {
				return "", err
			}
			str.WriteString("ok")
			return str.String(), nil

		} else {
			str.WriteString(fmt.Sprintf("invalid user: %s", username))
			return "", errors.New(str.String())
		}

	case "password":
		username := p[1]
		password := p[2]
		str.WriteString(fmt.Sprintf("set password for user= %s password= %s", username, password))
		if err := module.UserSetPassword(username, password); err != nil {
			return str.String(), err
		} else {
			str.WriteString("\nok")
			return str.String(), nil
		}

	case "list":
		uis, _ := database.UserGetAll()
		str.WriteString(printUsers(uis, 100))
		return str.String(), nil

	}
	return "", errors.New("unknown command")
}

func parseGroup(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		str.WriteString("group create g1,desc1 g2,desc2...\n" +
			"group remove g1\n" +
			"group add g1 role u1 u2...\n" +
			"group del g1 u1 u2\n" +
			"group set g1 role u1 u2...\n" +
			"group get g1 u1\n" +
			"group list\n" +
			"group sync\n" +
			"group view add g1 4ap")
		return "", errors.New(str.String())
	}

	switch p[0] {
	// group create gn,dispn[,memo,gtype,containtype]
	case "create":
		for _, given := range p[1:] {
			groupname := given
			displayname := given
			groupType := "label_dicom"
			containType := "instance_id"
			memo := "cli_create"

			s := strings.Split(given, ",")
			switch len(s) {
			case 2:
				// group create gn,dispn
				groupname = s[0]
				displayname = s[1]
			case 3:
				// group create gn,dispn,memo
				groupname = s[0]
				displayname = s[1]
				memo = s[2]
			case 4:
				// group create gn,dispn,memo,gtype
				groupname = s[0]
				displayname = s[1]
				memo = s[2]
				groupType = s[3]
			case 5:
				// group create gn,dispn,memo,gtype,containtype
				groupname = s[0]
				displayname = s[1]
				memo = s[2]
				groupType = s[3]
				containType = s[4]
			}

			if err := module.GroupCreate(groupname, displayname, groupType, containType, memo); err != nil {
				log("e", err.Error())
			}
		}
		return str.String(), nil

	case "list":
		if gis, err := database.GroupGetAll(); err != nil {
			return "", err
		} else {
			return printGroups(gis, 100), nil
		}

	case "remove":
		group := p[1]
		str.WriteString(fmt.Sprintf("Remove group [ %s ]", group))
		if err := module.GroupDel(group); err != nil {
			return "", err
		}
		return "ok", nil

	case "add":
		if len(p) < 4 {
			return "", errors.New("format check failed")
		}

		var err error

		group := p[1]
		role := p[2]
		users := p[3:]

		gid := module.GroupGetGid(group)
		if gid == 0 {
			gid, err = strconv.Atoi(group)
		}

		if err != nil {
			return "", err
		} else if gid == 0 {
			return "", errors.New("cannot find group")
		}

		for _, user := range users {
			uid := module.UserGetUid(user)
			if uid == 0 {
				uid, err = strconv.Atoi(user)
				if err != nil {
					log("e", err.Error())
					continue
				}
			}

			if uid == 0 {
				log("e", "cannot find user")
				continue
			}

			str.WriteString(fmt.Sprintf("add user %s [%d] into group %s [%d] as role [ %s ]\n", module.UserGetRealname(uid), uid, module.GroupGetGroupname(gid), gid, role))
			if err = module.GroupUserAdd(gid, uid, role); err != nil {
				log("e", err.Error())
			}
		}

		return str.String(), err

	case "set":
		group := p[1]
		gid := module.GroupGetGid(group)
		role := p[2]
		for _, u := range p[3:] {
			uid := module.UserGetUid(u)
			err := module.GroupUserSetPermissioin(gid, uid, role)
			if err != nil {
				log("e", err.Error())
			} else {
				str.WriteString(fmt.Sprintf("Set user [ %s,%d ] into group [ %s,%d ] as role [ %s ]\n", u, uid, group, gid, role))
			}
		}
		return str.String(), nil

	case "del":
		group := p[1]
		for _, u := range p[2:] {
			str.WriteString(fmt.Sprintf("Del user %s from group %s\n", u, group))
		}
		return str.String(), errors.New("not functioned")

	case "view":
		switch p[1] {
		case "add":
			group := p[2]
			gid := module.GroupGetGid(group)
			mids := module.GroupGetMedia(gid)
			view := p[3]
			for _, mid := range mids {
				err := database.MediaAddView(mid, view)
				if err != nil {
					log("e", err.Error())
				}
			}
		}
		return "ok", nil

	case "sync":
		//groupSyncUsers()
		return "", errors.New("non functioned")

	default:
		return "", errors.New("unknown command")

	}
}

func parseMedia(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		str.WriteString("media load from l to l1 l2 l3 l4\n" +
			"media list\n" +
			"media list l\n" +
			"media find hash h1\n" +
			"media label del m1 m2...\n" +
			"media label setByMid m1 progress 1-7\n" +
			"media label setByHash h1 progress 1-7\n" +
			"media move target\n" +
			"media check")
		return "", errors.New(str.String())
	}

	switch p[0] {
	//
	case "load":
		if p[1] == "from" && p[3] == "to" {
			gmaster := p[2]
			gsubs := p[4:]

			gmasterId := module.GroupGetGid(gmaster)
			if gmasterId <= 0 {
				panic("group is not existed:" + gmaster)
			}

			gsubsId := make([]int, 0)
			for _, gname := range gsubs {
				gid := module.GroupGetGid(gname)
				if gid == 0 {
					panic("group is not existed:" + gname)
				}
				gsubsId = append(gsubsId, gid)
			}
			fmt.Printf("from [%s]=[%d] to %v=%v\n", gmaster, gmasterId, gsubs, gsubsId)

			gSubsMids := make(map[int][]int, 0)
			mids := module.GroupGetMedia(gmasterId)
			sort.Ints(mids)
			gidx := 0
			for _, mid := range mids {
				if gidx == len(gsubsId) {
					gidx = 0
				}
				gsubId := gsubsId[gidx]
				gSubsMids[gsubId] = append(gSubsMids[gsubId], mid)
				gidx++
			}

			for gid, mids := range gSubsMids {
				for _, mid := range mids {
					err := database.GroupAddContain(gid, mid)
					if err != nil {
						fmt.Println("import:", err.Error())
					}
				}
			}

			break
		}

	case "list":
		switch len(p) {
		case 1:
			mis, _ := database.MediaGetAll()
			printMedias(mis)

		default:
			mids := make([]int, 0)
			gid := module.GroupGetGid(p[1])
			mids = module.GroupGetMedia(gid)

			var mis []database.MediaInfo
			for _, mid := range mids {
				mi, _ := database.MediaGet(mid)
				mis = append(mis, mi)
			}
			printMedias(mis)
		}
		return "ok", nil

	case "find":
		if len(p) < 3 {
			return "", errors.New("use: media find hash HASH")
		}

		switch p[1] {
		case "hash":
			mediaHash := p[2]
			mi, err := database.MediaGet(mediaHash)
			if err != nil {
				return "", err
			}

			if _, err = os.Stat(mi.Path); err != nil {
				return "", errors.New("file check: not found")
			} else {
				str.WriteString(fmt.Sprintf("dispname: %s\nfilepath: %s", mi.DisplayName, mi.Path))
				return str.String(), nil
			}

		default:

			mids := make([]int, 0)
			m1 := make(map[int]string) //mid到hash
			m2 := make(map[int]int)    //mid到uid
			m3 := make(map[int]string) //mid到name
			for i := 2; i < len(p); i++ {
				hash := p[i]
				mi, err := database.MediaGet(hash)
				if err != nil {
					log("e", err.Error())
					continue
				}

				mids = append(mids, mi.Mid)
				m1[mi.Mid] = hash
				var uid int
				uid = mi.LabelAuthorUid
				m2[mi.Mid] = uid
				user, err := database.UserGet(uid)
				if err != nil {
					log("e", err.Error())
					continue
				}
				m3[mi.Mid] = user.Realname
			}
			sort.Ints(mids)
			for i := 0; i < len(mids); i++ {
				mid := mids[i]
				str.WriteString(fmt.Sprintf("%d/%s/%d/%s\n", mid, m1[mid], m2[mid], m3[mid]))
			}
			return str.String(), nil
		}

	case "label":
		switch p[1] {
		case "del":
			for _, strmid := range p[2:] {
				mid, err := strconv.Atoi(strmid)
				if err != nil {
					fmt.Println("E;invalid mid:", strmid)
					continue
				}
				fmt.Println("removing label for mid:", mid)
				summary, err := module.MediaGetSummary(mid)
				fmt.Printf("summary:\n%v\nare you sure? (y/n):", summary)
				var confirm string
				fmt.Scanln(&confirm)
				switch confirm {
				case "y", "Y":
					err = module.MediaDeleteLabelAll(mid)
					if err != nil {
						fmt.Println("E:", err.Error())
					} else {
						module.UserSetMediaMemo(1, mid, "")
						fmt.Println("OK")
					}
				default:
					fmt.Println("user ignore")
				}
			}

		case "setByMid":
			switch p[len(p)-2] {
			case "progress":
				prog, _ := strconv.Atoi(p[len(p)-1])
				if prog < 1 {
					fmt.Println("E: progress=", prog)
					break
				}
				m := 0
				n := 0
				for i := 3; i < len(p)-2; i++ {
					n++
					mid, _ := strconv.Atoi(p[i])
					if mid < 1 {
						fmt.Println("E: mid=", mid)
						continue
					}
					mi, err := database.MediaGet(mid)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					li, err := database.LabelGet(mi.Hash)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					fmt.Printf("Mid = %d, Lid = %d\n", mid, li.Lid)
					if err = database.LabelUpdateProgress(li.Lid, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else if err = database.MediaUpdateLabelProgress(mid, li.AuthorUid, li.ReviewUid, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else {
						m++
						fmt.Println("OK")
					}
				}
				fmt.Printf("%d/%d finished\n", m, n)
			default:
				fmt.Println("unknown:", p[len(p)-2])
			}

		case "setByHash":
			switch p[len(p)-2] {
			case "progress":
				prog, _ := strconv.Atoi(p[len(p)-1])
				if prog < 1 {
					fmt.Println("E: progress=", prog)
					break
				}
				m := 0
				n := 0
				var i int
				for i = 3; i < len(p)-2; i++ {
					n++
					hash := p[i]
					mi, err := database.MediaGet(hash)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					li, err := database.LabelGet(hash)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					mid := mi.Mid
					fmt.Printf("Mid = %d, Lid = %d\n", mid, li.Lid)
					if err = database.LabelUpdateProgress(li.Lid, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else if err = database.MediaUpdateLabelProgress(mid, li.AuthorUid, li.ReviewUid, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else {
						m++
						fmt.Println("OK")
					}
				}
				fmt.Printf("%d/%d finished\n", m, n)
			default:
				fmt.Println("unknown:", p[len(p)-2])
			}
		}
		break

	}
	return "ok", nil
}

func parseLabel(p []string) (string, error) {
	if len(p) == 0 {
		return "", errors.New("label sync")

	} else if p[0] == "sync" {
		mis, err := database.MediaGetAll()
		if err != nil {
			return "", err
		}

		for _, mi := range mis {
			mid := mi.Mid
			aid := mi.LabelAuthorUid
			rid := mi.LabelReviewUid
			prog := mi.LabelProgress
			if aid != 0 || rid != 0 {
				li, err := database.LabelGet(mi.Hash)
				if err != nil {
					log("e", err.Error())
					continue
				}

				lid := li.Lid
				aid2 := li.AuthorUid
				rid2 := li.ReviewUid
				prog2 := li.Progress
				if aid2 != aid || rid2 != rid || prog2 != prog {
					log("i", fmt.Sprintf("Conflict! MID: %d, LID: %d, A1-A2: %d<-%d, R1-R2: %d<-%d, P1-P2: %d<-%d\n", mid, lid, aid, aid2, rid, rid2, prog, prog2))
					err = database.MediaUpdateLabelProgress(mid, aid2, rid2, prog2)
					if err != nil {
						log("e", err.Error())
					}
				}
			}
		}
		return "ok", nil
	} else {
		return "", errors.New("unknown command")
	}
}

func parseImport(p []string) (string, error) {
	if len(p) == 0 {
		return "", errors.New("import Folder1 Folder2")
	}

	destFolder, err := filepath.Abs(path.Join(global.GetAppSettings().PathMedia, "us", time.Now().Format("20060102-15H")))
	if err != nil {
		return "", err

	} else if err = os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return "", err

	}

	for _, srcFolder := range p {
		if err = importFolder(srcFolder, destFolder); err != nil {
			log("e", err.Error())
		}
	}

	return "ok", nil
}

func parseGenJson(p []string) (string, error) {
	FileJson := "data.json"
	FolderMedia := "media"

	switch len(p) {
	case 0:
		log("i", "genjson folder jsonfile")
	case 1:
		FolderMedia = p[0]
	default:
		FolderMedia = p[0]
		FileJson = p[1]
	}

	if s, err := os.Stat(FolderMedia); err != nil {
		return "", errors.New(fmt.Sprintf("folder not exist: %s", FolderMedia))
	} else if !s.IsDir() {
		return "", errors.New(fmt.Sprintf("need a folder:", FolderMedia))
	}

	if data, err := scanFolder(FolderMedia); err != nil {
		return "", err
	} else if fp, err := os.OpenFile(FileJson, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644); err != nil {
		return "", err
	} else if bs, err := json.Marshal(data); err != nil {
		return "", err
	} else if n, err := fp.Write(bs); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("json write %d bytes\n", n), nil
	}
}

func importFolder(srcFolder, destFolder string) error {
	var bs []byte
	var data []module.MediaImportJson

	fmt.Println("checking folder:", srcFolder)
	if _, err := os.Stat(srcFolder); err != nil {
		return err

	} else if fp, err := os.Open(filepath.Join(srcFolder, "data.json")); err != nil {
		return err

	} else {
		bs, _ = ioutil.ReadAll(fp)
		fp.Close()
	}

	if err := json.Unmarshal(bs, &data); err != nil {
		return err

	} else if len(data) == 0 {
		return errors.New("ignore empty data.json")

	}

	fmt.Printf("import media from %s => %s\nwaiting 5 seconds to continue", srcFolder, destFolder)

	return module.MediaImportFromJson(1, srcFolder, destFolder, data)
}

func scanFolder(srcFolder string) ([]module.MediaImportJson, error) {
	fmt.Println("scan folder:", srcFolder)

	return nil, nil
}

func printMedias(mis []database.MediaInfo) {
	termbox.Init()
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	termW, termH := termbox.Size()

	lineRemains := 6
	lineStart := 0
	lineTotal := len(mis)

	linePage := termH - lineRemains
	pageTotal := lineTotal / linePage

	pgb := progressbar.NewOptions(pageTotal, progressbar.OptionShowCount(), progressbar.OptionSetPredictTime(false), progressbar.OptionFullWidth())
	pgb.Set(0)

	termbox.Flush()

	printMediasPage(mis[0:linePage], termW)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown, termbox.KeySpace:
				if lineTotal-lineStart+1 > linePage {
					lineStart += linePage
				} else {
					continue
				}

			case termbox.KeyArrowUp:
				if lineStart == 0 {
					continue
				} else if lineStart-linePage < 0 {
					lineStart = 0
				} else {
					lineStart -= linePage
				}

			case termbox.KeyCtrlC, termbox.KeyEsc, termbox.KeyCtrlQ:
				return

			default:
				log("w", "unknown press", termbox.EventKey)
			}

		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			termW, termH = termbox.Size()

			linePage = termH - lineRemains
			pageTotal = lineTotal / linePage
			pgb.ChangeMax(pageTotal)
			termbox.Flush()
		}

		pageCurrent := lineStart * (pageTotal + 1) / lineTotal
		//tools.ScreenClear()
		pgb.Reset()
		pgb.Set(pageCurrent)

		if lineTotal-lineStart > linePage {
			printMediasPage(mis[lineStart:lineStart+linePage], termW)
		} else {
			printMediasPage(mis[lineStart:], termW)
		}
	}
}

func printMediasPage(mis []database.MediaInfo, termWidth int) {
	if len(mis) == 0 {
		return
	}

	var lines []string
	maxLineChars := termWidth - 2

	for _, mi := range mis {
		line := fmt.Sprintf("| %5d | %-4s | %s", mi.Mid, mi.IncludeViews, mi.Path)
		lineWidth := runewidth.StringWidth(line)
		lines = append(lines, line)

		if lineWidth > maxLineChars {
			maxLineChars = lineWidth
		}
	}

	title := tools.LineBuilder(maxLineChars+2, "-")
	fmt.Println("")
	fmt.Println(title)
	fmt.Printf("%s |\n", runewidth.FillRight("|  Mid  | View | PathCache", maxLineChars))
	fmt.Println(title)

	for _, line := range lines {
		fmt.Printf("%s |\n", runewidth.FillRight(line, maxLineChars))
	}

	fmt.Println(title)
}

func printUsers(uis []database.UserInfo, termWidth int) string {
	str := strings.Builder{}

	var lines []string
	maxLineChars := termWidth - 2

	widthUN := 10
	widthRN := 10

	for _, ui := range uis {
		un := runewidth.StringWidth(ui.Username)
		rn := runewidth.StringWidth(ui.Realname)
		if un > widthUN {
			widthUN = un
		}
		if rn > widthRN {
			widthRN = rn
		}
	}
	//fmt.Println(widthUN, widthRN)

	for _, ui := range uis {
		username := runewidth.FillRight(ui.Username, widthUN)
		realname := runewidth.FillRight(ui.Realname, widthRN)
		line := fmt.Sprintf("| %3d | %s |  %d  | %-20s | %-20s | %-3d | %s", ui.Uid, username, ui.Activate, ui.RegisterTime, ui.LoginTime, ui.LoginCount, realname)
		lines = append(lines, line)
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxLineChars {
			maxLineChars = lineWidth
		}
	}
	head := fmt.Sprintf("| UID | %s | Act | %-20s | %-20s | LCo | %s", runewidth.FillRight("Username", widthUN), "Reg", "Login", runewidth.FillRight("Realname", widthRN))
	title := tools.LineBuilder(maxLineChars+2, "-")
	//tools.ScreenClear()

	str.WriteString(title)
	str.WriteString("\n")
	str.WriteString(fmt.Sprintf("%s |\n", runewidth.FillRight(head, maxLineChars)))
	str.WriteString(title)
	str.WriteString("\n")
	for _, line := range lines {
		str.WriteString(fmt.Sprintf("%s |\n", runewidth.FillRight(line, maxLineChars)))
	}
	str.WriteString(title)
	str.WriteString("\n")
	return str.String()
}

func printGroups(gis []database.GroupInfo, termWidth int) (result string) {
	var lines []string
	maxLineChars := termWidth - 2

	widthGN := 10
	widthDN := 10

	for _, gi := range gis {
		gn := runewidth.StringWidth(gi.GroupName)
		dn := runewidth.StringWidth(gi.DisplayName)
		if gn > widthGN {
			widthGN = gn
		}
		if dn > widthDN {
			widthDN = dn
		}
	}

	for _, gi := range gis {
		var userAndPerms map[string]int
		var usernameStr string
		var groupType string

		json.Unmarshal([]byte(gi.Users), &userAndPerms)

		uidIndex := make([]int, 0)

		for k, _ := range userAndPerms {
			uid, _ := strconv.Atoi(k)
			uidIndex = append(uidIndex, uid)
		}

		sort.Ints(uidIndex)

		userRealnames := make([]string, 0)
		for _, uid := range uidIndex {
			if userinfo, err := database.UserGet(uid); err == nil {
				realname := userinfo.Realname
				perm := userAndPerms[strconv.Itoa(uid)]
				switch perm {
				case 127:
					realname += "*"
				case 0:
					realname += "?"
				}
				userRealnames = append(userRealnames, realname)
			}
		}
		usernameStr = strings.Join(userRealnames, " ")

		switch gi.GroupType {
		case "admin":
			groupType = "GAdmin"
		case "label_media":
			groupType = "LMedia"
		case "label_dicom":
			groupType = "LDicom"
		case "screen_dicom":
			groupType = "SDicom"
		}

		countMedia := 0

		switch gi.GroupType {
		case "label_media":
			if ids, _, err := database.GroupGetContains(gi.Gid); err == nil {
				countMedia = len(ids)
			} else {
				fmt.Println(err.Error())
			}

		case "label_dicom", "screen_dicom":
			if ids, _, err := database.GroupGetContains(gi.Gid); err == nil {
				countMedia = len(ids)
			} else {
				fmt.Println(err.Error())
			}
		}

		groupname := runewidth.FillRight(gi.GroupName, widthGN)
		dispname := runewidth.FillRight(gi.DisplayName, widthDN)
		line := fmt.Sprintf("| %3d | %s | %s | %-6s | %5d | %s", gi.Gid, groupname, dispname, groupType, countMedia, strings.TrimSpace(usernameStr))
		lines = append(lines, line)
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxLineChars {
			maxLineChars = lineWidth
		}
	}
	title := tools.LineBuilder(maxLineChars+2, "-")
	//tools.ScreenClear()
	head := fmt.Sprintf("| GID | %s | %s | GTypes | Count | Users", runewidth.FillRight("Groupname", widthGN), runewidth.FillRight("Dispname", widthDN))
	head = fmt.Sprintf("%s |", runewidth.FillRight(head, maxLineChars))

	var content string
	for _, line := range lines {
		content += fmt.Sprintf("%s |\n", runewidth.FillRight(line, maxLineChars))
	}

	result = fmt.Sprintf("%s\n%s\n%s\n%s%s", title, head, title, content, title)
	return result
}
