// $.jgrid.defaults.styleUI = 'Bootstrap4';
$.jgrid.defaults.iconSet = 'fontAwesome';
let def_show_btn_hidden = 1

function page_resize() {
    let boxWidth = $(".card").outerWidth()
    let boxHeight = $(window).height() - $(".card-body").offset().top - $(".card-header").outerHeight()
        - $(".ui-jqgrid-pager").outerHeight() - $(".main-footer").outerHeight()
    let target = $(".my-table")
    target.jqGrid('setGridHeight', boxHeight);
    target.jqGrid('setGridWidth', boxWidth);
}

class tbObj {
    mainObj;
    data = []
    rowIds = []
    goods_count = 0
    out_count = 0;
    current_gid;
    current_page;
    current_row;
    current_count;

    objTable;
    objPager;

    init(parentObj) {
        this.mainObj = $('<table id="screen-medialist" class="my-table"/>');
        this.objPager = $('<div id="jqGridPager"/>');

        parentObj.append(this.mainObj).append(this.objPager);

        this.mainObj.jqGrid({
            colModel: [
                {label: "PatientID", name: "patient_id", width: 100},
                {label: "SeriesID", name: "series_id", width: 100},
                {label: "StudiesID", name: "studies_id", width: 100, hidden: true},
                {label: "实例数", name: "instance_count", width: 25, summaryType: 'sum'},
                {label: "检查时间", name: "studies_datetime", width: 40},
                {label: "上传时间", name: "record_datetime", width: 40},
                {label: "进度", name: "progress", width: 30},
                {label: "标注", name: "author", width: 30},
                {label: "审阅", name: "reviewer", width: 30},
                {label: "备注", name: "memo", width: 40},
                {label: "操作", name: "oprRender", width: 60, formatter: this.oprRender}
            ],
            styleUI: 'Bootstrap4',
            datatype: 'local',
            rownumbers: true,
            height: 1200,
            rowList: [10, 20, 30, 50],
            rowNum: 10,
            autowidth: true,
            pager: this.objPager,
            grouping: true,
            groupingView: {
                groupField: ["patient_id"],
                groupColumnShow: [false],
                groupText: ["<b style='display: inline-block;width: 130px;'>{0}</b>"],
                // groupOrder: ["asc"],
                // groupSummary: [true],
                groupCollapse: false
            },
        });

        page_resize()
        $(window).resize(() => {
            setTimeout("page_resize()", 100);
        });
    }

    load(gid, page, row, count) {
        console.log("load", gid, page, row, count)
        this.current_gid = gid
        this.current_page = page
        this.current_row = row
        this.reload()
    }

    reload() {
        $.get(`/api/v1/screen?action=getlist&src=ui&gid=${this.current_gid}&page=${this.current_page}&row=${this.current_row}&count=${this.current_count}`).done(result => {
            if (result.code === 200) {
                this.data = result.data.data
                this.refresh()
            }
        })
    }

    refresh() {
        this.mainObj.jqGrid("clearGridData")
        this.mainObj.jqGrid('setGridParam', {data: this.data || []});
        this.mainObj.trigger('reloadGrid');
    }

    oprRender(data, options, row) {
        // console.log(options,row)
        let  isReadonly = (row['progress'] === '审核完成')

        let obj = $('<div class="row">')
        obj.append(`<div class='btn btn-info btn-xs' style="margin-right: 5px" onClick='OpSearchHis("${row.patient_id}")'>报告</div>`)
        obj.append(`<div class='btn btn-primary btn-xs' style="margin-right: 5px" onClick='OpScreenTool("${row.studies_id}","${row.series_id}","${row.patient_id}",${isReadonly})'>挑图</div>`)
        if (row['progress'] === '待审核' ||row['progress'] === '待重审' ) {
            obj.append(`<div class='btn btn-primary btn-xs' style="margin-right: 5px" onClick='OpScreenTool("${row.studies_id}","${row.series_id}","${row.patient_id}",1)'>审核</div>`)
        }
        obj.append(`<div class='btn btn-warning btn-xs' style="margin-right: 5px" onClick='OpClean("${row.studies_id}","${row.series_id}")'>清除</div>`)

        if (def_show_btn_hidden) {
            obj.append(`<div class='btn btn-danger btn-xs' style="margin-right: 5px" onClick='OpHideStudies("${row.studies_id}")'>隐藏</div>`)
        }
        // let html = "<table class='table-borderless table table-primary table-valign-middle d-table-row'><tbody><tr class='d-table d-table-row'><td class='d-table-cell'><div type='button' class='btn btn-primary btn-xs'>挑图</div></td><td><div type='button' class='btn btn-primary btn-xs'>2</div></td></tr></tbody></table>"

        return obj.html()
    }

    doExpand() {
        this.mainObj.jqGrid('groupingGroupBy', 'studies_id', {groupCollapse: false});
    }

    doCollapse() {
        this.mainObj.jqGrid('groupingGroupBy', 'studies_id', {groupCollapse: true});
    }
}

function OpClean(studiesId, seriesId) {
    Swal.fire({
        icon: 'warning',
        title: '清除人员信息',
        text: "将清空标注人员关联，并重置标注进度！",
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
    }).then((result) => {
        if (result.isConfirmed) {
            let u = `/api/v1/screen?studies_id=${studiesId}&series_id=${seriesId}`
            $.ajax({
                url: u,
                type: "delete",
                contentType: "application/json",
                dataType: "json",
                data: "",
                success: function (resp) {
                    if (resp.code !== 200) {
                        windowError(resp.msg)
                    } else {
                        windowMessage('已清除！', '标注与人员的关联信息已清除.')
                        tbo.reload()
                    }
                },
            });
        }
    })
}

function OpScreenTool(studiesId,seriesId,patientId,review,readonly) {
    let u = `/api/v1/screen?action=getlock&type=us&studies_id=${studiesId}&series_id=${seriesId}`
    $.get(u, result => {
        // console.log(result)
        if (result.code === 200) {
            let targetURL = `/ui/screen?action=screen&type=us&studies_id=${studiesId}&series_id=${seriesId}&patient_id=${patientId}&readonly=${!!readonly}`
            targetURL += (review) ? "&review=1" : "";
            window.open(targetURL, "", 'fullscreen, toolbar=no, menubar=no, scrollbars=no, resizable=no,location=no, status=no')
        } else {
            console.log(result)
            windowError("其它用户正在标注本视频，请等待或选择其它数据处理。",2000)
        }
    })
}

function OpHideStudies(studiesId) {
    Swal.fire({
        icon: 'warning',
        title: '隐藏本实例',
        text: "本实例将被隐藏",
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: '确认',
        cancelButtonText: '取消',
    }).then((result) => {
        if (result.isConfirmed) {
            let u = `/api/v1/screen/studies/${studiesId}/hidden`
            let data = {}
            data["value"] = true
            $.post(u, JSON.stringify(data), result => {
                console.log(result)
                if (result.code === 200) {
                    windowMessage('完成', '本关联实例已隐藏.',1000)
                    tbo.reload()
                } else {
                    windowError(result.msg)
                }
            });
        }
    });
}

class ScreenList {
    lastPageIndex = 1;
    cardContainer;
    cardHead;
    cardBody;
    cardFoot;
    cardGroupButtons;

    dataGroups;

    constructor() {
        console.log("screen list 1.0")
    }

    Start() {
        $.get("/api/v1/group?action=getlistfull&grouptype=pacs_studies_id", result => {
            if (result.code === 200 && result.data.length > 0) {
                this.dataGroups = result.data
                this.cardHead = $(`<div class="card-header d-flex p-0"><h3 class="card-title p-3">媒体列表</h3><ul id="group_ids" class="nav nav-pills ml-auto p-2" /></div>`);
                this.cardContainer = $('<div class="card" />').append(this.cardHead)
                $("#main-content").append(this.cardContainer)

                this.CreateCardHeadGroupButton();
                this.LoadUserLastStatus();
            }
        });
    }

    // Load favourite
    LoadUserLastStatus() {
        $.get("/api/v1/user?action=laststatus&grouptype=screen", result => {
            if (result.code === 200) {
                let data = JSON.parse(result.data);
                let gid = data['lastGroupId']
                let lastPage = data['lastPageIndex']
                if (!(gid in this.dataGroups)) {
                    gid = this.dataGroups[0]['Gid']
                    lastPage = 1
                }

                $(`#groupid_${gid}`).addClass("active");
                this.CreateCardBody(gid, lastPage);
            }
        })
    }

    CreateCardHeadGroupButton() {
        this.dataGroups.forEach((v) => {
            let gid = v['Gid']
            let gname = v['Name']
            let obj = $(`#groupid_${gid}`);
            if (obj.length > 0) {
                obj.remove()
            }

            gname = (gname === "") ? `G_${gid}` : gname;

            obj = $(`<li class="nav-item"><a id="groupid_${gid}" class="nav-link" href="#" data-toggle="tab">${gname}</a></li>`).click(() => {
                this.CreateCardBody(gid, 1);
            });
            $("#group_ids").append(obj)
        });
    }

    CreateCardBody(gid, lastPageIndex) {
        if (this.cardBody) {
            this.cardBody.remove()
        }
        this.cardBody = $("<div class='card-body p-0'/>");
        this.cardContainer.append(this.cardBody);
        tbo.init(this.cardBody)
        tbo.load(gid, lastPageIndex, -1,-1)
    }

}

let slo = new ScreenList()
let tbo = new tbObj()

slo.Start()