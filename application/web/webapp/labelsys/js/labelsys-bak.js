$(function(window, document) {
    var CRFButton_Def=[
        {id:"crfSetSSMQ",ctype:'set',ctool:'time',clevel:1,ccolor:'red',cid:"SSMQ",name:"收缩末期"},
        {id:"crfSetSZMQ",ctype:'set',ctool:'time',clevel:1,ccolor:'lightblue',cid:"SZMQ",name:"舒张末期"},
        {id:"crfSetSpecTime",ctype:'set',ctool:'time',clevel:1,ccolor:'lightblue',cid:"SpecTime",name:"特殊时间"},

        {id:"crfSetXG",ctype:'set',ctool:'area',clevel:1,ccolor:'Purple',cid:"XG",name:"胸骨"},
        {id:"crfSetAD",ctype:'set',ctool:'area',clevel:1,ccolor:'Plum',cid:"JZ",name:"脊柱"},
        {id:"crfSetDA",ctype:'set',ctool:'area',clevel:1,ccolor:'red',cid:"DA",name:"降主动脉"},

        {id:"crfSetFD",ctype:'set',ctool:'area',clevel:1,ccolor:'SeaShell',cid:"FD",name:"房顶"},
        {id:"crfSetXJ",ctype:'set',ctool:'area',clevel:1,ccolor:'Tomato',cid:"XJ",name:"心尖"},
        {id:"crfSetSZ",ctype:'set',ctool:'area',clevel:1,ccolor:'Wheat',cid:"SZ",name:"十字交叉"},

        {id:"crfSetXJWM",ctype:'set',ctool:'area',clevel:1,ccolor:'PaleTurquoise',cid:"XJWM",name:"心肌外膜"},
        {id:"crfSetXJNM",ctype:'set',ctool:'area',clevel:1,ccolor:'Olive',cid:"XJNM",name:"心肌内膜"},
        {id:"crfSetFJM",ctype:'set',ctool:'area',clevel:1,ccolor:'MintCream',cid:"FJM",name:"肺静脉"},

        {id:"crfSetLA",ctype:'set',ctool:'area',clevel:1,ccolor:'LightBlue',cid:"LA",name:"左房 LA"},
        {id:"crfSetLV",ctype:'set',ctool:'area',clevel:1,ccolor:'LightCoral',cid:"LV",name:"左室 LV"},
        {id:"crfSetMV",ctype:'set',ctool:'multi',clevel:1,ccolor:'LightCyan',cid:"MV",name:"二尖瓣 MV"},

        {id:"crfSetRA",ctype:'set',ctool:'area',clevel:1,ccolor:'SlateBlue',cid:"RA",name:"右房 RA"},
        {id:"crfSetRV",ctype:'set',ctool:'area',clevel:1,ccolor:'Lavender',cid:"RV",name:"右室 RV"},
        {id:"crfSetTV",ctype:'set',ctool:'multi',clevel:1,ccolor:'LightSlateGray',cid:"TV",name:"三尖瓣 TV"},

        {id:"crfSetRYKBM",ctype:'set',ctool:'area',clevel:1,ccolor:'DarkCyan',cid:"RYKBM",name:"卵圆孔瓣膜"},
        {id:"crfSetRYKKK",ctype:'set',ctool:'multi',clevel:1,ccolor:'Lime',cid:"RYKKK",name:"卵圆孔开口"},
        {id:"crfSetTSC",ctype:'set',ctool:'multi',clevel:1,ccolor:'DodgerBlue',cid:"TSC",name:"TSC"},


        {id:"crfSetVideoQ1",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"VQ1",name:"视频清晰"},
        {id:"crfSetVideoQ2",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"VQ2",name:"质量一般"},
        {id:"crfSetVideoQ3",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"VQ3",name:"视频模糊"},
        {id:"crfSetFrameQ1",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"FQ1",name:"本帧清晰"},
        {id:"crfSetFrameQ2",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"FQ2",name:"本帧一般"},
        {id:"crfSetFrameQ3",ctype:'set',ctool:'quality',clevel:1,ccolor:'yellow',cid:"FQ3",name:"本帧模糊"},
    ];

    // 系统类
    var Debug = true;
    var xmlns = "http://www.w3.org/2000/svg";
    // video系统
    var mediaVar= {
        duration: null, coffsetxy: null, ch: null, cw: null,
        containerRef:null, playerRef: null,
        progressBackRef: null, progressFrontRef: null, progressStrRef: null,
        ctrlBtnPlay: null, ctrlBtnStop: null, ctrlBtnPrev: null, ctrlBtnNext: null,
    };
    // draw 系统
    var ButtonsDef_Tools=[
        {id:"drawSelectBtn",cid:"select",name:"指针",func:null},
        {id:"drawPointBtn",cid:"point",name:"点",func:null},
        {id:"drawLineBtn",cid:"line",name:"直线",func:null},
        {id:"drawPolylineBtn",cid:"polyline",name:"折线",func:null},
        {id:"drawPolygonBtn",cid:"polygon",name:"多边形",func:null},
        {id:"drawTextBtn",cid:"text",name:"文字",func:null},
        {id:"drawDelete",cid:"remove",name:"删除",func:null},
        {id:"drawRemoveFrame",cid:"",name:"清空本帧",func:function(){LabelsPanel.clean();labelsData.RemoveFrame(timeFrame.getFrame(mediaVar.playerRef.currentTime));}},
        {id:"drawRemoveAll",cid:"",name:"清空全部",func:function(){LabelsPanel.clean();labelsData.Empty();}},
        {id:"drawSaveImage",cid:"savef",name:"保存图片",func:null},
        {id:"drawLoadJson",cid:"",name:"打开标注",func:function(){labelsData.LoadJsonfile();}},
        {id:"drawSaveJson",cid:"",name:"导出标注",func:function(){labelsData.SaveJsonfile();}}
    ];
    var drawVar={
        penColor:null,selectedBkColor:"#0080FF",unselectBkColor:"#E0E0E0",
    };
    var globalVar={
        disableSvgBgOnClickFunc:false,
        targetObj:null,
        targetElement:null,
        onDraw:false,
        svgRef:null,
    };

    var drawToolMode = "", drawPenDown = false, drawTargetID, svgPolyPathArray, drawFrameLabelModifyFlag = false;
    var targetElement;
    // label 系统
    var labelsData = new LabelsDataObj();
    // 存储CRF标注
    var LabelsObjData = new Array();
    var labelsShownType=preValue.LabelsShownType;

    // 初始化各类
    var timeFrame = new TimeFrameObj();
    var timer = new TimerObj();

    var Panels={
        init:function () {
            this.initToolsPanel();
            this.initWorkPanel();
            this.initInfoPanel();
            this.initCRFPanel();
        },
        initToolsPanel:function () {
            this.initToolButtons();
            this.initColorButtons();
        },
        initCRFPanel:function () {
            CRFPanel.init();
        },
        initWorkPanel:function () {
            switch(preValue.SourceType){
                case "video":
                    MediaPanel.init();
                    window.onresize = MediaPanel.adapt;
                    break;
                case "image":
                    ImagePanel.init();
                    break;
                default:
                    alert('数据源类型 "',preValue.SourceType,'" 无法识别。');
            }
            LabelsPanel.init();
        },
        initInfoPanel:function () {
        },
        initToolButtons: function () {
            ButtonsDef_Tools.forEach(ele => {
                ref = document.createElement("button");
                ref.id=authorData.id;
                ref.className= "toolBtn";
                $("#tools").append(ref);
                ref.cid = authorData.cid;
                ref.innerText = authorData.name;
                ref.onclick = function () {
                    toolSelect(authorData.cid);
                    if (authorData.func) {authorData.func();}
                }
            })
        },
        initColorButtons: function () {
            drawVar.penColor=preValue.PenColor;
            ButtonsDef_Colors.forEach(ele => {
                ref = document.createElement("button");
                ref.id=authorData.id;
                ref.cid=authorData.cid;
                ref.className="colorBtn";
                ref.style.setProperty("background-color",authorData.color,);
                ref.style.setProperty('border','1px none');
                $("#colors").append(ref);
                ref.onclick = function () {
                    drawVar.penColor = authorData.color;
                    Array.from($(".colorBtn")).forEach(v => {
                        if (v.cid == authorData.color) {
                            v.style.setProperty('border','5px black double')
                        } else {
                            v.style.setProperty('border','1px none')
                        }
                    });
                }
            })
        },
    };

    var MediaPanel={
        init:function () {
            var mediaContainer = document.createElement("div");
            mediaContainer.id = "mediaContainer";
            mediaContainer.className = "videoBackground";
            $("#WorkPanel").append(mediaContainer);

            var videoPlayer = document.createElement("video");
            videoPlayer.id = "videoContainer";
            videoPlayer.className = "workSpaceOverlayer";
            videoPlayer.setAttribute("style", "z-index: 1");
            videoPlayer.setAttribute("preload", "auto");
            videoPlayer.setAttribute("src", preValue.Source);
            videoPlayer.setAttribute("type", "video/mp4");
            videoPlayer.innerText = "您所使用的浏览器不支持 HTML5 视频播放，请换用Chrome或Firefox浏览器（国产浏览器请切换至“急速模式”。";
            mediaContainer.append(videoPlayer);

            globalVar.svgPlayer = document.createElementNS(xmlns, "svg");
            globalVar.svgRef.id = "svgContainer";
            globalVar.svgRef.className = "workSpaceOverlayer";
            globalVar.svgRef.setAttribute("class", 'workSpaceOverlayer');
            globalVar.svgRef.setAttribute("style", "z-index:2");
            mediaContainer.append(globalVar.svgRef);

            var vpBK = document.createElement("div");
            vpBK.setAttribute("id", "videoProgressBackground");
            vpBK.className = "videoProgressBackground";

            var vpF = document.createElement("div");
            vpF.setAttribute("id", "videoProgressFrontend");
            vpF.className = "videoProgressFrontend";

            var vpStr = document.createElement("span");
            vpStr.setAttribute("id", "videoProgressStr");
            vpStr.innerText = "0%";

            vpF.appendChild(vpStr);
            vpBK.appendChild(vpF);
            $("#WorkPanel").append(vpBK);

            var videoControlCenter = document.createElement("div");
            videoControlCenter.setAttribute("id", "videoControlCenter");
            videoControlCenter.className = "videoControlCenter";
            var videoPlaytime = document.createElement("span");
            videoPlaytime.id = "videoPlaytime";
            videoPlaytime.innerText = "0:00.000 / 0:00.000";
            videoControlCenter.appendChild(videoPlaytime);
            var videoLC = document.createElement("div");
            videoLC.id = "videoLC";
            this.createButtons(ButtonsDef_VideoControls, videoLC);
            videoControlCenter.appendChild(videoLC);
            var videoRC = document.createElement("div");
            videoRC.id = "videoRC";
            // 标注跟随模式
            this.createRadios(RadioDef_LabelShowType, videoRC, "labelShowntype");
            // 视频循环使能
            target = document.createElement("input");
            target.id = "videoLoopEnable";
            target.setAttribute('checked', true);
            target.setAttribute("type", "checkbox");
            target.setAttribute("name", "videoLoopEnable");
            target.setAttribute("value", "loop");
            targetL = document.createElement("label");
            targetL.appendChild(target);
            targetL.append("  循环  ");
            videoRC.appendChild(targetL);
            videoControlCenter.appendChild(videoRC);
            $("#WorkPanel").append(videoControlCenter);

            mediaVar.progressBackRef = vpBK;
            mediaVar.progressFrontRef = vpF;
            mediaVar.progressStrRef = vpStr;
            mediaVar.playerRef = videoPlayer;
            mediaVar.containerRef = mediaContainer;
            mediaVar.ctrlBtnPlay = document.getElementById('videoPlayBtn');
            mediaVar.ctrlBtnStop = document.getElementById('videoStopBtn');
            mediaVar.ctrlBtnNext = document.getElementById('videoNextFrameBtn');
            mediaVar.ctrlBtnPrev = document.getElementById('videoPrevFrameBtn');

            mediaVar.ctrlBtnPlay.onclick=this.videoPlay;
            mediaVar.ctrlBtnStop.onclick=this.videoStop;
            mediaVar.ctrlBtnPrev.onclick=this.videoPrevFrame;
            mediaVar.ctrlBtnNext.onclick=this.videoNextFrame;

            videoPlayer.onloadedmetadata = this.adapt;
        },
        adapt:function() {
            v = mediaVar.playerRef;
            c = mediaVar.containerRef;
            var vratio = v.videoWidth / v.videoHeight;
            var cratio = c.offsetWidth / c.offsetHeight;
            if (vratio > cratio) {
                v.width = c.offsetWidth;
            } else {
                v.height = c.offsetHeight;
            }
            mediaVar.coffsetxy= (Offset.fromDocument(v));
            mediaVar.cw=v.offsetWidth;
            mediaVar.ch=v.offsetHeight;
            //v.style.setProperty("margin-left",(c.offsetWidth-v.clientWidth)/2+"px");

            mySVG=document.getElementById('svgContainer');
            mySVG.setAttribute('width','100%');
            mySVG.setAttribute('height','100%');
            mySVG.style.setProperty('height',mediaVar.ch+'px');
            mySVG.style.setProperty('width',mediaVar.cw+'px');
            //mySVG.style.setProperty("margin-left",(c.offsetWidth-v.clientWidth)/2+"px");
            // 禁用右键
            mySVG.oncontextmenu = function (e) {
                e.preventDefault();
            };
            // 处理鼠标移动
            mySVG.onmousemove = function (e) {
                doMouse.move(e);
            };
            // 处理鼠标按键
            mySVG.onmousedown = function (e) {
                switch (e.button) {
                    case 0:
                        doMouse.downLeft(e);
                        break;
                    case 1:
                        doMouse.downMiddle(e);
                        break;
                    case 2:
                        doMouse.downRight(e);
                        break;
                }
            };
            mediaVar.duration=v.duration;
            timeFrame.init(preValue.FPS, v.duration);
        },
        createButtons:function (buttonCollections,parent) {
            buttonCollections.forEach(ele=>{
                target = document.createElement("button");
                target.setAttribute("id",authorData.id);
                target.className=authorData.class;
                target.innerText=authorData.text;
                parent.appendChild(target);
            })
        },
        createRadios:function (collections,parent,strName) {
            collections.forEach(ele=>{
                targetRadio=document.createElement("input");
                targetRadio.setAttribute("type","radio");
                targetRadio.id=authorData.id;
                targetRadio.setAttribute("name",strName);
                targetRadio.setAttribute("value",authorData.value);
                if (authorData.checked) {
                    targetRadio.setAttribute("checked","true");
                }
                targetL=document.createElement("label");
                targetL.appendChild(targetRadio);
                targetL.append(authorData.text);
                parent.appendChild(targetL);
            })

        },
        videoPlay:function () {
            v = mediaVar.playerRef;
            if (v.paused || v.ended) {
                if (v.ended) {
                    v.currentTime = 0;
                }
                v.play();
                timer.set();
            } else {
                v.pause();
                timer.clear();
            }
            MediaPanel.updateFast();
        },
        videoStop:function () {
            v = mediaVar.playerRef;
            v.pause();
            v.currentTime = 0;
            timer.clear();
            MediaPanel.updateFast();
        },
        videoPrevFrame:function(){
            v = mediaVar.playerRef;
            v.pause();
            var time = (v.currentTime - 1 / preValue.FPS).toFixed(6);
            if (Debug) {
                console.log("PrevTime: ", v.currentTime, "; PrevFrame:", timeFrame.getFrame(v.currentTime));
                console.log("Set Time: ", time, "; Set Frame:", timeFrame.getFrame(time));
            }
            if (time >= 0) {
                v.currentTime = time;
            } else {
                v.currentTime = 0;
            }
            MediaPanel.updateFast();
        },
        videoNextFrame:function(){
            v = mediaVar.playerRef;
            v.pause();
            var time = (v.currentTime + (1 / preValue.FPS)).toFixed(6);
            if (Debug) {
                console.log("PrevTime: ", v.currentTime, "; PrevFrame:", timeFrame.getFrame(v.currentTime));
                console.log("Set Time: ", time, "; Set Frame:", timeFrame.getFrame(time));
            }
            if (time >= v.duration) {
                v.currentTime = v.duration;
            } else {
                v.currentTime = time;
            }
            MediaPanel.updateFast();
        },
        updateFast:function () {
            var vP = mediaVar.playerRef;
            var percent = vP.currentTime / vP.duration;
            mediaVar.progressStrRef.innerText = ((percent * 100).toFixed(0) + "%");//进度条文字进度
            $('#videoPlaytime').html(convert.formatTime(vP.currentTime) + " / " + convert.formatTime(mediaVar.duration));//控制条时间

            mediaVar.progressFrontRef.style.width = percent * (mediaVar.progressBackRef.offsetWidth) + "px";//调整控制条长度

            document.getElementById("videoPlayBtn").innerText= (vP.paused || vP.ended) ? "Play" : "Pause";// 调整按键文字

            switch (labelsShownType) {
                case "0":// 禁用
                    LabelsPanel.clean();
                    break;
                case "1":// 跟随
                    LabelsPanel.clean();
                    LabelsPanel.ShowFrame(timeFrame.getFrame(vP.currentTime));
                    break;
                case "2":// 固定
                    break;
                default:
                    break;
            }
            // 标注修改切换时提交
            if (drawFrameLabelModifyFlag) {
                LabelsPanel.SaveFrame();
            }
            drawFrameLabelModifyFlag = false;
        }
    };
    var LabelsPanel = {
        init:function(){
            if (preValue.SourceType == "video"){
                this.initVideoLabelTool();
            } else {
                this.initImageLabelTool();
            }
        },
        initImageLabelTool:function() {

        },
        initVideoLabelTool: function () {
            var lst_follow = document.getElementById('labelsShownType_Follow');
            var lst_fixed=document.getElementById('labelsShownType_Fixed');
            var lst_none=document.getElementById('labelsShownType_None');

            if (lst_follow.checked) {
                labelsShownType = lst_follow.value;
            }
            if (lst_fixed.checked) {
                labelsShownType = lst_fixed.value;
            }
            if (lst_none.checked) {
                labelsShownType = lst_none.value;
            }
            // 按键响应绑定
            lst_follow.onclick = function () {
                if(Debug){console.log("LabelsShownType",this.value)}
                labelsShownType = this.value;
            };
            lst_fixed.onclick = function () {
                if(Debug){console.log("LabelsShownType",this.value)}
                labelsShownType = this.value;
            };
            lst_none.onclick = function () {
                if(Debug){console.log("LabelsShownType",this.value)}
                labelsShownType = this.value;
            };
        },
        clean: function () {
            var childs = mySVG.childNodes;
            for (var i = childs.length-1; i > -1; i--) {
                mySVG.removeChild(childs[i]);
            }
        },
        ShowFrame: function (iframe) {
            labelsData.drawFrame(iframe);
        },
        SaveFrame: function (iframe) {
            labelsData.SaveFrame(iframe);
        }
    };
    var CRFPanel = {
        init: function () {
            var b=null;
            b=this.addBlock('blkTime');
            this.addText('>> 时间点标注',b);
            this.addButtons('time',b);

            b=this.addBlock('blkArea');
            this.addText('>> 单结构标注',b);
            this.addButtons('area',b);

            b=this.addBlock('blkMulti');
            this.addText('>> 多结构标注',b);
            this.addButtons('multi',b);

            b=this.addBlock('blkQ');
            this.addText('>> 质量标注',b);
            this.addButtons('quality',b);
        },
        addButtons: function (ctool,ref) {
            CRFButton_Def.forEach(ele => {
                if (authorData.ctool != ctool) return;
                var obj = document.createElement("button");
                obj.id=authorData.id;
                obj.style.setProperty("height","50px");
                obj.style.setProperty("width","80px");
                obj.style.setProperty("font-size","14px");
                obj.style.setProperty("border","1px solid gray");
                obj.style.setProperty("background-color",'white');//authorData.ccolor);
                obj.innerText = authorData.name;

                obj.cid = authorData.cid;
                obj.ctype=authorData.ctype;
                obj.ctool=authorData.ctool;
                obj.ccolor=authorData.ccolor;
                obj.clevel=authorData.clevel;

                ref.appendChild(obj);
                obj.onclick = CRFBtnOnClick;
            });
        },
        addText:function (string,ref) {
            var p=document.createElement('p');
            p.innerHTML=string;
            ref.appendChild(p)
        },
        addBlock:function (strDiv) {
            var crfPanelRef = document.getElementById('RPanel');
            var d=document.createElement('div');
            d.id=strDiv;
            d.className="labelBlk";
            crfPanelRef.appendChild(d);
            return d;
        },
        reflushAllBtn:function () {
            this.reflushTimeBtn();
            this.reflushAreaBtn();
            this.reflushMultiBtn();
            this.reflushQBtn()
        },
        reflushTimeBtn:function() {
            ref=document.getElementById('blkTime');
            console.log(ref.childNodes);
        },
        reflushAreaBtn:function() {
            ref=document.getElementById('blkArea');
            console.log(ref.childNodes);

        },
        reflushMultiBtn:function () {
            ref=document.getElementById('blkMulti');
            console.log(ref.childNodes);

        },
        reflushQBtn:function(){
            ref=document.getElementById('blkQ');
            console.log(ref.childNodes);

        }
    };
    
    var mainObject = {
        init: function () {
            var that = this;
            //videoVar.playerRef.removeAttribute("controls");
            this.bindFunctions();
            this.videoOperateControls();
        },
        videoOperateControls: function () {
            bindEvent(mediaVar.progressBackRef, "mousedown", videoObj.onSeekbar);
        },
        bindFunctions: function () {
            bindEvent(mediaVar.playerRef, "ended", videoObj.onEnded);
            document.onkeydown = function (evt) {
                var theEvent = window.event || evt;
                var code = theEvent.keyCode || theEvent.which;
                if (Debug) {
                    console.log(code.toString());
                }
                parseKeys(code);
                return;
            };
        }
    };
    var videoObj = {
        // 进度条鼠标点击动作
        onSeekbar: function (ele) {
            var pN=mediaVar.progressFrontRef;
            var pA=mediaVar.progressBackRef;
            var vP=mediaVar.playerRef;
            timer.clear();
            var length = authorData.pageX - Offset.fromDocument(pA).left;
            var percent = length / pA.offsetWidth;
            vP.currentTime = percent * vP.duration;
            MediaPanel.updateFast();
            if (!vP.ended) {
                timer.set();
            }
        },
        // 视频结束动作
        onEnded: function () {
            if (document.getElementById('videoLoopEnable').checked) {
                vc=document.getElementById("videoContainer");
                vc.currentTime = 0;
                vc.play();
            }
            MediaPanel.updateFast();
        },
    };

    var svgPoint = {
        fun: function (x, y) {
            drawTargetID = svgElement.GenID();
            this.create(drawTargetID, x, y, preValue.PointR, preValue.PenWidth, drawVar.penColor);
            svgElement.addTargetEventListeners(drawTargetID);
            this.submit(drawTargetID);
        },
        drawNow: function (strID, x, y, color) {
            this.create(strID, x, y, preValue.PointR, preValue.PenWidth, color);
            svgElement.addTargetEventListeners(strID);
        },
        create: function (strID, x, y, r, width, color) {
            let target = document.createElementNS(xmlns, "circle");
            target.setAttribute("id", strID);
            target.setAttribute("cx", x);
            target.setAttribute("cy", y);
            target.setAttribute("r", r);
            target.setAttribute("fill", color);
            target.setAttribute("stroke", "black");
            target.setAttribute("stroke-width", width);
            mySVG.appendChild(target);
        },
        submit: function (id) {
            svgElement.submit(id, "point");
        },
    };
    var svgLine = {
        fun: function (x, y) {
            if (!drawPenDown) {
                // 起点
                drawPenDown = true;
                drawTargetID = svgElement.GenID();
                this.create(drawTargetID, x, y, x, y, drawVar.penColor);
            } else {
                // 终点
                drawPenDown = false;
                this.modify(x, y);
                svgElement.addTargetEventListeners(drawTargetID);
                this.submit(drawTargetID);
            }
            return drawTargetID;
        },
        drawNow: function (strID, x1, y1, x2, y2, color) {
            this.create(strID, x1, y1, x2, y2, color);
            svgElement.addTargetEventListeners(strID);
        },
        create: function (id, x1, y1, x2, y2, color) {
            let target = document.createElementNS(xmlns, "line");
            target.setAttribute("id", id);
            target.setAttribute("x1", x1.toString());
            target.setAttribute("y1", y1.toString());
            target.setAttribute("x2", x2.toString());
            target.setAttribute("y2", y2.toString());
            target.setAttribute("stroke", color);
            target.setAttribute("stroke-width", preValue.PenWidth);
            mySVG.appendChild(target);
        },
        modify: function (x2, y2) {
            var target = document.getElementById(drawTargetID);
            if (target) {
                target.setAttribute("x2", x2);
                target.setAttribute("y2", y2);
            } else {
                console.log("Target not found.")
            }
        },
        submit: function (strID) {
            svgElement.submit(strID, "line");
        },
    };
    var svgPolyline = {
        fun: function (x, y) {
            if (!drawPenDown) {
                // 起点
                drawTargetID = svgElement.GenID();
                drawPenDown = true;
                this.create(drawTargetID, x, y, drawVar.penColor);
            } else {
                // 续点
                this.addNode(x, y);
            }
        },
        drawNow: function (strID, cxy, color) {

        },
        create: function (id, x, y, color) {
            svgPolyPathArray = [[x, y]];
            let target = document.createElementNS(xmlns, "polyline");
            target.setAttribute("id", id);
            target.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
            target.setAttribute("fill", "none");
            target.setAttribute("stroke-width", preValue.PenWidth);
            target.setAttribute("stroke", color);
            mySVG.appendChild(target);
        },
        modify: function (x, y) {
            let targetPolyline = document.getElementById(drawTargetID);
            let tempArray = svgPolyPathArray.slice();
            tempArray.push([x, y]);
            targetPolyline.setAttribute("points", svgPath.ArrayToStr(tempArray));
        },
        addNode: function (x, y) {
            let targetPolyline = document.getElementById(drawTargetID);
            svgPolyPathArray.push([x, y]);
            targetPolyline.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
        },
        removeNode: function () {
            let targetPolyline = document.getElementById(drawTargetID);
            svgPolyPathArray.pop();
            targetPolyline.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
        },
        complete: function () {
            let targetPolyline = document.getElementById(drawTargetID);
            drawPenDown = false;
            if (svgPolyPathArray.length > 1) {
                targetPolyline.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
                svgElement.addTargetEventListeners(drawTargetID);
                this.submit(drawTargetID);
            } else {
                svgElement.remove(drawTargetID);
                svgElement.RollbackID();
            }
        },
        submit: function (strID) {
            svgElement.submit(strID, "polyline");
        },
    };
    var svgPolygon = {
        fun: function (x, y) {
            if (!drawPenDown) {
                // 起点
                drawTargetID = svgElement.GenID();
                drawPenDown = true;
                this.create(drawTargetID, x, y, drawVar.penColor);
            } else {
                // 续点
                this.addNode(x, y);
            }
        },
        drawNow: function (strID, cxy, color) {

        },
        create: function (id, x, y, color) {
            svgPolyPathArray = [[x, y]];
            let target = document.createElementNS(xmlns, "polygon");
            target.setAttribute("id", id);
            target.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
            target.setAttribute("fill", color);
            target.setAttribute("stroke-width", preValue.PenWidth);
            target.setAttribute("stroke", color);
            target.setAttribute("opacity", preValue.Opacity);
            mySVG.appendChild(target);
            target.addEventListener('mouseover', svgElement.onMouseOver);
            target.addEventListener('mouseout', svgElement.onMouseOut);
        },
        modify: function (x, y) {
            let targetPolygon = document.getElementById(drawTargetID);
            let tempArray = svgPolyPathArray.slice();
            tempArray.push([x, y]);
            targetPolygon.setAttribute("points", svgPath.ArrayToStr(tempArray));
        },
        addNode: function (x, y) {
            let targetPolygon = document.getElementById(drawTargetID);
            svgPolyPathArray.push([x, y]);
            targetPolygon.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
        },
        removeNode: function () {
            let targetPolygon = document.getElementById(drawTargetID);
            svgPolyPathArray.pop();
            targetPolygon.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
        },
        complete: function () {
            let targetPolygon = document.getElementById(drawTargetID);
            if (svgPolyPathArray.length > 1) {
                drawPenDown = false;
                targetPolygon.setAttribute("points", svgPath.ArrayToStr(svgPolyPathArray));
                svgElement.addTargetEventListeners(drawTargetID);
                this.submit(drawTargetID);
            } else {
                svgElement.remove(targetPolygon);
                svgElement.RollbackID();
            }
        },
        submit: function (id) {
            svgElement.submit(id);
        },
    };
    var svgText = {
        fun: function (x, y) {
            var str = prompt("请输入文字内容");
            drawTargetID = svgElement.GenID();
            this.create(drawTargetID, x, y, drawVar.penColor, str);
            svgElement.addTargetEventListeners(drawTargetID);
            this.submit(drawTargetID);
        },
        drawNow: function (strID, x, y, color, str) {
            this.create(strID, x, y, color, str);
            svgElement.addTargetEventListeners(drawTargetID);
        },
        create: function (strID, x, y, color, str) {
            let target = document.createElementNS(xmlns, "text");
            target.setAttribute("id", strID);
            target.setAttribute("x", x.toString());
            target.setAttribute("y", y.toString());
            target.setAttribute("fill", color);
            target.innerHTML = str;
            mySVG.appendChild(target);
            target.addEventListener('mouseover', svgElement.onMouseOver);
            target.addEventListener('mouseout', svgElement.onMouseOut);
        },
        submit: function (strID) {
            svgElement.submit(strID, "text");
        },
    };
    var svgElement = {
        GenID: function () {
            let strID = "svg" + preValue.SVGNumber.toString();
            while (document.getElementById(strID)) {
                preValue.SVGNumber++;
                strID = "svg" + preValue.SVGNumber.toString();
            }
            console.log("Gen SVG Element ID: " + strID);
            return strID;
        },
        CheckID: function (strID) {
            if (document.getElementById(strID)) {
                return this.GenID();
            } else {
                return strID;
            }
        },
        RollbackID: function () {
            preValue.SVGNumber--;
            //console.log("Roll ID Back to ", SVGNumber.toString());
        },
        remove: function (strID) {
            mySVG.removeChild(document.getElementById(strID));
        },
        submit: function (strID, type) {
            if(Debug){console.log("Submit ",type+' @ID="' + strID);}
            let target = document.getElementById(strID);
            labelsData.AddLabelHtml(timeFrame.getFrame(mediaVar.playerRef.currentTime), target.outerHTML, strID, type);
        },
        addTargetEventListeners(strID) {
            var target = document.getElementById(strID);
            target.addEventListener('mouseover', this.onMouseOver);
            target.addEventListener('mouseout', this.onMouseOut);
            // 添加效果后标注本帧已被修改
            drawFrameLabelModifyFlag = true;
        },
        onMouseOver: function (evt) {
            var obj = evt.target;
            obj.setAttribute("stroke-width", "5");
            //obj.setAttribute("fill","red");
            targetElement = obj.id;

        },
        onMouseOut: function (evt) {
            var obj = evt.target;
            obj.setAttribute("stroke-width", preValue.PenWidth.toString());
            //obj.setAttribute("fill",penColor);

        },
    };
    var svgPath = {
        ArrayToStr: function (svgArray) {
            let str = "";
            for (let i = 0; i < svgArray.length; i++) {
                str += (i > 0) ? " " : "";
                str += svgArray[i][0].toString() + "," + svgArray[i][1].toString();
            }
            return str;
        },
        ArrayToJson: function (svgArray) {
            let jsonObj = [];
            for (var i = 0; i < svgArray.length; i++) {
                var c = {
                    x: svgArray[i][0],
                    y: svgArray[i][1]
                };
                jsonObj.push(c);
            }
            return JSON.stringify(jsonObj);
        },
        JsonToArray: function (jsonStr) {
            let jsonObj = eval(jsonStr);
            var svgArray = [];
            for (var i = 0; i < jsonObj.length; i++) {
                var c = [jsonObj[i].x, jsonObj[i].y];
                svgArray.push(c);
            }
            return svgArray;
        },
    };

    var doMouse = {
        move: function (ele) {
            if (drawPenDown) {
                var offset = Offset.svgElement(authorData);
                var x = offset.x, y = offset.y;
                switch (drawToolMode) {
                    case "line":
                        svgLine.modify(x, y);
                        break;
                    case "polyline":
                        svgPolyline.modify(x, y);
                        break;
                    case "polygon":
                        svgPolygon.modify(x, y);
                        break;
                    default:
                        break;
                }
            }
        },
        downLeft: function (evt) {
            if (drawToolMode == 'LabelArea') {
                CRFLabelAreaMouseDownL(evt);
                return;
            }
            if (drawToolMode) {
                var offset = Offset.svgElement(evt);
                var x = offset.x, y = offset.y;
                switch (drawToolMode) {
                    case "select":
                        break;
                    case "point":
                        svgPoint.fun(x, y);
                        break;
                    case "line":
                        svgLine.fun(x, y);
                        break;
                    case "polyline":
                        svgPolyline.fun(x, y);
                        break;
                    case "text":
                        svgText.fun(x, y);
                        break;
                    case "polygon":
                        svgPolygon.fun(x, y);
                        break;
                    case "remove":
                        svgElement.remove(targetElement);
                        break;
                    default:
                        return;
                }
            }
        },
        downMiddle: function (ele) {
            if (Debug) {
                console.log("Mouse: Middle key pressed, but nothing to do :)");
            }
        },
        downRight: function (evt) {
            if (drawToolMode == 'LabelArea') {
                CRFLabelAreaMouseDownR(evt);
                return;
            }

            if (drawPenDown) {
                var offset = Offset.svgElement(evt);
                var x = offset.x, y = offset.y;
                switch (drawToolMode) {
                    case "polyline":
                        svgPolyline.complete();
                        break;
                    case "polygon":
                        svgPolygon.complete();
                        break;
                    default:
                        svgElement.remove(drawTargetID);
                        svgElement.RollbackID();
                       toolSelect("");
                }
            } else {
               toolSelect("");
            }
        },
    };
    function parseKeys(code) {
        switch (code) {
            case 27://esc
                doKeyPress.esc();
                break;
            case 32://space
            case 80://p
                doKeyPress.space();
                break;
            case 13://enter
                doKeyPress.enter();
                break;
            case 83:
                doKeyPress.s();
                break;
            case 68:
                doKeyPress.d();
                break;
            case 65:
                doKeyPress.a();
                break;
            case 84:
                doKeyPress.t();
                break;
            case 39://right
                MediaPanel.videoNextFrame();
                break;
            case 37://left
                MediaPanel.videoPrevFrame();
                break;
            case 90://z
                doKeyPress.undo();
                break;
            case 88://x
            case 85://u
                doKeyPress.redo();
                break;
            case 46://del
                doKeyPress.del();
                break;
            default:
                break;
        }
    }
    var doKeyPress = {
        esc: function () {
            if (drawPenDown) {
                switch (drawToolMode) {
                    case "polyline":
                        if (svgPolyPathArray.length > 1) {
                            svgPolyline.removeNode();
                            break;
                        }
                    case "polygon":
                        if (svgPolyPathArray.length > 1) {
                            svgPolygon.removeNode();
                            break;
                        }
                    default:
                        svgElement.remove(drawTargetID);
                        svgElement.RollbackID();
                       toolSelect("");
                }
            } else {
               toolSelect("");
            }
        },
        space: function () {
            MediaPanel.videoPlay();
        },
        enter: function () {
            var iFrame = timeFrame.getFrame(mediaVar.playerRef.currentTime);
            LabelsPanel.SaveFrame(iFrame);
            if (Debug) {
                console.log("手动提交，F=", iFrame);
            }
        },
        s: function () {
            let f = timeFrame.getFrame(mediaVar.playerRef.currentTime);
            if (Debug) {
                console.log("Frame:", f, "; Time:", mediaVar.playerRef.currentTime);
            }
            labelsData.PrintLabelsByFrame(f);
        },
        d: function () {
            Debug = !Debug;
            console.log("Debug Flag:", Debug.toString());
        },
        a: function () {
            if (Debug) {
                labelsData.PrintAllLabels();
            }
        },
        t: function () {
            if (Debug) {
                timeFrame.Print();
            }
        },
        undo:function() {
            if (globalVar.onDraw) {
                globalVar.targetObj.undo();
            }
        },
        redo:function () {
            if (globalVar.onDraw) {
                globalVar.targetObj.redo();
            }

        },
        del:function () {
            if (globalVar.onDraw) {
                globalVar.targetObj.cancel();
            }
        }

    };

    var Offset = {
        svgElement: function (ele) {
            var x = authorData.pageX, y = authorData.pageY;
            var offset = this.fromDocument(mySVG);
            x -= offset.left;
            y -= offset.top;
            return {
                x: x,
                y: y
            }
        },
        fromDocument: function (ele) {
            var totalLeft = 0, totalTop = 0, par = authorData.parentNode;
            // 首先加自己本身的左偏移和上偏移
            if (authorData.offsetTop) {
                totalTop += authorData.offsetTop;
            }
            if (authorData.offsetLeft) {
                totalLeft += authorData.offsetLeft;
            }
            while (par) {
                if (navigator.userAgent.indexOf("MSIE 8.0") === -1) {
                    totalLeft += par.clientLeft;
                    totalTop += par.clientTop;
                }
                totalLeft += par.offsetLeft;
                totalTop += par.offsetTop;
                par = par.offsetParent;
            }
            return {
                left: totalLeft,
                top: totalTop
            }
        }
    };

    function TimeFrameObj() {
        this.videoFrames = 0;
        this.videoFrameTime = [0];
        this.init = function (fps, duration) {
            this.videoFrames = (fps * duration);
            var frameStep = parseFloat((1 / fps).toFixed(6));
            var frameTime = 0;
            for (var i = 1; i <= this.videoFrames; i++) {
                frameTime += frameStep;
                //this.videoFrameTime.push((i / this.videoFrames * duration).toFixed(6));
                this.videoFrameTime.push(parseFloat(frameTime.toFixed(6)));
            }
        };
        this.getFrame = function (time) {
            for (var i = 0; i < this.videoFrames; i++) {
                if (this.videoFrameTime[i] == time) {
                    return i;
                }
                if (this.videoFrameTime[i] > time) {
                    return (i - 1);
                }
            }
            return -1;
        };
        this.getTime = function (frame) {
            return this.videoFrameTime[frame];
        };
        this.Print = function () {
            console.log("VideoFrames: ", this.videoFrames);
            console.log("VideoFrameTime: ", this.videoFrameTime);
        };
        this.getVideoTime = function () {


        }
    }
    function TimerObj() {
        var timerRef = null;
        // 定时器启动
        this.set = function () {
            this.clear();
            timerRef = setInterval(MediaPanel.updateFast, 1000 / preValue.FPS);
        };
        // 定时器结束
        this.clear = function () {
            if (timerRef) {
                clearInterval(timerRef);
            }
        }
    }
    function LabelsDataObj() {
        let frameData = []; // 标注帧通用对象
        this.Empty = function () {
            frameData = [];
        };
        this.RemoveFrame = function (iframe) {
            frameData[iframe] = [];
        };
        this.SetFrame = function (iframe, frameData) {
            frameData[iframe] = frameData;
        };
        this.GetFrame = function (iframe) {
            return frameData[iframe];
        };
        this.SaveFrame = function (iframe) {
            // 提交给服务器
            console.log("Upload Frame Label Data to Server.")
        };
        this.AddToFrame = function (iframe, frameData) {
            frameData.forEach(ele => {
                frameData.push(authorData);
            });
        };
        this.GetAllData = function () {
            return frameData;
        };
        this.SetAllData = function (data) {
            frameData = data;
        };
        this.AddLabelHtml = function (iframe, labelHtmlElement, strID, strType) {
            if (!frameData[iframe]) {
                frameData[iframe] = [];
            }
            let lbe = new LabelElement(strID, labelHtmlElement, strType, preValue.author, mediaVar.playerRef.currentTime);
            frameData[iframe].push(lbe);
        };
        this.AddLabelData = function (iframe, labelData) {
            if (!frameData[iframe]) {
                frameData[iframe] = [];
            }
            let lbe = new LabelElement(labelData.id, "", labelData, preValue.author, mediaVar.playerRef.currentTime);
            frameData[iframe].push(lbe);
        };
        this.RemoveLabelByFrame = function (iframe, labelID) {
            let ret = -1;
            let fd = frameData[iframe];
            if (fd) {
                ret = 0;
                for (let i = 0; i < fd.length; i++) {
                    if (fd[i].id == labelID) {
                        fd.splice(i, 1);
                        ret = 1;
                    }
                }
            }
            return ret;
        };
        this.RemoveLabel = function (labelID) {
            if (frameData) {
                for (let i = 0; i < frameData.length; i++) {
                    this.RemoveLabelByFrame(i, labelID)
                }
            }
        };
        this.PrintLabelsByFrame = function (iframe) {
            let fd = frameData[iframe];
            if (fd) {
                fd.forEach(ele => {
                    console.log(authorData);
                })
            } else {
                console.log("This cFrame has no label.");
            }
        };
        this.PrintAllLabels = function () {
            if (frameData) {
                frameData.forEach((ele, iframe, arr) => {
                    this.PrintLabelsByFrame(iframe);
                })
            }
        };
        this.drawFrame = function (iframe) {
            let fd = frameData[iframe];
            if (fd) {
                fd.forEach(ele => {
                    strID = svgElement.CheckID(authorData.id);
                    if (authorData.html) {
                        //console.log("DrawType:",authorData.type,"; HTML",authorData.html);
                        var target = null;
                        if (authorData.type == "point") {
                            target = document.createElementNS(xmlns, "circle");
                        } else {
                            target = document.createElementNS(xmlns, authorData.type);
                        }
                        mySVG.appendChild(target);
                        target.outerHTML = authorData.html;
                    }
                })
            }
        };
        this.drawElement = function (iframe, ele) {
            var strID = svgElement.CheckID(authorData.id);
            switch (authorData.type) {
                case "point":
                    svgPoint.drawNow(strID, authorData.c[0], authorData.c[1], authorData.color);
                    break;
                case "line":
                    svgLine.drawNow(strID, authorData.c[0][0], authorData.c[0][1], authorData.c[1][0], authorData.c[1][1], authorData.color);
                    break;
                case "polyline":
                    break;
                case "polygon":
                    break;
                case "text":
                    svgText.drawNow(strID, authorData.x, authorData.y, authorData.color, authorData.str);
                    break;
                default:
            }
        };
        this.ExportFrameToJson = function (iframe) {
            alert("On developing.");
        };
        this.LoadJsonfile = function () {
            var file=document.createElement('input');
            file.id='filereader';
            file.setAttribute('type','file');



        };
        this.SaveJsonfile = function () {
            if (frameData.length) {
                var data = JSON.stringify(frameData);
                var blob = new Blob([data], {type: "text/plain;charset=utf-8"});
                saveAs(blob, preValue.JsonFile);
            } else {
                alert("本内容无待保存标记信息！");
            }
        };
        this.ImportFromJson = function (strJson) {
            if (strJson=="") {
                return;
            }
            let temp=[];
            try {
                temp = eval(strJson);
            } catch (e) {
                alert(e);
                return;
            }
            if (!frameData) {
                frameData = [];
            }
            temp.forEach(ele => {
                var vTime = authorData[0].videoTime;
                var iFrame = timeFrame.getFrame(vTime);
                var data = authorData.slice();
                if (!frameData[iFrame]) {
                    frameData[iFrame] = [];
                }
                frameData[iFrame].push(data);
            });
        }
    }
    function LabelElement(strID,strHtml,strType,strAuthor,iVideoTime) {
        this.id=strID;
        this.html=strHtml;
        this.type=strType;
        this.data=[];
        this.author=strAuthor;
        this.createTime=new Date();
        this.videoTime=iVideoTime;
    }

    function CRFBtnOnClick(evt) {
        var obj=evt.target;
        // 暂停未完成工作
        if (globalVar.onDraw) {
            console.log('onPause:',obj.cid);
            globalVar.targetObj.onPause();
        }
        // 激活当前Target
        globalVar.TargetElement = document.getElementById(obj.cid);
        // 处理区域类标注
        if (obj.ctool=='area') {
            toolSelect('LabelArea');
            // 激活SVG背景响应
            globalVar.disableSvgBgOnClickFunc=false;
            if (globalVar.TargetElement) {
                // 非第一次创建本cid区域
                globalVar.targetObj=LabelsObjData[obj.cid];
                console.log(LabelsObjData,globalVar.targetObj);
                globalVar.targetObj.onContinue();
            } else {
                // 第一次创建本cid区域
                globalVar.targetObj = new LabelAreaObj(obj.cid,mediaVar.playerRef.currentTime,obj.ccolor);
                globalVar.targetObj.create(obj.cid);
            }
            return;
        }
        if (obj.ctool=='time') {
            toolSelect('LabelTime');

            return;
        }
        if (obj.ctool=='multi') {
            toolSelect('LabelMulti');

            return;
        }
        if (obj.ctool=='quility') {
            toolSelect('LabelQ');

            return;
        }
    }
    function CRFLabelAreaMouseDownL(evt) {
        if (globalVar.disableSvgBgOnClickFunc) {
            if(Debug){console.log('SVG Background onClick Function Disabled.')}
        } else {
            switch (globalVar.targetObj.mode) {
                case 'point':
                    var offset = Offset.svgElement(evt);
                    globalVar.targetObj.addPoint(offset.x, offset.y);
                    break;
                case 'mask':
                    break;
                default:
                    break;
            }
        }
    }
    function CRFLabelAreaMouseDownR(evt) {
        var obj=globalVar.targetObj;
        switch (obj.mode) {
            case 'point':
                obj.mode = 'mask';
                obj.maskArea();
                obj.hidePoints();
                obj.submit();
                break;
            case 'mask':
                obj.mode = 'modify';
                obj.showPoints();
                globalVar.onDrawing=true;
                break;
            case 'modify':
                obj.mode = 'mask';
                obj.hidePoints();
                globalVar.onDrawing=false;
                obj.submit();
                break;
            default:
                alert('特殊右键模式，已禁止:');
                console.log(obj);
                break;
        }
    }

    function LabelAreaObj (cid,time,color) {
        this.cTime = timeFrame.getTime();
        this.cFrame = timeFrame.getFrame();
        this.cSvgRef = null;
        this.mode='point';
        this.cPointNum=1;
        this.cR=3;
        this.cid=cid;
        this.cFillColor=color;
        this.cStrokeColor="black";
        this.cStrokeWidth=1;
        this.cPoints=new Array();
        this.undoList=[];
        this.create =function (id) {
            var obj = document.createElementNS(xmlns,"svg");
            obj.id=id;
            obj.cid=this.cid;
            this.cSvgRef=obj;
            globalVar.onDrawing=true;
            mySVG.appendChild(obj);
        };
        this.addPoint = function (x,y) {
            var id=this.cSvgRef.id + "_" + this.cPointNum;
            this.cPointNum++;
            var obj = this.drawPoint(x,y,id);
            this.cPoints[id]=convert.AbsoluteToRelative(x,y);
            this.undoList=[];
        };
        this.drawPoint=function(x,y,id){
            var obj =document.createElementNS(xmlns,"circle");
            obj.id=id;
            obj.setAttribute("cx", x);
            obj.setAttribute("cy", y);
            obj.setAttribute("r", this.cR);
            obj.setAttribute("fill", this.cFillColor);
            obj.setAttribute("stroke", this.cStrokeColor);
            obj.setAttribute("stroke-width", this.cStrokeWidth);
            obj.onmouseenter=doMouseEnter;
            obj.onmouseleave=doMouseLeave;
            obj.onmousedown=doMouseDown;
            obj.onmouseup=doMouseUp;
            obj.Dad=this;
            this.cSvgRef.appendChild(obj);
            return obj
        };
        this.delPointByID = function (id) {
            target = document.getElementById(id);

        };
        this.delLatestPoint = function() {
            var nodelist = this.cSvgRef.childNodes;
            if (nodelist.length > 0) {
                var target = nodelist[nodelist.length - 1];
                var data=this.cPoints[target.id];

                this.undoList.push([target,data]);

                delete this.cPoints[target.id];
                this.cSvgRef.removeChild(target);
            } else {
                console.log('无法再撤销。');
            }
        };
        this.redo=function () {
            var evt=this.undoList.pop();
            if (evt==null) {
                console.log('已重做至最新状态。')
                return;
            }
            var target = evt[0];
            var data = evt[1];

            this.cSvgRef.appendChild(target);
            this.cPoints[target.id]=data;
        };
        this.undo=function () {
            console.log(this.mode)
            if (this.mode=='point'){
                this.delLatestPoint();
            } else {
                alert('仅在点模式下允许撤销。')
            }
        };
        this.maskArea = function () {
            var id = this.cSvgRef.id + "_mask";
            let obj = document.getElementById(id);
            if (obj==null){
                obj = document.createElementNS(xmlns, "polygon");
                obj.id=id;
            }
            obj.setAttribute("points", convert.PointsArrayToStr(this.cPoints));
            obj.setAttribute("fill", this.cFillColor);
            obj.setAttribute("stroke-width", this.cStrokeWidth);
            obj.setAttribute("stroke", this.cStrokeColor);
            obj.setAttribute("opacity", preValue.Opacity);
            this.cSvgRef.insertBefore(obj,this.cSvgRef.firstChild)
        };
        this.hidePoints = function() {
            var nodelist = this.cSvgRef.childNodes;
            var maskID = this.cSvgRef.id + '_mask';
            for (var i=nodelist.length-1;i>=0;i--) {
                if (nodelist[i].id !=maskID){
                    //console.log(this.cPoints);
                    //delete this.cPoints[nodelist[i].id];
                    this.cSvgRef.removeChild(nodelist[i]);
                }
            }
        };
        this.showPoints = function() {
            for (var key in this.cPoints) {
                var p=convert.RelativeToAbsolute(this.cPoints[key][0],this.cPoints[key][1]);
                this.drawPoint(p[0],p[1],key);
                if(Debug){console.log(p)};
            }
        };
        this.submit = function () {
            globalVar.onDrawing=false;
            this.onPause();
            console.log('submit, coding uncomplete.');
            console.log('Time',this.cTime,'Points',this.cPoints);
            CRFPanel.reflushAllBtn();
        };
        this.destory = function() {
            console.log('send message to server to destory label information.');
        };
        this.onPause = function () {
            console.log('Pause label:',this.cSvgRef.cid);
            LabelsObjData[this.cSvgRef.cid]=this;
            this.hidePoints();
            console.log('LabelsData',LabelsObjData);
        };
        this.onContinue = function () {
            console.log('Continue label:',this.cSvgRef.cid);
        };
        this.cancel = function () {
            console.log('cancel objects.')
            if (this.mode == 'point') {
                // 点模式下未提交
            } else {
                // 其他模式已经submit，需要撤回
                this.destory();
            }
            console.log(LabelsObjData[this.cid]);
            this.cSvgRef.parentNode.removeChild(this.cSvgRef);
        }
    }

    // 鼠标键按下
    function doMouseDown (evt) {
        // 禁止SVGPlayer响应
        globalVar.disableSvgBgOnClickFunc=true;
        var obj=evt.target;
        switch (evt.button) {
            case 0:
                if (obj.selected) {
                    // 已选择点，放下
                    pointPutdown(obj);
                } else {
                    // 未选择点，拾起
                    pointPickup(obj);
                }
                break;
            case 2:
                if (obj.selected) {
                    // 已选择点，放弃
                    pointCancelMove(obj);
                }
                break;
        }
    }
    function doMouseUp(evt) {
        globalVar.disableSvgBgOnClickFunc=false;
    }
    function doMouseMove(evt) {
        var obj=evt.target;
        var offset = Offset.svgElement(evt);
        if (obj.selected) {
            pointMove(obj,offset.x,offset.y);
        }
    }
    // 鼠标进出效果
    function doMouseEnter (evt) {
        var obj = evt.target;
        obj.setAttribute('fill',this.Dad.cStrokeColor);
        obj.setAttribute('stroke',this.Dad.cFillColor);
    }
    function doMouseLeave(evt) {
        var obj = evt.target;
        obj.setAttribute("fill", this.Dad.cFillColor);
        obj.setAttribute("stroke", this.Dad.cStrokeColor);
    }

    function pointPickup (obj) {
        obj.selected = true;
        obj.undoX = obj.getAttribute('cx');
        obj.undoY = obj.getAttribute('cy');
        obj.onmousemove = doMouseMove;
        obj.onmouseout = doMouseMove;

    }
    function pointMove(obj,x,y) {
        obj.setAttribute('cx',x);
        obj.setAttribute('cy',y);
        obj.Dad.cPoints[obj.id]=convert.AbsoluteToRelative(x,y);
    }
    // 放弃点移动
    function pointCancelMove(obj) {
        obj.selected = false;
        obj.onmousemove = null;
        obj.onmouseout = null;
        obj.setAttribute('cx', obj.undoX);
        obj.setAttribute('cy', obj.undoY);
    }
    function pointPutdown(obj) {
        obj.selected = false;
        obj.onmousemove = null;
        obj.onmouseout = null;
        // 如处于修改模式，需要重绘mask
        if (obj.mode='mask') {
            obj.Dad.maskArea();
        }
    }

    // 转换工具集
    var convert = {
        RelativeToAbsolute: function (pW,pH) {
            var x=pW*mediaVar.cw/100;
            var y=pH*mediaVar.ch/100;
            return [x,y]
        },
        AbsoluteToRelative: function (x,y) {
            var pW=x*100/mediaVar.cw;
            var pH=y*100/mediaVar.ch;
            return [pW,pH]
        },
        PointsArrayToStr: function (points) {
            let str = '';
            for (var key in points) {
                var p=this.RelativeToAbsolute(points[key][0],points[key][1]);
                str += p[0].toFixed(0).toString()+','+p[1].toFixed(0).toString()+' ';
            }
            console.log(str);
            return str;
        },
        // 按n要求为num补零
        pad: function (num, n) {
            var len = num.toString().length;
            while (len < n) {
                num = "0" + num;
                len++;
            }
            return num;
        },
        // 形成模仿时间字符串
        formatTime: function (time) {
            time = time.toFixed(3);
            var Hour, Minute, Second, Msecond, lastTime, ret;
            if (time < 3600) {
                Hour = 0
            } else {
                Hour = parseInt(time / 3600)
            }
            lastTime = time - Hour * 3600;
            if (lastTime < 60) {
                Minute = 0
            } else {
                Minute = parseInt(lastTime / 60)
            }
            lastTime = lastTime - Minute * 60;
            if (lastTime < 1) {
                Second = 0
            } else {
                Second = parseInt(lastTime)
            }
            Msecond = parseInt(1000 * (lastTime - Second));
            ret = Hour.toString() + ":" + Minute.toString() + ":" + Second.toString() + "." + this.pad(Msecond, 3).toString();
            return ret;
        },
    };
    // 选择工具
    function toolSelect(selectTool) {
        drawPenDown = false;
        drawToolMode = selectTool;
        console.log('mode:',drawToolMode);
        if (!selectTool || selectTool === "") {
            Array.from($(".toolBtn")).forEach(v => {
                v.style.backgroundColor = 'white';
            })
        } else {
            Array.from($(".toolBtn")).forEach(v => {
                v.style.backgroundColor = (v.cid == selectTool) ? 'lightblue' : 'white';
            })
        }
    }
    // 事件绑定
    function bindEvent(ele, eventName, func) {
        if (window.addEventListener) {
            authorData.addEventListener(eventName, func);
        } else {
            authorData.attachEvent('on' + eventName, func);
        }
    }

    // 初始化界面
    Panels.init();
    // 初始化对象
    mainObject.init();
    // 载入预定义标注信息
    labelsData.ImportFromJson(preValue.JsonLabelData);
}(this, document));