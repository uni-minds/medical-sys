// medialist
$(function() {
    class InstanceViewer {
        constructor() {
            this.data = {}
            this.data.title = "序列浏览器"
        }

        set dataset(d) {
            this.data.dicom = d
        }

        get dataset() {
            return this.data.dicom
        }

        studyUID;
        seriesUID;
        baseUrl;

        load(db,node,studyuid,seriesuid) {
            toolWindows.frozen("正在加载数据……")
            this.studyUID=studyuid
            this.seriesUID=seriesuid
            this.baseUrl = `/api/v1/pacs/${db}/${node}/${studyuid}/${seriesuid}`
            $.get(this.baseUrl).done(resp => {
                toolWindows.unfrozen()
                if (resp.code !== 200) {
                    toolWindows.autoWarning(resp.msg)
                    return
                }
                this.dataset = resp.data
                this.pageEnd = Math.floor((this.dataset.length + 11) / 12);
                this.setup();
                $("#instance-viewer-btn-show").click();
            });
        }

        showPage(page) {
            if (page < 1) {
                this.pageCurrent = 1;
            } else if (page > this.pageEnd) {
                this.pageCurrent = this.pageEnd;
            } else {
                this.pageCurrent = page
            }
            $("#instance-viewer-title").text(`${this.data.title} ${this.pageCurrent} / ${this.pageEnd}`)
            let count = (this.pageCurrent - 1) * 12;
            for (let i = 0; i < 12; i++) {
                let imgURL = "/dist/img/boxed-bg.png"
                $(`#instance-img${i}`).attr("src",imgURL).attr("data",count + i + 1).unbind("click")
            }

            for (let i = 0; i < 12; i++) {
                let imgURL = "";
                if (count + i < this.dataset.length) {
                    let value = this.dataset[count + i];
                    let ObjectUID = getDICOMValue(value, '00080018');
                    imgURL = `${this.baseUrl}/${ObjectUID}`
                } else {
                    imgURL = "/dist/img/boxed-bg.png"
                }

                $(`#instance-img${i}`).attr("src", imgURL).attr("data", count + i + 1).click(function (e) {
                    this.imageOnClick(e)
                }.bind(this));
            }
        }

        imageOnClick(e) {
            let url = $(e.target).attr("src")
            imageViewer.show(url)
        }

        setup() {
            let p = document.getElementById("instance-viewer-content");
            while (p.hasChildNodes()) {
                p.removeChild(p.firstChild)
            }
            p.innerHTML = '<table><tr>' +
                '  <td rowspan="3" id="instance-viewer-btn-prev" width="10px">⬅️</td>' +
                '  <td><img id="instance-img0" height="256px" width="256px" src="" /></td>' +
                '  <td><img id="instance-img1" height="256px" width="256px" src="" /></td>' +
                '  <td><img id="instance-img2" height="256px" width="256px" src="" /></td>' +
                '  <td><img id="instance-img3" height="256px" width="256px" src="" /></td>' +
                '  <td rowspan="3" id="instance-viewer-btn-next" width="10px">➡️</td>' +
                '</tr><tr>' +
                '  <td><img id="instance-img4" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img5" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img6" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img7" height="256px" width="256px" src="" /></td>\n' +
                '</tr><tr>' +
                '  <td><img id="instance-img8" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img9" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img10" height="256px" width="256px" src="" /></td>\n' +
                '  <td><img id="instance-img11" height="256px" width="256px" src="" /></td>\n' +
                '</tr></table>'

            $("#instance-viewer-btn-next").unbind("click")
                .click(function () {
                    this.showPage(this.pageCurrent + 1);
                }.bind(this))

            $("#instance-viewer-btn-prev").unbind("click")
                .click(function () {
                    this.showPage(this.pageCurrent - 1);
                }.bind(this));

            this.showPage(1)
        }
    }
    class ImageViewer {
        show(url, width, height) {
            $("#image-viewer-title").text("实例预览")
            $("#image-viewer-img-content").attr("src", url).click(function () {
                this.hide()
            }.bind(this))
            $("#image-viewer-btn-show").click()
        }

        selector(url, callbackFunc, width, height) {
            $("#image-viewer-title").text("坐标选择")
            $("#image-viewer-img-content").attr("src", url)
                .unbind("click").click(function (e) {
                if (!!callbackFunc) {
                    let offset = $(e.target).offset();
                    let top = e.pageY - offset.top - 2;
                    let left = e.pageX - offset.left - 2;
                    let X = Math.round(left * 512 / 798);
                    let Y = Math.round(top * 512 / 798);
                    callbackFunc({X, Y})
                }
                this.hide()
            }.bind(this))
            $("#image-viewer-btn-show").click()
        }

        hide() {
            $("#image-viewer-btn-close").click()
        }
    }
    class DeepAnalyser extends InstanceViewer {
        constructor() {
            super();
            this.reset()
        }

        reset() {
            this.data.pointL = true
            this.data.pointR = false
            this.data.submit = false
            this.data.pointdata = []
            this.data.page = 1
            this.status = false
            this.mode = ""
        }

        get status() {
            return (!!this.data.status)
        }

        set status(s) {
            this.data.status = !!s
        }

        set mode(m) {
            this.data.mode = m
        }

        get mode() {
            return this.data.mode
        }

        load(db,node,studyuid,seriesuid) {
            this.data.pointL = true
            this.data.pointR = false
            this.data.title = "请选择左冠脉开口"
            return super.load(db,node,studyuid,seriesuid);
        }

        next() {
            switch (this.data.mode) {
                case "ccta":
                    if (this.data.pointL) {
                        this.data.pointL = false
                        this.data.pointR = true
                        this.data.title = "请选择右冠脉开口"
                        this.showPage(this.data.page)
                    } else if (this.data.pointR) {
                        this.data.pointL = false
                        this.data.pointR = false
                        this.submit(this.data.p)
                    }
            }
        }

        setPoint(index, p) {
            index = parseInt(index)
            switch (this.data.mode) {
                case "ccta":
                    if (this.data.pointL) {
                        this.data.pointdata[0] = {i: index, x: p.X, y: p.Y}
                    } else {
                        this.data.pointdata[1] = {i: index, x: p.X, y: p.Y}
                    }
            }
            this.next()
        }

        imageOnClick(e) {
            let o = $(e.target)
            let url = o.attr("src")
            let index = o.attr("data")
            imageViewer.selector(url, function (point) {
                this.setPoint(index, point)
            }.bind(this))
        }

        submit() {
            let d = {}
            d["StudiesUID"] = this.studyUID;
            d["SeriesUID"] = this.seriesUid;
            d["L"] = this.data.pointdata[0];
            d["R"] = this.data.pointdata[1];
            d["T"] = this.mode;
            toolWindows.autoWarning("正在进行分析，完成后页面将跳转")
            $.post(`/api/v1/analysis/ct/${this.mode}/deepbuild`,JSON.stringify(d)).fail(()=>{
                toolWindows.autoWarning("无法访问，请确认合约权限")
            }).done(resp=>{
                if (resp.code !== 200) {
                    toolWindows.autoWarning(resp.message)
                } else {
                    toolWindows.autoMessage("完成分析",2000)
                    setTimeout(()=>{
                        window.location.href = `analysis?type=deepbuild&mode=ccta&pipe=${resp.data}`;
                    },500)
                }
            })
            $("#instance-viewer-btn-close").click()
        }
    }
    class ToolWindows {
        constructor() {
            this.messageRef = $(`#messagebox-content`)
            this.messageBtnShow = $(`#messagebox-btn-show`)
            this.messageBtnHide = $(`#messagebox-btn-close`)

            this.warningRef = $("#warningbox-content")
            this.warningBtnShow = $("#warningbox-btn-show")
            this.warningBtnHide = $("#warningbox-btn-close")

            this.frozenRef = $("#frozenbox-content")
            this.frozenBtnShow = $("#frozenbox-btn-show")
            this.frozenBtnHide = $("#frozenbox-btn-close")
        }

        frozen(message) {
            if (!!message) {
                this.frozenRef.text(message);
                this.frozenBtnShow.click();
            } else {
                setTimeout(function () {
                    this.frozenBtnHide.click();
                }.bind(this), 500)
            }
        }

        unfrozen() {
            this.frozen(null)
        }

        showMessage(message) {
            this.messageRef.text(message);
            this.messageBtnShow.click();
        }

        hideMessage() {
            setTimeout(function () {
                this.messageBtnHide.click()
            }.bind(this), 500)
        }

        autoMessage(message, time) {
            this.messageRef.text(message);
            this.messageBtnShow.click();
            setTimeout(function () {
                this.messageBtnHide()
            }.bind(this), time)
        }

        showWarning(message) {
            this.warningRef.text(message)
            this.warningBtnShow.click()
        }

        hideWarning() {
            setTimeout(function () {
                this.warningBtnHide.click()
            }.bind(this), 500)
        }

        autoWarning(message, time) {
            if (!time) {
                time = 3000
            }
            this.warningRef.text(message)
            this.warningBtnShow.click()
            setTimeout(function () {
                this.warningBtnHide.click()
            }.bind(this), time)
        }
    }

    function getDICOMValue(data,key) {
        let c = data[key];
        if (c) {
            if (c.Value) {
                return c.Value[0];
            }
        }
        return '';
    }
    function formatTime(str) {
        return str.substr(0,2)+":"+str.substr(2,2)+":"+str.substr(4,2)
    }
    function formatDate(str) {
        return str.substr(0,4)+"/"+str.substr(4,2)+"/"+str.substr(6,2)
    }

    function deepSearch(db,node,studyuid,seriesuid,objectuid) {
        data = {}
        data.Db = db
        data.Node = node
        data.studyUID = studyuid
        data.seriesUID = seriesuid
        data.objectUID = objectuid
        toolWindows.frozen("正在启动图检索服务")
        let postUrl = '/api/v1/analysis/ct/ccta/deepsearch'
        let mode = 'ccta'
        switch (data.Db) {
            case "db_cta":
                postUrl = '/api/v1/analysis/ct/cta/deepsearch'
                mode = "cta"
                break
            case "db_ccta":
                postUrl = '/api/v1/analysis/ct/ccta/deepsearch'
                mode = "ccta"
                break
        }
        $.post(postUrl,JSON.stringify(data)).fail(()=>{
            toolWindows.autoWarning("图检索服务未相应，请确认权限")
        }).done(resp =>{
            toolWindows.unfrozen()
            // toolWindows.autoWarning("未检测到匹配项，且影像数量少于20例，请扩充特征池")
            if (resp.code !==200) {
                console.log("R",resp)
                toolWindows.autoWarning(resp.msg)
            } else {
                toolWindows.autoMessage("完成分析", 2000)
                setTimeout(() => {
                    window.location.href = `analysis?type=deepsearch&mode=${mode}&StudiesUID=${studyuid}&SeriesUID=${seriesuid}&ObjectUID=${objectuid}&pipe=${resp.data}`;
                }, 500)
            }
        })
    }
    function deepAnalysis(db,node,studyuid,seriesuid,mode) {
        deepAnalyser.mode = mode
        deepAnalyser.load(db, node, studyuid, seriesuid)
    }

    class SearchEngine {
        dblist;
        nodelist;
        btnSearch;
        useDb;
        useNode;
        urlPrefix;
        dataSearchParams;
        refTableBody;
        targetStudyUid;
        targetSeriesUid;

        TitleStudies = ["序号", "所属节点", "患者ID", "研究ID", "拍摄日期", "拍摄时间", "影像类型", "访问编号", "影像数量", "描述"];
        TitleSeries = ["序号", "站名", "序列编号", "部位", "序列描述", "影像数量", "任务"];
        TitleInstances = ["序号", "SOP Class UID", "Object UID", "长", "宽", "位", "任务"];

        constructor() {
            this.dblist = $("#dblist")
                .bind("change", () => {
                    this.selectDB(this.dblist.val())
                })

            this.nodelist = $("#nodelist")
                .bind("change", () => {
                    this.selectNode(this.nodelist.val())
                })

            this.btnSearch = $('#search')
                .click(() => {
                    this.search()
                })
            this.refTableBody = $('#resultTableBody');

            this.btnNextPage = $('#btnNextPage')
                .click(() => {
                    this.dataSearchParams.Offset = this.dataSearchParams.Offset + this.dataSearchParams.Limit
                    this.search(this.dataSearchParams);
                });
            $('#btnPrevPage')
                .click(() => {
                    let offset = this.dataSearchParams.Offset - this.dataSearchParams.Limit
                    if (offset < 0) {
                        offset = 0
                    }
                    this.dataSearchParams.Offset = offset
                    this.search(this.dataSearchParams);
                });
        }

        init() {
            $.get("/api/v1/pacs").done(resp => {
                if (resp.code === 200) {
                    resp.data.forEach((db) => {
                        let o = $("<option>").text(db)
                        this.dblist.append(o)
                    })
                    this.selectDB(resp.data[0])
                }
            })
        }

        selectDB(db) {
            this.useDb = db;
            $.get(`/api/v1/pacs/${this.useDb}`).done(resp => {
                if (resp.code === 200) {
                    // this.nodelist.empty().append($("<option>").text("All"))
                    // resp.data.forEach(node => {
                    //     let o = $("<option>").text(node)
                    //     this.nodelist.append(o)
                    // })
                    // console.log("nodelist",resp.data)
                    // this.selectNode("all")
                    this.nodelist.empty()
                    resp.data.forEach(node => {
                        let o = $("<option>").text(node)
                        this.nodelist.append(o)
                    })
                    console.log("nodelist",resp.data)
                    this.selectNode(resp.data[0])
                }
            })
        }

        selectNode(node) {
            this.useNode = node
        }

        search(data) {
            if (!data) {
                data = {}
                data.PatientID = $('#PatientID').val();
                data.StudyInstanceUID = $('#StudyInstanceUID').val();
                data.StudyDate = $('#StudyDate').val();
                data.Offset = parseInt($('#SearchOffset').val()) - 1;
                data.Limit = parseInt($('#SearchLimit').val());
                this.dataSearchParams = data
            }

            toolWindows.frozen("正在加载数据……")
            $.post(`/api/v1/pacs/${this.useDb}/${this.useNode}`, JSON.stringify(data)).done(resp => {
                toolWindows.unfrozen()
                if (resp.code !== 200) {
                    console.log(resp)
                    toolWindows.autoWarning(`交互异常，请检查合约权限`)
                } else {
                    this.createStudiesTable(resp.data)
                }
            })
        }

        createStudiesTable(data) {
            if (!data) {
                data = this.dataStudies
            } else {
                this.dataStudies = data
            }
            this.refTableBody.empty()
            this.refTableBody.append(this.createTableHead(this.TitleStudies));
            data.length >= this.dataSearchParams.Limit ? $('#btnNextPage').removeClass("disabled") : $('#btnNextPage').addClass("disabled")
            this.dataSearchParams.Offset > 1 ? $('#btnPrevPage').removeClass("disabled") : $('#btnPrevPage').addClass("disabled")

            data.forEach((value, index) => {
                if (index >= this.dataSearchParams.Limit) return;
                let StudyUID = getDICOMValue(value, '0020000D');
                let nodename = getDICOMValue(value, "nodename");
                let lineObj = this.createTableLine([this.dataSearchParams.Offset + index + 1,
                    nodename,
                    getDICOMValue(value, '00100020'),
                    StudyUID,
                    formatDate(getDICOMValue(value, '00080020')),
                    formatTime(getDICOMValue(value, '00080030')),
                    getDICOMValue(value, '00080061'),
                    getDICOMValue(value, '00080050'),
                    getDICOMValue(value, '00201208'),
                    getDICOMValue(value, '00081030')]);
                lineObj.click(() => {
                    toolWindows.frozen("正在加载数据……")
                    this.targetStudyUid = StudyUID
                    this.useNode = nodename
                    $.get(`/api/v1/pacs/${this.useDb}/${this.useNode}/${this.targetStudyUid}`).done(resp => {
                        toolWindows.unfrozen()
                        if (resp.code !== 200) {
                            console.log(resp)
                            toolWindows.autoWarning(`节点（${nodename}）交互异常，请检查合约权限`)
                        } else {
                            this.createSeriesTable(resp.data);
                        }
                    })
                });

                this.refTableBody.append(lineObj);
            });
        }

        createSeriesTable(data) {
            if (!data) {
                data = this.dataSeries
            } else {
                this.dataSeries = data
            }
            this.refTableBody.empty()
            this.refTableBody.append(this.createTableHead(this.TitleSeries));

            //创建顶部返回
            let lineObj = $(`<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>`);
            lineObj.click(() => {
                this.createStudiesTable();
            });
            this.refTableBody.append(lineObj);

            $('#btnNextPage').addClass("disabled");
            $('#btnPrevPage').addClass("disabled");

            data.forEach((value, index) => {
                let SeriesUID = getDICOMValue(value, '0020000E');
                let dicomCount = getDICOMValue(value, '00201209');
                let deepSupport = false

                let btns = []
                let btnView = $("<span/>").addClass("btn btn-primary seriesPreviewer").text("浏览序列")
                    .attr("studyuid", this.targetStudyUid)
                    .attr("seriesuid", SeriesUID).click(() => {
                        window.event.stopPropagation();
                        instanceViewer.load(this.useDb, this.useNode, this.targetStudyUid, SeriesUID)
                    })
                btns.push(btnView)

                if (parseInt(dicomCount) >= 50) {
                    let btnAnalysis = $("<span/>").addClass("btn btn-warning deepAnalysis").text("特征分析")
                        .attr("studyuid", this.targetStudyUid)
                        .attr("seriesuid", SeriesUID).click(() => {
                            window.event.stopPropagation();
                            let mode = (this.useDb == "db_cta")?"cta":"ccta"
                            deepAnalysis(this.useDb, this.useNode, this.targetStudyUid, SeriesUID, mode);
                        })
                    btns.push(btnAnalysis)
                    deepSupport = true
                }


                console.log(value)
                let lineObj = this.createTableLine([index + 1,
                    // getDICOMValue(value, '00100020'),
                    //StudyID,
                    getDICOMValue(value, '00081010'),
                    getDICOMValue(value, '00200011'),
                    getDICOMValue(value, '00180015'),
                    getDICOMValue(value, '0008103E'),
                    getDICOMValue(value, '00201209'),
                ], btns);

                lineObj.click(() => {
                    toolWindows.frozen("正在加载数据……")
                    this.targetSeriesUid = SeriesUID
                    $.get(`/api/v1/pacs/${this.useDb}/${this.useNode}/${this.targetStudyUid}/${this.targetSeriesUid}`).done(resp => {
                        toolWindows.unfrozen()
                        if (resp.code !== 200) {
                            console.log(resp)
                            toolWindows.autoWarning(`节点（${nodename}）交互异常，请检查合约权限`)
                        } else {
                            this.createInstancesTable(resp.data, deepSupport);
                        }
                    })
                    // URLInstance = `${dcm4cheeUrl}/${StudyUID}/series/${SeriesUID}/instances?includefield=all&offset=0&orderby=InstanceNumber`

                });
                this.refTableBody.append(lineObj);
            });

            //创建底部返回
            lineObj = $('<tr><td colspan="8" "><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(() => {
                this.createStudiesTable();
            });
            this.refTableBody.append(lineObj);
        }

        createInstancesTable(data, ShowDeep) {
            this.refTableBody.empty();
            this.refTableBody.append(this.createTableHead(this.TitleInstances));

            //创建返回
            let lineObj = $('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(() => {
                this.createSeriesTable();
            });
            this.refTableBody.append(lineObj);

            data.forEach((value, index) => {
                let ObjectUID = getDICOMValue(value, '00080018');
                let btns = []
                if (ShowDeep) {
                    let obj = $("<span/>").addClass("btn btn-success deepSearch").text("以图检图")
                        .click(() => {
                            window.event.stopPropagation();
                            deepSearch(this.useDb, this.useNode, this.targetStudyUid, this.targetSeriesUid, ObjectUID)
                        })
                    btns.push(obj)
                }
                let lineObj = this.createTableLine([index + 1,
                    getDICOMValue(value, '77771052'),
                    ObjectUID,
                    getDICOMValue(value, '00280010'),
                    getDICOMValue(value, '00280011'),
                    getDICOMValue(value, '00280100'),
                ], btns);

                lineObj.click(() => {
                    imageViewer.show(`/api/v1/pacs/${this.useDb}/${this.useNode}/${this.targetStudyUid}/${this.targetSeriesUid}/${ObjectUID}`);
                });
                this.refTableBody.append(lineObj);
            });

            //创建返回
            lineObj = $('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(() => {
                this.createSeriesTable();
            });
            this.refTableBody.append(lineObj);

            $('#btnNextPage').addClass("disabled");
            $('#btnPrevPage').addClass("disabled");
        }

        createTableHead(title) {
            let strTitle = "<tr>";
            title.forEach(function (str) {
                strTitle += '<th>' + str + '</th>';
            });
            strTitle += '</tr>';
            return $(strTitle);
        }

        createTableLine(data, btn) {
            let obj = $("<tr>")
            data.forEach(text => {
                // if (count === 2 && !!node) {
                //     strTitle = `${strTitle}<td valign="middle">${node}</td><td valign="middle">${str}</td>`;
                // } else {
                obj.append($(`<td valign="middle"/>`).text(text));
                // }
            });

            if (!!btn && btn.length > 0) {
                let t = $(`<td valign="middle"/>`)
                btn.forEach(b => {
                    t.append(b)
                })
                obj.append(t)
            }

            return obj
        }
    }

    let se = new SearchEngine()
    let instanceViewer = new InstanceViewer()
    let imageViewer = new ImageViewer()
    let toolWindows = new ToolWindows()
    let deepAnalyser = new DeepAnalyser()

    se.init()
});