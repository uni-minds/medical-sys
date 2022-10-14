/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpClient.go
 */

package tools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const timeout = 10 * time.Second
const clientMax = 100

type ConnectData struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

type ClientPool struct {
	Clients chan *http.Client
	Used    int
}

var clientPool ClientPool

func init() {
	clientPool.Clients = make(chan *http.Client, clientMax)
}

func clientGet() (client *http.Client) {
	t := time.NewTicker(timeout)
	for {
		select {
		case client = <-clientPool.Clients:
			clientPool.Used += 1
			return client

		case <-t.C:
			log("e", "no more http clients available")
			return nil

		default:
			if clientPool.Used < clientMax {
				tr := &http.Transport{
					TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
					TLSHandshakeTimeout:   timeout,
					ResponseHeaderTimeout: timeout,
					ExpectContinueTimeout: timeout,
				}

				client = &http.Client{
					Transport: tr,
					Timeout:   timeout,
				}

				clientPool.Used += 1
				log("w", fmt.Sprintf("client new, %d use / %d max", clientPool.Used, clientMax))
				return client
			}
		}
	}
}
func clientRecycle(client *http.Client) {
	if client != nil {
		clientPool.Used -= 1
		clientPool.Clients <- client
	}
}

func HttpGet(url string) (data ConnectData, rs []byte, err error) {
	client := clientGet()
	defer clientRecycle(client)

	if resp, err := client.Get(url); err != nil {
		return data, nil, err
	} else if rs, err = ioutil.ReadAll(resp.Body); err != nil {
		return data, nil, err
	} else if err = resp.Body.Close(); err != nil {
		return data, rs, err
	} else {
		if err = json.Unmarshal(rs, &data); err != nil {
			log("t", "http get: not a standard response")
		}
		return data, rs, nil
	}
}

func HttpPost(url string, data interface{}, contentType string) (rdata ConnectData, bs []byte, err error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return rdata, nil, err
	}
	req.Header.Add("content-type", contentType)

	client := clientGet()
	defer clientRecycle(client)

	resp, err := client.Do(req)
	req.Body.Close()
	if err != nil {
		return rdata, nil, err
	}
	defer resp.Body.Close()

	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		return rdata, bs, err
	} else {
		if err = json.Unmarshal(bs, &rdata); err != nil {
			log("t", "http post: not a standard response")
		}
		return rdata, bs, nil
	}
}
