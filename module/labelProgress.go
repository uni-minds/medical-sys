package module

var progressData = map[int]string{
	0: "待领取",
	1: "标注中",
	2: "待审核",
	3: "审核中",
	4: "审核已退回",
	5: "待重审",
	6: "重审中",
	7: "审核完成",
}

func ProgressQuery(status int) string {
	value, ok := progressData[status]
	if ok {
		return value
	} else {
		return "状态异常"
	}
}
