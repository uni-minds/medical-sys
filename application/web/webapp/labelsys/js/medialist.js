$(function(window, document) {
    var headers = ['文件名', '类型', '标注人员','更新日期','视频质量','标注总帧数','标注总数量','所属组','上传日期','备注'];
    var indexes = ['fn',    'tp',   'aus',   'udt',     'vq',    'lfs',        'lnum',      'og'    ,'ult'   ,'comment'];

    var from,count,isHaveNext=true,isHavePrev=false;
    var btnNext,btnPrev;
    updateMediaList(MediaReq);

    var context=document.getElementById('context');
    var pagectrl=document.getElementById('pagectrl');
    createButton(pagectrl);

    function updateTable() {
        if (document.getElementById('mediatb')) removeTable();
        createTable(context, headers, MediaList);
    }

    function createTable(parent, headers, datas) {
        var table = document.createElement("table");
        table.id = "mediatb";
        parent.appendChild(table);

        var thead = document.createElement("thead");
        table.appendChild(thead);
        var tr = document.createElement("tr");
        tr.className="info";
        thead.appendChild(tr);
        for (var i = 0; i < headers.length; i++) {
            var th = document.createElement("th");
            th.innerHTML = headers[i];
            th.className=indexes[i];
            tr.appendChild(th);
        }

        var tbody = document.createElement("tbody");
        table.appendChild(tbody);
        for (i = 0; i < datas.length; i++) {
            tr = document.createElement("tr");
            tr.className="info";
            tbody.appendChild(tr);

            for (var j=0;j<indexes.length;j++){
                var td=document.createElement("td");
                if (indexes[j]=='fn') {
                    var a=document.createElement("a");
                    td.appendChild(a);
                    a.innerText=datas[i][indexes[j]];
                    a.setAttribute('href',"/labelsys?mid="+datas[i]._id);
                } else {
                    td.innerHTML=datas[i][indexes[j]];
                }
                tr.appendChild(td);
            }
        }

    }

    function createButton(parent) {
        btnNext =document.createElement("button");
        btnNext.id='btnNext';
        btnNext.innerText='Next';
        btnNext.onclick=function(){
            var req=MediaReq;
            req.from+=req.count;
            updateMediaList(req);
        };

        btnPrev =document.createElement("button");
        btnPrev.id='btnPrev';
        btnPrev.innerText='Prev';
        btnPrev.onclick=function(){
            var req=MediaReq;
            req.from-=req.count;
            if (req.from<0) req.from=0;
            updateMediaList(req)
        };

        parent.appendChild(btnPrev);
        parent.appendChild(btnNext);
    }

    function updateButton() {
        btnPrev.disabled = !isHavePrev;
        btnNext.disabled = !isHaveNext;
    }

    function removeTable() {
        var tb=document.getElementById('mediatb');
        tb.parentNode.removeChild(tb);
    }

    function updateMediaList(req) {
        $.ajax({
            type:'POST',
            url:'/medialist',
            data:JSON.stringify(req),
            dataType:'json',
            contentType:'application/json',
            success:function(data){
                MediaList=JSON.parse(data.data);
                MediaReq=JSON.parse(data.req);
                isHaveNext=data.isHaveNext;
                isHavePrev=data.isHavePrev;
                updateTable();
                updateButton();
            },
        })
    }

}(this,document));