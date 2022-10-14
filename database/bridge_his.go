package database

import (
	"errors"
	"fmt"
	"gitee.com/uni-minds/bridge_his"
	_ "gitee.com/uni-minds/bridge_his"
	"gitee.com/uni-minds/medical-sys/global"
	"strings"
)

var hisPort bridge_his.HisManager

func BridgeHisInit() {
	app := global.GetAppSettings()
	table := "version2"
	index := "产妇入院登记号号码"
	hisPort.Init(app.DbFileHis, table, index)
	log("t", "use his db:", app.DbFileHis, table, index)
}

func BridgeGetHisDatabaseRetrieve(code string) (data []map[string]string, err error) {
	s1 := strings.Split(code, "-")
	if len(s1) > 0 {
		code = s1[len(s1)-1]
	}

	fmt.Println("HIS <-", code)

	result, err := hisPort.Query(code)
	if len(result) > 10 {
		return nil, errors.New(fmt.Sprintf("存在过多的索引结果，请尝试完善病例号：[%s]", code))
	} else if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("未找到相关病例：[%s]", code))
	} else {
		return result, nil
	}
}
