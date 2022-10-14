btnValuesView = ["未标注", "无法识别","-", "4AP", "LVOT","RVOT","3V","3VT","AC", "AA","VC","-","股骨","头围","双顶径"]
btnValuesDiagnose = ["未标注", "无法判断", "切面正常", "切面异常"]
btnValuesInterfere = ["未标注", "无干扰",  "存在测量线"]

ImgPerPage = 8

class ScreenPlot {
    studies_id;
    series_id;
    instance_ids;
    readonly;
    data_loaded;
    data;
    maxPage
    maxInstance;
    page = 0;
    pageLink =[];

    refHead;
    refBody;
    refFoot;

    objPager;

    pic_per_line = 4;
    nowRows = 0;
    refTrs = [];
    refTdImg = [];
    refTdAct = [];
    btnsView = {};
    btnsDiagnose = {};
    btnsInterfere = {};

    constructor(parent) {
        let u = analysisURL(window.location.href)
        this.studies_id = u['studies_id']
        this.series_id = u['series_id']
        this.readonly = u['readonly']

        let btnSubmit = $(`<button class="btn btn-sm btn-block btn-primary " type="button">提交审核</button>`).click(() => {
            comm.submitProgress()
        })

        this.refHead = $('<div id="container_head" class="card-header"/>').append($(`<div class="float-right"/>`).append(btnSubmit))
        this.refBody = $('<div id="container_body" class="table text-center" />').width("100%")
        this.refFoot = $(`<div id="container_foot" class="card-footer"/>`)

        let main = $('<div class="card card-primary card-outline"/>').append(this.refHead).append(this.refBody).append(this.refFoot)
        parent.append(main)
    }

    load() {
        $.get(`/api/v1/screen?action=getdetails&src=ui&studies_id=${this.studies_id}&series_id=${this.series_id}`, result => {
            if (result.code === 200) {
                this.data = result.data
                this.maxInstance = this.data['instance_details'].length
                this.maxPage = Math.ceil(this.maxInstance / ImgPerPage)
                this.pagerInit(this.maxPage)
                this.pageSelect(0)
            }
        })
    }

    reload() {
        $.get(`/api/v1/screen?action=getdetails&src=ui&studies_id=${this.studies_id}&series_id=${this.series_id}`, result => {
            if (result.code === 200) {
                this.data = result.data
            }
        })
    }

    pagerInit(maxPage) {
        let obj = $(`<div class="clearfix float-right"/>`)
        this.objPager = $(`<ul class = "pagination m-0"/>`)

        let pgLink = $(`<a class="page-link">«</a>`).click(() => {
            this.pagePrev()
        })
        let pgItem = $(`<li class="page-item"/>`).append(pgLink)
        this.objPager.append(pgItem)

        for (let i = 1; i <= maxPage; i++) {
            pgLink = $(`<a class="page-link"/>`).text(i).click(() => {
                this.pageSelect(i-1)
            })
            pgItem = $(`<li class="page-item"/>`).append(pgLink)
            this.objPager.append(pgItem)
            this.pageLink.push(pgItem)
        }

        pgLink = $(`<a class="page-link">»</a>`).click(() => {
            this.pageNext()
        })
        pgItem = $(`<li class="page-item"/>`).append(pgLink)
        this.objPager.append(pgItem)

        obj.append(this.objPager)
        this.refFoot.append(obj)
    }

    pageNext() {
        if (this.page < this.maxPage-1) {
            this.pageSelect(this.page+1)
        }
    }

    pagePrev() {
        if (this.page > 0) {
            this.pageSelect(this.page-1)
        }
    }

    pageSelect(page) {
        console.log("page:",this.page,"->", page)
        let start = page * ImgPerPage
        let data = this.data['instance_details'].slice(start, start + ImgPerPage)

        this.pageLink[this.page].removeClass("active")
        this.pageLink[page].addClass("active")
        this.page = page

        this.figureClean()
        data.forEach((v, i) => {
            this.figurePlot(i, v['instance_id'], !!v['frames'], v)
        })
    }

    addRow(count) {
        let row = this.nowRows
        for (; row < this.nowRows + count; row++) {
            let rowObj1 = $('<div class="row" />').attr("id", `instance_${row}`)
            let rowObj2 = $('<div class="row" />').attr("id", `instance_${row}_act`).attr("style", "padding-top:5px;padding-bottom:3px")
            for (let col = 0; col < this.pic_per_line; col++) {
                let colObj1 = $('<div class="col-3"/>').attr("id", `instance_${row}_${col}`)
                let colObj2 = $('<div class="col-3"/>').attr("id", `instance_${row}_${col}`)
                rowObj1.append(colObj1);
                rowObj2.append(colObj2);
                this.refTdImg.push(colObj1);
                this.refTdAct.push(colObj2);
            }
            this.refTrs.push(rowObj1) // push img line
            this.refTrs.push(rowObj2) // puah act line
            this.refBody.append(rowObj1) // add img line
            this.refBody.append(rowObj2) // add act line
        }
        this.nowRows = row
    }

    figureClean() {
        this.nowRows = 0
        this.refTrs = []
        this.refTdImg = []
        this.refTdAct = []
        this.refBody.children().remove()
    }

    figurePlot(num, instance_id, isVideo, data) {
        // console.log("plot",num,instance_id,isVideo,data)

        let imgUrl = `/api/v1/screen?action=getmedia&instance_id=${instance_id}&src=ui`
        let col = num % this.pic_per_line
        let row = Math.floor(num / this.pic_per_line)
        if (row >= this.nowRows) {
            this.addRow(row - this.nowRows + 1)
        }
        let objImg = this.refTdImg[num]
        let objBtn = this.refTdAct[num]

        let img = $(`<img/>`).attr("src", imgUrl + "&thumb=300").width("100%")
        if (isVideo) {
            img.click(() => {
                imgViewer.show(imgUrl+"&video=true")
            })
        } else {
            img.click(() => {
                imgViewer.show(imgUrl)
            })

        }
        objImg.append(img)
        let acts = this.action(instance_id,isVideo,data)
        objBtn.append(acts)
    }

    action(instance_id,isVideo,data) {
        let obj = $('<div class="row"/>')

        let btnV = this.createBtnActionView(instance_id, data['label_view'])
        let btn1 = $('<div class="col-4" />').append(btnV)
        let btnD = this.createBtnActionDiagnose(instance_id,  data['label_diagnose'])
        let btn2 = $('<div class="col-4" />').append(btnD)
        let btnI = this.createBtnActionInterfere(instance_id,  data['label_interfere'])
        let btn3 = $('<div class="col-4" />').append(btnI)
        if (isVideo) {
            btnV.children().addClass("bg-warning")
            btnI.children().addClass("bg-warning")
            btnD.children().addClass("bg-warning")
        }

        obj.append(btn1).append(btn2).append(btn3)
        this.btnsView[instance_id] = btnV
        this.btnsDiagnose[instance_id] = btnD
        this.btnsInterfere[instance_id] = btnI
        return obj
    }

    createBtnActionView(instance_id, value) {
        if (!value || value === "0") {
            value = btnValuesView[0]
        }

        let btnDisp = $('<button type="button" class="btn btn-sm btn-flat"/>').addClass(value !== btnValuesView[0]?"btn-primary":"btn-info").text(value)
        let btnDrop = $('<button type="button" class="btn btn-sm btn-info btn-flat dropdown-toggle dropdown-icon" data-toggle="dropdown"/>')
        let btnContext = $('<div class="dropdown-menu" role="menu"/>')
        btnValuesView.forEach((txt, index) => {
            let obj = {}
            if (txt === "-") {
                obj = $('<div class="dropdown-divider"/>')
            } else {
                obj = $('<a class="dropdown-item"/>').text(txt).click(() => {
                    comm.submitValue(instance_id, "view", index)
                    if (txt !== btnValuesView[0]) {
                        btnDisp.removeClass("btn-info").addClass("btn-primary").text(txt)
                    } else {
                        btnDisp.removeClass("btn-primary").addClass("btn-info").text(txt)
                    }
                })
            }
            btnContext.append(obj)
        })

        let group = $('<div class="btn-group"/>').append(btnDisp).append(btnDrop).append(btnContext)

        let title = $('<div/>').text("切面标识")
        return $('<div/>').css("font-size", "80%").append(title).append(group)
    }

    createBtnActionDiagnose(instance_id, value) {
        if (!value || value === "0") {
            value = btnValuesDiagnose[0]
        }

        let btnDisp = $('<button type="button" class="btn btn-sm btn-flat"/>').addClass(value !== btnValuesDiagnose[0]?"btn-primary":"btn-info").text(value)
        let btnDrop = $('<button type="button" class="btn btn-sm btn-info btn-flat dropdown-toggle dropdown-icon" data-toggle="dropdown"/>')
        let btnContext = $('<div class="dropdown-menu" role="menu"/>')
        btnValuesDiagnose.forEach((txt, index) => {
            let obj = $('<a class="dropdown-item"/>').text(txt).click(() => {
                comm.submitValue(instance_id, "diagnose", index)
                if (txt !== btnValuesDiagnose[0]) {
                    btnDisp.removeClass("btn-info").addClass("btn-primary").text(txt)
                } else {
                    btnDisp.removeClass("btn-primary").addClass("btn-info").text(txt)
                }
            })
            btnContext.append(obj)
        })

        let group = $('<div class="btn-group"/>').append(btnDisp).append(btnDrop) .append(btnContext)

        let title = $('<div/>').text("诊断标识")
        return $('<div/>').css("font-size", "80%").append(title).append(group)
    }

    createBtnActionInterfere(instance_id, value) {
        if (!value || value === "0") {
            value = btnValuesInterfere[0]
        }

        let btnDisp = $('<button type="button" class="btn btn-sm btn-flat"/>').addClass(value !== btnValuesInterfere[0]?"btn-primary":"btn-info").text(value)
        let btnDrop = $('<button type="button" class="btn btn-sm btn-info btn-flat dropdown-toggle dropdown-icon" data-toggle="dropdown"/>')
        let btnContext = $('<div class="dropdown-menu" role="menu"/>')
        btnValuesInterfere.forEach((txt, index) => {
            let obj = $('<a class="dropdown-item"/>').text(txt).click((e) => {
                comm.submitValue(instance_id, "interfere", index)
                if (txt !== btnValuesInterfere[0]) {
                    btnDisp.removeClass("btn-info").addClass("btn-primary").text(txt)
                } else {
                    btnDisp.removeClass("btn-primary").addClass("btn-info").text(txt)
                }
            })
            btnContext.append(obj)
        })

        let group = $('<div class="btn-group"/>').append(btnDisp).append(btnDrop) .append(btnContext)

        let title = $('<div/>').text("干扰项")
        return $('<div/>').css("font-size", "80%").append(title).append(group)
    }

    updateBtnAction(instance_id,selector,value) {
        let obj = {};
        switch (selector) {
            case "view":
                obj = this.btnsView[instance_id]
                break
            case "diagnose":
                obj = this.btnsDiagnose[instance_id]
                break
            case "interfere":
                obj = this.btnsInterfere[instance_id]
                break
            default:
                alert("E1")
        }
        if (!!obj) {
            console.log(obj)
        }
    }
}

class ImageViewer {
    objViewer;
    objTitle;
    objImage;

    constructor() {
        let objTitleTxt = $(`<h4 class="modal-title" />`).text("影像预览")
        let objTitleBtn = $('<button type="button" class="close">').text("×").click(() => {
            this.hide()
        })
        let objTitle = $('<div class="modal-header"/>').append(objTitleTxt).append(objTitleBtn)

        let objImage = $('<img/>').width(800).attr("src", "")
        let objBodyContent = $('<div class="text-center"/>').append(objImage)
        let objBody = $('<div class="modal-body"/>').append(objBodyContent).click(() => {
            this.hide()
        })

        let obj1 = $('<div class="modal-content bg-primary"/>').append(objTitle).append(objBody);
        let obj2 = $('<div class="modal-dialog modal-lg" style="max-width: 1000px"/>').append(obj1);
        let obj = $('<div class="modal fade" />').hide().append(obj2)

        this.objViewer = obj
        this.objTitle = objTitleTxt
        this.objImage = objImage

        $('.wrapper').append(obj)
    }

    show(url, isVideo, width, height) {
        this.objImage.attr("src", url)
        this.objViewer.addClass("show").show()
    }

    hide() {
        this.objImage.attr("src", "")
        this.objViewer.removeClass("show").hide()
    }
}

class Communicator {
    studies_id;
    series_id;

    constructor() {
        let u = analysisURL(window.location.href)
        this.series_id = u['series_id']
        this.studies_id = u['studies_id']
    }

    submitValue(instance_id, selector, value) {
        console.log(this.studies_id, this.series_id, instance_id, selector, value)
        let u = `/api/v1/screen?action=author&studies_id=${this.studies_id}&series_id=${this.series_id}&instance_id=${instance_id}`
        let info = {}
        info['selector'] = selector

        switch (selector) {
            case "view":
                info['value'] = btnValuesView[value]
                break
            case "diagnose":
                info['value'] = btnValuesDiagnose[value]
                break
            case "interfere":
                info['value'] = btnValuesInterfere[value]
                break
            default:
                alert("无效选择器" + selector + value)
        }
        let data = JSON.stringify(info)
        $.post(u, data, (resp) => {
            if (resp.code === 200) {
                sp.updateBtnAction(instance_id, selector, resp.data)
                sp.reload()
            } else {
                ui.message(resp.msg, true)
            }
        })
    }

    submitProgress() {
        console.log("submit progress")
        let u = `/api/v1/screen?action=author&selector=submit&studies_id=${this.studies_id}&series_id=${this.series_id}`
        $.post(u, (resp) => {
            if (resp.code === 200) {
                alert("完成提交")
            } else {
                alert("提交失败，请重试")
            }
        })
    }
}


let comm = new Communicator()
let imgViewer = new ImageViewer()
let sp = new ScreenPlot($("#main-content"))
sidebar("hide")
sp.load()