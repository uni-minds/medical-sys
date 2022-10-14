/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: main.go
 */

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/global"
	"uni-minds.com/liuxy/medical-sys/logger"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/liuxy/medical-sys/tools"
)

func main() {
	logger.Init("", true)
	if len(os.Args) > 1 {
		var p []string
		for _, arg := range os.Args[1:] {
			p = append(p, arg)
		}
		run(p)
	} else {
		command()
	}

	os.Exit(0)
}

func command() {
	var input string
	for {
		fmt.Printf("#> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = strings.TrimSpace(scanner.Text())
		p := strings.Split(input, " ")
		if len(p) == 0 {
			continue
		} else if p[0] == "exit" {
			os.Exit(0)
		} else {
			run(p)
		}
	}
}

func run(p []string) {
	switch strings.ToLower(p[0]) {
	case "group":
		parseGroup(p[1:])

	case "user":
		parseUser(p[1:])

	case "media":
		parseMedia(p[1:])

	case "label":
		parseLabel(p[1:])

	case "import":
		parseImport(p[1:])

	case "progress":
		fmt.Println("1:标注中\n" +
			"2:标注完成\n" +
			"3:审阅中\n" +
			"4:审阅完成，拒绝\n" +
			"5:标注修改中\n" +
			"6:标注完成修改，提交审阅\n" +
			"7:审阅接受，最终状态")

	default:
		fmt.Printf("unsupported command: %s\n"+
			"support user | group | media | label | import | progress", p[0])
	}
}

func parseUser(p []string) {
	if len(p) == 0 {
		fmt.Println("user create u1 p1 (description)\n" +
			"user remove u1 u2...\n" +
			"user activate u1\n" +
			"user password u1 p1\n" +
			"user merge u_old into u_new" +
			"user list")
		return
	}

	switch p[0] {
	case "create":
		break

	case "remove":
		for _, u := range p[1:] {
			uid := module.UserGetUid(u)
			var answer string
			fmt.Printf("DANGER!! User %s [UID=%d] will be delete.", u, uid)
			fmt.Printf("Confirm? [Y/n]:")
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" {
				database.UserDelete(uid)
			}
		}
		break

	case "merge":
		uo := p[1]
		un := p[3]
		uido := module.UserGetUid(uo)
		uidn := module.UserGetUid(un)
		if uido > 0 || uidn > 0 {
			var answer string
			fmt.Printf("DANGER!! User %s [UID=%d] will controll all label made by User %s [UID=%d]\n", un, uidn, uo, uido)
			fmt.Printf("Confirm? [Y/n]:")
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" {
				fmt.Println("User confirmed.")
				lis, _ := database.LabelGetAll()
				for _, li := range lis {
					change := false
					aid := li.AuthorUid
					rid := li.ReviewUid
					if aid == uido {
						aid = uidn
						change = true
					}
					if rid == uido {
						rid = uidn
						change = true
					}
					if change {
						li.AuthorUid = aid
						li.ReviewUid = rid
						database.LabelUpdate(li)
						fmt.Printf("Label owner: A<-%d, R<-%d\n", aid, rid)
					}
				}
				return
			}
		} else {
			fmt.Println("User is not existed or user is admin.")
		}
		fmt.Println("Canceled.")
		break

	case "activate":
		username := p[1]
		uid := module.UserGetUid(username)
		if uid > 0 {
			fmt.Println(module.UserSetActive(uid))
		} else {
			fmt.Println("invalid user:", username)
		}
		break

	case "password":
		username := p[1]
		password := p[2]
		fmt.Println("set password for user:", username, " password=", password)
		if err := module.UserSetPassword(username, password); err != nil {
			fmt.Println("E:", err.Error())
		} else {
			fmt.Println("OK")
		}
		break

	case "list":
		uis, _ := database.UserGetAll()
		printUsers(uis, 100)
	}
	return
}

func parseGroup(p []string) {
	if len(p) == 0 {
		fmt.Println("group create g1,desc1 g2,desc2...\n" +
			"group remove g1\n" +
			"group add g1 role u1 u2...\n" +
			"group del g1 u1 u2\n" +
			"group set g1 role u1 u2...\n" +
			"group get g1 u1\n" +
			"group list\n" +
			"group sync\n" +
			"group view add g1 4ap")
		return
	}

	switch p[0] {
	case "create":
		for _, given := range p[1:] {
			var groupname, displayname string
			memo := ""
			if strings.Contains(given, ",") {
				s := strings.Split(given, ",")
				groupname = s[0]
				displayname = s[1]
				if len(s) > 2 {
					memo = s[3]
				}
			} else {
				groupname = given
				displayname = given
			}
			if groupname != "" && displayname != "" {
				module.GroupCreate(groupname, displayname, memo)
			}
		}
		break
	case "list":
		gis, _ := database.GroupGetAll()
		printGroups(gis, 100)
	case "remove":
		group := p[1]
		fmt.Println("Remove group [", group, "]")
		if err := module.GroupDel(group); err != nil {
			fmt.Println("E:", err.Error())
		}
		break
	case "add":
		groupAddUsers(p[1], p[2], p[3:])
		break
	case "set":
		group := p[1]
		gid := module.GroupGetGid(group)
		role := p[2]
		for _, u := range p[3:] {
			uid := module.UserGetUid(u)
			err := module.GroupUserSetPermissioin(gid, uid, role)
			if err != nil {
				fmt.Println("E:", err.Error())
			} else {
				fmt.Printf("Set user [ %s,%d ] into group [ %s,%d ] as role [ %s ]\n", u, uid, group, gid, role)
			}
		}
		break
	case "del":
		group := p[1]
		for _, u := range p[2:] {
			fmt.Printf("Del user %s from group %s\n", u, group)
		}
		break
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
					fmt.Println("E;media add view:", err.Error())
				}
			}
		}
		break
	case "sync":
		groupSyncUsers()
	}
	return
}

func parseMedia(p []string) {
	if len(p) == 0 {
		fmt.Println("media load from l to l1 l2 l3 l4\n" +
			"media list\n" +
			"media list l\n" +
			"media find hash h1\n" +
			"media label del m1 m2...\n" +
			"media label setByMid m1 progress 1-7\n" +
			"media label setByHash h1 progress 1-7\n" +
			"media move target\n" +
			"media check")
		return
	}

	switch p[0] {
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
					err := database.GroupAddMedia(gid, mid)
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
		break

	case "find":
		if len(p) < 3 {
			fmt.Println("media find hash HASH")
			return
		}

		switch p[1] {
		case "hash":
			mediaHash := p[2]
			mi, err := database.MediaGet(mediaHash)
			if err != nil {
				fmt.Println("E:", err.Error())
				return
			}
			printMedias([]database.MediaInfo{mi})

			if _, err = os.Stat(mi.Path); err != nil {
				fmt.Println("File check: not found")
			} else {
				fmt.Println("File check: OK")
			}
			return

		default:

			mids := make([]int, 0)
			m1 := make(map[int]string) //mid到hash
			m2 := make(map[int]int)    //mid到uid
			m3 := make(map[int]string) //mid到name
			for i := 2; i < len(p); i++ {
				hash := p[i]
				mi, err := database.MediaGet(hash)
				if err != nil {
					fmt.Println("E:", err.Error())
					continue
				}
				mids = append(mids, mi.Mid)
				m1[mi.Mid] = hash
				var uid int
				uid = mi.LabelAuthorUid
				m2[mi.Mid] = uid
				user, err := database.UserGet(uid)
				if err != nil {
					fmt.Println("E:", err.Error())
					continue
				}
				m3[mi.Mid] = user.Realname
			}
			sort.Ints(mids)
			for i := 0; i < len(mids); i++ {
				mid := mids[i]
				fmt.Printf("%d/%s/%d/%s\n", mid, m1[mid], m2[mid], m3[mid])
			}
			break
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

	case "move":
		target := global.GetAppSettings().SystemMediaPath
		if len(p) > 1 && p[1] != "" {
			target = p[1]
		}

		if target, err := filepath.Abs(target); err != nil {
			fmt.Println("E:", err.Error())
			return

		} else if _, err = os.Stat(target); err != nil {
			fmt.Println("E:", err.Error())
			return
		}

		fmt.Printf("move media to %s\nwait 5 seconds to continue", target)
		sleep(0)

		mis, err := database.MediaGetAll()
		if err != nil {
			fmt.Println("E:", err.Error())
			return
		}

		fmt.Println("total media:", len(mis))
		for _, mi := range mis {
			rel, err := filepath.Rel(target, mi.Path)
			if err != nil {
				fmt.Println(mi.Path, err.Error())
			} else if rel[:1] == "." {
				p := strings.Split(mi.Path, "media/us/")
				if len(p) > 1 {
					// e.g. /application/media/us/20200820-21H/27818-3675637-4C.ogv

					srcFile := mi.Path
					destFolder := path.Join(target, path.Dir(p[1]))
					os.MkdirAll(destFolder, os.ModePerm)
					destFile := path.Join(destFolder, path.Base(p[1]))

					switch filepath.Ext(srcFile) {
					case ".jpg":
						fmt.Println("ignore jpg file:", srcFile)
						continue

					default:
						if err := tools.MoveFile(srcFile, destFile); err != nil {
							panic(err.Error())

						}
					}

					if err := database.MediaUpdatePath(mi.Mid, destFile); err != nil {
						fmt.Println(err.Error())
						panic(err.Error())
					}

					fmt.Println("OK")
				}
			}
		}

	case "check":
		fmt.Println("checking files in database")
		mis, _ := database.MediaGetAll()
		for _, mi := range mis {
			file := mi.Path
			_, err := os.Stat(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			switch mi.IncludeViews {
			case "[]", "null":
				database.MediaUpdateViews(mi.Mid, "")

			}
		}
	}
}

func parseLabel(p []string) {
	if len(p) == 0 {
		fmt.Println("label sync")
		return
	} else if p[0] == "sync" {
		mis, err := database.MediaGetAll()
		if err != nil {
			fmt.Println("E:", err.Error())
			return
		}
		for _, mi := range mis {
			mid := mi.Mid
			aid := mi.LabelAuthorUid
			rid := mi.LabelReviewUid
			prog := mi.LabelProgress
			if aid != 0 || rid != 0 {
				li, err := database.LabelGet(mi.Hash)
				if err != nil {
					fmt.Println("E:", err.Error())
					continue
				}

				lid := li.Lid
				aid2 := li.AuthorUid
				rid2 := li.ReviewUid
				prog2 := li.Progress
				if aid2 != aid || rid2 != rid || prog2 != prog {
					fmt.Printf("Conflict! MID: %d, LID: %d, A1-A2: %d<-%d, R1-R2: %d<-%d, P1-P2: %d<-%d\n", mid, lid, aid, aid2, rid, rid2, prog, prog2)
					database.MediaUpdateLabelProgress(mid, aid2, rid2, prog2)
				}
			}
		}

	}
}

func parseImport(p []string) {
	if len(p) == 0 {
		fmt.Println("import Folder1 Folder2")
		return
	}

	destFolder, err := filepath.Abs(path.Join(global.GetAppSettings().SystemMediaPath, "us", time.Now().Format("20060102-15H")))
	if err != nil {
		fmt.Println("E:", err.Error())
		return

	} else if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		fmt.Println("E:", err.Error())
		return
	}

	for _, folder := range p {
		if err := importFolder(folder, destFolder); err != nil {
			fmt.Println("E:", err.Error())
		}
	}
}

func groupAddUsers(group string, role string, users []string) {
	var err error
	if group == "" {
		fmt.Println("group name is empty")
		return
	}

	gid := module.GroupGetGid(group)
	if gid == 0 {
		gid, err = strconv.Atoi(group)
	}

	if err != nil {
		fmt.Println("E:", err.Error())
		return
	} else if gid == 0 {
		fmt.Println("cannot find group")
		return
	}

	for _, user := range users {
		uid := module.UserGetUid(user)
		if uid == 0 {
			uid, err = strconv.Atoi(user)
			if err != nil {
				fmt.Println("E:", err.Error())
				continue
			}
		}

		if uid == 0 {
			fmt.Println("cannot find user")
			continue
		}

		fmt.Printf("add user %s [%d] into group %s [%d] as role [ %s ]\n", module.UserGetRealname(uid), uid, module.GroupGetGroupname(gid), gid, role)
		module.GroupUserAdd(gid, uid, role)
	}
	//groupSyncUsers()
}

func groupSyncUsers() {
	gis, _ := database.GroupGetAll()
	for _, gi := range gis {
		var us map[int]int
		json.Unmarshal([]byte(gi.Users), &us)
		for uid, perm := range us {
			ui, _ := database.UserGet(uid)
			database.UserAddGroup(uid, gi.Gid)
			fmt.Println("G", gi.Gid, ui, perm)
		}
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

	sleep(5)

	return module.MediaImportFromJson(1, srcFolder, destFolder, data)
}

func sleep(sec int) {
	for i := 0; i < sec; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
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
		tools.ScreenClear()
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
	fmt.Printf("%s |\n", runewidth.FillRight("|  Mid  | View | Path", maxLineChars))
	fmt.Println(title)

	for _, line := range lines {
		fmt.Printf("%s |\n", runewidth.FillRight(line, maxLineChars))
	}

	fmt.Println(title)
}

func printUsers(uis []database.UserInfo, termWidth int) {
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
	fmt.Println(widthUN, widthRN)

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
	tools.ScreenClear()
	fmt.Println(title)
	fmt.Printf("%s |\n", runewidth.FillRight(head, maxLineChars))
	fmt.Println(title)
	for _, line := range lines {
		fmt.Printf("%s |\n", runewidth.FillRight(line, maxLineChars))
	}
	fmt.Println(title)
}

func printGroups(gis []database.GroupInfo, termWidth int) {
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
		var uidstr string

		json.Unmarshal([]byte(gi.Users), &userAndPerms)

		for uid, _ := range userAndPerms {
			if uidstr == "" {
				uidstr = uid
			} else {
				uidstr = fmt.Sprintf("%s %s", uidstr, uid)
			}
		}

		groupname := runewidth.FillRight(gi.GroupName, widthGN)
		dispname := runewidth.FillRight(gi.DisplayName, widthDN)
		line := fmt.Sprintf("| %3d | %s | %s | %s", gi.Gid, groupname, dispname, uidstr)
		lines = append(lines, line)
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxLineChars {
			maxLineChars = lineWidth
		}
	}
	head := fmt.Sprintf("| GID | %s | %s | Users", runewidth.FillRight("Groupname", widthGN), runewidth.FillRight("Dispname", widthDN))
	title := tools.LineBuilder(maxLineChars+2, "-")
	tools.ScreenClear()
	fmt.Println(title)
	fmt.Printf("%s |\n", runewidth.FillRight(head, maxLineChars))
	fmt.Println(title)
	for _, line := range lines {
		fmt.Printf("%s |\n", runewidth.FillRight(line, maxLineChars))
	}
	fmt.Println(title)
}

func log(level string, message ...interface{}) {
	msg := tools.ExpandInterface(message)
	logger.Write("TOOL", level, msg)
}
