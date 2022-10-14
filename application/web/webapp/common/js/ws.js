/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: ws.js
 * Description:
 */

class WebClient {
    constructor() {
        let ws = new WebSocket('ws://192.168.2.101/ws');
        this.ws = ws

        ws.onmessage = (e) => {
            this.recv(e.data)
        }
        ws.onopen = (e) => {
            console.log("Connection...");
            this.send("Hello");
        };
        ws.onclose = (e) => {
            console.log("Connection closed.")
        }
    }

    recv(msg) {
        console.log("R:", msg)
    }

    send(msg) {
        this.ws.send(msg)
    }

    close() {
        this.ws.close()
    }
}

function msgUserMessageSend() {

}

function msgUserMessageRecv() {

}

function msgUserMessageInit() {
    console.log("Message system initialized.");
}

function msgSetIconPopupNumber(i) {
    $("#user-message-remain-count").text(i)
}

var message_count