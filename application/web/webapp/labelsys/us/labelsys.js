/*
 * THIS FILE IS PART OF MEDICAL-SYS (MSYS)
 * Filename: labelsys.js
 * Author: Liu Xiangyu (ansilxy@163.com)
 * Copyright (c) 2020.
 */

// region UI and Timer

/**
 * 用户交互类
 */
class UI {
    refs
    data
    flag

    constructor() {
        this.refs = {}
        this.data = {}
        this.flag = {}

        this.refs.popupMessageBox = $(`
            <div class="modal fade" id="modal-overlay" style="display: none;">
              <div class="modal-dialog"><div class="modal-content">
                <div class="overlay"><i class="fas fa-2x fa-sync fa-spin"></i></div>
                <div class="modal-body"><p></p></div>
            </div></div></div>`)
        $("#ui-message-window").append(this.refs.popupMessageBox)

        let mainObj = $('<div id="ui-constructor" class="row" />').height(48)
        $("#main-footer").css("padding", 0).append(mainObj)
        this.mainObj = mainObj

        this.createInfoPanel()
        this.createCopyright()

        this.message("系统正在初始化", false)
    }

    createCopyright() {
        let obj = $('<div class="col-md-2 text-center align-self-center bg-white" />').css("padding", 0)
        let t1 = `core-sys 2.1<br/>(c) 2018 - 2021`
        let t2 = `Liuxy [BUAA]<br/>Uni-Minds.com`
        let t = $('<div />').html(t1).hover(() => {
            this.refs.copyright.html(t2)
        }, () => {
            this.refs.copyright.html(t1)
        })
        obj.append(t)

        this.refs.copyright = t
        this.mainObj.append(obj)
    }

    createInfoPanel() {
        let obj = $('<div class="col-md-10 text-center align-self-center" />').css("padding", 0)
        let info = $('<div/>')
        obj.append(info)

        this.refs.infolog = info
        this.mainObj.append(obj)
    }

    message(msg, warn) {
        this.refs.infolog.text(msg)
        console.info("UI MSG:",msg)
        if (warn) {
            this.mainObj.addClass("bg-gradient-red").removeClass("bg-gradient-yellow")
        } else {
            this.mainObj.addClass("bg-gradient-yellow").removeClass("bg-gradient-red")
        }
    }

    confirm(msg) {
        return confirm(msg)
    }

    prompt(msg, def) {
        return prompt(msg, def)
    }

    alert(msg) {
        return alert(msg)
    }

    set popupContent(message) {
        console.info("UI POP:",message)
        let ref = this.refs.popupMessageBox
        if (message) {
            if (!this.flag.popupMessageShow) {
                ref.addClass("show").css("display", "block")
                ref.find("p").text(message)
                this.flag.popupMessageShow = true
            } else {
                ref.find("p").text(message)
            }
        } else {
            ref.removeClass("show").css("display", "none")
            this.flag.popupMessageShow = false
        }
    }

    backgroundFrozen(show) {
        if (!this.refs.backdrop) {
            this.refs.backdrop = $('<div class="modal-backdrop fade"></div>')
        }

        let backdrop = this.refs.backdrop

        if (show) {
            $(document.body).append(backdrop)
            backdrop.addClass("show")
        } else {
            backdrop.removeClass("show")
            backdrop.remove()
        }
    }

    backgroundUnfrozen() {
        this.backgroundFrozen(false)
    }
}

/**
 * 定时器父类
 * @param: Class
 */
class Timer {
    constructor(f) {
        this.f = f
    }

    set time(t) {
        this.t = t
    }

    get time() {
        return this.t
    }

    start(t) {
        this.stop()
        if (t) {
            this.time = t
        }

        if (this.time > 0) {
            this.timer = setInterval(this.f, this.time)
            return this.timer
        } else {
            return false
        }
    }

    stop() {
        if (this.timer) {
            clearInterval(this.timer)
            this.timer = null
            return true
        } else {
            return false
        }
    }
}

/**
 * 媒体锁
 */
class MediaLockerObj extends Timer {
    apiLock = `${apiRoot}/lock`
    lockInterval = 0

    constructor(lockInterval) {
        super(() => {
            if (this.lockInterval > 0) {
                $.post(this.apiLock).done(resp => {
                    if (resp.code !== 200) {
                        ui.message(resp.msg, true)
                    }
                })
            }
        })
        if (lockInterval > 0) {
            this.lockInterval = lockInterval
        } else {
            console.warn("media autolock disabled")
        }
    }

    lock() {
        if (this.lockInterval > 0) {
            return super.start(this.lockInterval)
        }
    }

    unlock(func) {
        if (this.lockInterval > 0) {
            let r = super.stop()
            $.ajax({
                url: this.apiLock,
                type: 'DELETE',
            }).done(resp => {
                if (resp.code === 200 && !!func) {
                    console.log("post-unlock")
                    func()
                }
            });
            return r
        }
    }
}

//endregion

//region Containers
/**
 * 容器父类
 */
class BasicContainer {
    constructor(parent) {
        let container = $("<div/>")
        $(parent.mainObj).append(container)

        this.parent = parent
        this.mainObj = container
        this.data = {}
        this.refs = {}
    }

    set window(data) {
        let obj = this.mainObj
        if (data.w) {
            obj.width(`${data.w}px`)
        }
        if (data.h) {
            obj.height(`${data.h}px`)
        }
        if (data.t === 0) {
            obj.css("top", 0)
        } else if (data.t) {
            obj.css("top", `${data.t}px`)
        }
        if (data.l === 0) {
            obj.css("left", 0)
        } else if (data.l) {
            obj.css("left", `${data.l}px`)
        }
        this.onResize()
    }

    get window() {
        let obj = this.mainObj
        return {
            w: obj.width().toFixed(3),
            h: obj.height().toFixed(3),
            l: obj.position().left.toFixed(3),
            t: obj.position().top.toFixed(3),
        }
    }

    onResize() {
    }
}

/**
 * 视频容器
 */
class VideoContainer extends BasicContainer {
    constructor(parent, id) {
        super(parent, id);
        this.mainObj.css("z-index", 1).height("100%").width("100%").addClass("lsWorkspaceOverlay").css("margin", "0 auto")
        this.data.loop = true
        this.data.duration = 0
        this.data.metaloaded = false
        this.data.currentFrame = 0
        this.timer = {}
    }

    set framesData(data) {
        this.data.frames = data
    }

    get framesData() {
        return this.data.frames
    }

    framesInit(frames) {
        let d = this.data.duration
        console.log(`frames init duration= ${d} / ${frames}f`)
        if (!frames) {
            ui.message("图片标注模式")

            console.groupCollapsed("image data prepare")

            this.framesData = [0]
            this.data.frameTime = 0
            this.data.frameCount = 0
            console.log(this.data)

            console.groupEnd()

        } else {
            ui.message("视频标注模式")

            console.groupCollapsed("video data prepare")

            let step = parseFloat((d / frames).toFixed(6));
            let t = 0;
            let data = []
            for (let i = 0; i < frames; i++) {
                data.push(parseFloat(t.toFixed(6)));
                t += step;
                //this.videoFrameTime.push((i / this.videoFrames * duration).toFixed(6));
            }
            this.framesData = data
            this.data.frameTime = step
            this.data.frameCount = frames;
            console.log(this.data)

            console.groupEnd()

        }
    }

    get videoSize() {
        let p = this.refs.player
        return {"width": p.videoWidth, "height": p.videoHeight}
    }

    createPlayer(url) {
        let player = $("<video />").css("height", "100%").css("width", "100%").attr("preload", "auto")
            .bind("ended", () => {
                this.data.loop ? this.play() : this.stop()
            })
            .bind("loadedmetadata", () => {
                console.log("loading media: target duration=", mediaDuration)
                this.timer.meta = setInterval(() => {
                    let dur = this.refs.player.duration
                    if (dur === Infinity) {
                        ui.popupContent = `载入进度：0%`
                        return
                    }

                    this.data.duration = dur
                    let progress = dur / mediaDuration * 100
                    ui.popupContent = `载入进度：${progress.toFixed(1)}%`
                    console.log("loading media: duration=", dur, "progress=", progress)

                    // 判断加载进度
                    if (progress > 99) {
                        console.log("loading media: finished")
                        clearInterval(this.timer.meta)
                        this.data.metaloaded = true
                        this.framesInit(mediaFrames)
                        // mainPanel.vc.onMetaLoaded()
                    }
                }, 100)
            });
        player.append($(`<source src=/api/v1/media/${mediaIndex}/video?type=ogv type='video/ogg; codecs="theora, vorbis"'/>`))
        player.append($(`<source src=/api/v1/media/${mediaIndex}/video?type=mp4 type="video/mp4" />`))
        // player.append($(`<source src="${url}&type=webm" type='video/webm'/>`))
        player.append($(`<p>错误：您所使用的浏览器不支持 HTML5 视频播放，请换用Chrome或Firefox浏览器（国产浏览器请切换至“急速模式”）。</p>`))

        this.refs.player = player[0];
        this.mainObj.append(player);
    }

    set loop(b) {
        // console.log("video loop status", b)
        this.data.loop = b
    }

    get loop() {
        return this.data.loop
    }

    get progress() {
        return this.refs.player.currentTime / this.refs.player.duration
    }

    set progress(p) {
        if (p < 0) {
            return
        }
        this.currentFrame = Math.floor(p * this.data.frameCount)
    }

    set currentFrame(f) {
        if (f < 0) {
            f = this.data.frameCount - 1
        } else if (f >= this.data.frameCount) {
            f = 0
        }
        this.data.currentFrame = f
        if (this.framesData) {
            this.refs.player.currentTime = this.framesData[f]
        }
    }

    get currentFrame() {
        return this.data.currentFrame
    }

    set currentTime(t) {
        // if (t < 0) {
        //     t = 0
        //
        // } else if (t > this.refs.player.duration) {
        //     t = this.refs.player.duration
        //
        // }
        // this.refs.player.currentTime = t
        // this.currentFrameUpdate()
        ui.alert("set current time 0")

    }

    get currentTime() {
        return this.refs.player.currentTime
    }

    currentFrameUpdate() {
        let p = this.refs.player
        this.data.currentFrame = Math.floor(p.currentTime * this.data.frameCount / p.duration)
    }

    get current() {
        return {frame: this.currentFrame, time: this.currentTime, progress: this.progress}
    }

    get duration() {
        return this.data.metaloaded ? this.data.duration : false
    }

    play() {
        let p = this.refs.player;
        if (p.paused || p.ended) {
            if (p.ended) {
                p.currentTime = 0;
            }
            // this.timer.play = setInterval(() => {
            //     this.currentFrameUpdate()
            //     mainPanel.cc.pageLoad(this.currentFrame)
            //     // console.log("vc:", this.currentFrame)
            // }, this.data.frameTime * 500)
            p.play();
            return true

        } else {
            this.pause();
            return false
        }
    }

    pause() {
        if (this.timer.play != null) {
            clearInterval(this.timer.play)
        }
        this.refs.player.pause()
        this.currentFrameUpdate()

    }

    stop() {
        if (this.timer.play != null) {
            clearInterval(this.timer.play)
        }
        this.refs.player.pause()
        this.currentFrame = 0

    }

    next() {
        if (!this.refs.player.paused) {
            this.pause()
        }
        this.currentFrame += 1
        return this.current
    }

    prev() {
        if (!this.refs.player.paused) {
            this.pause()
        }
        this.currentFrame -= 1
        return this.current
    }

    jumpTo(frame) {
        if (!this.refs.player.paused) {
            this.pause()
        }
        this.currentFrame = frame
        return this.current
    }

    update() {
        console.log("video update")
    }

    videoGetBufferedPercentage() {
        const buffered = this.refs.player.buffered;
        if (buffered.length) {
            ui.message(`载入进度：${buffered.end(0) * 100}%`, true)
            return buffered.end(0) * 100 / this.duration;
        } else {
            return 0
        }
    }
}

/**
 * 媒体容器，包含视频及标注容器，主要用于自适应尺寸变更及一致性保护
 */
class MediaContainer extends BasicContainer {
    constructor(parent, height) {
        super(parent, "mediaContainer");
        this.mainObj.height(height).addClass("lsWorkspaceBG")
    }
}

// endregion
/**
 * 标注容器
 */
class CanvasContainer extends BasicContainer {
    constructor(parent, id) {
        super(parent, id);
        this.init()
        this.mode = this.def.ModeDisable
        this.mainObj.css("z-index", 2).addClass("lsWorkspaceOverlay").css("left", 0).css("top", 0)
    }

    init() {
        this.page = {}
        this.pageCurrent = 0

        this.def = {}

        this.def.ModeDisable = "d"
        this.def.ModeCreate = "c"
        this.def.ModeModify = "m"
        this.def.ModePoint = "p"
        this.def.ModeEnable = "e"

        this.pageNew()
    }

    // region Pages
    get isModified() {
        return this.page.modify
    }

    set isModified(d) {
        if (this.page) {
            this.page.modify = !!d
        } else {
            alert("PANIC: page is not defined")
        }
    }

    get isHide() {
        return this.mainObj.css("display") === "none"
    }

    get pageData() {
        let data = {}
        data["clabels"] = {}
        this.page.data.forEach((value, id) => {
            switch (id) {
                case "ctime":
                    data["ctime"] = mainPanel.vc.currentTime
                    break
                case "cdescribe":
                    data["cdescribe"] = value
                    break

                case "cid":
                    data["cid"] = value
                    break

                default:
                    data.clabels[id] = value
            }
        })
        return data
    }

    set pageData(d) {
        console.log("page set data:",d)
        let data = new Map
        let label = d.clabels
        for (let id in label) {
            data.set(id, label[id])
        }
        data.set("ctime", d.ctime)
        data.set("cid", d.cid)
        data.set("cdescribe", d.cdescribe)
        this.page.data = data
    }

    get pageCurrent() {
        if (!this.page.p) {
            this.pageCurrent = 0
        }
        return this.page.p
    }

    set pageCurrent(p) {
        this.page.p = p
    }

    get pageIds() {
        let t = []
        let k = this.page.parts.keys()
        while (1) {
            let d = k.next()
            t.push(d.value)
            if (d.done) {
                break
            }
        }
        return t
    }

    pageNew() {
        this.svg.empty()
        this.isModified = false
        this.page.parts = new Map
    }

    /**
     * 加载页数据
     * @param page 页号
     */
    pageLoad(page) {
        // pageLoad(null)
        if (!page && page !== 0) {
            page = this.pageCurrent
        }
        console.groupCollapsed("canvas cont.: load page:",page)

        this.pageNew()
        this.pageCurrent = page
        this.pageData = ldata.getPage(page)
        console.groupCollapsed("page draw part",this.page.data)
        this.page.data.forEach((value, id) => {
            switch (id) {
                case "cid":
                case "ctime":
                case "cdescribe":
                    break
                default:
                    console.info("draw part:",id,value)
                    let p = new PolyPart(id, this)
                    p.pointData = value
                    p.redraw()
                    this.setPart(id, p)
            }
        })
        console.groupEnd()
        console.groupEnd()
    }

    pageSave() {
        console.groupCollapsed("page save:", this.pageCurrent, this.pageData)
        ldata.setPage(this.pageCurrent, this.pageData)
        ldata.uploadFull()
        this.pageLoad()
        console.groupEnd()
    }

    pageSetTimeLabel(id, describe, time) {
        console.warn("page set time", id, describe, time)
        this.page.data.set("cid", id)
        this.page.data.set("ctime", time)
        this.page.data.set("cdescribe", (id === "SPEC") ? describe : id)
        this.isModified = true
        this.pageSave()
    }

    pageGetTimeLabel() {
        if (this.page.data.has("cid")) {
            return [this.page.data.get("cid"), this.page.data.get("cdescribe")]
        }
    }

    pageSavePart(id, data) {
        console.log("page save part",id,data)
        this.page.data.set(id, data)
        this.pageSave()
        mainPanel.cp.setButton(id, "on")

        this.pageLoad()
    }

    hideParts() {
        console.log("hide all parts")
        this.mainObj.hide()
    }

    showParts() {
        console.log("show all parts")
        this.mainObj.show()
    }

    hasPart(id) {
        return this.page.parts.has(id)
    }

    setPart(id, obj) {
        this.page.parts.set(id, obj)
    }

    getPart(id) {
        return this.page.parts.get(id)
    }

    delPart(id) {
        console.log("cc remove part", id)
        if (this.hasPart(id)) {
            let obj = this.getPart(id)
            obj.deactivate()
            obj.remove()
            this.page.parts.delete(id)
            this.page.data.delete(id)
            this.isModified = true
            this.pageSave()
        }
    }

    getActivates() {
        let ids = new Map
        this.page.parts.forEach((v, id) => {
            if (v.isActivate) {
                ids.set(id, true)
            } else {
                ids.set(id, false)
            }
        })
        return ids
    }

    // endregion

    redo() {
        switch (this.mode) {
            case this.def.ModeCreate:
                console.log("CanvasContainer.redo()")
                console.log("this(cc).page.parts.get(this.activateId) = ", this.page.parts.get(this.activateId))
                let tmp = this.page.parts.get(this.activateId)//tmp的type为PolyPart
                let flg = tmp.pointRedo()
                if (flg === false) {
                    break
                }
                let keys = tmp.data.points.keys()
                let id = ""
                while (true) {
                    let o = keys.next()
                    console.log("o = ", o)
                    if (o.done) {
                        break
                    }
                    id = o.value
                }
                let p = tmp.data.points.get(id)
                p = tmp.WHtoXY(p)

                let obj = tmp.newCircle.attr("cx", `${p.x}`).attr("cy", `${p.y}`).attr("id", id)
                    .attr("r", 3.2).attr("fill", "red").attr("stroke", "black").attr("stroke-width", 0.5)
                    .hover(tmp.onAttention)
                    .click(tmp.pointOnClick.bind(tmp))
                    .contextmenu(tmp.pointOnContext.bind(tmp))
                tmp.mainObj.parent().append(obj);
                break
            default:
                ui.message('仅在creat模式下允许重做,当前模式为：' + this.mode, true)
        }
    }

    undo() {
        switch (this.mode) {
            case this.def.ModeCreate:
                console.log("CanvasContainer.undo()")
                console.log("this(cc).page.parts = ", this.page.parts)
                let flg = this.page.parts.get(this.activateId).pointUndo()
                if (flg) {
                    this.mainObj[0].firstElementChild.lastElementChild.remove()
                }
                break
            default:
                ui.message('仅在creat模式下允许撤销,当前模式为：' + this.mode, true)
        }
    }

    remove(id) {
        let p = this.page.parts
        switch (id) {
            case "activate":
                p.forEach((v, id) => {
                    if (v.isActivate) {
                        this.delPart(id)
                    }
                })
                break

            case "all":
                p.forEach((v, id) => {
                    v.remove()
                    p.delete(id)
                })
                break

            default:
                if (p.has(id)) {
                    p.get(id).remove()
                    p.delete(id)
                }
        }
    }

    // region Main
    get svg() {
        if (!this.refs.svg) {
            let obj = document.createElementNS(xmlns, "svg")
            let o = $(obj).attr("width", "100%").attr("height", "100%")
            this.mainObj.append(o)
            this.refs.svg = o
        }
        return this.refs.svg
    }

    set mode(m) {
        this.data.mode = m
    }

    get mode() {
        return this.data.mode
    }

    createCommon(id, color) {
        let obj = new PolyPart(id, this)
        this.page.parts.set(id, obj)
    }

    doOnClick(e) {
        switch (this.data.mode) {
            case this.def.ModeCreate:
                this.page.parts.get(this.activateId).pointCreate(this.getPosition(e))
                break
            case this.def.ModeModify:

                break
            default:
                console.log("canvas is not ready")
        }
    }

    doOnContextMenu() {
        console.log("cc right click mode:", this.mode, this)
        switch (this.data.mode) {
            case this.def.ModeCreate:
            case this.def.ModeModify:
                let obj = this.page.parts.get(this.activateId)
                if (obj.isModified) {
                    obj.confirm()
                } else {
                    obj.cancel()
                }
                this.data.mode = this.def.ModeDisable
                break

            default:
        }
    }

    onResize() {
        this.pageLoad()
    }

    getPosition(e) {
        let o = this.mainObj.offset()
        let d = this.window

        let x = (e.pageX - o.left)
        let y = (e.pageY - o.top)

        let w = parseFloat((x * 100 / d.w).toFixed(3))
        let h = parseFloat((y * 100 / d.h).toFixed(3))

        return {x, y, w, h}
    }

    activate(id, type, color) {
        this.enable()
        this.activateId = id
        if (this.hasPart(id)) {
            let obj = this.getPart(id)
            obj.activate()
            this.mode = this.def.ModeModify
        } else {
            console.log("creator part:", id, type, color)
            switch (type) {
                case "com":
                    this.createCommon(id, color)
                    break

                default:
                    alert(`Unknown type: ${type}`)
                    return
            }
            this.mode = this.def.ModeCreate
        }
        this.data.ready = true
    }

    deactivate(id) {
        this.activateId = ""
        console.log("canvas container deactivate:", id)
        if (this.hasPart(id)) {
            let obj = this.getPart(id)
            console.log("cc found part", obj)
            if (!obj.deactivate()) {
                console.log("cc remove part", id)
                this.delPart(id)
            }
        }
        this.disable()
    }

    set activateId(id) {
        this.data.activateId = id
    }

    get activateId() {
        return this.data.activateId
    }

    // endregion
    enable() {
        this.disable()
        this.mainObj.click(this.doOnClick.bind(this)).contextmenu(this.doOnContextMenu.bind(this))
        console.log("canvas enable")
        this.mode = this.def.ModeEnable
    }

    disable() {
        console.log("canvas disable")
        this.mainObj.off("click").off("contextmenu")
        this.mode = this.def.ModeDisable
    }

    hideVideo() {
        this.mainObj.css("background-color", "white")
    }
}

// region Parts
/**
 * 标注结构父类，基础描点功能
 */
class BasicPart {
    constructor(id, width, height, container) {
        console.log(`new basic, ${id}:`,container)
        this.id = id
        this.mainObj = {}
        this.container = container

        this.data = {}
        this.data.color = ldata.crfId(id).color
        this.data.width = width
        this.data.height = height
        this.data.offset = container.mainObj.offset()
        this.isActivate = true

        this.data.undo = []
        this.data.activate = false
        this.data.modified = false
        this.data.pointPick = ""
        this.data.pointCount = 0

        this.data.points = new Map
        this.data.pointRefs = new Map
    }

    // region 数据操作
    get isModified() {
        return this.data.modified
    }

    set isModified(b) {
        this.data.modified = !!b
    }

    get isActivate() {
        return this.data.activate
    }

    set isActivate(b) {
        this.data.activate = b
    }

    set pointData(d) {
        d = d.cPoints
        for (const id in d) {
            let a = d[id]
            this.data.points.set(id, {w: parseFloat(a[0].toFixed(3)), h: parseFloat(a[1].toFixed(3))})
        }
    }

    get pointData() {
        console.log(`get points data`, this.data.points)
        let count = 0
        let clabel = {}
        clabel.cid = this.id
        clabel.cPoints = {}
        this.data.points.forEach((p) => {
            let id = `${this.id}_${count}`
            clabel.cPoints[id] = [p.w, p.h]
            count++
        })
        return clabel
    }

    set pointPick(id) {
        this.data.pointPick = id
    }

    get pointPick() {
        return this.data.pointPick
    }

    get pointString() {
        let str = ""
        this.data.points.forEach((p) => {
            p = this.WHtoXY(p)
            str += `${p.x},${p.y} `
        })
        return str
    }

    set resolution(r) {
        this.data.width = r.w
        this.data.height = r.h
        this.onResize()
    }

    // endregion

    activate() {
        console.log(`basic part activate: ${this.id}`)
        if (!this.isActivate) {
            this.isActivate = true
            this.showPoints()
        }
    }

    /**
     *
     * @returns {boolean} obj是否存在
     */
    deactivate() {
        console.log(`basic part deactivate: ${this.id}, activate: ${this.data.activate}`)
        this.isActivate = false
        this.hidePoints()
        console.log("basic part data:", this.data)
        if (this.data.points.size > 0) {
            this.isModified ? this.confirm() : this.cancel()
            return true
        } else {
            this.destroy()
            return false
        }
    }

    showPoints() {
        this.data.points.forEach((d, id) => {
            this.data.pointRefs.set(id, this.pointCreate(d, id))
        })
    }

    hidePoints() {
        if (this.isModified) {
            this.confirm()
        }
        // console.log("hide points")
        if (this.data.pointRefs.size > 0) {
            this.data.pointRefs.forEach((obj, id) => {
                obj.remove()
                this.data.pointRefs.delete(id)
            })
        }
    }

    moveTop() {
        this.mainObj.parent().children().last().after(this.mainObj)
    }

    moveBottom() {
        this.mainObj.parent().children().first().before(this.mainObj)
    }

    confirm() {
        console.log("basic part confirm:", this.id)
        this.isModified = false
        this.container.pageSavePart(this.id, this.pointData)
    }

    cancel() {
        console.log(`basic part cancel:`, this.id)
    }

    destroy() {
        console.log("basic part destroy:", this.id)
        this.mainObj.remove()
    }

    redraw() {
        // console.log("basic redraw")
    }

    remove() {
        this.hidePoints()
        this.mainObj.remove()
    }

    // region 点操作函数
    /**
     * 创建锚点
     * @param p 坐标数据 {w:p.w,h:p.h}
     * @param id ID
     * @returns {*} object
     */
    pointCreate(p, id) {
        this.isModified = true
        while (!id || this.data.pointRefs.has(id)) {
            id = `${this.id}_${this.data.pointCount}`
            this.data.pointCount++
        }
        console.log(`point create: ${id}`)
        p = this.WHtoXY(p)
        let obj = this.newCircle.attr("cx", `${p.x}`).attr("cy", `${p.y}`).attr("id", id)
            .attr("r", 3.2).attr("fill", "red").attr("stroke", "black").attr("stroke-width", 0.5)
            .hover(this.onAttention)
            .click(this.pointOnClick.bind(this))
            .contextmenu(this.pointOnContext.bind(this))
        this.mainObj.parent().append(obj);
        this.data.points.set(id, {w: p.w, h: p.h})
        this.data.undo = [];
        return obj
    }

    pointMove(p, id) {
        this.isModified = true
        if (!id) {
            id = this.pointPick
        }
        // console.log(`point move: ${id} ->`, p)
        this.data.pointRefs.get(id).attr("cx", p.x).attr('cy', p.y)
    }

    pointSave(p, id) {
        this.isModified = true
        if (!id) {
            id = this.pointPick
        }
        this.data.points.set(id, {w: p.w, h: p.h})
    }

    pointRemove(id) {
        console.log("BasicPart.pointRemove()")
        this.isModified = true
        if (!id) {
            id = this.pointPick
        }
        console.log(`point remove: ${id}`)
        let c = this.data
        //c.pointRefs.get(id).remove()
        //c.points.get(id).remove()
        //c.pointRefs.delete(id)
        c.undo.push([id, c.points.get(id)])
        c.points.delete(id)
        console.log("c = this(bp).data = ", this.data)
    }

    pointCancel(id) {
        if (!id) {
            id = this.pointPick
        }
        let p = this.WHtoXY(this.data.points.get(id))
        this.mainObj.parent().off("mousemove")
        this.data.pointRefs.get(id).attr("cx", p.x).attr("cy", p.y)
            .off("mousemove").off("mouseleave")
            .hover(this.onAttention);
        this.pointPick = ""
    }

    pointUndo() {
        console.log("BasicPart.pointUndo()")
        //let keys = this.data.pointRefs.keys()
        let keys = this.data.points.keys()
        let id = ""
        while (true) {
            let o = keys.next()
            if (o.done) {
                break
            }
            id = o.value
        }
        console.log("point latest id:", id)
        if (id !== "") {
            console.log("id = ", id)
            this.pointRemove(id)
            return id
        } else {
            ui.message('无法再撤销。', false);
            return false
        }
    }

    pointRedo() {
        console.log("BasicPart.pointUndo()")
        let evt = this.data.undo.pop()
        console.log("evt = ", evt)
        if (evt == null) {
            ui.message('无法再重做。', false);
            return false;
        }
        this.data.points.set(evt[0], evt[1])
        console.log("this(bp).data = ", this.data)

    }

    //endregion
    // region 事件响应
    onResize() {
        this.data.offset = this.container.mainObj.offset()
    }

    onAttention() {
        let obj = $(this)
        let cs = obj.attr("fill")
        let cf = obj.attr("stroke")
        obj.attr('fill', cf).attr('stroke', cs);
    }

    pointOnClick(e) {
        e.stopPropagation();
        const obj = $(e.target);
        if (this.pointPick) {
            // 已选择点，放下
            obj.off("mousemove").off("mouseleave")
                .hover(this.onAttention);
            this.mainObj.parent().off("mousemove")

            this.pointSave(this.getPosition(e))
            this.pointPick = ""
            this.redraw()

        } else {
            // 未选择点，拾起
            this.pointPick = e.target.id
            obj.off("mouseenter").off("mouseleave")
                .on("mousemove", (e) => {
                    this.pointMove(this.getPosition(e))
                })
            this.mainObj.parent().on("mousemove", (e) => {
                this.pointMove(this.getPosition(e))
            })
        }
    }

    pointOnContext(e) {
        e.stopPropagation();
        e.preventDefault();
        this.pointCancel(e.target.id)
    }

    //endregion
    // region 工具函数
    WHtoXY(p) {
        p.x = parseFloat((p.w * this.data.width / 100).toFixed(3))
        p.y = parseFloat((p.h * this.data.height / 100).toFixed(3))
        return p
    }

    XYtoWH(p) {
        p.w = parseFloat((p.x * 100 / this.data.width).toFixed(3))
        p.h = parseFloat((p.y * 100 / this.data.height).toFixed(3))
        return p
    }

    getPosition(e) {
        let x = (e.pageX - this.data.offset.left)
        let y = (e.pageY - this.data.offset.top)
        return this.XYtoWH({x, y})
    }

    get newCircle() {
        let obj = document.createElementNS(xmlns, "circle")
        return $(obj)
    }

    // endregion
}

/**
 * 多边形标注结构
 */
class PolyPart extends BasicPart {
    constructor(id, container) {
        let d = container.window
        super(id, d.w, d.h, container);
        this.mainObj = this.newPolygon.attr("id", id).attr("fill", this.data.color)
            .attr("stroke-width", 1).attr("stroke", "black").attr("opacity", 0.5)
            .hover(this.onAttention)
            .click(this.moveBottom.bind(this))
            .contextmenu((e) => {
                e.stopPropagation()
                e.preventDefault()
                this.isActivate ? this.deactivate() : this.activate()
            })
        container.svg.append(this.mainObj)
    }

    get newPolygon() {
        let obj = document.createElementNS(xmlns, "polygon")
        return $(obj)
    }

    redraw() {
        super.redraw();
        this.mainObj.attr("points", this.pointString)
    }
}

// endregion

//region Panels
/**
 * 面板父类
 */
class BasicPanel {
    constructor(parent) {
        this.mainObj = $("<div />")
        this.parent = parent
        parent.append(this.mainObj)
    }
}

/**
 * 系统面板，右侧隐藏
 */
class SystemPanel extends BasicPanel {
    constructor(parent) {
        super(parent)
        let ul = $('<ul class="nav nav-pills nav-sidebar flex-column nav-flat" data-widget="treeview" role="menu" data-accordion="false"/>')
        let nav = $('<nav class="mt-0" />').append(ul)
        this.mainList = ul
        this.mainObj.addClass("sidebar").append(nav)
        // 用户备注
        this.memo = new ButtonGroup(this, "备注内容", true, true)
        this.memo.addMemo(5, "输入备注信息……");
        this.system = new SystemButtonGroup(this, "系统工具", true)
        this.admin = new AdminButtonGroup(this, "管理工具", false)
    }
}

/**
 * 主面板，功能主体
 */
class MainPanel extends BasicPanel {
    constructor(parent) {
        super(parent)
        this.data = {}
        this.data.select = new Map
        this.data.duration = 0
        this.data.videomode = !!mediaFrames
        this.refs = {}
        this.timer = {}
        this.def = {}
        this.mainObj.height("100%").width("100%")
        this.initResizeListeners()
        this.initKeyPressListeners()

        this.cp = new CanvasPanel($("#creator"))
        this.sp = new SystemPanel($("#system"))
        this.mc = new MediaContainer(this, "calc(100% - 70px)")

        this.def.ModeInput = "i"
        this.def.ModeControl = "c"
        this.setControlMode()
    }

    close() {
        console.log("main panel:exit")
    }

    setInputMode() {
        this.mode = this.def.ModeInput
    }

    setControlMode() {
        this.mode = this.def.ModeControl
    }

    set mode(m) {
        this.data.mode = m
    }

    get mode() {
        return this.data.mode
    }

    setValue(group, value, crf) {
        // console.log("button set value", group, value, crf)
        switch (crf.domain) {
            case "global":
                ldata.setGlobal(group, crf.value)
                break
            case "frame":
                switch (group) {
                    case "t":
                        this.cc.pageSetTimeLabel(crf.id, value, this.vc.currentTime)
                        break
                    default:
                }
                console.log("set value/F")
                this.currentPageSave()
                break
            default:

        }
        this.data.select.set(group, value)
    }

    getValue(group) {
        if (this.data.select.has(group)) {
            return this.data.select.get(group)
        } else {
            return null
        }
    }

    hasValue(group) {
        return this.data.select.has(group)
    }

    delValue(group) {
        this.data.select.delete(group)
    }

    getPart(id) {
        return this.cc.getPart(id)
    }

    hasPart(id) {
        return this.cc.hasPart(id)
    }

    delPart(id) {
        console.log('remove part', id)
        switch (id) {
            case "all":
                if (ui.confirm(`确认删除本帧全部标注结构？`)) {
                    this.cc.getActivates().forEach((isAct, id) => {
                        this.cc.delPart(id)
                        this.cp.setButton(id, "off")

                    })
                }
                break
            case "activate":
                this.cc.getActivates().forEach((isAct, id) => {
                    if (isAct && ui.confirm(`确认删除标注结构：${id}？`)) {
                        this.cc.delPart(id)
                        this.cp.setButton(id, "off")
                    }
                })
                break
            default:
                if (this.hasPart(id)) {
                    this.cc.delPart(id)
                    this.cp.setButton(id, "off")
                }
        }
    }

    buttonOnClickL(btnId) {
        console.group("L_Click:",btnId)
        let crf = ldata.crfId(btnId)
        let group = ldata.groupId(crf.group)
        let gid = crf.group

        let vOld = this.getValue(gid)
        let vNew = btnId

        console.log("Value:",vOld,"->",vNew)

        if (group.gradio) {
            let custom = false
            if (crf.value === "INPUT") {
                custom = true
                vNew = ui.prompt("请输入状态", (vOld) ? vOld : vNew)
                if (!vNew) {
                    return
                }
            }

            if (vOld) {
                if (vNew && vNew !== vOld && ui.confirm("确认调整")) {
                    this.cp.setButton(ldata.crfId(vOld) ? vOld : "SPEC", "off")
                    this.cp.setButton(btnId, "on", (custom) ? vNew : "")
                    this.setValue(gid, vNew, crf)
                }
            } else {
                this.cp.setButton(btnId, "on", (custom) ? vNew : "")
                this.setValue(gid, vNew, crf)
            }
        } else {
            console.log(`main panel ${vOld} => ${vNew}`)
            if (!this.cc.pageGetTimeLabel()) {
                ui.alert("请先进行时间标注")
                return
            }

            if (vOld) {
                this.cc.deactivate(vOld)
                this.cp.setButton(vOld, (this.hasPart(vOld)) ? "on" : "off")

                if (vOld === vNew) {
                    this.delValue(gid)
                    return
                }
            }
            this.cp.setButton(btnId, "hold")
            this.cc.activate(btnId, crf.type, crf.color)
            this.setValue(gid, btnId, crf)
        }
        console.groupEnd()
    }

    buttonOnClickR(id) {
        console.group("R_click:", id)
        console.groupEnd()
    }

    buttonSetStatus(id, status) {

    }

    set page(page) {
        this.skipToFrame(page)
    }

    /**
     * 窗口调整监听
     */
    initResizeListeners() {
        const callback = (mutationsList, observer) => {
            for (let mutation of mutationsList) {
                if (mutation.type === 'childList') {
                    console.log('A child node has been added or removed.');
                } else if (mutation.type === 'attributes') {
                    // console.log('The ' + mutation.attributeName + ' attribute was modified.');
                    if (mutation.attributeName === "style") {
                        this.onResize()
                    }
                }
            }
        };
        const observer = new MutationObserver(callback);
        observer.observe($(".content-wrapper")[0], {attributes: true, childList: false, subtree: false});

        window.addEventListener("resize", (e) => {
            this.onResize()
        })
    }

    /**
     * 按键监听器
     */
    initKeyPressListeners() {
        $(document).on("keydown", (e) => {
            switch (this.mode) {
                case this.def.ModeControl:
                    this.onKeyDown(e.keyCode || e.which || e.charCode)
                    break

                case this.def.ModeInput:
                    break
            }
        })
    }

    onCRFFinishDownload() {
        this.cp.init()
    }

    onGlobalFinishDownload() {
        let q = ldata.getGlobal("q")
        if (q && q > 0) {
            let id = `FQ${q}`
            this.cp.setButton(id, "on")
        }
        this.skipToFrame(0)
    }

    onPageFinishDownload() {

    }

    onMouseScroll(e) {
        let delta = (e.originalEvent.wheelDelta && (e.originalEvent.wheelDelta > 0 ? 1 : -1)) ||  // chrome & ie
            (e.originalEvent.detail && (e.originalEvent.detail > 0 ? -1 : 1));// firefox
        if (delta > 0) {
            mainPanel.prevFrame()
        } else if (delta < 0) {
            mainPanel.nextFrame()
        }
    }

    onMetaLoaded() {
        console.groupCollapsed("main panel on metaloaded")
        console.info(`m.panel set duration: ${this.data.duration} -> ${this.vc.duration}`)
        this.data.duration = this.format(this.vc.duration)
        console.log("create progress bar")
        this.progressBarCreate()
        console.log("create progress tag")
        this.progressTagCreate()
        console.log("create control bar")
        this.controlBarCreate()

        let vcsize = this.vc.videoSize
        this.scale = vcsize.width / vcsize.height
        console.log("resize window")
        this.onResize()
        $(this.mainObj).on("mousewheel DOMMouseScroll", this.onMouseScroll)
        console.info("Timer.buffer start")
        this.timer.buffer = new Timer(() => {
            if (mediaDuration === 0 || this.progressBufferMonitor()) {
                console.info("Timer.buffer stop")
                this.timer.buffer.stop()
                ui.message("完成视频载入")
                ui.popupContent = ""
                ui.backgroundUnfrozen()
            }
        })
        this.timer.buffer.start(100)
        console.groupEnd()
    }

    onResize() {
        let d = this.vc.window
        let vScale = this.scale

        let mcs = this.mc.window
        let containerH = mcs.h
        let containerW = mcs.w
        let cScale = containerW / containerH

        if (vScale >= cScale) {
            d.w = containerW
            d.h = containerW / vScale
            d.l = 0
            d.t = (containerH - d.h) / 2
        } else {
            d.h = containerH
            d.w = containerH * vScale
            d.l = (containerW - d.w) / 2
            d.t = 0
        }
        this.vc.window = d
        this.cc.window = d
    }

    onKeyDown(code) {
        switch (code) {
            // ESC
            case 27:
                this.cc.pageLoad()
                break;
            // P
            case 80:
                this.play()
                break;
            // Enter
            case 13:
                break;
            // Z
            case 90://z
                this.cc.undo();
                break;
            // X, U
            case 88:
            case 85:
                this.cc.redo();
                break;
            // Del, D, Backspace
            case 46:
            case 68:
            case 8:
                mainPanel.delPart("activate")
                break;
            // Q
            case 81:
                mainPanel.delPart("all");
                break;
            // H
            case 72:
                if (this.cc.isHide) {
                    this.cc.showParts()
                } else {
                    this.cc.hideParts()
                }
                break;
            // Home
            case 36:
                this.skipToFrame(0)
                break;
            // End
            case 35:
                this.skipToFrame(-1)
                break;
            // left
            case 37:
                this.prevFrame()
                break;
            // right
            case 39:
                this.nextFrame();
                break;
            // <
            case 188:
                this.prevLabel()
                break;
            // >
            case 190:
                this.nextLabel()
                break;
            default:
            // console.log(code.toString());
        }
    }

    init() {
        // init video
        this.vc = new VideoContainer(this.mc)
        this.cc = new CanvasContainer(this.mc)

        this.vc.createPlayer(mediaIndex)

        ui.backgroundFrozen("数据载入中……")

        this.timer.metaload = new Timer(() => {
            if (this.vc.data.metaloaded) {
                console.info("Timer.metaload stop")
                this.timer.metaload.stop()
                this.onMetaLoaded()
                ui.popupContent = "数据预载结束"
                ui.message("数据预载结束")
            } else {
                ui.message("数据载入未完成...", true)
            }
        })
        console.info("Timer.metaload start")
        this.timer.metaload.start(100)

        // setTimeout(this.cp.init.bind(this.cp), 100)
    }

    play() {
        this.currentPageSave()
        if (this.vc.play()) {
            this.refs.buttonPlay.html('<i class="fas fa-pause-circle"></i>')
            this.cc.pageNew()
            this.osdAutoUpdateStart()

        } else {
            this.refs.buttonPlay.html('<i class="fas fa-play-circle"></i>')
            this.vc.currentFrameUpdate()
            // console.log("current frame",this.vc.currentTime,this.vc.currentFrame)
            this.update(this.vc.current)
            this.osdAutoUpdateStop()

        }
    }

    stop() {
        this.refs.buttonPlay.html('<i class="fas fa-play-circle"></i>')
        this.vc.stop()
        this.osdAutoUpdateStop();
        this.osdManualUpdate();
    }

    skipToFrame(frame) {
        this.currentPageSave()
        this.update(this.vc.jumpTo(frame))
    }

    prevFrame() {
        this.currentPageSave()
        this.update(this.vc.prev())
    }

    nextFrame() {
        this.currentPageSave()
        this.update(this.vc.next())
    }

    update(d) {
        if (d) {
            console.groupCollapsed("main panel: frame to", d.frame)
            this.progressPlay = d.progress
            this.currentPageLoad(d.frame)
            console.groupEnd()
        }
    }

    prevLabel() {
        this.currentPageSave()
        let vc = this.vc
        let f = ldata.before(vc.currentFrame)
        if (!!f) {
            vc.currentFrame = f
            this.update(vc.current)
        } else {
            ui.message("无标注信息")
        }
    }

    nextLabel() {
        this.currentPageSave()
        let vc = this.vc
        let f = ldata.after(vc.currentFrame)
        if (!!f) {
            vc.currentFrame = f
            this.update(vc.current)
        } else {
            ui.message("无标注信息")
        }
    }

    currentPageSave() {
        let cc = this.cc
        if (cc.isModified) {
            ldata.setPage(cc.pageCurrent, cc.pageData)
            console.log("current page save")
        }
    }

    currentPageLoad(page) {
        this.cc.pageLoad(page)
        this.cp.pageLoad(page)
    }

    format(time) {
        time = time.toFixed(3);
        let h = time < 3600 ? 0 : parseInt(time / 3600)
        time = time - h * 3600;
        let m = time < 60 ? 0 : parseInt(time / 60)
        time = time - m * 60;
        let s = time < 1 ? 0 : parseInt(time)
        return `${h}:${m}:${s}.${this.pad(parseInt(1000 * (time - s)), 3)}`;
    }

    pad(num, n) {
        let len = num.toString().length;
        while (len < n) {
            num = "0" + num;
            len++;
        }
        return num;
    }

    osdAutoUpdateStart() {
        let t = this.timer.osdUpdate
        if (t) {
            // console.info("Timer.osd stop 1")
            t.stop()
        } else {
            t = new Timer(this.osdManualUpdate.bind(this))
            this.timer.osdUpdate = t
        }
        // console.info("Timer.osd start")
        t.start(this.vc.data.frameTime*500)
    }

    osdAutoUpdateStop() {
        // console.info("Timer.osd stop")
        this.timer.osdUpdate.stop()
    }

    osdManualUpdate() {
        this.vc.currentFrameUpdate()
        this.cc.pageLoad(this.vc.currentFrame)
        this.progressPlay = this.vc.progress
    }

    controlBarUpdate() {
        let t = this.format(this.vc.refs.player.currentTime)
        this.refs.videoPlaytime.text(`${t} / ${this.data.duration}`)
    }

    controlBarCreate() {
        let videoControlCenter = $("<div class='row'/>").height(40);
        let videoPlaytime = $("<div/>").text(`0:00.000 / ${this.data.duration}`);
        // Left control
        let btnPlay = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-play-circle"></i>')
        let btnStop = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-stop-circle"></i>')
        let btnPFra = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-arrow-alt-circle-left"></i>')
        let btnNFra = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-arrow-alt-circle-right"></i>')
        let btnPLab = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-chevron-circle-left"></i>')
        let btnNLab = $("<button class='btn btn-info videoBtn'/>").html('<i class="fas fa-chevron-circle-right"></i>')
        let LA = $("<div class='col-md-5 btn-group videoControl'/>").append(btnPlay).append(btnStop).append(btnPFra).append(btnNFra).append(btnPLab).append(btnNLab)
        // Right control
        let btnLoop = $("<input />").attr("type", "checkbox").attr("name", "videoLoopEnable").attr("value", "loop")
        let labelLoop = $("<div />").append(btnLoop).append(`  循环`)
        let RA = $("<div class='col-md-1 align-self-center text-center'/>").append(labelLoop);
        // Middle
        let MA = $("<div class='col-md-6 videoControl align-self-center text-center' />").append(videoPlaytime);

        videoControlCenter.append(LA).append(MA).append(RA);

        if (this.data.videomode) {
            btnPlay.click(this.play.bind(this))
            btnStop.click(this.stop.bind(this))
            btnPFra.click(this.prevFrame.bind(this))
            btnNFra.click(this.nextFrame.bind(this))
            btnPLab.click(this.prevLabel.bind(this))
            btnNLab.click(this.nextLabel.bind(this))
            btnLoop.attr("checked", true).click(() => {
                this.vc.loop = !this.vc.loop
            })
        } else {
            btnPlay.addClass("disabled")
            btnStop.addClass("disabled")
            btnPFra.addClass("disabled")
            btnNFra.addClass("disabled")
            btnPLab.addClass("disabled")
            btnNLab.addClass("disabled")
            btnLoop.addClass("disabled")

        }

        this.mainObj.append(videoControlCenter);

        this.refs.buttonPlay = btnPlay
        this.refs.buttonStop = btnStop
        this.refs.buttonPrev = btnPFra
        this.refs.buttonNext = btnNFra
        this.refs.btnNextLabel = btnNLab
        this.refs.btnPrevLabel = btnPLab
        this.refs.buttonLoop = btnLoop
        this.refs.videoPlaytime = videoPlaytime
    }

    set progressBuffer(percent) {
        this.refs.pgBuffer.css("width", `${percent}%`);
    }

    set progressPlay(percent) {
        let p = `${(percent * 100).toFixed(0)}%`
        if (this.refs.pgPlayStr && this.refs.pgPlay) {
            this.refs.pgPlayStr.text(p);//进度条文字进度
            this.refs.pgPlay.css("width", p);//调整控制条长度
            this.controlBarUpdate()
        }
    }

    progressBarCreate() {
        let bg = $("<div class='progress videoProgressBackground'/>")
        let pp = $("<div class='progress-bar bg-primary progress-bar-striped videoProgressFront'/>").css("width", 0);
        let pb = $("<div class='progress-bar bg-warning videoProgressFront'/>").css("width", 0);
        let str = $("<span />").addClass('videoProgressStr').text("0%");

        pp.append(str);
        bg.append(pb).append(pp).click((e) => {
            let obj = this.refs.bg
            let x = e.pageX - obj.offset().left
            let w = x / obj.width()

            this.progressPlay = w
            this.vc.progress = w
        })

        this.mainObj.append(bg);
        this.refs.bg = bg
        this.refs.pgPlay = pp
        this.refs.pgPlayStr = str
        this.refs.pgBuffer = pb
    }

    progressTagCreate() {
        let bg = $('<div class="bg-gray"/>').height(10)
        this.mainObj.append(bg)
    }

    progressBufferMonitor() {
        let buf = this.vc.videoGetBufferedPercentage();
        this.progressBuffer = buf;
        return (buf >= 100)
    }

    progressTagInsert(frame) {
        console.log("insert progress tag:", frame)
    }

    progressTagRemove(frame) {
        console.log("remove progress tag:", frame)
    }
}

//endregion
/**
 * 标注功能面板，左侧
 */
class CanvasPanel extends BasicPanel {
    constructor(parent) {
        super(parent)
        this.groups = new Map
        let ul = $('<ul class="nav nav-pills nav-sidebar flex-column nav-flat" data-widget="treeview" role="menu" data-accordion="false"/>')
        let nav = $('<nav class="mt-0" />')
        nav.append(ul)
        this.mainList = ul
        this.mainObj.addClass("sidebar").append(nav)
        this.data = {}
        this.data.buttonStatus = new Map
    }

    init() {
        // 标签信息
        ldata.crfGroups.forEach(v => {
            let color = ""
            switch (v.name) {
                case "通用标签":
                    color = "yellow"
                    break
                case "异常标签":
                    color = "red"
                    break
            }
            let obj = new ButtonGroup(this, v.name, v.gopen, color)
            obj.addButtons(ldata.crfGroup(v.group));
            this.groups.set(v.group, obj)
        })
        this.pageLoad(0)
    }

    pageLoad(page) {
        console.groupCollapsed("canvas panel: load page:",page)
        this.pageClear()
        let d = ldata.getPage(page)
        console.info("page data:",d)
        if (d && d.clabels) {
            for (const id in d.clabels) {
                console.log("set button", id)
                this.setButton(id, "on")
            }

            if (d.cid === "SPEC") {
                console.log("set button", d.cid, d.describe)
                this.setButton(d.cid, "on", d.describe)

            } else {
                console.log("set button", d.cid)
                this.setButton(d.cid, "on",)

            }
        }
        console.groupEnd()
    }

    pageClear() {
        console.groupCollapsed("page clear: button status")
        this.data.buttonStatus.forEach((v, id) => {
            console.log("check:", id, v)
            switch (v) {
                case "on":
                case "hold":
                    let crf = ldata.crfId(id)
                    if (crf.domain !== "global") {
                        console.info(`clear: ${id}`, crf)
                        this.setButton(id, "off", null)
                    } else {
                        console.log("ignore global:", id)
                    }
            }
        })
        console.groupEnd()
    }

    setButton(id, status, text) {
        let textOverride = text?`text<-${text}`:""
        console.log(`cp set button id<-${id} status<-${status}`+textOverride)
        if (!id) {
            return
        }
        let crf = ldata.crfId(id)
        if (!crf) {
            return
        }
        let target = this.groups.get(crf.group)
        switch (status) {
            case "hold":
                target.hold(id)
                this.data.buttonStatus.set(id, status)
                break

            case "off":
                target.off(id)
                target.text(id, crf.name)
                this.data.buttonStatus.delete(id)
                break

            case "on":
            default:
                target.on(id)
                this.data.buttonStatus.set(id, status)
                if (text) {
                    target.text(id, text)
                }
                break
        }
    }
}

//region Buttons
/**
 * 按钮组父类
 */
class ButtonGroup {
    constructor(parent, title, open, circleColor) {
        let nav = $('<li class="nav-item has-treeview" />')
        if (open) {
            nav.addClass("menu-open")
        }
        let ccolor = "text-info"
        switch (circleColor) {
            case "red":
                ccolor = "text-danger"
                break
            case "yellow":
                ccolor = "text-warning"
                break
        }
        let header = $(`<a href="#" class="nav-link"><i class="nav-icon far fa-circle ${ccolor}"></i><i class="nav-icon far fa-circle-thin "></i><p>${title}<i class="fas fa-angle-left right"></i></p></a>`)
        let tree = $('<ul class="nav nav-treeview"></ul>')

        nav.append(header).append(tree)
        parent.mainList.append(nav)

        this.parent = parent
        this.mainObj = tree
        this.data = {}
        this.crfs = new Map
        this.refs = new Map
        // this.data.radio = !!radio
        this.data.select = ""
        this.data.c = 0
    }

    addButtons(dat) {
        dat.forEach((d) => {
            let obj = this.newButton(d.name, d.id).click(() => {
                mainPanel.buttonOnClickL(d.id)
            }).contextmenu(() => {
                mainPanel.buttonOnClickR(d.id)
            })

            if (this.data.c > 2 || this.data.c === 0) {
                this.data.c = 0
                this.mainLine = this.newLine()
                this.mainObj.append(this.mainLine)
            }

            this.mainLine.append(obj)
            this.refs.set(d.id, obj)
            this.data.c++
        })
    }

    addMemo(rows, placeholder) {
        rows = (rows) ? rows : 3
        let urlMemo = apiRoot + '/label/memo'
        let objMemo = $(`<textarea class="form-control" id="usermemo" rows="${rows}" placeholder="${placeholder}"></textarea>`)
        let btn = this.newButton("保存备注", "save memo", 12, "#ff851b").click(() => {
            let d = {}
            d.media = mediaIndex;
            d.data = objMemo.val();
            d.direction = 'upload';
            let data = JSON.stringify(d)

            console.warn("memo -> server")

            $.post(urlMemo, data).done(resp => {
                console.groupCollapsed("memo detail")
                console.log("data:", data)
                console.log("resp:", resp)
                console.groupEnd()
                if (resp.code === 200) {
                    btn.css("background-color", "#28a745")
                    setTimeout(() => {
                        btn.css("background-color", "#ff851b")
                    }, 1000)

                    $.get(urlMemo).done(resp => {
                        if (resp.code === 200) {
                            objMemo.val(resp.data)
                        }
                    })

                } else {
                    ui.message(resp.msg, true)
                }
            })

        })

        $.get(urlMemo).done(resp => {
            if (resp.code === 200) {
                objMemo.val(resp.data)
            }
        })

        this.mainObj.append(objMemo).append(btn)
        this.refs.set("memo", objMemo)
    }

    newLine() {
        return $('<div />').addClass("row").css("margin", 0)
    }

    newButton(name, id, buttonLength, bgColor, wordCut) {
        buttonLength = (buttonLength) ? buttonLength : 4
        wordCut = (wordCut) ? wordCut : 5
        const content = (name.length > wordCut) ? `${name.substring(0, wordCut - 1)}…` : name
        const obj = $(`<button class="crfBtn crfBtnNone col-md-${buttonLength}" title="${name}"/>`)
            .html(`<div id="name">${content}</div><div id="id" class="crfBtnId">${id}</div>`)
        if (bgColor) {
            obj.css("background-color", bgColor)
        }

        return obj
    }

    off(id) {
        if (id) {
            // this.selectId = ""
            this.refs.get(id).removeClass("crfBtnHold")
                .addClass("crfBtnNone").css("background-color", "")
        }
    }

    on(id) {
        if (id) {
            // this.selectId = id
            let crf = ldata.crfId(id)
            this.refs.get(id).removeClass("crfBtnNone").removeClass("crfBtnHold")
                .css("background-color", (crf.color) ? crf.color : ldata.groupId(crf.group).color)
        }
    }

    hold(id) {
        if (id) {
            // this.selectId = id
            let crf = ldata.crfId(id)
            this.refs.get(id).removeClass("crfBtnNone")
                .addClass("crfBtnHold").css("background-color", (crf.color) ? crf.color : ldata.groupId(crf.group).color)
        }
    }

    text(id, txt) {
        if (id) {
            this.refs.get(id).children("#name").text(txt)
        }
    }
}

/**
 * 系统按钮组
 */
class SystemButtonGroup extends ButtonGroup {
    constructor(parent, title, open) {
        super(parent, title, open);
        const senddata = JSON.stringify({"media": mediaIndex})
        switch (urlPort) {
            case "author":
                let obj = this.newButton("提交审核", "submit", 12, "#67afe5").click(function () {
                    if (ui.confirm("确认提交审核？")) {
                        $.post(apiRoot + '/label/author?do=submit', senddata).done(resp => {
                            if (resp.code === 200 && resp.data === "exit") {
                                LabelToolSystemExit()
                            } else {
                                ui.message(resp.msg, resp.code !== 200)
                            }
                        })
                    }
                })
                this.mainObj.append(obj)
                break

            case "review":
                let btnReject = this.newButton("驳回", "reject", 6, "#cc6666").click(function () {
                    if (ui.confirm("确认驳回？")) {
                        $.post(apiRoot + '/label/review?do=reject', senddata).done(resp => {
                            if (resp.code === 200 && resp.data === "exit") {
                                LabelToolSystemExit()
                            } else {
                                ui.message(resp.msg, resp.code !== 200)
                            }
                        })
                    }
                })
                let btnConfirm = this.newButton("通过", "confirm", 6, "#99ffcc").click(function () {
                    if (ui.confirm("确认通过审核？")) {
                        $.post(apiRoot + '/label/review?do=confirm', senddata).done(resp => {
                            if (resp.code === 200 && resp.data === "exit") {
                                LabelToolSystemExit()
                            } else {
                                ui.message(resp.msg, resp.code !== 200)
                            }
                        })
                    }
                })

                this.mainObj.append(btnConfirm).append(btnReject)
                break
        }
    }
}

/**
 * 管理员按钮组
 */
class AdminButtonGroup extends ButtonGroup {
    constructor(parent, title, open) {
        super(parent, title, open);
        let obj = this.newButton("清空全部标注", "drop", 12, "#dc3545").click(function () {
            if (ui.confirm("警告！确认后将清空本媒体对应的全部标注数据")) {
                ui.message("未授权操作", true)
            }
        })
        this.mainObj.append(obj)

        obj = this.newButton("导出标注", "export", 6,).click(function () {
            ui.prompt("标注内容", ldata.raw)
        })
        this.mainObj.append(obj)

        obj = this.newButton("导入标注", "import", 6,).click(function () {
            let data = ui.prompt("数据内容", ldata.raw)
            if (data) {
                ldata.raw = data

                console.log("up 2")
                ldata.uploadFull()
            } else {
                ui.message("导入已取消", true)
            }
        })
        this.mainObj.append(obj)

        obj = this.newButton("至无标注状态", "release", 6, "", 6).click(function () {
            let data = ui.prompt("需要管理员提权，删除后本窗口将关闭")
            if (data) {
                $.ajax({
                    url: apiRoot + '/label/full',
                    type: 'DELETE',
                    data: JSON.stringify({"admin": data}),
                }).done(resp => {
                    if (resp.code === 200 && resp.data === "exit") {
                        LabelToolSystemExit()
                    } else {
                        ui.message(resp.msg, resp.code !== 200)
                    }
                })
            }
        })
        this.mainObj.append(obj)

        obj = this.newButton("至无审阅状态", "revoke", 6, "", 6).click(function () {
            let data = ui.prompt("调整审阅状态需要管理员提权")
            if (data) {
                const senddata = JSON.stringify({"admin": data})
                console.log("senddata",senddata)
                // /api/v1/media/MEDIA/label/:op
                $.post(apiRoot + '/label/review?do=revoke', senddata).done(resp => {
                    if (resp.code === 200 && resp.data === "exit") {
                        LabelToolSystemExit()
                    } else {
                        ui.message(resp.msg, resp.code !== 200)
                    }
                })
            }
        })
        this.mainObj.append(obj)
    }
}

//endregion
/**
 * 标注数据存储器
 */
class LabelData {
    constructor() {
        this.data = {}
        this.data.crfs = new Map
        this.data.groups = new Map
        this.data.initCRF = false
        this.data.initFull = false
        this.data.urlCrf = `/api/v1/raw?action=getviewjson&view=${urlPara.crf}`
        this.data.urlFull = apiRoot + `/label/${urlPort}?do=full`

        this.downloadCrf()
    }

    set crfData(raw) {
        raw.forEach((v) => {
            switch (v.type) {
                case "group":
                    this.data.groups.set(v.group, v)
                    break

                default:
                    this.data.crfs.set(v.id, v)
            }
        })
    }

    get crfData() {
        let data = []
        this.data.crfs.forEach(value => {
            data.push(value)
        })
        return data
    }

    set raw(json) {
        console.groupCollapsed("set label raw data")
        try {
            console.log("JSON:", json)
            let d = JSON.parse(json)
            this.data.q = !!d.q ? d.q : null
            this.data.c = !!d.c ? d.c : null
            this.data.frames = !!d.frames ? d.frames : []
        } catch (e) {
            this.data.q = null
            this.data.c = null
            this.data.frames = []
        }
        console.groupEnd()
    }

    get raw() {
        let d = {}
        if (this.data.q) {
            d.q = this.data.q
        }

        d.c = this.data.frames.length
        d.frames = this.data.frames
        return JSON.stringify(d)
    }

    get crfGroups() {
        let data = []
        this.data.groups.forEach((v) => {
            data.push(v)
        })
        return data
    }

    setGlobal(group, value) {
        console.log("set global value:", group, "<-", value)
        switch (group) {
            case "q":
                this.data.q = parseInt(value)
                break

            default:

        }

        console.log("up 3")
        this.uploadFull()

    }

    setPage(page, data) {
        console.log(`set page=${page}:`, this.data.frames[page], "=>", data)
        this.data.frames[page] = data
    }

    getPage(page) {
        let data = {}
        if (this.has(page)) {
            data = this.data.frames[page]
        }
        return data
    }

    getGlobal(group) {
        return this.data[group]
    }

    crfGroup(group) {
        let data = []
        this.data.crfs.forEach((v) => {
            if (v.group === group)
                data.push(v)
        })
        return data
    }

    groupId(id) {
        return this.data.groups.get(id)
    }

    crfId(id) {
        // console.log("CRF",this.data.crfs)
        return this.data.crfs.get(id)
    }

    updateId(page, id, data) {
        let c = this.data.frames[page].clabels
        c[id] = data
        this.data.frames[page].clabels = c
    }

    has(page) {
        if (this.data.initFull && this.data.frames) {
            return (page < this.data.frames.length) ? (!!this.data.frames[page]) : false
        } else {
            return false
        }
    }

    after(page) {
        let t = this.data.frames.length
        page = (page < t) ? page : 0
        for (let i = 0; i < t; i++) {
            page = (page === (t - 1)) ? 0 : page + 1
            if (this.has(page)) return page
        }
        return false
    }

    before(page) {
        let t = this.data.frames.length
        page = (page < t) ? page : t - 1
        for (let i = 0; i < t; i++) {
            page = (page === 0) ? t - 1 : page - 1
            if (this.has(page)) return page
        }
        return false
    }

    genPostData(data) {
        let d = {}
        d.media = mediaIndex;
        d.data = data;
        d.direction = 'upload';
        return JSON.stringify(d)
    }

    uploadFull() {
        let data = this.genPostData(this.raw)
        console.warn("upload -> server")
        $.post(this.data.urlFull, data)
            .done(resp => {
                console.groupCollapsed("Upload_Full")
                console.log("upload data:", data)
                console.log("server resp:", resp)
                if (resp.code === 200) {
                    ui.message(resp.data, false)
                } else {
                    ui.message(resp.msg, true)
                }
                console.groupEnd()
            })
    }

    downloadCrf() {
        $.get(this.data.urlCrf, (result) => {
            if (result.code === 200) {
                this.crfData = result.data
                if (!this.data.initCRF) {
                    mainPanel.onCRFFinishDownload()
                    this.downloadFull()
                    this.data.initCRF = true
                }
            }
        })
    }

    downloadFull() {
        $.get(this.data.urlFull, (result) => {
            if (result.code === 200) {
                let data = result.data

                switch (urlPara.crf) {
                    case "4ap":
                    case "van":
                        data = data.replace(/DA"/g, 'DAO"')
                        data = data.replace(/DA_/g, 'DAO_')
                        break
                }
                this.raw = data
                // 统一标注数据

                ui.message("标注数据同步完成", false)
                this.data.initFull = true
                mainPanel.onGlobalFinishDownload()
            } else {
                ui.message(`标注数据同步错误：${result.msg}`, true)
            }
        })
    }
}

// region Main
/**
 * 初始化函数
 * @constructor
 */
function LabelToolSystemInit() {
    // 退出响应
    window.onbeforeunload = function () {
        mlocker.unlock()
        return "确认关闭标注界面么？"
    };
    // 右键响应
    document.body.oncontextmenu = function () {
        return false
    }
    // 媒体定时锁
    mlocker.lock()
    // 初始化界面
    mainPanel.init()
}

/**
 * 退离函数，处理交互退出相关任务
 * @constructor
 */
function LabelToolSystemExit() {
    window.onbeforeunload = null
    mainPanel.close()
    mlocker.unlock(()=>{window.close()})
}

function analysisLabelsysUrl() {
    let u1 = window.location.href.split("?")
    let u2 = u1[0].split('/')
    urlPort = u2[u2.length - 1]

    u1[1].split('&').forEach(e => {
        let p = e.split('=');
        urlPara[p[0]] = p[1];
    })
}

// 基本结构
// urlData.type .media .crf .action

class Logger {
    isShow;

    constructor() {
        this.isShow = true;
    }

    show() {
        this.isShow = true;
    }

    hide() {
        this.isShow = false;
    }

    write(msg) {
        if (this.isShow) {
            console.log(msg)
        }
    }
}
console.info("medi-sys version 2")
let urlPara = {}
let urlPort = ""
analysisLabelsysUrl()

const mlockInterval = 0//10000
const apiRoot = `/api/v1/media/${mediaIndex}`
const xmlns = "http://www.w3.org/2000/svg";

const ui = new UI()
const ldata = new LabelData()
const mainPanel = new MainPanel($("#main-content"));
const mlocker = new MediaLockerObj(mlockInterval)

LabelToolSystemInit();
// endregion
