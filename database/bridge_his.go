package database

import (
	"errors"
	"fmt"
	his_mgr "gitee.com/uni-minds/bridge-his/manager"
	"gitee.com/uni-minds/medical-sys/global"
	"strings"
)

var his his_mgr.HisManager

func BridgeHisInit() {
	dbfile, err := global.GetDbFile("his")
	if err != nil {
		log.Error(err.Error())
	}
	table := "version2"
	index := "产妇入院登记号号码"
	his.Init(dbfile, table, index)
	log.Println("his db ->", dbfile, table, index)
}

func BridgeGetHisDatabaseRetrieve(code string) (data []map[string]string, err error) {
	s1 := strings.Split(code, "-")
	if len(s1) > 0 {
		code = s1[len(s1)-1]
	}

	log.Debug(fmt.Sprintf("Search HIS:", code))

	result, err := his.Query(code)
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
