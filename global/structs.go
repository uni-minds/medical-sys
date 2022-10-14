/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: structs.go
 * Date: 2022/08/15 13:21:15
 */

package global

import "time"

type NodeInfo struct {
	Name        string
	PubKey      string
	IP          string
	Port        int
	Mode        string
	LastTalk    time.Time
	BlockHeight int
	BlockHash   string

	SystemNodes int
	MinVotes    int
}

type NodeStatus struct {
	Name   string
	Alive  bool
	IP     string
	Height int
}
