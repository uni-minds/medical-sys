// $.jgrid.defaults.styleUI = 'Bootstrap4';
$.jgrid.defaults.iconSet = 'fontAwesome';

class tbObj {
    mainObj;
    data = []
    rowIds = []
    goods_count = 0
    out_count = 0;

    objTable;
    objPager;

    init(parentObj) {
        this.mainObj = $('<table id="screen-medialist" />');
        this.objPager = $('<div id="jqGridPager"/>')

        parentObj.append(this.mainObj).append(this.objPager)

        this.mainObj.jqGrid({
            colModel: [
                {label: "SeriesID", name: "series_id", width: 100},
                {label: "StudiesID", name: "studies_id", width: 100},
                {label: "实例数", name:"instance_count",width:50,summaryType:'sum'},
                {label: "PatientID", name: "patient_id", width: 100, hidden:true},
                {label: "进度", name: "progress", width: 30},
                {label: "标注", name: "author", width: 50},
                {label: "审阅", name: "reviewer", width: 50},
                {label: "备注", name: "memo", width: 40},
                {label: "操作", name: "oprRender", width: 60, formatter: this.oprRender}
            ],
            styleUI:'Bootstrap4',
            datatype: 'local',
            rownumbers: true,
            height: 1200,
            rowList:[10,20,30],
            rowNum:10,
            autowidth: true,
            pager: this.objPager,
            grouping: true,
            groupingView: {
                groupField: ["studies_id"],
                groupColumnShow: [false],
                groupText: ["<b style='display: inline-block;width: 130px;'>{0}</b>"],
                // groupOrder: ["asc"],
                // groupSummary: [true],
                groupCollapse: false
            },
        });
    }

    load(gid,page,count){
        $.get(`/api/v1/screen?action=getlist&src=ui&gid=${gid}&page=${page}&count=${count}`, result => {
            if (result.code === 200) {
                console.log(result.data)
                this.data = result.data.data
                this.reflush()
            }
        })
    }

    reflush() {
        this.mainObj.jqGrid("clearGridData")
        this.mainObj.jqGrid('setGridParam', {data: this.data || []});
        this.mainObj.trigger('reloadGrid');
    }

    oprRender(data, options, row) {
        return `<button type='button' class='btn btn-primary btn-xs' onClick='screen("${row.studies_id}","${row.series_id}")'>挑图</button>`
    }

    doExpand() {
        this.mainObj.jqGrid('groupingGroupBy', 'studies_id', {groupCollapse: false});
    }

    doCollapse() {
        this.mainObj.jqGrid('groupingGroupBy', 'studies_id', {groupCollapse: true});
    }

}

function screen(studiesId,seriesId) {
    OpenScreenTool(studiesId,seriesId,false)
}

function OpenScreenTool(studiedId,seriesId,readonly) {
    let u = `/api/v1/screen?action=getlock&type=us&studies_id=${studiedId}&series_id=${seriesId}`
    $.get(u, result => {
        console.log(result)
        if (result.code === 200) {
            let targetURL = `/ui/screen?action=screen&type=us&studies_id=${studiedId}&series_id=${seriesId}&readonly=${!!readonly}`;
            window.open(targetURL, "", 'fullscreen, toolbar=no, menubar=no, scrollbars=no, resizable=no,location=no, status=no')
        } else {
            alert("其它用户正在标注本视频，请等待或选择其它数据处理。")
        }
    })
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
        tbo.load(gid, lastPageIndex, 400)
    }

}

let sl = new ScreenList()
let tbo = new tbObj()

sl.Start()