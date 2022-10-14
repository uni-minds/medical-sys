jsGrid.locale("zh-cn");

$(function() {
    let selectorTag = new Map([["selectLabelGroupId","标注组"],["selectScreenGroupId","挑图组"],["selectView","切面"],["showReviewedOnly","仅显示挑图通过"]])
    let obj = new UltrasonicMedialist("main-content",selectorTag)
    obj.Start()
});


class UltrasonicMedialist {
    mainObj = {}
    screenClass = new Map()

    objLabelGroupSelector = {}
    objScreenGroupSelector = {}
    objViewSelector = {}

    userSelectGroupId = 0
    userSelectViewId = ""
    userIgnoreProgressCheck = false

    urlPullData = ""


    constructor(mainId, selectorTag) {
        this.screenClass = selectorTag
        this.mainObj = $('<div class="card"></div>')
        $(`#${mainId}`).append(this.mainObj);
    }

    async Start() {
        let objH1 = $(`<div class="card-header"><div class="card-title">媒体列表</div></div>`);
        let objH2 = $('<div class="p-0 row text-center"/>')
        this.mainObj.append(objH1).append($(`<div class="card-header"/>`).append(objH2))

        for (let [selector, describeCn] of this.screenClass) {
            let selectGroup = $('<div class="col-md-3"></div>')
            let optionTitle = $(`<div class="p-0 text-info">${describeCn}</div>`)
            let optionGroup = $('<div class="pb-0" style="width:100%"></div>')
            selectGroup.append(optionTitle).append(optionGroup)

            switch (selector) {
                case "selectLabelGroupId":
                    let labelGroupData = await this.getdata("/api/v1/group?action=getlistfull&type=label")

                    if (labelGroupData.code === 200) {
                        let data = labelGroupData.data

                        let objSelector = $('<select class="form-control" style="width: 100%;"/>').change(e => {
                            let gid = $(e.target).val()
                            this.UserSelect({"selectLabelGroupId": gid})
                        })
                        optionGroup.append(objSelector)
                        this.objLabelGroupSelector = objSelector
                        this.SelectorRenewData(objSelector, data)
                        objSelector.select2({theme: 'bootstrap4'})
                    }
                    break

                case "selectScreenGroupId":
                    let screenGroupData = await this.getdata("/api/v1/group?action=getlistfull&type=screen")

                    if (screenGroupData.code === 200) {
                        let data = screenGroupData.data

                        let objSelector = $('<select class="form-control" style="width: 100%;"/>').change(e => {
                            let gid = $(e.target).val()
                            this.UserSelect({"selectScreenGroupId": gid})
                        })
                        optionGroup.append(objSelector)
                        this.objScreenGroupSelector = objSelector
                        this.SelectorRenewData(objSelector, data)
                        objSelector.select2({theme: 'bootstrap4'})
                    }
                    break

                case "selectView":
                    let data = []

                    let objSelector = $('<select class="form-control" style="width: 100%;"/>').change(e => {
                        let gid = $(e.target).val()
                        this.UserSelect({"selectViewId": gid})
                    })
                    optionGroup.append(objSelector)
                    this.objViewSelector = objSelector
                    this.SelectorRenewData(objSelector, data)
                    objSelector.select2({theme: 'bootstrap4'})

                    break

                case "showReviewedOnly":
                    let viewdata = [{Gid: 1, Name: "是"}, {Gid: 0, Name: "否"}]
                    let objReviewSelector = $('<select class="form-control" style="width: 100%;"/>').change(e => {
                        let val = $(e.target).val()
                        this.UserSelect({"showReviewedOnly": val})
                    })
                    optionGroup.append(objReviewSelector)
                    this.SelectorRenewData(objReviewSelector, viewdata, true)
                    objReviewSelector.select2({theme: 'bootstrap4'})

                    break

            }
            objH2.append(selectGroup)
        }
    }

    SelectorRenewData(objSelector, data, nodefault) {
        // console.log("renew dropmenu",data)
        objSelector.children().remove()
        if (!nodefault) {
            let objOption = $(`<option value="-1">未选择</option>`)
            objSelector.append(objOption)
        }
        if (data == null || data.length === 0) {
            objSelector.attr("disabled", true)
        } else {
            objSelector.removeAttr("disabled")
            data.forEach(v => {
                let objOption = $(`<option value="${v.Gid}">${v.Name}</option>`)
                objSelector.append(objOption)
            })
        }

        return objSelector
    }

    async UserSelect(data) {
        if (data === null) {
            console.warn("user select empty data")
            return
        }

        let groupId = 0

        for (let k in data) {
            switch (k) {
                case "showReviewedOnly":
                    let val = data["showReviewedOnly"]
                    this.userIgnoreProgressCheck = (val !== "1")
                    break


                case "selectLabelGroupId":
                    groupId = data["selectLabelGroupId"]
                    if (groupId < 0) {
                        console.warn("user select label gid -1")
                        return
                    }
                    this.userSelectGroupId = groupId
                    this.userSelectViewId = ""
                    break

                case "selectScreenGroupId":
                    groupId = data["selectScreenGroupId"]
                    if (groupId < 0) {
                        console.warn("user select screen gid -1")
                        return
                    }

                    this.userSelectGroupId = groupId
                    this.userSelectViewId = ""

                    let getViewUrl = `/api/v1/media?action=getview&gid=${groupId}`
                    if (this.userIgnoreProgressCheck) {
                        getViewUrl += '&ignoreProgressCheck=1'
                    }

                    let groupView = await this.getdata(getViewUrl)
                    if (groupView.code === 200) {
                        let data = []
                        if (groupView.data !== null && groupView.data.length > 0) {
                            groupView.data.forEach(v => {
                                data.push({Gid: v, Name: v})
                            })
                        }
                        this.SelectorRenewData(this.objViewSelector, data)
                    }
                    break

                case "selectViewId":
                    this.userSelectViewId = data["selectViewId"]
                    break

            }
        }

        if (this.userSelectGroupId > 0) {
            this.urlPullData = `/api/v1/media?action=getlist&gid=${this.userSelectGroupId}`
            if (this.userSelectViewId !== "-1") {
                this.urlPullData += `&view=${this.userSelectViewId}`
            }
        }

        if (this.userIgnoreProgressCheck) {
            this.urlPullData += '&ignoreProgressCheck=1'
        }

        // console.log("update data pull url:", this.urlPullData)
        this.CreateCardBody(data.Gid, 1);
    }

    // Load favourite
    async LoadUserLastStatus() {
        let response = await this.getdata("/api/v1/user?action=laststatus")
        if (response.code === 200) {
            return JSON.parse(response.data);
        }
    }


    async CreateSelectors() {
        let objMargin = $('<div class="card-body"/>').addClass("margin row")
        this.mainObj.append(objMargin)
        for (let [selectorEn, selectorCn] of this.screenClass) {

            /*
            <div class="btn-group">
                    <button type="button" class="btn btn-info">Action</button>
                    <button type="button" class="btn btn-info dropdown-toggle dropdown-icon" data-toggle="dropdown">
                      <span class="sr-only">Toggle Dropdown</span>
                    </button>
                    <div class="dropdown-menu" role="menu">
                      <a class="dropdown-item" href="#">Action</a>
                      <a class="dropdown-item" href="#">Another action</a>
                      <a class="dropdown-item" href="#">Something else here</a>
                      <div class="dropdown-divider"></div>
                      <a class="dropdown-item" href="#">Separated link</a>
                    </div>
                  </div>
             */

            let obj1 = $('<button/>').attr("type", "button").addClass("btn btn-info").text(selectorCn)
            let obj2 = $('<button/>').attr("type", "button").addClass("btn btn-info dropdown-toggle dropdown-hover dropdown-icon").attr("data-toggle", "dropdown")
            let obj3 = $("<span class='sr-only'>DropDown</span>")
            let btnGroup = $('<div class="btn-group col-md-2"/>')
            obj2.append(obj3)
            btnGroup.append(obj1).append(obj2)

            if (selectorEn === "labelGroupId") {
                let response = await this.getdata("/api/v1/group?action=getlistfull")

                if (response.code === 200) {
                    let objMenu = $('<div class="dropdown-menu" role="menu" style="height:300px;overflow:scroll"></div>')
                    response.data.forEach(v => {
                        let objItem = $('<a class="dropdown-item" href="#">').text(v.Name)
                            .attr("Gtype", v.GType).attr("Gid", v.Gid).click(() => {
                                this.select_labelGroupId_type = v.GType
                                console.log("LOG", v.GType)
                                this.CreateCardBody(v.Gid, 1);
                            });
                        objMenu.append(objItem)
                    })
                    btnGroup.append(objMenu)
                }
            }
            objMargin.append(btnGroup)
        }
    }

    CreateCardBody(gid, lastPageIndex) {
        $("#media-table").remove()
        let obj = $("<div id='media-table' />").addClass("card-body").css("padding", 0);
        let mediaTableObj = $("<div />");
        obj.append(mediaTableObj);
        this.mainObj.append(obj);
        this.CreateMediaTable(gid, lastPageIndex, mediaTableObj)
    }

    CreateMediaTable(gid, pageIndex, tableobj) {
        let fields = [
            {
                name: "mid", type: "number", title: "ID", align: "left", width: 20,
                itemTemplate: FormatId
            },
            {
                name: "name", type: "string", title: "名称", align: "left", width: 100,
                itemTemplate: FormatName
            },
            {
                name: "duration", type: "number", title: "时长", align: "center", width: 40,
                itemTemplate: FormatDurationContent
            },
            {name: "frames", type: "number", title: "总帧", align: "center", width: 30},
            {name: "view", type: "string", title: "切面", align: "center", width: 40, itemTemplate: FormatViewContent},
            {
                name: "authors",
                type: "string",
                title: "标注人",
                align: "center",
                width: 50,
                itemTemplate: LabelAuthorRender
            },
            {
                name: "reviews",
                type: "string",
                title: "审阅人",
                align: "center",
                width: 50,
                itemTemplate: LabelReviewRender
            },
            {name: "memo", type: "string", title: "备注"},
        ];

        /**
         * @return {string}
         */
        function FormatViewContent(value) {
            // console.log("View:",value)
            if (value.startsWith('[')) {
                let v = JSON.parse(value);
                let t = "";
                v.forEach(e => {
                    t += e + "; "
                });
                t = t.substring(0, t.length - 2);
                return t

            } else {
                return value
            }
        }


        function FormatId(value) {
            let ids = value.split(".")
            if (ids.length <= 1) {
                return value
            } else if (ids.length >= 12) {

                return `D${ids[11]}`
            } else {
                console.log(value)
                return "None"
            }
        }

        function FormatName(value) {
            const dicom_us_id = "1.2.276.0.26.1.1.1.2."
            const SOP_US_IMAGE = "1.2.840.10008.5.1.4.1.1.6.1"
            const SOP_US_ENHANCE_VOLUME = "1.2.840.10008.5.1.4.1.1.6.2"
            const SOP_MULTI_FRAME = "1.2.840.10008.5.1.4.1.1.3.1"
            const SOP_SECONDARY_SCREEN = "1.2.840.10008.5.1.4.1.1.7"
            const SOP_COMPREHENSIVE_SR = "1.2.840.10008.5.1.4.1.1.88.33"

            if (value.indexOf(dicom_us_id) >= 0) {
                // dicom us
                return value.replace(dicom_us_id, "d.us.")
            } else {
                // console.log("miss",value)
                return value
            }
        }

        /**
         * @return {string}
         */
        function FormatDurationContent(value) {
            if (value === 0) {
                return "0:00.000"
            }

            let min = Math.floor(value / 60);
            let str = min + ":";
            value = value - min * 60;
            value = value.toFixed(3);
            str += value < 10 ? "0" + value : value;
            return str
        }

        /**
         * @return {null}
         */
        function LabelAuthorRender(value, dataCol) {
            let view;
            try {
                view = FormatViewContent(dataCol.view).toLowerCase()
            } catch (err) {
                view = "unknown"
            }

            let btn = $("<div>").addClass("btn btn-sm btn-block row col-sm-10 offset-sm-1").text(value.realname)
            switch (value.status) {
                case "using":
                    btn.addClass("btn-info")
                    break

                case "submit":
                    btn.addClass("btn-warning")
                    break

                case "a_reject":
                    btn.addClass("btn-danger")
                    break

                case "":
                case "free":
                default:
                    btn.addClass("btn-default").text("未标注")
                    break
            }

            btn.hover(function () {
                // 移入
                const mediaIndex = dataCol.media
                $.get(`/api/v1/media/${mediaIndex}/lock`).done(resp => {
                    if (resp.code === 200 && !!resp.msg) {
                        $(this).removeClass('btn-warning btn-info btn-default')
                            .addClass('btn-danger').text("已锁定").click(() => {
                            ui.message("其它用户正在使用本视频，请等待或选择其它数据处理。", true)
                        });
                        return
                    }

                    $.get(`/api/v1/media/${mediaIndex}/label/summary?do=author`).done(resp => {
                        let obj = $(this).removeClass('btn-warning btn-info btn-default btn-danger btn-secondary')
                        switch (resp.code) {
                            case 200:
                                //存在标注信息
                                let summary = resp.data
                                switch (summary.ReviewProgress) {
                                    case "using":
                                        obj.attr('title', summary.AuthorTips).addClass("btn-danger").click(function () {
                                            alert("审阅中，禁止修改标注")
                                        })
                                        break;
                                    default:
                                        obj.attr('title', summary.AuthorTips).addClass("btn-info").text("开始标注").click(function () {
                                            OpenLabelTool(dataCol.media, 'author', view)
                                        })
                                }
                                break
                            case 30001:
                                //不存在标注信息
                                obj.attr('title', "未标注").text("开始标注").addClass("btn-info").click(() => {
                                    OpenLabelTool(mediaIndex, 'author', view)
                                })
                                break
                            case 403:
                                //已分配他人
                                obj.attr('title', '已由他人标注').addClass("btn-secondary")
                                break
                        }
                    })
                })
            }, function () {
                // 移出
                const mediaIndex = dataCol.media

                setTimeout(() => {
                    let obj = $(this).unbind('click').removeClass('btn-default btn-info btn-warning btn-danger btn-secondary')
                    $.get(`/api/v1/media/${mediaIndex}/label/summary`).done(resp => {
                        let summary = {}
                        switch (resp.code) {
                            case 200:
                                //存在标注信息
                                summary = resp.data
                                break
                            case 403:
                                summary = resp.msg
                                break
                            case 30001:
                                //不存在标注信息
                                obj.attr('title', "未标注").addClass("btn-default").text("未标注")
                                return
                        }

                        if (!!summary) {
                            switch (summary.AuthorProgress) {
                                case "using":
                                    obj.addClass("btn-info").text(summary.AuthorRealname)
                                    break

                                case "submit":
                                    obj.addClass("btn-warning").text(summary.AuthorRealname)
                                    break

                                case "a_reject":
                                    obj.addClass("btn-danger").text(summary.AuthorRealname)
                                    break

                                case "":
                                case "free":
                                default:
                                    obj.addClass("btn-default").text("未标注");
                                    break
                            }
                        }
                    })
                }, 500)
            });
            return btn
        }

        /**
         * @return {null}
         */
        function LabelReviewRender(value, dataCol) {
            // console.log(2, value, dataCol);

            let view;
            try {
                view = FormatViewContent(dataCol.view).toLowerCase()
                // view = JSON.parse(dataCol.view)[0].toLowerCase()
            } catch (err) {
                view = "unknown"
            }

            let btn = $("<div>").addClass("btn btn-sm btn-block row col-sm-10 offset-sm-1").text(value.realname)
            switch (value.status) {
                case "using":
                    btn.addClass("btn-info")
                    break

                case "submit":
                    btn.addClass("btn-warning")
                    break

                case "r_warning":
                    btn.addClass("btn-danger")
                    break

                case "r_confirm":
                    btn.addClass("btn-primary")
                    break

                case "free":
                    btn.addClass("btn-default").text("待审核")
                    break

                default:
                    return null
            }

            btn.hover(function () {
                const mediaIndex = dataCol.media
                $.get(`/api/v1/media/${mediaIndex}/lock`).done((resp) => {
                    if (resp.code === 200 && !!resp.msg) {
                        $(this).removeClass('btn-warning btn-info btn-default')
                            .addClass('btn-danger').text("已锁定").click(() => {
                            alert("其它用户正在使用本视频，请等待或选择其它数据处理。")
                        });
                        return
                    }

                    $.get(`/api/v1/media/${mediaIndex}/label/summary?do=review`).done(resp => {
                        let summary = {}
                        switch (resp.code) {
                            case 200:
                                summary = resp.data
                                break
                            case 403:
                                summary = resp.msg
                                break
                            default:
                                return
                        }

                        if (!!summary) {
                            let obj = $(this).attr('title', summary.ReviewTips).removeClass('btn-default btn-info btn-primary btn-warning btn-danger')
                            switch (summary.AuthorProgress) {
                                case "using":
                                    obj.addClass("btn-danger").click(function () {
                                        alert("作者修改中，尚未提交审阅")
                                    })
                                    break;

                                default:
                                    obj.addClass("btn-info").text("开始审阅").click(function () {
                                        OpenLabelTool(dataCol.media, 'review', view)
                                    })
                            }
                        }

                    })
                })
            }, function () {
                const mediaIndex = dataCol.media

                setTimeout(() => {
                    $.get(`/api/v1/media/${mediaIndex}/label/summary?do=review`).done(resp => {
                        let summary = {}
                        switch (resp.code) {
                            case 200:
                                summary = resp.data
                                break
                            case 403:
                                summary = resp.msg
                                break
                            default:
                                return
                        }

                        if (!!summary) {
                            let obj = $(this).unbind('click').removeClass('btn-default btn-info btn-primary btn-warning btn-danger');
                            switch (summary.ReviewProgress) {
                                case "using":
                                    obj.addClass("btn-info").text(summary.ReviewRealname)
                                    break


                                case "submit":
                                    obj.addClass("btn-warning").text(summary.ReviewRealname)
                                    break


                                case "r_warning":
                                    obj.addClass("btn-danger").text(summary.ReviewRealname)
                                    break

                                case "r_confirm":
                                    btn.addClass("btn-primary").text(summary.ReviewRealname)
                                    break


                                default:
                                    obj.addClass("btn-default").text("待审阅")
                                    break

                            }
                        }
                    })
                }, 500)
            });
            return btn;
        }

        tableobj.jsGrid({
            height: "auto",
            width: "100%",

            fields: fields,

            sorting: true,
            paging: true,
            autoload: true,
            pageLoading: true,
            pageSize: 20,
            pageIndex: pageIndex,

            controller: {
                loadData: (e) => {
                    let d = $.Deferred();
                    let url = this.urlPullData
                    if (e.pageIndex) {
                        url += "&page=" + e.pageIndex;
                    }
                    if (e.pageSize) {
                        url += "&count=" + e.pageSize;
                    }
                    if (e.sortField) {
                        url += "&field=" + e.sortField;
                    }
                    if (e.sortOrder) {
                        url += "&order=" + e.sortOrder;
                    }
                    $.ajax({
                        url: url,
                        dataType: "json",
                        type: "GET",
                    }).done((response) => {
                        let data = {};
                        if (response.code === 200) {
                            data = response.data
                        }
                        d.resolve(data);
                    });
                    return d.promise();
                },
            },

            rowClick: function (args) {
                //console.log(args)
            },
        });
    }


    getdata(url) {
        return new Promise((resolve, reject) => {
            $.get(url, html => {
                resolve(html)
            })
        })
    }
}

function OpenLabelTool(mediaIndex, userType, viewCrf) {
    let targetURL = `/ui/labeltool/media/${mediaIndex}/${userType}?crf=${viewCrf}`;
    window.open(targetURL, "", 'fullscreen, toolbar=no, menubar=no, scrollbars=no, resizable=no,location=no, status=no')
}