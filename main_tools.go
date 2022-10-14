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
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"uni-minds.com/liuxy/medical-sys/module"
	"uni-minds.com/medical-sys/database"
	"uni-minds.com/medical-sys/upgrade"
)

//L1 张烨 郭勇 陈倬 周晓雪
//L2 谷孝艳 李烨 王静怡 王斯宇
//L3 刘晓伟 满婷婷 赵映 薛超 [王欣]
//L4 韩建成 郑敏 杨旭 武玉多

func main() {
	for {
		fmt.Printf("\n#:")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		if len(data) == 0 {
			continue
		}
		input := strings.Split(string(data), " ")
		inputlen := len(input)
		switch input[0] {
		case "exit":
			return
		case "group":
			if inputlen == 1 {
				fmt.Printf("\n" +
					"group create g1 (description)\n" +
					"group remove g1\n" +
					"group add g1 u1 (role)\n" +
					"group del g1 u1\n" +
					"group set g1 u1 role\n" +
					"group get g1 u1\n" +
					"group list\n" +
					"group view add g1 4ap")
				continue
			}

			switch input[1] {
			case "create":
				var groupname, displayname string
				switch inputlen {
				case 3:
					groupname = input[2]
					displayname = groupname
				case 4:
					groupname = input[2]
					displayname = input[3]
				default:

				}
				module.GroupCreate(groupname, displayname, "")
				continue

			case "list":
				di := module.GroupGetAll()
				var keys = make([]int, 0)
				for i, _ := range di {
					keys = append(keys, i)
				}
				sort.Ints(keys)
				for _, idx := range keys {
					fmt.Printf("%3d, %v\n", idx, di[idx])
				}
				continue
			case "remove":
				continue
			case "add":
				var username, role string
				if inputlen == 4 {
					username = input[3]
					role = "guest"
				} else if inputlen > 4 {
					username = input[3]
					role = input[4]
				} else {
					break
				}
				module.GroupUserAddFrendly(input[2], username, role)
				continue

			case "set":
				groupname := input[2]
				username := input[3]
				role := input[4]
				gid := module.GroupGetGid(groupname)
				uid := module.UserGetUid(username)
				err := module.GroupUserSetPermissioin(gid, uid, role)
				if err != nil {
					fmt.Println("E:", err.Error())
				} else {
					fmt.Println("OK")
				}
				continue

			case "del":
				continue
			case "view":
				switch input[2] {
				case "add":
					gid := module.GroupGetGid(input[3])
					view := input[4]
					mids := module.GroupGetMedia(gid)
					for _, mid := range mids {
						err := database.MediaAddView(mid, view)
						if err != nil {
							fmt.Println("E;media add view:", err.Error())
						}
					}
				}
				continue

			}
		case "user":
			if inputlen == 1 {
				fmt.Printf("\n" +
					"user create u1 p1 (description)\n" +
					"user remove u1\n" +
					"user activate u1\n" +
					"user password u1 p1\n" +
					"user list\n")
				continue
			}
			switch input[1] {
			case "create":
				continue

			case "remove":
				continue

			case "activate":
				uid := module.UserGetUid(input[2])
				if uid > 0 {
					fmt.Println(module.UserSetActive(uid))
				} else {
					fmt.Println("invalid user:", input[2])
				}
				continue

			case "password":
				username := input[2]
				password := input[3]
				fmt.Println("set password for user:", username, " password=", password)
				if err := module.UserSetPassword(username, password); err != nil {
					fmt.Println("E:", err.Error())
				} else {
					fmt.Println("OK")
				}
				continue

			case "list":
				uis := module.UserGetAll()
				keys := make([]int, 0)
				for k, _ := range uis {
					keys = append(keys, k)
				}
				sort.Ints(keys)
				for _, uid := range keys {
					ui := uis[uid]
					fmt.Printf("%-6d | %-10s | %-10s | %d | %-20s | %-20s | %d\n", uid, ui.Username, ui.Realname, ui.Activate, ui.RegisterTime, ui.LoginTime, ui.LoginCount)
				}
				continue
			}
		case "media":
			if inputlen == 1 {
				fmt.Printf("\n" +
					"media load from l to l1 l2 l3 l4\n" +
					"media list\n" +
					"media list l\n" +
					"media label del m1\n" +
					"media label set m1 progress 1-7\n")
				continue
			}
			switch input[1] {
			case "load":
				if input[2] == "from" && input[4] == "to" {
					gmaster := input[3]
					gsubs := make([]string, 0)
					for i := 5; i < inputlen; i++ {
						gsubs = append(gsubs, input[i])
					}

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

					continue
				}
			case "list":
				mids := make([]int, 0)
				switch inputlen {
				case 2:
					mis := module.MediaGetAll()
					for k, _ := range mis {
						mids = append(mids, k)
					}
					sort.Ints(mids)
				default:
					gid := module.GroupGetGid(input[2])
					mids = module.GroupGetMedia(gid)
					sort.Ints(mids)
				}

				for _, mid := range mids {
					mi, _ := database.MediaGet(mid)
					fmt.Printf("%-6d| %-32s | %-40s | %-20s\n", mi.Mid, mi.Hash, mi.DisplayName, mi.IncludeViews)
				}
				continue
			case "label":
				switch input[2] {
				case "del":
					for _, strmid := range input[3:] {
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
				case "set":
					switch input[4] {
					case "progress":
						prog, _ := strconv.Atoi(input[5])
						mid, _ := strconv.Atoi(input[3])
						if prog < 1 || mid < 1 {
							fmt.Println("E: mid=", mid, ", progress=", prog)
							continue
						}
						mi, _ := database.MediaGet(mid)
						fmt.Println(mi)
						li, err := database.LabelGet(mi.Hash)
						fmt.Println(li)
						if err != nil {
							fmt.Println("E:", err.Error())
							continue
						}
						if err = database.MediaUpdateLabelProgress(mid, li.AuthorUid, li.ReviewUid, prog); err != nil {
							fmt.Println("E:", err.Error())
						} else {
							fmt.Println("OK")
						}
					default:
						fmt.Println("unknown:", input[4])
					}
				}
				continue
			}
		case "progress":
			fmt.Println("\n" +
				"1:标注中\n" +
				"2:标注完成\n" +
				"3:审阅中\n" +
				"4:审阅完成，拒绝\n" +
				"5:标注修改中\n" +
				"6:标注完成修改，提交审阅\n" +
				"7:审阅接受，最终状态")
		case "sync":
			upgrade.Run3()
		}
		fmt.Println("参数错误", input)
	}
}
