$(function(window, document) {
    document.body.oncontextmenu=function () {return false};
    var Debug = true;
    var xmlns = "http://www.w3.org/2000/svg";

    let objTmp=Object.create(null);
    //let strTmp="";

    var xmlReqJson = new XMLHttpRequest();
    xmlReqJson.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            objTmp = JSON.parse(this.responseText);
        }
    };
    /*
    var xmlReqText = new XMLHttpRequest();
    xmlReqText.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            strTmp = this.responseText;
        }
    };
    */

    // 载入 正常标注按钮定义
    xmlReqJson.open("GET", "/js/crf/normal.json", false);xmlReqJson.send();
    let CRFButton_Common=objTmp;
    console.log("1",CRFButton_Common);

    // 载入 异常标注按钮定义
    xmlReqJson.open("GET", "/js/crf/tsc.json", false);xmlReqJson.send();
    let CRFButton_Spec=objTmp;
    console.log("2",CRFButton_Common);


    // 载入 时间节点定义
    xmlReqJson.open("GET", "/js/crf/time.json", false);xmlReqJson.send();
    let CRFButton_T=objTmp;


    // 载入 视频质量定义
    xmlReqJson.open("GET", "/js/crf/quility.json", false);xmlReqJson.send();
    let CRFButton_Q=objTmp;


    var mediaVar= {
        duration: null, coffsetxy: null, ch: null, cw: null,
        containerRef:null, playerRef: null,progressRef:null,
        ctrlBtnPlay: null,
    };
    var drawVar={
        penColor:null,
        selectedBkColor:"#0080FF",
        unselectBkColor:"#E0E0E0",
        mode:'',
        penDown:false,
        onDraw:false,
        onModify:false,
        svgRef:null,
    };
    var disableSvgBgOnClickFunc=false;

    var frameData=null;
    var frameObjs=[],btnObjs=[];
    var targetObj=null,targetBtn=null;
    var mySVG=null;

    var hiddenFrameObjs=false;

    // 初始化各类
    var timeFrame = new TimeFrameObj();
    var timer = new TimerObj();
    var lData=new LDataObj();
    mediaVar.progressRef=new VideoProgressObj();

    var Panels={
        init:function () {
            this.initWorkPanel();
            this.initCRFPanel();
        },
        initCRFPanel:function () {
            CRFPanel.init();
        },
        initWorkPanel:function () {
            switch(preValue.SourceType){
                case "video":
                case ".mp4":
                case ".ogv":
                    MediaPanel.init();
                    window.onresize = MediaPanel.adapt;
                    break;
                case "image":
                case ".jpg":
                    ImagePanel.init();
                    break;
                default:
                    alert('数据源类型 "'+preValue.SourceType+'" 无法识别。');
            }
        },
    };
    var MediaPanel={
        init:function () {
            var mediaContainer = document.createElement("div");
            mediaContainer.id = "mediaContainer";
            mediaContainer.className = "workSpaceBackground";
            document.getElementById('Workspace').append(mediaContainer);

            var videoPlayer = document.createElement("video");
            videoPlayer.id = "videoContainer";
            videoPlayer.className = "workSpaceOverlayer";
            videoPlayer.setAttribute("style", "z-index: 1");
            videoPlayer.setAttribute("preload", "auto");
            videoPlayer.setAttribute("src", preValue.Source);

            if (preValue.SourceType==".mp4") videoPlayer.setAttribute("type", "video/mp4");
            else  videoPlayer.setAttribute("type", "video/ogv");
            videoPlayer.innerText = "您所使用的浏览器不支持 HTML5 视频播放，请换用Chrome或Firefox浏览器（国产浏览器请切换至“急速模式”。";
            mediaContainer.append(videoPlayer);

            mySVG = document.createElementNS(xmlns, "svg");
            mySVG.id = "svgContainer";
            mySVG.className = "workSpaceOverlayer";
            mySVG.setAttribute("class", 'workSpaceOverlayer');
            mySVG.setAttribute("style", "z-index:2");
            mediaContainer.append(mySVG);

            mediaVar.progressRef.createBar(document.getElementById('Workspace'));

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
            document.getElementById('Workspace').append(videoControlCenter);

            mediaVar.playerRef = videoPlayer;
            mediaVar.containerRef = mediaContainer;
            mediaVar.ctrlBtnPlay = document.getElementById('videoPlayBtn');

            mediaVar.ctrlBtnPlay.onclick=this.videoPlay;
            document.getElementById('videoStopBtn').onclick=this.videoStop;
            document.getElementById('videoPrevFrameBtn').onclick=this.videoPrevFrame;
            document.getElementById('videoNextFrameBtn').onclick=this.videoNextFrame;

            videoPlayer.onloadedmetadata = this.adapt;
        },
        adapt:function() {
            let v = mediaVar.playerRef;
            let c = mediaVar.containerRef;
            var vratio = preValue.Width / preValue.Height;
            var cratio = c.offsetWidth / c.offsetHeight;
            if (vratio > cratio) {v.width = c.offsetWidth;}
            else {v.height = c.offsetHeight;}
            mediaVar.coffsetxy= (Offset.fromDocument(v));
            mediaVar.cw=v.offsetWidth;
            mediaVar.ch=v.offsetHeight;
            mySVG=document.getElementById('svgContainer');
            mySVG.setAttribute('width','100%');
            mySVG.setAttribute('height','100%');
            mySVG.style.setProperty('height',mediaVar.ch+'px');
            mySVG.style.setProperty('width',mediaVar.cw+'px');
            mediaVar.duration=v.duration;
            timeFrame.init(preValue.Frames, v.duration);
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
            let v = mediaVar.playerRef;
            if (v.paused || v.ended) {
                if (v.ended) {
                    v.currentTime = 0;
                }
                v.play();
                timer.set();
                MediaPanel.updateFast();
            } else {
                v.pause();
                timer.clear();
                MediaPanel.updateFull();
            }
        },
        videoStop:function () {
            let v = mediaVar.playerRef;
            v.pause();
            v.currentTime = 0;
            timer.clear();
            MediaPanel.updateFull();
        },
        videoPrevFrame:function(){
            let v = mediaVar.playerRef;
            v.pause();
            f=timeFrame.getVideoCurrentFrame(v.currentTime)-1;
            v.currentTime=(f<0)?0:timeFrame.getTime(f);
            MediaPanel.updateFull();
        },
        videoNextFrame:function(){
            let v = mediaVar.playerRef;
            v.pause();
            f=timeFrame.getVideoCurrentFrame(v.currentTime)+1;
            v.currentTime=(f<preValue.Frames)?timeFrame.getTime(f):0;
            MediaPanel.updateFull();
        },
        videoNextLabel:function(){
            console.log("Next label.");
            for (i=timeFrame.getVideoCurrentFrame()+1;i<=lData.frames.length;i++) {
                if(lData.frames[i]!=null){
                    MediaPanel.videoSkipToFrame(i);
                    break;
                }
            }
        },
        videoPrevLabel:function(){
            console.log("Prev label.");
            console.log("Next label.");
            for (i=timeFrame.getVideoCurrentFrame()-1;i>=0;i--) {
                if(lData.frames[i]!=null){
                    MediaPanel.videoSkipToFrame(i);
                    break;
                }
            }
        },
        videoSkipToFrame:function(f){
            let v=mediaVar.playerRef;
            v.currentTime=timeFrame.getTime(f);
            console.log(v.currentTime);
            MediaPanel.updateFull();
        },
        videoGetBufferedPercentage:function(){
            let v=mediaVar.playerRef;
            let buffered=v.buffered;
            if (buffered.length) {
                return buffered.end(0) / mediaVar.duration;
            } else {
                return null
            }
        },
        updateFast:function () {
            var vP = mediaVar.playerRef;
            var percent = vP.currentTime / vP.duration;
            var f = timeFrame.getVideoCurrentFrame();
            mediaVar.progressRef.setPlayProgress(percent);
            document.getElementById('videoPlaytime').innerText=convert.formatTime(vP.currentTime) + " / " + convert.formatTime(mediaVar.duration);//控制条时间
            mediaVar.ctrlBtnPlay.innerText= (vP.paused || vP.ended) ? "Play" : "Pause";// 调整按键文字
            disableSvgBgOnClickFunc=false;
            hiddenFrameObjs=false;
            frameData =lData.frames[f];
            CRFPanel.refreshTimeBtn(f);
        },
        updateFull:function () {
            console.log('Upd Full');
            MediaPanel.updateFast();
            CRFPanel.refreshAllBtn();
            CRFRedrawAllLabels();
        }
    };
    var CRFPanel = {
        init: function () {
            var b=null;
            b=this.addBlock('blkTime');
            this.addText('>> 时间点标注',b);
            this.addButtons(CRFButton_T,'time',b);

            
            b=this.addBlock('blkCommon');
            this.addText('>> 通用信息标注',b);
            this.addButtons(CRFButton_Common,'com',b);
            this.addButtons(CRFButton_Common,'text',b);


            b=this.addBlock('blkSpec');
            this.addText('>> 病理影像标注',b);
            this.addButtons(CRFButton_Spec,'com',b);
			
			b=this.addBlock('blkQ');
            this.addText('>> 质量评估',b);
            this.addButtons(CRFButton_Q,'q',b);

        },
        addButtons: function (dat,tool,ref) {
            for (var cid in dat){
                authorData=dat[cid];
                if (authorData.tool != tool) continue;
                var obj = document.createElement("button");
                //obj.id=authorData.id;
                obj.className='crfBtn crfBtnNone';
                obj.innerText = authorData.name;

                obj.cid = authorData.cid;
                obj.tool=authorData.tool;

                ref.appendChild(obj);
                obj.onclick = CRFBtnOnClick;

                btnObjs[authorData.cid]=obj;
            }
        },
        addText:function (string,ref) {
            var p=document.createElement('p');
            p.innerHTML=string;
            ref.appendChild(p)
        },
        addBlock:function (strDiv) {
            var crfPanelRef = document.getElementById('Tools');
            var d=document.createElement('div');
            d.id=strDiv;
            d.className="labelBlk";
            crfPanelRef.appendChild(d);
            return d;
        },
        refreshAllBtn:function () {
            let f=timeFrame.getVideoCurrentFrame();
            this.refreshTimeBtn(f);
            this.refreshAreaBtn(f);
            this.refreshSpecBtn(f);
            this.refreshQBtn(f);
        },
        refreshTimeBtn:function() {
            let ref = document.getElementById('blkTime');
            if (frameData) {
                ref.childNodes.forEach(ele => {
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        if (frameData.cid == authorData.cid) authorData.className = 'crfBtn crfBtnTime';
                        else authorData.className = 'crfBtn crfBtnNone';

                        if (authorData.cid == 'SPEC'){
                            if (frameData.cdescribe=='收缩末期'||frameData.cdescribe=='舒张末期') authorData.innerText='特殊时间';
                            else authorData.innerText='特殊时间\r' + frameData.cdescribe;
                        }
                    }
                })
            } else {
                ref.childNodes.forEach(ele => {
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        authorData.className = 'crfBtn crfBtnNone';
                        if (authorData.cid=='SPEC') authorData.innerText='特殊时间';
                    }
                })
            }
        },
        refreshAreaBtn:function() {
            let ref = document.getElementById('blkCommon');
            if (frameData) {
                ref.childNodes.forEach(ele => {
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        if (drawVar.onDraw && targetObj && targetObj.cid == authorData.cid) authorData.className = 'crfBtn crfBtnLabeling';
                        else {
                            let data = frameData.clabels[authorData.cid];
                            authorData.className = (data) ? 'crfBtn crfBtnArea' : 'crfBtn crfBtnNone';
                            authorData.style.setProperty('background-color', (data) ? CRFButton_Common[authorData.cid].color : CRFButton_Common['default'].color);
                        }
                    }
                })
            } else {
                ref.childNodes.forEach(ele => {
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        authorData.className = 'crfBtn crfBtnNone';
                        authorData.style.setProperty('background-color', CRFButton_Common['default'].color);
                    }
                })
            }
        },
        refreshSpecBtn:function () {
            let ref=document.getElementById('blkSpec');
            if (frameData) {
                ref.childNodes.forEach(ele => {
                    console.log("refreshSpecBtn.authorData",authorData);
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        if (drawVar.onDraw && targetObj && targetObj.cid == authorData.cid) authorData.className = 'crfBtn crfBtnLabeling';
                        else {
                            let data = frameData.clabels[authorData.cid];
                            authorData.className = (data) ? 'crfBtn crfBtnArea' : 'crfBtn crfBtnNone';
                            authorData.style.setProperty('background-color', (data) ? CRFButton_Spec[authorData.cid].color : CRFButton_Common['default'].color);
                        }
                    }
                })
            } else {
                ref.childNodes.forEach(ele => {
                    if (authorData.nodeName.toLowerCase() == 'button') {
                        authorData.className = 'crfBtn crfBtnNone';
                        authorData.style.setProperty('background-color', CRFButton_Common['default'].color);
                    }
                })
            }
        },
        refreshQBtn:function() {
            let ref = document.getElementById('blkQ');
            let vqcid = 'VQ0', vccid = 'VC0', fqcid = 'FQ0';
            if (lData && lData.q != 0) vqcid = 'VQ' + lData.q;
            if (lData && lData.c != 0) vccid = 'VC' + lData.c;
            if (frameData && frameData.q != 0) fqcid = 'FQ' + frameData.q;
            ref.childNodes.forEach(ele => {
                if (authorData.nodeName.toLowerCase() == 'button') {
                    if (authorData.cid == vqcid || authorData.cid == vccid || authorData.cid == fqcid) authorData.className = 'crfBtn crfBtnQ';
                    else authorData.className = 'crfBtn crfBtnNone';
                }
            })
        }
    };
    var mainObj = {
        init: function () {
            var that = this;
            //videoVar.playerRef.removeAttribute("controls");
            this.bindFunctions();
            this.videoOperateControls();
        },
        videoOperateControls: function () {
            ref=mediaVar.progressRef.backgroundRef;
            bindEvent(ref, "mousedown", videoObj.onSeekbar);
        },
        bindFunctions: function () {
            bindEvent(mediaVar.playerRef, "ended", videoObj.onEnded);
            mySVG.onmousemove = function (e) {doMouse.move(e);};// 处理鼠标移动
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
            };// 处理鼠标按键
            document.onkeydown = function (evt) {
                var theEvent = window.event || evt;
                var code = theEvent.keyCode || theEvent.which;
                parseKeys(code);
            };
            if(document.addEventListener){document.addEventListener('DOMMouseScroll',doMouseScroll,false);} //W3C
            window.onmousewheel=document.onmousewheel=doMouseScroll;//IE/Opera/Chrome/Safari
        }
    };
    var videoObj = {
        // 进度条鼠标点击动作
        onSeekbar: function (ele) {
            var pA=mediaVar.progressRef.backgroundRef;
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
                timer.clear();
                timer.set();
            } else {
                timer.clear()
            }
            MediaPanel.updateFast();
        },
    };

    function parseKeys(code) {
        //console.log(code.toString());
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
            case 68://d
                doKeyPress.del();
                break;
            case 81:
                doKeyPress.q();
                break;
            case 72:
                doKeyPress.h();
                break;
            case 36:
                doKeyPress.home();
                break;
            case 35:
                doKeyPress.end();
                break;
            case 188:// <
                doKeyPress.less();
                break;
            case 190:// >
                doKeyPress.great();
                break;
            default:
                break;
        }
    }
    var doMouse = {
        move: function (ele) {
            if (drawVar.penDown) {
                var offset = Offset.svgElement(authorData);
                var x = offset.x, y = offset.y;
                switch (drawVar.mode) {
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
            if (drawVar.mode) {
                var offset = Offset.svgElement(evt);
                var x = offset.x, y = offset.y;
                switch (drawVar.mode) {
                    case 'LabelArea':
                        CRFLabelAreaMouseDownL(evt);
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
            if (drawVar.mode) {
                var offset = Offset.svgElement(evt);
                var x = offset.x, y = offset.y;
                switch (drawVar.mode) {
                    case 'LabelArea':
                        CRFLabelAreaMouseDownR(evt);
                        break;
                    default:
                        console.log('Switch default:', drawVar.penDown)
                }
            } else {
                toolSelect("");
            }
        },
    };
    var doKeyPress = {
        esc: function () {
            if (targetObj) {
                targetObj.onPause();
                targetObj=null;
            }
            drawVar.penDown=false;
            targetBtn=null;
            MediaPanel.updateFull();
        },
        space: function () {
            MediaPanel.videoPlay();
        },
        enter: function () {
            var f = timeFrame.getFrame(mediaVar.playerRef.currentTime);
            LDataUpload();
        },
        s: function () {
            var f = timeFrame.getVideoCurrentFrame();
            if (Debug) {
                console.log('frameData',frameData);
                console.log('frameObjs',frameObjs,'targetObj',targetObj);
                console.log('targetBtn',targetBtn);
                console.log('video',mediaVar.playerRef.buffer)
            }
        },
        d: function () {
            Debug = !Debug;
            console.log("Debug Flag:", Debug.toString());
        },
        a: function () {
            if (Debug) {
                console.log('LData:',lData);
            }
        },
        t: function () {
            if (Debug) {
                timeFrame.Print();
            }
        },
        undo:function() {
            if (drawVar.onDraw) {
                targetObj.undo();
            }
        },
        redo:function () {
            if (drawVar.onDraw) {
                targetObj.redo();
            }

        },
        del:function () {
            drawVar.onDraw=false;
            drawVar.penDown=false;
            let ans =confirm("确认删除本结构标注？");
            if (ans) {
                if (targetBtn) {
                    cid=targetBtn.cid;
                    targetObj=frameObjs[cid];
                    delete frameObjs[cid];
                    delete frameData.clabels[cid];
                    if (targetObj) targetObj.remove();
                }
            }
            LDataUpload();
            MediaPanel.updateFull();
        },
        q:function() {
            ans=confirm('将清除本帧全部标注信息，是否确认？');
            if (ans) {
                f=timeFrame.getVideoCurrentFrame();
                CRFResetFrame(f);
                LDataUpload();
            }
        },
        h:function() {
            if (hiddenFrameObjs) {
                hiddenFrameObjs=false;
                //show all
                UserMessage('显示标注信息',false);
                MediaPanel.updateFull()
            } else {
                hiddenFrameObjs=true;
                //hidden all
                UserMessage('隐藏标注信息',false);
                for (var cid in frameObjs) {
                    if (frameObjs[cid]) {
                        frameObjs[cid].remove();
                        delete frameObjs[cid];
                    }
                }
            }
        },
        home:function(){
            MediaPanel.videoSkipToFrame(0);
        },
        end:function() {
            MediaPanel.videoSkipToFrame(timeFrame.getLastFrame());
        },
        less:function(){
            MediaPanel.videoPrevLabel();
        },
        great:function(){
            MediaPanel.videoNextLabel();
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
        this.init = function (frames, duration) {
            this.videoFrames = frames;
            var frameStep =   parseFloat((duration/frames).toFixed(6));
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
        this.getVideoCurrentTime = function () {
            return mediaVar.playerRef.currentTime;
        };
        this.getVideoCurrentFrame = function () {
            return this.getFrame(mediaVar.playerRef.currentTime);
        };
        this.getLastFrame = function() {
            return this.videoFrames;
        }

    }
    function TimerObj() {
        var timerRef = null;
        // 定时器启动
        this.set = function () {
            this.clear();
            sleeptime=1000/2*preValue.Duration/preValue.Frames;
            timerRef = setInterval(MediaPanel.updateFull, sleeptime );
        };
        // 定时器结束
        this.clear = function () {
            if (timerRef) {
                clearInterval(timerRef);
            }
        }
    }
    function VideoProgressObj() {
        this.backgroundRef=null;
        this.bufferedRef=null;
        this.playProgressRef=null;
        this.strRef=null;
        var monitorTimer=null;
        this.createBar=function (div) {
            this.backgroundRef = document.createElement("div");
            this.backgroundRef.id="videoProgressBackground";
            this.backgroundRef.className = "videoProgressBackground";

            this.playProgressRef = document.createElement("div");
            this.playProgressRef.id="videoProgressFrontend";
            this.playProgressRef.className = "videoProgressFrontend";

            this.bufferedRef=document.createElement('div');
            this.bufferedRef.id='videoProgressBuffered';
            this.bufferedRef.className='videoProgressBuffered';

            this.strRef = document.createElement("span");
            this.strRef.id="videoProgressStr";
            this.strRef.className='videoProgressStr';
            this.strRef.innerText = "0%";

            this.playProgressRef.appendChild(this.strRef);
            this.backgroundRef.appendChild(this.bufferedRef);
            this.backgroundRef.appendChild(this.playProgressRef);
            div.append(this.backgroundRef);
        };
        this.setBufferProgress=function(percent) {
            let width = percent * (this.backgroundRef.offsetWidth) + "px";//调整控制条长度
            this.bufferedRef.style.width = width;
                console.log('Set Buffer:',percent,"width:",width);
        };
        this.setPlayProgress=function(percent) {
            console.log('Set PlayProgress',percent);
            this.strRef.innerText = ((percent * 100).toFixed(0) + "%");//进度条文字进度
            this.playProgressRef.style.width = percent * (this.backgroundRef.offsetWidth) + "px";//调整控制条长度
        };
        this.monitorBufferStatus=function() {
            let buf=MediaPanel.videoGetBufferedPercentage();
            mediaVar.progressRef.setBufferProgress(buf);
            if (buf>=1) {
                console.log('Buf End.');
                clearInterval(monitorTimer);
            }
        };
        this.startMonitorBuffer=function(){
            monitorTimer = setInterval(this.monitorBufferStatus,1000);
        };
        this.adapt=function () {
            this.backgroundRef.width=mediaVar.cw;
            console.log("wid",this.backgroundRef.width);
        }
    }

    function CRFBtnOnClick(evt) {
        targetBtn=evt.target;
        let cid=targetBtn.cid;
        if (drawVar.onDraw) {targetObj.onPause();targetObj=null;}
        let btnCID={};
        switch (targetBtn.tool) {
            case 'time':
                toolSelect('LabelTime');
                let t = timeFrame.getVideoCurrentFrame();
                CRFAddTimeLabel(cid,t,targetBtn.innerText);
                break;
            case 'q':
                toolSelect('LabelQ');
                CRFLabelQ(targetBtn.cid);
                CRFPanel.refreshQBtn();
                break;
            case 'com':
				if (CRFButton_Common[cid]) {
					btnCID=CRFButton_Common[cid];				
				} else {
					btnCID=CRFButton_Spec[cid];
				}
				if (!frameData) {alert('请先进行>>时间点标注<<');break;}
                toolSelect('LabelArea');
                // 激活SVG背景响应
                disableSvgBgOnClickFunc = false;
				targetBtn.style.backgroundColor=btnCID.color;
                targetBtn.className='crfBtn crfBtnAreaLabeling';
                // 获取按钮 cid 对应元素引用
                if (document.getElementById(cid)) {
                    // 非第一次创建本cid区域
                    targetObj = frameObjs[cid];
                    if (targetObj) targetObj.onContinue();
                } else {
                    // 第一次创建本cid区域
                    targetObj = new LFrameAreaObj(cid, btnCID.color);
                    targetObj.create(cid);
                    frameObjs[cid]=targetObj;
                }
                break;
            case 'imp':
                if (!frameData) {alert('请先进行>>时间点标注<<');break;}
                toolSelect('LabelMulti');
                break;
            case 'text':
                alert('功能未开放，请等待。');
                return;




                console.log("markUp.");

                //Download
                let data={};
                data.mid=preValue.Mid;
                data.direction='download';
                console.log(data);
                xmlReqText.open("POST", "/comment", false);
                xmlReqText.send(data);
                let strResult=strTmp;
                console.log('DownloadResult',strResult);

                //Upload
                data={};
                data.mid=preValue.Mid;
                data.direction='upload';
                data.data="Hello";
                xmlReqText.open("POST", "/comment", false);
                xmlReqText.send(data);
                strResult=strTmp;
                console.log('UploadResult',strResult);



                break;
        }
        lData.frames[timeFrame.getVideoCurrentFrame()]=frameData;
        LDataUpload();
    }
    function CRFRedrawAllLabels() {
        frameObjs = [];
        targetObj=null;
        var childs = mySVG.childNodes;
        for (let i = childs.length - 1; i > -1; i--) mySVG.removeChild(childs[i]);
        if (frameData && frameData.clabels) {
            for (let cid in frameData.clabels) {
                if (frameData.clabels[cid]) {
                    var color;
                    if (CRFButton_Common[cid]) {
                        color = CRFButton_Common[cid].color;
                    } else {
                        color = CRFButton_Spec[cid].color;
                    }
                    let obj = new LFrameAreaObj(cid, color);
                    obj.create(cid);
                    obj.redraw(frameData.clabels[cid].cPoints);
                    frameObjs[cid] = obj;
                }
            }
        }
        drawVar.onDraw = false;
    }
    function CRFAddTimeLabel (cid,frame,cdescribe) {
        if (frameData == null) {
            cdescribe = (cid=='SPEC')?prompt('请输入事件描述',cdescribe):cdescribe;
            frameData=new LFrameObj(cid,frame,cdescribe);
        } else {
            if (frameData.cid != cid) {
                let ans = confirm('确认调整时间点标注？');
                if (ans) {
                    frameData.cdescribe = (cid=='SPEC')?prompt('请输入事件描述',frameData.cdescribe):cdescribe;
                    frameData.cid=cid;
                }
            }
        }
        lData.frames[timeFrame.getVideoCurrentFrame()]=frameData;
        console.log(lData);
        CRFPanel.refreshTimeBtn(frame);
    }
    function CRFLabelAreaMouseDownL(evt) {
        if (disableSvgBgOnClickFunc == false) {
            switch (targetObj.mode) {
                case 'point':
                    var offset = Offset.svgElement(evt);
                    targetObj.addPoint(offset.x, offset.y);
                    break;
                case 'mask':
                    break;
                default:
                    break;
            }
        }
    }
    function CRFLabelAreaMouseDownR(evt) {
        let obj=frameObjs[targetBtn.cid];
        if (obj) {
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
                    drawVar.onDraw = true;
                    break;
                case 'modify':
                    obj.mode = 'mask';
                    obj.hidePoints();
                    drawVar.onDraw = false;
                    obj.submit();
                    break;
                default:
                    alert('特殊右键模式，已禁止:');
                    break;
            }
        }
    }
    // 鼠标键按下
    function doMouseDown (evt) {
        // 禁止SVGPlayer响应
        disableSvgBgOnClickFunc=true;
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
        disableSvgBgOnClickFunc=false;
    }
    function doMouseMove(evt) {
        var obj=evt.target;
        var offset = Offset.svgElement(evt);
        if (obj.selected) {
            pointMove(obj,offset.x,offset.y);
        }
    }
    function doMouseScroll(evt) {
        evt = evt || window.event;
        if (evt.wheelDelta) {//IE/Opera/Chrome
            if (evt.wheelDelta < 0) {
                MediaPanel.videoNextFrame();
            } else {
                MediaPanel.videoPrevFrame();
            }
        } else if (evt.detail) {//Firefox
            if (evt.detail < 0) {
                MediaPanel.videoPrevFrame();
            } else {
                MediaPanel.videoNextFrame();
            }
        }
    }
    // 鼠标进出效果
    function doMouseEnter(evt) {
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
        let mode=obj.Dad.mode;
        obj.selected = false;
        obj.onmousemove = null;
        obj.onmouseout = null;
        // 如处于修改模式，需要重绘mask
        if (mode=='mask'||mode == 'modify') {
            obj.Dad.maskArea();
        }
    }

    function CRFResetFrame(f) {
        delete lData.frames[f];
        frameData=null;
        frameObjs=[];
        MediaPanel.updateFull();
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
        drawVar.penDown = false;
        drawVar.mode = selectTool;
    }
    // 事件绑定
    function bindEvent(ele, eventName, func) {
        if (window.addEventListener) {
            authorData.addEventListener(eventName, func);
        } else {
            authorData.attachEvent('on' + eventName, func);
        }
    }
    function CRFLabelQ(cid) {
        let q=cid.substring(2,3);
        switch (cid.substring(0,2)){
            case 'VQ':
                lData.q=q;
                break;
            case 'VC':
                lData.c=q;
                break;
            case 'FQ':
                if (!frameData) {alert('请先进行>>时间点标注<<');break;}
                frameData.q=q;
                break;
        }
    }
    function LDataObj() {
        this.q=0;
        this.c=0;
        this.frames=[];
    }
    function LFrameObj(cid,frame,describe) {
        this.cid=cid;
        this.cframe=frame;
        this.cdescribe=describe;
        this.clabels={};
        this.q=0;
    }
    function LFrameAreaObj (cid,color) {
        this.cTime = timeFrame.getVideoCurrentTime();
        this.cFrame = timeFrame.getVideoCurrentFrame();
        this.cSvgRef = null;
        this.mode='point';
        this.cPointNum=1;
        this.cR=preValue.PointR;
        this.cid=cid;
        this.cFillColor=color;
        this.cStrokeColor="black";
        this.cStrokeWidth=1;
        this.cPoints={};
        this.undoList=[];
        this.create =function (id) {
            var obj = document.createElementNS(xmlns,"svg");
            obj.id=id;
            obj.cid=this.cid;
            this.cSvgRef=obj;
            drawVar.onDraw=true;
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
        this.redraw=function(cPoints){
            this.cPoints=cPoints;
            this.maskArea();
            this.mode='mask';
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
                UserMessage('无法再撤销。',false);
            }
        };
        this.redo=function () {
            var evt=this.undoList.pop();
            if (evt==null) {
                UserMessage('已重做至最新状态。',false);
                return;
            }
            var target = evt[0];
            var data = evt[1];

            this.cSvgRef.appendChild(target);
            this.cPoints[target.id]=data;
        };
        this.undo=function () {
            if (this.mode=='point'){
                this.delLatestPoint();
            } else {
                UserMessage('仅在point模式下允许撤销,当前模式为：'+this.mode,true)
            }
        };
        this.maskArea = function () {
            let id=this.cSvgRef.id + "_mask";
            let obj=document.getElementById(id);
            if (obj==null){
                obj=document.createElementNS(xmlns, "polygon");
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
                    this.cSvgRef.removeChild(nodelist[i]);
                }
            }
        };
        this.showPoints = function() {
            for (let key in this.cPoints) {
                var p=convert.RelativeToAbsolute(this.cPoints[key][0],this.cPoints[key][1]);
                this.drawPoint(p[0],p[1],key);
            }
        };
        this.submit = function () {
            console.log('submit:',this.cid);
            this.onPause();
            LDataUpload();
        };
        this.onPause = function () {
            let cid = this.cid;
            drawVar.onDraw = false;
            this.hidePoints();
            if (Object.keys(this.cPoints).length > 1) {
                // 已标记两个以上的点
                btnObjs[cid].className = 'crfBtn crfBtnArea';
                frameObjs[cid] = this;
                frameData.clabels[cid] = new FrameLabelDataObj(this.cid, this.cPoints);
                lData.frames[timeFrame.getVideoCurrentFrame()] = frameData;
            } else {
                btnObjs[this.cid].className = 'crfBtn crfBtnNone';
                btnObjs[this.cid].style.backgroundColor= '';
                frameObjs[cid] = null;
                target=document.getElementById(cid);
                target.parentElement.removeChild(target);
            }
        };
        this.onContinue = function () {
            targetObj=frameObjs[this.cid];
            console.log('Continue label:',this.cid,targetObj);
            parEle = targetObj.cSvgRef.parentNode;
            parEle.appendChild(targetObj.cSvgRef, parEle.lastChild);
        };
        this.cancel = function () {
            cid=this.cid;
            console.log('Removing cid:',cid);
            this.cSvgRef.parentNode.removeChild(this.cSvgRef);
        };
        this.remove = function () {
            this.cSvgRef.parentNode.removeChild(this.cSvgRef);
        }
    }
    function FrameLabelDataObj(cid,cPoints,cType) {
        this.cPoints=cPoints;
        this.cid=cid;
        this.ctype=cType;
        this.cTime=timeFrame.getVideoCurrentTime();
        this.cFrame=timeFrame.getVideoCurrentFrame();
    }
    function UserMessage(strMessage,isWarning){
        document.getElementById('infoLog').innerText=strMessage;
        document.getElementById('InfoPanel').style.backgroundColor=isWarning?"red":"#ffffaa";
    }

    function LDataUpload() {
        let data={};
        data.mid=preValue.Mid;
        data.data=JSON.stringify(lData);
        data.direction='upload';
        $.ajax({
            type:'post',
            url:'/labeldata',
            data:JSON.stringify(data),
            contentType:'application/json',
            dataType:'json',
            success:function(data){
                if (data.success) UserMessage("上传成功",false);
                else UserMessage("上传失败",true);
                console.log('上传成功：',lData);
            }
        });
    }
    function LDataDownload() {
        let data={};
        data.mid=preValue.Mid;
        data.direction='download';
        $.ajax({
            type:'post',
            url:'/labeldata',
            data:JSON.stringify(data),
            contentType:'application/json',
            dataType:'json',
            success:function(data) {
                if (data.success) {
                    lData=JSON.parse(data.data);
                    frameData=lData.frames[timeFrame.getVideoCurrentFrame()];
                    MediaPanel.updateFull();
                    UserMessage("已载入您的历史标注信息。",false);
                    console.log(lData)
                }
            }
        });
    }
    function CalcSection(data) {
        let maxlen=data.x.length-1;
        let ret =0;
        if (maxlen>1) { // 满足3点
            for (let i=0;i<maxlen;i++) {
                ret += data.x[i + 1] * data.y[i] - data.x[i] * data.y[i + 1];
            }
            ret += data.x[0]*data.y[maxlen]- data.x[maxlen]*data.y[0];
            ret = ret/2;
            if (ret<0) {ret=-ret};
            return ret
        }
    }
    function CRFGetCalc(labelData) {
        
    }
    // 初始化界面
    Panels.init();
    // 初始化对象
    mainObj.init();
    // 视频预载进度
    mediaVar.progressRef.startMonitorBuffer();
    // 数据预载
    LDataDownload();

    //data={};
    //data.x=[1,2,1];
    //data.y=[1,1,2];
    //console.log(CRFGetLabelSection(data));

}(this, document));