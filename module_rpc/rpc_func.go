/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: rpc_func.go
 * Date: 2022/09/04 10:14:04
 */

package module_rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/uni-minds/medical-sys/database"
	"gitee.com/uni-minds/medical-sys/logger"
	"gitee.com/uni-minds/medical-sys/module"
	"gitee.com/uni-minds/utils/tools"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/schollz/progressbar/v3"
	"os"
	"sort"
	"strconv"
	"strings"
)

var log *logger.Logger

func (this *RpcFunc) ParsePacs(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		return "", errors.New(`pacs commands:
pacs export view 4ap to group GN1
pacs list view 4CV group GN1
pacs list views
pacs split GroupSrc to GDesc1,GDesc2
pacs import Group ["A","B"]`)
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
			case "views":
				return "", errors.New("incomplete")
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
		var gi database.DbStructGroup

		if srcGid, err = strconv.Atoi(srcGroup); err != nil {
			gi, err = database.GroupGet(srcGroup)
			if err != nil {
				return "", errors.New(fmt.Sprintf("Unknown group name: %s", err.Error()))
			} else {
				srcGid = gi.Id
			}
		}

		if dstGid, err = strconv.Atoi(dstGroup); err != nil {
			gi, err = database.GroupGet(dstGroup)
			if err != nil {
				return "", errors.New(fmt.Sprintln("Unknown group name:", err.Error()))
			} else {
				dstGid = gi.Id
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
	// pacs split pacs_label_l to pacs_label_l_1 pacs_label_l_2 pacs_label_l_3 pacs_label_l_4
	case "split":
		if len(p) >= 3 && p[2] == "to" {
			gSrc := p[1]
			gDesc := strings.Split(p[3], ",")

			gSrcId := module.GroupGetGid(gSrc)
			if gSrcId <= 0 {
				return "", errors.New(fmt.Sprintf("source group is not exist: %s", gSrc))
			}

			gDescIds := make([]int, 0)
			for _, gname := range gDesc {
				gid := module.GroupGetGid(gname)
				if gid <= 0 {
					log.Error(fmt.Sprintf("subs group is not exist:", gname))
					continue
				}
				gDescIds = append(gDescIds, gid)
			}

			if len(gDescIds) == 0 {
				return "", errors.New("target group is not existed")
			}

			gSrcDicomIds, gSrcDicomType := module.GroupGetContainMedia(gSrcId)
			gDescDicomIds := make([]string, 0)
			groupDicomIds := make(map[int][]string, 0)

			// 检查是否目录组中已经导入对应id
			for _, groupId := range gDescIds {
				containIds, containType := module.GroupGetContainMedia(groupId)
				if gSrcDicomType != containType {
					eStr := fmt.Sprintf("组间DICOM类型不符：%s != %s", gSrcDicomType, containType)
					log.Error(eStr)
					return "", errors.New(eStr)
				}

				groupDicomIds[groupId] = containIds
				gDescDicomIds = append(gDescDicomIds, containIds...)
			}

			log.Log("i", fmt.Sprintf("split group %s from %d to %v", gSrcDicomType, gSrcId, gDescIds))

			str.WriteString(fmt.Sprintf("src_group:\t%s [%d]\ndest_groups:\t%v %v\n", gSrc, gSrcId, gDesc, gDescIds))

			log.Log("i", fmt.Sprintf("src group contain %d ids.", len(gSrcDicomIds)))
			log.Log("i", fmt.Sprintf("desc groups contain %d ids.", len(gDescDicomIds)))

			// 排除已经导入的id
			remainIds := tools.StringsExcept(gSrcDicomIds, gDescDicomIds)

			log.Log("i", "remain ids to import:", len(remainIds))

			// 轮询并分组
			gidx := 0
			for _, id := range remainIds {
				if gidx == len(gDescIds) {
					gidx = 0
				}
				gsubId := gDescIds[gidx]
				groupDicomIds[gsubId] = append(groupDicomIds[gsubId], id)
				gidx++
			}

			// 导入
			for gid, ids := range groupDicomIds {
				err := database.GroupAddContains(gid, ids)
				if err != nil {
					log.Log("e", "group import ids:", err.Error())
				} else {
					str.WriteString(fmt.Sprintf("group %d has %d id(s) imported.\n", gid, len(ids)))
				}
			}
			str.WriteString("ok")
			return str.String(), nil
		}
		return "", errors.New("unrecognized format")
	// pacs import pacs_label_l []
	case "import":
		if len(p) < 3 {
			return "", errors.New("more params needed")
		}

		groupName := p[1]
		groupId := module.GroupGetGid(groupName)
		if groupId < 1 {
			return "", errors.New("no group")
		}

		ids, _ := module.GroupGetContainMedia(groupId)
		countBefore := len(ids)

		strBuilder := strings.Builder{}
		var studiesIds []string

		for _, str := range p[2:] {
			strBuilder.WriteString(str)
		}

		fmt.Println("import json string:", strBuilder.String())
		if err := json.Unmarshal([]byte(strBuilder.String()), &studiesIds); err != nil {
			return "", err
		} else if err = module.GroupAddMedia(groupId, studiesIds); err != nil {
			return "", err
		} else {
			ids, _ := module.GroupGetContainMedia(groupId)
			countAfter := len(ids)
			return fmt.Sprintf("studies %d ids imported: %d -> %d", len(studiesIds), countBefore, countAfter), nil
		}

	default:
		return "", errors.New(fmt.Sprintf("unknown params: %s", p))
	}
}

func (this *RpcFunc) ParseUser(p []string) (string, error) {
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
					log.Error(err.Error())
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
	//		fmt.printf("DANGER!! User %s [UID=%d] will controll all label made by User %s [UID=%d]\n", un, uidn, uo, uido)
	//		fmt.printf("Confirm? [Y/n]:")
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
	//					database.UpdateAll(li)
	//					fmt.printf("Label owner: A<-%d, R<-%d\n", aid, rid)
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

func (this *RpcFunc) ParseGroup(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		return "", errors.New(`group help:
group create g1,desc1 g2,desc2...
group remove g1
group add g1 role u1 u2...
group del g1 u1 u2
group set g1 role u1 u2...
group get g1 u1
group list
group sync
group view add g1 4ap`)
	}

	switch p[0] {
	// group create gn,dispn[,memo,gtype,containtype]
	case "create":
		for _, given := range p[1:] {
			groupname := given
			groupType := "label_dicom"
			containType := "instance_id"
			memo := "cli_create"

			s := strings.Split(given, ",")
			switch len(s) {
			case 2:
				// group create gn,dispn
				groupname = s[0]
			case 3:
				// group create gn,dispn,memo
				groupname = s[0]
				memo = s[2]
			case 4:
				// group create gn,dispn,memo,gtype
				groupname = s[0]
				memo = s[2]
				groupType = s[3]
			case 5:
				// group create gn,dispn,memo,gtype,containtype
				groupname = s[0]
				memo = s[2]
				groupType = s[3]
				containType = s[4]
			}

			if err := module.GroupCreate(groupname, groupType, containType, memo); err != nil {
				log.Error(err.Error())
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
					log.Error(err.Error())
					continue
				}
			}

			if uid == 0 {
				log.Error("cannot find user")
				continue
			}

			str.WriteString(fmt.Sprintf("add user %s [%d] into group %s [%d] as role [ %s ]\n", module.UserGetRealname(uid), uid, module.GroupGetGroupname(gid), gid, role))
			if err = module.GroupUserAdd(gid, uid, role); err != nil {
				log.Error(err.Error())
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
				log.Error(err.Error())
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
			mediaUUID, _, _ := module.GroupGetContains(gid)
			view := p[3]
			for _, mid := range mediaUUID {
				err := database.GetMedia().AddCrfView(mid, view)
				if err != nil {
					log.Error(err.Error())
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

func (this *RpcFunc) ParseMedia(p []string) (string, error) {
	str := strings.Builder{}

	if len(p) == 0 {
		return "", errors.New(`media help:
media load from l to l1 l2 l3 l4
media list
media list l
media find hash h1
media label del m1 m2...
media label setByMid m1 progress 1-7
media label setByHash h1 progress 1-7
media move target
media check`)
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

			gSubsMids := make(map[int][]string, 0)
			mediaUUIDs, _, _ := module.GroupGetContains(gmasterId)
			gidx := 0
			for _, mediaUUID := range mediaUUIDs {
				if gidx == len(gsubsId) {
					gidx = 0
				}
				gsubId := gsubsId[gidx]
				gSubsMids[gsubId] = append(gSubsMids[gsubId], mediaUUID)
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
			mis, _ := database.GetMedia().GetAll()
			printMedias(mis)

		default:
			mediaUUIDs := make([]string, 0)
			gid := module.GroupGetGid(p[1])
			mediaUUIDs, _, _ = module.GroupGetContains(gid)

			var mis []database.DbStructMedia
			for _, mid := range mediaUUIDs {
				mi, _ := database.GetMedia().Get(mid)
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
			mi, err := database.GetMedia().Get(mediaHash)
			if err != nil {
				return "", err
			}

			if _, err = os.Stat(mi.Path); err != nil {
				return "", errors.New("file check: not found")
			} else {
				str.WriteString(fmt.Sprintf("dispname: %s\nfilepath: %s", mi.DisplayName, mi.Path))
				return str.String(), nil
			}

		}

	case "label":
		switch p[1] {
		case "del":
			for _, mid := range p[2:] {
				fmt.Println("removing label for mid:", mid)
				summary, err := module.GetSummaryMedia(mid)
				fmt.Printf("summary:\n%v\nare you sure? (y/n):", summary)
				var confirm string
				fmt.Scanln(&confirm)
				switch confirm {
				case "y", "Y":
					err = module.MediaRemoveLabelAll(mid)
					if err != nil {
						fmt.Println("E:", err.Error())
					} else {
						module.MediaSetMemo(mid, "")
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
					mediaUUID := p[i]
					if mediaUUID != "" {
						fmt.Println("E: mediaUUID=", mediaUUID)
						continue
					}
					info, err := database.GetMedia().Get(mediaUUID)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					li, err := database.GetLabel().Get(info.MediaUUID)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					fmt.Printf("Mid = %d, Id = %d\n", mediaUUID, li.Id)
					if err = database.GetLabel().UpdateProgress(li.Id, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else if err = database.GetMedia().LabelUpdateProgress(mediaUUID, prog); err != nil {
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
					mi, err := database.GetMedia().Get(hash)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					li, err := database.GetLabel().Get(hash)
					if err != nil {
						fmt.Println("E:", err.Error())
						continue
					}
					mid := mi.Id
					fmt.Printf("Mid = %d, Id = %d\n", mid, li.Id)
					if err = database.GetLabel().UpdateProgress(li.Id, prog); err != nil {
						fmt.Println("E:", err.Error())
					} else if err = database.GetMedia().LabelUpdateProgress(mid, prog); err != nil {
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

func (this *RpcFunc) ParseGenJson(p []string) (string, error) {
	FileJson := "data.json"
	FolderMedia := "media"

	switch len(p) {
	case 0:
		log.Log("i", "genjson folder jsonfile")
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

func scanFolder(srcFolder string) ([]module.MediaImportJson, error) {
	fmt.Println("scan folder:", srcFolder)

	return nil, nil
}

func printMedias(mis []database.DbStructMedia) {
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
				log.Warn(fmt.Sprintf("unknown press", termbox.EventKey))
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

func printMediasPage(mis []database.DbStructMedia, termWidth int) {
	if len(mis) == 0 {
		return
	}

	var lines []string
	maxLineChars := termWidth - 2

	for _, mi := range mis {
		line := fmt.Sprintf("| %5d | %-4s | %s", mi.Id, mi.CrfDefine, mi.Path)
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

func printUsers(uis []database.DbStructUser, termWidth int) string {
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

func printGroups(gis []database.DbStructGroup, termWidth int) (result string) {
	var lines []string
	maxLineChars := termWidth - 2

	widthGN := 10
	widthDN := 10

	for _, gi := range gis {
		gn := runewidth.StringWidth(gi.Name)
		dn := runewidth.StringWidth(gi.Name)
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

		for k := range userAndPerms {
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

		switch gi.Type {
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

		switch gi.Type {
		case "label_media":
			if ids, _, err := database.GroupGetContains(gi.Id); err == nil {
				countMedia = len(ids)
			} else {
				fmt.Println(err.Error())
			}

		case "label_dicom", "screen_dicom":
			if ids, _, err := database.GroupGetContains(gi.Id); err == nil {
				countMedia = len(ids)
			} else {
				fmt.Println(err.Error())
			}
		}

		groupname := runewidth.FillRight(gi.Name, widthGN)
		dispname := runewidth.FillRight(gi.Name, widthDN)
		line := fmt.Sprintf("| %3d | %s | %s | %-6s | %5d | %s", gi.Id, groupname, dispname, groupType, countMedia, strings.TrimSpace(usernameStr))
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
