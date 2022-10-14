// medialist
$(function() {
    let SearchOffset = 0, SearchLimit = 10;
    let URLStudiesBase = "", URLStudies = "", URLSeries = "", URLInstance = "";
    let UseDB = "";

    let dcm4cheeUrl = "/api/v1/database/ct/CTA/rs/studies";
    let dcm4cheeWado = "/api/v1/database/ct/CTA/wado?requestType=WADO";
    let optionCSV = '&accept=text/csv;delimiter=semicolon';
    let optionZIP = '?accept=application/zip';

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

        set studiesUid(id) {
            this.StudiesUID = id
        }

        get studiesUid() {
            return this.StudiesUID
        }

        set seriesUid(id) {
            this.SeriesUID = id
        }

        get seriesUid() {
            return this.SeriesUID
        }

        load(id) {
            let tmpStr = id.split("&");
            this.id = id
            this.studiesUid = tmpStr[1].slice(9);
            this.seriesUid = tmpStr[2].slice(10);

            toolWindows.frozen("正在加载数据……")
            $.getJSON(`${dcm4cheeUrl}/${this.studiesUid}/series/${this.seriesUid}/instances?includefield=all&offset=0&orderby=InstanceNumber`, function (resp) {
                toolWindows.unfrozen()
                if (resp.code !== 200) {
                    toolWindows.autoWarning(resp.msg)
                    return
                }
                this.dataset = JSON.parse(resp.data)
                this.pageEnd = Math.floor((this.dataset.length + 11) / 12);
                this.setup();
                $("#instance-viewer-btn-show").click();
            }.bind(this));
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
                    imgURL = `${dcm4cheeWado}&studyUID=${this.studiesUid}&seriesUID=${this.seriesUid}&objectUID=${ObjectUID}&contentType=image/jpeg&frameNumber=1`
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

        load(id) {
            this.data.pointL = true
            this.data.pointR = false
            this.data.title = "请选择左冠脉开口"
            return super.load(id);
        }

        next() {
            switch (this.data.mode) {
                case "CCTA":
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
                case "CCTA":
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
            d["StudiesUID"] = this.studiesUid;
            d["SeriesUID"] = this.seriesUid;
            d["L"] = this.data.pointdata[0];
            d["R"] = this.data.pointdata[1];
            d["T"] = this.mode;
            let ds = JSON.stringify(d);
            $.ajax({
                type: "POST",
                url: "/api/v1/analysis/ct/cta/deepbuild",
                contentType: "application/json; charset=utf-8",
                data: ds,
                dataType: "json",
                success: function (message) {
                    if (message > 0) {
                    }
                },
                error: function(XMLHttpRequest, textStatus, errorThrown) {
                    if (XMLHttpRequest.status === 404) {
                        toolWindows.autoWarning("分析服务器无法链接，请稍后再试")
                    }
                },
            });
            $("#instance-viewer-btn-close").click()
            // alert("正在进行分析，完成后页面将跳转");
            // //sleep(5000);
            // window.location.href = "analysis?type=ccta";
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

    let instanceViewer = new InstanceViewer()
    let imageViewer = new ImageViewer()
    let toolWindows = new ToolWindows()
    let deepAnalyser = new DeepAnalyser()

    function UpdateURLStudiesBase(Database,PatientID,StudyInstanceUID,StudyDate) {
        switch (Database) {
            case  "CCTA":
                dcm4cheeUrl = "/api/v1/database/ct/CCTA/rs/studies";
                dcm4cheeWado = "/api/v1/database/ct/CCTA/wado?requestType=WADO";
                break;
            case "CTA":
                dcm4cheeUrl = "/api/v1/database/ct/CTA/rs/studies";
                dcm4cheeWado = "/api/v1/database/ct/CTA/wado?requestType=WADO";
                break;
            default:
                console.log("Unknown database");
        }

        URLStudiesBase = dcm4cheeUrl+ "?includefield=all";
        if (PatientID) URLStudiesBase += "&PatientID=" + PatientID;
        if (StudyInstanceUID) URLStudiesBase += "&StudyInstanceUID=" + StudyInstanceUID;
        if (StudyDate) URLStudiesBase += "&StudyDate=" + StudyDate;
    }
    function UpdateURLStudies(offset,limit) {
        if (offset < 0) {
            SearchOffset = 0;
        } else {
            SearchOffset = offset;
        }
        offset = SearchOffset;

        SearchLimit = limit;
        let tmpURL = `${URLStudiesBase}&offset=${offset}&limit=${SearchLimit}`
        URLStudies = tmpURL;
        return tmpURL;
    }

    $('#search').click(function () {
        UseDB = $('#dataCluster').val();
        let PatientID = $('#PatientID').val();
        let StudyInstanceUID = $('#StudyInstanceUID').val();
        let StudyDate = $('#StudyDate').val();
        let StrOffset = $('#SearchOffset').val();
        let StrLimit = $('#SearchLimit').val();
        UpdateURLStudiesBase(UseDB,PatientID,StudyInstanceUID,StudyDate);

        if (StrOffset) {SearchOffset = parseInt(StrOffset) - 1;}
        if (StrLimit) {SearchLimit = parseInt(StrLimit) + 1;}
        UpdateURLStudies(SearchOffset,SearchLimit);

        createStudiesTable();
    });
    $('#btnNextPage').click(function(){
        UpdateURLStudies(SearchOffset+SearchLimit-1,SearchLimit);
        createStudiesTable();
    });
    $('#btnPrevPage').click(function () {
        UpdateURLStudies(SearchOffset-SearchLimit+1,SearchLimit);
        createStudiesTable();
    });

    function createTableHead(title) {
        let strTitle = "<tr>";
        title.forEach(function (str) {
            strTitle += '<th>' + str + '</th>';
        });
        strTitle += '</tr>';
        return $(strTitle);
    }
    function makeTableLine(data,node) {
        let strTitle = "<tr>";
        let count = 1
        data.forEach(function (str) {
            if (count === 2 && !!node ) {
                strTitle = `${strTitle}<td valign="middle">${node}</td><td valign="middle">${str}</td>`;
            } else {
                strTitle = `${strTitle}<td valign="middle">${str}</td>`;
            }
            count++
        });
        strTitle += '</tr>';
        return $(strTitle);
    }

    let StudiesTitle=["序号","所属节点","患者ID","研究ID","拍摄日期","拍摄时间","影像类型","访问编号","影像数量","描述"];
    function createStudiesTable() {
        if (URLStudies) {
            toolWindows.frozen("正在加载数据……")
            $('#resultTableBody').empty();
            $.getJSON(URLStudies, (resp) => {
                toolWindows.unfrozen()
                if (resp.code !== 200) {
                    console.log(resp)
                    toolWindows.autoWarning((resp.msg) ? `${resp.msg} (${resp.code})` : `无数据，请检查合约权限`)
                } else {
                    let data = JSON.parse(resp.data)
                    createStudiesTableRunner(data)
                }
            });
        }
    }
    function createStudiesTableRunner(data) {
        let tb=$('#resultTableBody');
        tb.empty();
        tb.append(createTableHead(StudiesTitle));
        if (data) {
            if (data.length >= SearchLimit) {
                $('#btnNextPage').removeClass("disabled");
            } else {
                $('#btnNextPage').addClass("disabled");
            }

            if (SearchOffset>1) {
                $('#btnPrevPage').removeClass("disabled");
            } else {
                $('#btnPrevPage').addClass("disabled");
            }

            data.forEach(function (value, index) {
                if (index>=SearchLimit-1) return;
                let StudyUID=getDICOMValue(value, '0020000D');
                let nodename = getNodeNameByStudyUID(StudyUID)
                let lineObj=makeTableLine([SearchOffset + index+1,
                    getDICOMValue(value, '00100020'),
                    StudyUID,
                    formatDate(getDICOMValue(value, '00080020')),
                    formatTime(getDICOMValue(value, '00080030')),
                    getDICOMValue(value, '00080061'),
                    getDICOMValue(value, '00080050'),
                    getDICOMValue(value, '00201208'),
                    getDICOMValue(value, '00081030')],nodename);
                lineObj.click(function () {
                    switch (nodename) {
                        case "BN-1":
                        case "BN-2":
                            URLSeries = `${dcm4cheeUrl}/${StudyUID}/series?includefield=all&offset=0&orderby=SeriesNumber`
                            toolWindows.frozen("正在加载数据……")
                            $.getJSON(URLSeries, (resp) => {
                                toolWindows.unfrozen()
                                if (resp.code !== 200) {
                                    toolWindows.autoWarning(resp.msg)
                                } else {
                                    createSeriesTable(JSON.parse(resp.data), StudyUID);
                                }
                            });
                            break

                        default:
                            toolWindows.autoWarning(`无本节点（${nodename}）访问权限，请检查合约许可`)

                    }
                });

                tb.append(lineObj);
            });
        }
    }

    let SeriesTitle=["序号","站名","序列编号","部位","序列描述","影像数量","任务"];
    function getNodeNameByStudyUID(studyUID) {
        let id = studyUID.split(".")[5]
        return `BN-${parseInt(id) % 3 + 1}`
    }

    function createSeriesTable(data,StudyUID) {
        let tb = $('#resultTableBody');
        tb.empty();
        tb.append(createTableHead(SeriesTitle));

        //创建顶部返回
        let lineObj = $(`<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>`);
        lineObj.click(function () {
            createStudiesTable();
        });
        tb.append(lineObj);

        $('#btnNextPage').addClass("disabled");
        $('#btnPrevPage').addClass("disabled");

        data.forEach(function (value, index) {
            let SeriesUID = getDICOMValue(value, '0020000E');
            let baseInfo = "&studyUID=" + StudyUID + "&seriesUID=" + SeriesUID;
            let dicomCount = getDICOMValue(value, '00201209');
            let buttonHTML = '<span id="b_' + baseInfo + '" class="btn btn-primary seriesPreviewer">浏览序列</span>    ';
            let isDataStor = false;

            if (parseInt(dicomCount) >= 100) {
                isDataStor = true;
            }

            if (isDataStor) {
                buttonHTML += '<span id="' + baseInfo + '" class="btn btn-warning deepAnalysis">特征分析</span>';
            }

            let lineObj = makeTableLine([index + 1,
                // getDICOMValue(value, '00100020'),
                //StudyID,
                getDICOMValue(value, '00081010'),
                getDICOMValue(value, '00200011'),
                getDICOMValue(value, '00180015'),
                getDICOMValue(value, '0008103E'),
                getDICOMValue(value, '00201209'),
                buttonHTML,
            ]);

            lineObj.click(function () {
                URLInstance = `${dcm4cheeUrl}/${StudyUID}/series/${SeriesUID}/instances?includefield=all&offset=0&orderby=InstanceNumber`
                createInstancesTable(StudyUID, SeriesUID, isDataStor);
            });
            tb.append(lineObj);
        });

        $('.deepAnalysis').unbind("click").click(function () {
            window.event.stopPropagation();
            deepAnalysis("CCTA",this.id);
        });

        $('.seriesPreviewer').unbind("click").click(function () {
            window.event.stopPropagation();
            instanceViewer.load(this.id)
        });

        //创建底部返回
        lineObj = $('<tr><td colspan="8" "><i class="fas fa-level-up-alt"></i></td></tr>');
        lineObj.click(function () {
            createStudiesTable();
        });
        tb.append(lineObj);
    }

    let InstancesTitle=["序号","SOP Class UID","Object UID","长","宽","位","任务"];
    function createInstancesTable(StudyUID,SeriesUID,ShowDeep) {
        toolWindows.frozen("正在加载数据……")
        $.getJSON(URLInstance, (resp) => {
            toolWindows.unfrozen()
            if (resp.code !== 200) {
                toolWindows.autoWarning(resp.msg)
                return
            }
            let tb = $('#resultTableBody');
            tb.empty();
            tb.append(createTableHead(InstancesTitle));

            //创建返回
            let lineObj = $('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(function () {
                createSeriesTable(StudyUID);
            });
            tb.append(lineObj);

            let data = JSON.parse(resp.data)
            if (data) {
                data.forEach(function (value, index) {
                    let ObjectUID = getDICOMValue(value, '00080018');
                    let baseInfo = `&studyUID=${StudyUID}&seriesUID=${SeriesUID}&objectUID=${ObjectUID}`;
                    let buttonHTML = '';
                    if (ShowDeep) {
                        buttonHTML = `<span id="${baseInfo}" class="btn btn-success deepSearch">以图捡图</span>`;
                    }
                    let lineObj = makeTableLine([index + 1,
                            getDICOMValue(value, '77771052'),
                            ObjectUID,
                            getDICOMValue(value, '00280010'),
                            getDICOMValue(value, '00280011'),
                            getDICOMValue(value, '00280100'),
                            buttonHTML
                        ]
                    );

                    let imageURL = `${dcm4cheeWado}&studyUID=${StudyUID}&seriesUID=${SeriesUID}&objectUID=${ObjectUID}&contentType=image/jpeg&frameNumber=1`
                    lineObj.click(function () {
                        imageViewer.show(imageURL);
                    });
                    tb.append(lineObj);
                });
            }

            $(".deepSearch").click(function () {
                window.event.stopPropagation();
                console.log(this.id);
                let id = this.id;

                let StudiesUID = getQueryString(id, "studyUID");
                let SeriesUID = getQueryString(id, "seriesUID");
                let ObjectUID = getQueryString(id, "objectUID");
                deepSearch(StudiesUID, SeriesUID, ObjectUID)

            });

            //创建返回
            lineObj = $('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(function () {
                createSeriesTable(StudyUID);
            });
            tb.append(lineObj);

            $('#btnNextPage').addClass("disabled");
            $('#btnPrevPage').addClass("disabled");
        })
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
    function getQueryString(url,key) {
        var reg = new RegExp("(^|&)" + key + "=([^&]*)(&|$)", "i");
        var r = url.match(reg);
        if (r != null) {
            return unescape(r[2]);
        } else {
            return null;
        }
    }
    function formatTime(str) {
        return str.substr(0,2)+":"+str.substr(2,2)+":"+str.substr(4,2)
    }
    function formatDate(str) {
        return str.substr(0,4)+"/"+str.substr(4,2)+"/"+str.substr(6,2)
    }

    function deepSearch(StudiesUID,SeriesUID,ObjectUID) {
        toolWindows.frozen("正在启动图检索服务")
        setTimeout(function (){
            toolWindows.unfrozen()
            toolWindows.autoWarning("图检索服务未相应，请确认权限")
        },2000)
    }
    function deepAnalysis(mode,id) {
        deepAnalyser.mode=mode
        deepAnalyser.load(id)
    }
    createStudiesTable();
});