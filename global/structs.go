/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: structs.go
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
