jsGrid.locale("zh-cn");

$(function() {
    let fields = [
        {name: "uid", type: "number", title: "UID",align:"center", width: 20},
        {name: "loginenable", type: "checkbox", title: "激活",align:"center", width: 22},
        {name: "username", type: "string", title: "账号",align:"center", width: 40},
        {name: "realname", type: "string", title: "姓名",align:"center", width: 40},
        {name: "email", type: "string", title: "邮箱",align:"center", width: 80},
        {name: "groups", type: "string", title: "分组",align:"center"},
        {name: "logincount", type: "number", title: "登录次数",align:"center", width: 35},
        {name: "logintime", type: "string", title: "上次登录",align:"center", width: 70},
        {name: "remark", type: "string", title: "备注"},
        {
            type: "control",
            modeSwitchButton: false,
            editButton: false,
            headerTemplate: function () {
                return $("<button>").attr("type", "button").text("添加").on("click", function () {
                    showDetailsDialog("Add", {})
                })
            }
        }
    ];

    $("#user-table").jsGrid({
        height: "auto",
        width: "100%",

        fields: fields,

        sorting: true,
        paging: true,
        autoload: true,


        //noDataContent:"等待数据加载……",

        controller: {
            loadData: function () {
                let d = $.Deferred();
                $.ajax({
                    url: "/api/user?action=getlist",
                    dataType: "json",
                    type: "GET",
                }).done(function (response) {
                    d.resolve(JSON.parse(response.data));
                });
                return d.promise();
            },
        },

        deleteConfirm: function (item) {
            return '确认要删除用户"' + item.username + '"么？'
        },

        rowClick: function (args) {
            showDetailsDialog("Edit", args.item)
        },


    });

    $("#user-editor-dialog").dialog({
        autoOpen: false,
        width: 600,
        close: function () {
            console.log("Close")
        }
    });

    $("#user-editor-form").validate({
        rules: {
            username: "required",
            realname: "required",
        },
        submitHandler: function () {
            formSubmitHandler();
        }
    });

    var formSubmitHandler = $.noop;

    var showDetailsDialog = function (dialogType, client) {
        $("#username").val(client.username);
        $("#realname").val(client.realname);
        $("#email").val(client.email);

        formSubmitHandler = function () {
            saveClient(client, dialogType === "Add");
        };

        $("#user-editor-dialog").dialog("option", "title", dialogType + " Client").dialog("open");
    };

    var saveClient = function (client, isNew) {
        $.extend(client, {
            username: $("#username").val(),
            realname: $("#realname").val(),
        });

        $("#user-table").jsGrid(isNew ? "insertItem" : "updateItem", client);
        $("#user-editor-dialog").dialog("close");
    }
});




    //
    // $("#user-table").jsGrid({
    //     height: "100%",
    //     width: "100%",
    //
    //     filtering: true,
    //     sorting: true,
    //     paging: true,
    //     //
    //     // pageSize: 15,
    //     // pageButtonCount: 5,
    //
    //     // rowClick: function (item) {
    //     //     console.log(item)
    //     // },
    //
    //     // loadData: function (filter) {
    //     //     console.log("LoadData");
    //     //     return $.ajax({
    //     //         type: "GET",
    //     //         url: "/api/v1/user?action=getlist",
    //     //         data: data
    //     //     });
    //     // },
    //
    //     // updateItem: function (item) {
    //     //     console.log("Update");
    //     //     return $.ajax({
    //     //         type: "PUT",
    //     //         url: "/items",
    //     //         data: item
    //     //     });
    //     // },
    //     //
    //     // deleteItem: function (item) {
    //     //     console.log("Delete");
    //     //     return $.ajax({
    //     //         type: "DELETE",
    //     //         url: "/items",
    //     //         data: item
    //     //     });
    //     // },
    //
    //     data: db.clients,
    //
    //     fields: db.fields,
    // });
