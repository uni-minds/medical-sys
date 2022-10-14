jsGrid.locale("zh-cn");

$(function() {
    let lastPageIndex = 1;
    $.get("/api/v1/group?action=getlist", result=> {
        if (result.code === 200 && result.data.length > 0) {
            let cc = CreateCardContainer("main-content");
            CreateCardHead(cc);
            CreateCardHeadGroupButton(result.data, cc);
            LoadUserFav(cc);
        }
    });

    // Load favourite
    function LoadUserFav(cc) {
        $.get("/api/v1/user?action=laststatus", result => {
            if (result.code === 200) {
                let data = JSON.parse(result.data);
                let gid = data.lastGroupId
                $(`#groupid_${gid}`).addClass("active");
                CreateCardBody(cc, gid, data.lastPageIndex);
            }
        })
    }

    function CreateCardContainer(mainContentId) {
        let obj = $('<div class="row col-12 card"></div>')
        $(`#${mainContentId}`).append(obj);
        return obj
    }

    function CreateCardHead(cardContainer) {
        let obj = $(`<div class="card-header d-flex p-0">
<h3 class="card-title p-3">媒体列表</h3>
<ul id="group_ids" class="nav nav-pills ml-auto p-2" /></div>`);
        cardContainer.append(obj)
    }

    function CreateCardHeadGroupButton(gids, cardContainer) {
        gids.forEach(function (gid) {
            $.get(`/api/v1/group?action=getname&gid=${gid}`, resp => {
                if (resp.code === 200) {
                    CreateGroupButtonObj(gid, resp.data, cardContainer)
                }
            });
        });

        function CreateGroupButtonObj(gid, gname, cardContainer) {
            let obj = $(`#groupid_${gid}`);
            if (obj.length > 0) {
                obj.remove()
            }

            gname = (gname === "") ? `G_${gid}` : gname;

            obj = $(`<li class="nav-item"><a id="groupid_${gid}" class="nav-link" href="#" data-toggle="tab">${gname}</a></li>`)
                .click(function() {CreateCardBody(cardContainer, gid, 1);});
            $("#group_ids").append(obj)
        }
    }

    function CreateCardBody(cardContainer, gid, lastPageIndex) {
        $("#media-table").remove()
        let obj = $("<div id='media-table' />").addClass("card-body").css("padding",0);
        let mediaTableObj = $("<div />");
        obj.append(mediaTableObj);
        cardContainer.append(obj);
        CreateMediaTable(gid, lastPageIndex, mediaTableObj)
    }

    function CreateMediaTable(gid, pageIndex, tableobj) {
        let fields = [
            {name: "mid", type: "number", title: "ID", align: "center", width: 20},
            {name: "name", type: "string", title: "名称", align: "center", width: 100},
            {
                name: "duration",
                type: "number",
                title: "时长",
                align: "center",
                width: 40,
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
            console.log("View:",value)
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

        /**
         * @return {string}
         */
        function FormatDurationContent(value) {
            if (value == 0) {
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
            console.log(1, value)

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
                let status = GetMediaLockStatus(dataCol.media);
                let summary = GetLabelSummary(dataCol.media);
                // console.log(summary)
                if (!!status) {
                    $(this).removeClass('btn-warning btn-info btn-default').addClass('btn-danger').text("已锁定").click(function () {
                        alert("其它用户正在使用本视频，请等待或选择其它数据处理。")
                    });
                } else {
                    let obj = $(this).attr('title', summary.AuthorTips).removeClass('btn-warning btn-info btn-default')
                    switch (summary.ReviewProgress) {
                        case "using":
                            obj.addClass("btn-danger").click(function () {
                                alert("审阅中，禁止修改标注")
                            })
                            break;
                        default:
                            obj.addClass("btn-info").text("开始标注").click(function () {
                                OpenLabelTool(dataCol.media, 'author', view)
                            })
                    }
                }
            }, function () {
                // 移出
                let summary = GetLabelSummary(dataCol.media);
                let obj = $(this).unbind('click').removeClass('btn-default btn-info btn-warning btn-danger')
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
            });
            return btn
        }

        /**
         * @return {null}
         */
        function LabelReviewRender(value, dataCol) {
            console.log(2, value, dataCol);

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
                let status = GetMediaLockStatus(dataCol.media);
                if (!!status) {
                    $(this).removeClass('btn-default btn-info btn-warning').addClass('btn-danger').text("已锁定").click(function () {
                        alert("其它用户正在使用本视频，请等待或选择其它数据处理。")
                    });
                } else {
                    let summary = GetLabelSummary(dataCol.media)
                    let obj = $(this).attr('title',summary.ReviewTips).removeClass('btn-default btn-info btn-primary btn-warning btn-danger')
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
            }, function () {
                let summary = GetLabelSummary(dataCol.media)
                let obj= $(this).unbind('click').removeClass('btn-default btn-info btn-primary btn-warning btn-danger');
                switch(summary.ReviewProgress) {
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
                loadData: (e)=> {
                    let d = $.Deferred();
                    let url = "/api/v1/media?action=getlist&gid=" + gid;
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
    function GetLabelSummary(hash) {
        let url = '/api/v1/label?action=summary&media='+hash;
        let data = {}
        $.ajax({
            url: url,
            type: "GET",
            async: false,
        }).done(resp=>{
            if (resp.code === 200) {
                data = resp.data
            }
        });
        return data
    }

    /**
     * @return {string}
     */
    function GetLabelRealname(hash) {
        let url = "/api/v1/label?action=getrealname&label=" + hash;
        let realname = "";
        $.ajax({
            url: url,
            dataType: "json",
            type: "GET",
            async: false,
        }).done(function (resp) {
            if (resp.code === 200) {
                realname = resp.data
            }
        });
        return realname
    }

    /**
     * @return {null}
     */
    function GetMediaLockStatus(hash) {
        let url = "/api/v1/media?action=getlock&media=" + hash;
        let data = null;
        $.ajax({
            url: url,
            type: "GET",
            async: false,
        }).done(function (resp) {
            if (resp.code === 200)
                data = resp.msg
        });
        return data
    }

    function OpenLabelTool(mediaHash, labelType, labelCrf, labelHash,readonly) {
        let targetURL = "/ui/labeltool?type=us&media=" + mediaHash + "&crf=" + labelCrf + "&action=" + labelType;
        if (labelHash && labelHash.length > 31) {
            targetURL += "&label=" + labelHash;
        }
        if (readonly) {
            targetURL += "&readonly=true"
        }
        window.open(targetURL, "", 'fullscreen, toolbar=no, menubar=no, scrollbars=no, resizable=no,location=no, status=no')
    }

});