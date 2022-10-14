jsGrid.locale("zh-cn");

$(function() {
    let fields = [
        {name: "algo-index", type: "number", title: "序号",align:"center",width:20},
        {name: "algo-name", type: "string", title: "名称",align:"center"},
        {name: "algo-ref", type: "string", title: "引用地址",align:"center"},
        {
            type: "control",
            modeSwitchButton: false,
            editButton: false,
            deleteButton: false,
            headerTemplate: function () {
                return $("<button>").attr("type", "button").text("添加").on("click", function () {
                    $("#user-editor-dialog").dialog("open");
                })
            }
        }
    ];

    $("#algo-table").jsGrid({
        height: "auto",
        width: "100%",
        fields: fields,
        sorting: true,
        paging: true,
        autoload: true,
        controller: {
            loadData: function () {
                let d = $.Deferred();
                $.ajax({
                    url: "/api/v1/algo?action=getlist",
                    dataType: "json",
                    type: "GET",
                }).done(function (response) {
                    console.log("S",response)
                    d.resolve(response.data);
                }).fail(function (resp) {
                    console.log(resp)
                    d.resolve({})
                });
                return d.promise();
            },
        },
    });

    $("#user-editor-dialog").dialog({
        autoOpen: false,
        width: 300,
        modal: true,
        title: "登记算法",
        close: function () {
            location.reload()
        }
    });

    $("#user-editor-form").validate({
        rules: {
            "algo-name": "required",
            "algo-ref": "required",
        },
        submitHandler: function () {
            let data = {}
            data["algo-name"] = $("#algo-name").val()
            data["algo-ref"] = $("#algo-ref").val()
            $.post("/api/v1/algo",JSON.stringify(data),"json").done((resp)=>{
                console.log("S",resp)
            }).fail((resp)=>{
                console.log("F",resp)
            })
            $("#user-editor-dialog").dialog("close");
        }
    });
});