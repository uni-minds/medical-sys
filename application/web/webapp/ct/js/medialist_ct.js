
// medialist
$(function() {
    let SearchOffset = 0, SearchLimit = 10;
    let URLStudiesBase="",URLStudies="", URLSeries="", URLInstance="";
    let UseDB="";

    let dcm4cheeUrl = DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/rs/studies";
    let dcm4cheeWado = DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/wado?requestType=WADO";

    let optionCSV='&accept=text/csv;delimiter=semicolon';
    let optionZIP='?accept=application/zip';

    function UpdateURLStudiesBase(Database,PatientID,StudyInstanceUID,StudyDate) {
        switch (Database) {
            case  "CCTA":
                dcm4cheeUrl=DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/rs/studies";
                dcm4cheeWado = DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/wado?requestType=WADO";
                break;
            case "CTA":
                dcm4cheeUrl=DB_CTA+"/dcm4chee-arc/aets/AS_RECEIVED/rs/studies";
                dcm4cheeWado = DB_CTA+"/dcm4chee-arc/aets/AS_RECEIVED/wado?requestType=WADO";
                break;
            default:
                console.log("Unknown database.");
        }

        URLStudiesBase = dcm4cheeUrl+ "?includefield=all";
        if (PatientID) URLStudiesBase += "&PatientID=" + PatientID;
        if (StudyInstanceUID) URLStudiesBase += "&StudyInstanceUID=" + StudyInstanceUID;
        if (StudyDate) URLStudiesBase += "&StudyDate=" + StudyDate;
    }

    /**
     * @return {number}
     */
    function NNStartPage(p) {
        console.log("NN",p);
        if (p<0) {
            return 0;
        } else if (p>20000) {
            return (p-14000);
        } else {
            for (;p>6000;p-=5000) {
                console.log(p);
            }
            return p;
        }
    }

    function UpdateURLStudies(offset,limit) {
        var FakeSearchOffset=0;
        if (UseDB==="CCTA") {
            FakeSearchOffset=NNStartPage(offset);
            SearchOffset=offset;
        } else {
            if (offset<0) {
                SearchOffset=0;
            } else {
                SearchOffset = offset;
            }
            FakeSearchOffset=SearchOffset;
        }

        SearchLimit = limit;
        let tmpURL = URLStudiesBase + "&offset=" + FakeSearchOffset + "&limit=" + SearchLimit;
        URLStudies=tmpURL;
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

    function makeTableLine(data) {
        let strTitle = "<tr>";
        data.forEach(function (str) {
            strTitle += '<td valign="middle">' + str + '</td>';
        });
        strTitle += '</tr>';
        return $(strTitle);
    }

    let StudiesTitle=["序号","患者ID","研究ID","拍摄日期","拍摄时间","影像类型","访问编号","影像数量","描述"];
    function createStudiesTable() {
        console.log("Studies URL:",URLStudies);
        if (URLStudies==="") {
            return;
        }
        var data=getJson(URLStudies);
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
                let baseURL=getDICOMValue(value, '00081190');
                let StudyUID=getDICOMValue(value, '0020000D');
                let lineObj=makeTableLine([SearchOffset + index+1,
                    getDICOMValue(value, '00100020'),
                    StudyUID,
                    formatDate(getDICOMValue(value, '00080020')),
                    formatTime(getDICOMValue(value, '00080030')),
                    getDICOMValue(value, '00080061'),
                    getDICOMValue(value, '00080050'),
                    getDICOMValue(value, '00201208'),
                    getDICOMValue(value, '00081030')]);
                lineObj.click(function () {
                    URLSeries = baseURL + '/series?includefield=all&offset=0&orderby=SeriesNumber';
                    //let data=getJson(dcm4cheeUrl + '/' + StudyUID + params);
                    createSeriesTable(StudyUID);
                });

                tb.append(lineObj);
            });
        }
    }

    let SeriesTitle=["序号","患者ID","站名","序列编号","部位","序列描述","影像数量","任务"];
    function createSeriesTable(StudyUID) {
        let data=getJson(URLSeries);
        let tb=$('#resultTableBody');
        tb.empty();
        tb.append(createTableHead(SeriesTitle));

        //创建顶部返回
        var lineObj=$('<tr><td colspan="8" "><i class="fas fa-level-up-alt"></i></td></tr>');
        lineObj.click(function(){
            createStudiesTable();
        });
        tb.append(lineObj);

        $('#btnNextPage').addClass("disabled");
        $('#btnPrevPage').addClass("disabled");

        if (data) {
            data.forEach(function (value, index) {
                let baseURL=getDICOMValue(value, '00081190');
                let SeriesUID=getDICOMValue(value, '0020000E');
                let baseInfo = "&studyUID=" + StudyUID + "&seriesUID=" + SeriesUID;
                let dicomCount=getDICOMValue(value, '00201209');
                let buttonHTML = '<span id="b_' + baseInfo + '" class="btn btn-primary seriesPreviewer">浏览序列</span>    ';
                let isDataStor=false;
                if (parseInt(dicomCount) >=100) {
                    isDataStor=true;
                }

                if (isDataStor) {
                    buttonHTML += '<span id="' + baseInfo + '" class="btn btn-warning deepAnalysis">特征分析</span>';
                }

                let lineObj=makeTableLine([index + 1,
                    getDICOMValue(value, '00100020'),
                    //StudyID,
                    getDICOMValue(value, '00081010'),
                    getDICOMValue(value, '00200011'),
                    getDICOMValue(value, '00180015'),
                    getDICOMValue(value, '0008103E'),
                    getDICOMValue(value, '00201209'),
                    buttonHTML,
                ]);
                lineObj.click(function () {
                    URLInstance = baseURL + '/instances?includefield=all&offset=0&orderby=InstanceNumber';
                    createInstancesTable(StudyUID,SeriesUID,isDataStor);
                });
                tb.append(lineObj);
            });

            $('.deepAnalysis').click(function(){
                window.event.stopPropagation();
                analysisSeries(this.id);
            });

            $('.seriesPreviewer').click(function () {
                window.event.stopPropagation();
                previewSeries(this.id);
            });

            //创建底部返回
            var lineObj=$('<tr><td colspan="8" "><i class="fas fa-level-up-alt"></i></td></tr>');
            lineObj.click(function(){
                createStudiesTable();
            });
            tb.append(lineObj);
        }
    }

    let InstancesTitle=["序号","SOP Class UID","Object UID","长","宽","位","任务"];
    function createInstancesTable(StudyUID,SeriesUID,ShowDeep) {
        console.log(URLInstance);
        let data=getJson(URLInstance);
        let tb = $('#resultTableBody');
        tb.empty();
        tb.append(createTableHead(InstancesTitle));

        //创建返回
        var lineObj=$('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
        lineObj.click(function(){
            createSeriesTable(StudyUID);
        });
        tb.append(lineObj);

        if (data) {
            data.forEach(function (value, index) {
                //let baseURL = getDICOMValue(value, '00081190');
                let ObjectUID = getDICOMValue(value, '00080018');
                let baseInfo = "&studyUID=" + StudyUID + "&seriesUID=" + SeriesUID + "&objectUID=" + ObjectUID;
                let buttonHTML = '';
                if (ShowDeep) {
                    buttonHTML = '<span id="' + baseInfo + '" class="btn btn-success deepSearch"">以图捡图</span>';
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


                let jpegURL = getWadoJpegUrl(StudyUID, SeriesUID, ObjectUID);
                lineObj.click(function () {
                    previewImage(jpegURL);
                });
                tb.append(lineObj);
            });
        }

        $(".deepSearch").click(function () {
            window.event.stopPropagation();
            console.log(this.id);
            let id=this.id;

            let StudiesUID=getQueryString(id,"studyUID");
            let SeriesUID=getQueryString(id,"seriesUID");
            let ObjectUID=getQueryString(id,"objectUID");
            deepSearch(StudiesUID,SeriesUID,ObjectUID)

        });


        //创建返回
        var lineObj=$('<tr><td colspan="8"><i class="fas fa-level-up-alt"></i></td></tr>');
        lineObj.click(function(){
            createSeriesTable(StudyUID);
        });
        tb.append(lineObj);

        $('#btnNextPage').addClass("disabled");
        $('#btnPrevPage').addClass("disabled");
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

    function getWadoJpegUrl(studyUID,seriesUID,objectUID){
        //console.log(studyUID,seriesUID,objectUID);
        let wadoAddURL = "&studyUID=" + studyUID + "&seriesUID=" + seriesUID + "&objectUID=" + objectUID;
        let jpegParams = '&contentType=image/jpeg&frameNumber=1';
        return dcm4cheeWado + wadoAddURL + jpegParams;
    }

    function formatTime(str) {
        return str.substr(0,2)+":"+str.substr(2,2)+":"+str.substr(4,2)
    }
    function formatDate(str) {
        return str.substr(0,4)+"/"+str.substr(4,2)+"/"+str.substr(6,2)
    }

    function deepSearch(StudiesUID,SeriesUID,ObjectUID) {
        $("#messagebox-content").text("正在进行搜索");
        $("#btn-showMessage").click();

        setInterval(function () {
            $("#messagebox-content").text("库内数据未经预训练。")
        },2000);

        setInterval(function () {
            $("#messagebox-content").text("准备跳转至结果");
        },4000);
        setInterval(function () {
            window.location.href = "analysis?type=deepsearch&StudiesUID="+StudiesUID+"&SeriesUID="+SeriesUID+"&ObjectUID="+ObjectUID;
        },6000);
    }

    function sleep(delay) {
        var start = (new Date()).getTime();
        while ((new Date()).getTime() - start < delay) {
            continue;
        }
    }

    function getQueryString(url,key) {
        var reg = new RegExp("(^|&)" + key + "=([^&]*)(&|$)", "i");
        var r = url.match(reg);
        if ( r != null ){
            return unescape(r[2]);
        }else{
            return null;
        }
    }

    function previewImage(url) {
        console.log("Preview image ",url);
        $('#preview-image').attr('src',url);
        $('#btn-showImage').click();
    }

    /**
     * @return {string}
     */
    function GetInstanceURL(StudiesUID,SeriesUID) {
        let base=DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/rs/studies";
        return base + "/" + StudiesUID + "/series/" + SeriesUID + "/instances?includefield=all&offset=0&orderby=InstanceNumber";
    }
    
    function previewSeries(id){
        let [data,StudiesUID,SeriesUID] = previewSeriesPrepareData(id);

        let pageCurrent=1;
        let PageEnd=Math.floor((data.length+11)/12);
        gridView(pageCurrent);

        function gridView(page) {
            if (page < 1) {
                page = 1;
            } else if (page > PageEnd) {
                page = PageEnd;
            }
            document.getElementById("instance-preview-title").textContent = "序列浏览器 [ " + page + " / " + PageEnd + " ]";
            previewSeriesGridMode();
            document.getElementById("instance-preview-next").onclick=function () {
                pageCurrent++;
                if (pageCurrent > PageEnd) {
                    pageCurrent = PageEnd
                }
                gridView(pageCurrent);
            };
            document.getElementById("instance-preview-prev").onclick=function () {
                pageCurrent--;
                if (pageCurrent<1) {
                    pageCurrent = 1
                }
                gridView(pageCurrent);
            };
            let count = (page - 1) * 12;
            for (let i = 0; i < 12; i++) {
                let imgURL = "";
                if (count + i < data.length) {
                    let value = data[count + i];
                    let ObjectUID = getDICOMValue(value, '00080018');
                    imgURL = getWadoJpegUrl(StudiesUID, SeriesUID, ObjectUID);
                }
                PreviewInstanceSetImageURL("instance-img" + i, imgURL, count+i+1,focusView);
            }
            return page
        }

        function focusView() {
            previewSeriesFocusMode();
            let obj = document.getElementById("instance-img-focus");
            obj.src = this.src;
            obj.data=this.data;
            obj.onclick = function (e) {
                var offset = $(this).offset();
                var top = e.pageY - offset.top-2;
                var left = e.pageX - offset.left-2;
                var Y=Math.round(top*512/798);
                var X=Math.round(left*512/798);
                //console.log("offset",offset.top,offset.left);
                //console.log("events",e.pageY,e.pageX);
                //console.log("relati",top,left);
                console.log(Y,X, this.data);

                gridView(pageCurrent)
            }
        }

        $("#btn-showInstancePreviewer").click();
    }

    function previewSeriesPrepareData(id) {
        let strs=id.split("&");
        let StudiesUID=strs[1].slice(9);
        let SeriesUID=strs[2].slice(10);
        let u = GetInstanceURL(StudiesUID,SeriesUID);
        let data=getJson(u);

        return [data,StudiesUID,SeriesUID];
    }

    function previewSeriesGridMode() {
        let p=document.getElementById("instance-preview-content");
        while (p.hasChildNodes()) {
            p.removeChild(p.firstChild)
        }

        p.innerHTML='<table>\n' +
            '                <tr>\n' +
            '                  <td rowspan="3" id="instance-preview-prev" width="10px">⬅️</td>\n' +
            '                  <td><img id="instance-img0" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img1" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img2" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img3" height="256px" width="256px" src="" /></td>\n' +
            '                  <td rowspan="3" id="instance-preview-next" width="10px">➡️</td>\n' +
            '                </tr>\n' +
            '                <tr>\n' +
            '                  <td><img id="instance-img4" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img5" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img6" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img7" height="256px" width="256px" src="" /></td>\n' +
            '                </tr>\n' +
            '                <tr>\n' +
            '                  <td><img id="instance-img8" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img9" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img10" height="256px" width="256px" src="" /></td>\n' +
            '                  <td><img id="instance-img11" height="256px" width="256px" src="" /></td>\n' +
            '                </tr>\n' +
            '              </table>'
    }

    function previewSeriesFocusMode() {
        let p=document.getElementById("instance-preview-content");
        while (p.hasChildNodes()) {
            p.removeChild(p.firstChild)
        }

        p.innerHTML='<div class="text-center"><img id="instance-img-focus" height="800px" width="800px" src=""></div>'

    }

    function analysisSeries(id) {
        let [data,StudiesUID,SeriesUID] = previewSeriesPrepareData(id);
        let pageCurrent=1;
        let PageEnd=Math.floor((data.length+11)/12);
        let pointsLeft=[];
        let pointsRight=[];
        let pointMode="left";

        gridView(pageCurrent);

        function gridView(page) {
            if (page < 1) {
                page = 1;
            } else if (page > PageEnd) {
                page = PageEnd;
            }

            let requestStr="";
            switch (pointMode) {
                case "left":
                    requestStr="请选择左冠脉开口";
                    break;
                case "right":
                    requestStr="请选择右冠脉开口";
                    break;
            }

            document.getElementById("instance-preview-title").textContent = requestStr +" [ " + page + " / " + PageEnd + " ]";
            previewSeriesGridMode();
            document.getElementById("instance-preview-next").onclick=function () {
                pageCurrent++;
                if (pageCurrent > PageEnd) {
                    pageCurrent = PageEnd
                }
                gridView(pageCurrent);
            };
            document.getElementById("instance-preview-prev").onclick=function () {
                pageCurrent--;
                if (pageCurrent<1) {
                    pageCurrent = 1
                }
                gridView(pageCurrent);
            };
            let count = (page - 1) * 12;
            for (let i = 0; i < 12; i++) {
                let imgURL = "";
                if (count + i < data.length) {
                    let value = data[count + i];
                    let ObjectUID = getDICOMValue(value, '00080018');
                    imgURL = getWadoJpegUrl(StudiesUID, SeriesUID, ObjectUID);
                }
                PreviewInstanceSetImageURL("instance-img" + i, imgURL, count+i+1,focusView);
            }
            return page
        }

        function focusView() {
            previewSeriesFocusMode();
            let obj = document.getElementById("instance-img-focus");
            obj.src = this.src;
            obj.data=this.data;
            obj.onclick = function (e) {
                var offset = $(this).offset();
                var top = e.pageY - offset.top-2;
                var left = e.pageX - offset.left-2;
                var Y=Math.round(top*512/798);
                var X=Math.round(left*512/798);
                //console.log("offset",offset.top,offset.left);
                //console.log("events",e.pageY,e.pageX);
                //console.log("relati",top,left);
                console.log(Y,X, this.data);
                if (pointMode=="left") {
                    pointsLeft=[Y,X,this.data];
                    pointMode="right";
                } else {
                    pointsRight=[Y,X,this.data];
                    submit();
                }

                gridView(pageCurrent)
            }
        }

        $("#btn-showInstancePreviewer").click();

        function submit() {
            console.log("L=",pointsLeft," R=",pointsRight);
            let d={};
            d["StudiesUID"]=StudiesUID;
            d["SeriesUID"]=SeriesUID;
            d["L"]=pointsLeft;
            d["R"]=pointsRight;
            d["T"]="ccta";
            let ds=JSON.stringify(d);
            console.log(ds);
            $.ajax({
                type: "POST",
                url: "analysis",
                contentType: "application/json; charset=utf-8",
                data: ds,
                dataType: "json",
                success: function (message) {
                    if (message > 0) {
                    }
                },
            });
            alert("正在进行分析，完成后页面将跳转");
            //sleep(5000);
            window.location.href = "analysis?type=ccta";
        }

    }

    function PreviewInstanceSetImageURL(id,url,data,func) {
        let obj=document.getElementById(id);
        obj.src=url;
        obj.data=data;
        obj.onclick=func;
    }
    
    createStudiesTable();

});