/**
 * @Author: Liu Xiangyu
 * @Description:
 * @File:  main-tools
 * @Version: 1.0.0
 * @Date: 2020/5/6 13:14
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"uni-minds.com/liuxy/medical-sys/database"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/liuxy/medical-sys/upgrade"
)

func main() {
	var input string
	for {
		fmt.Printf("#:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()
		if len(input) == 0 {
			continue
		}

		params := strings.Split(input, " ")
		if len(params) > 0 {
			sel := params[0]
			if strings.ToLower(sel) == "exit" {
				return
			} else {
				process(sel, params[1:])
			}
		}
	}
}

func process(sel string, p []string) {
	if len(p) == 0 {
		switch sel {
		case "group":
			fmt.Println("\n" +
				"group create g1,desc1 g2,desc2...\n" +
				"group remove g1\n" +
				"group add g1 role u1 u2...\n" +
				"group del g1 u1 u2\n" +
				"group set g1 role u1 u2...\n" +
				"group get g1 u1\n" +
				"group list\n" +
				"group sync\n" +
				"group view add g1 4ap")
			break
		case "user":
			fmt.Println("\n" +
				"user create u1 p1 (description)\n" +
				"user remove u1 u2...\n" +
				"user activate u1\n" +
				"user password u1 p1\n" +
				"user merge u_old into u_new" +
				"user list")
			break
		case "media":
			fmt.Println("\n" +
				"media load from l to l1 l2 l3 l4\n" +
				"media list\n" +
				"media list l\n" +
				"media getByHash h1\n" +
				"media label del m1 m2...\n" +
				"media label setByMid m1 progress 1-7\n" +
				"media label setByHash h1 progress 1-7")
			break
		case "label":
			fmt.Println("\n" +
				"label sync")
			break
		case "progress":
			fmt.Println("\n" +
				"1:标注中\n" +
				"2:标注完成\n" +
				"3:审阅中\n" +
				"4:审阅完成，拒绝\n" +
				"5:标注修改中\n" +
				"6:标注完成修改，提交审阅\n" +
				"7:审阅接受，最终状态")
			break
		case "sync":
			upgrade.Run3()
			break
		default:
			break
		}
	} else {
		switch sel {
		case "group":
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
				di := module.GroupGetAll()
				var keys = make([]int, 0)
				for i, _ := range di {
					keys = append(keys, i)
				}
				sort.Ints(keys)
				for _, idx := range keys {
					data := di[idx]
					fmt.Printf("%-2d | %-15s | %s\n", idx, data[0], data[1])
					fmt.Printf(" - %s\n", data[2])
				}
				break
			case "remove":
				group := p[1]
				fmt.Println("Remove group [", group, "]")
				if err := module.GroupDel(group); err != nil {
					fmt.Println("E:", err.Error())
				}
				break
			case "add":
				group := p[1]
				role := p[2]
				for _, u := range p[3:] {
					fmt.Printf("Add user [ %s ] into group [ %s ] as role [ %s ]\n", u, group, role)
					module.GroupUserAddFrendly(group, u, role)
				}
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
				break
			}
			return

		case "user":
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
				uis := module.UserGetAll()
				keys := make([]int, 0)
				for k, _ := range uis {
					keys = append(keys, k)
				}
				sort.Ints(keys)
				fmt.Printf("%-3s | %-12s | %-4s | %-20s | %-20s | %-3s | %s\n", "UID", "Username", "Act", "Reg", "Login", "LCo", "Realname")
				for _, uid := range keys {
					ui := uis[uid]
					fmt.Printf("%-3d | %-12s | %-4d | %-20s | %-20s | %-3d | %s\n", uid, ui.Username, ui.Activate, ui.RegisterTime, ui.LoginTime, ui.LoginCount, ui.Realname)
				}
				break
			}
			return

		case "media":
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
				mids := make([]int, 0)
				switch len(p) {
				case 2:
					mis := module.MediaGetAll()
					for k, _ := range mis {
						mids = append(mids, k)
					}
					sort.Ints(mids)
				default:
					gid := module.GroupGetGid(p[2])
					mids = module.GroupGetMedia(gid)
					sort.Ints(mids)
				}

				for _, mid := range mids {
					mi, _ := database.MediaGet(mid)
					fmt.Printf("%-6d| %-32s | %-40s | %-20s\n", mi.Mid, mi.Hash, mi.DisplayName, mi.IncludeViews)
				}
				break

			case "getByHash":
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
							fmt.Println("user ignore\n")
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

		case "label":
			switch p[0] {
			case "sync":
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
				break
			}
			return
		}
		fmt.Println("参数错误", p)
	}
}
