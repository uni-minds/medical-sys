$(function(){
    function uiCopyright() {
        console.log("Copyright @ Uni-Minds 2019-2020");
        $(".main-footer").html(
            "<strong>Copyright &copy; 2020 <a href=\"http://uni-minds.com\">Uni-Minds</a> / <a href=\"http://www.buaa.edu.cn\">Beihang University</a> / <a href=\"http://www.anzhen.org\">Beijing Anzhen Hospital, CMU</a>. </strong>" +
            "All rights reserved.\n" +
            "<div class=\"float-right d-none d-sm-inline-block\">\n" +
            "  <b>Rev.</b> 20200330 Remastered\n" +
            "</div>");
    }
    function menuInit() {
        menuSetRealname();
        menuCreate();

        function menuCreate(){
            $.get("/webapp/common/static/menu.json", function (result) {
                menuLoadData(result);
                if (navMenuActive) {
                    menuActive(navMenuActive)
                } else {
                    setTimeout(menuCreate,3000);
                }
            });
        }
        function menuLoadData(menuJson) {
            const root = $(".mt-2 .nav");
            for (var i = 0; i < menuJson.length; i++) {
                let menu = menuJson[i];
                if (menu.child) {
                    let parent = menuParent(menu.name, menu.icon);
                    let tree = $('<ul class="nav nav-treeview"></ul>');
                    parent.append(tree);
                    root.append(parent);

                    let childlen = menu.child.length;
                    for (var j = 0; j < childlen; j++) {
                        let child = menu.child[j];
                        let obj = menuChild(child.id, child.name, child.icon, child.controller);
                        tree.append(obj);
                    }
                } else {
                    let obj = menuChild(menu.id, menu.name, menu.icon, menu.controller);
                    root.append(obj);
                }
            }
        }
        function menuParent(name, icon) {
            return $('<li class="nav-item has-treeview"><a href="#" class="nav-link"><i class="nav-icon ' +
                icon + '"></i>' +
                '<p>' + name + '<i class="right fas fa-angle-left"></p></i></a></li>');
        }
        function menuChild(id, name, icon, controller) {
            if (id == null) {
                return $('<li class="nav-item"><a href="' +
                    controller + '" class="nav-link"><i class="nav-icon ' +
                    icon + '"></i><p>' +
                    name + '</p></a>');
            } else {
                return $('<li class="nav-item"><a id="' +
                    id + '" href="' +
                    controller + '" class="nav-link"><i class="nav-icon ' +
                    icon + '"></i><p>' +
                    name + '</p></a>');
            }
        }
        function menuSetRealname() {
            $.get("/api/user?action=getrealname",function(result){
                let obj = $('.sidebar .info .d-block');
                obj.attr('href','/ui/logout');
                if (result.code === 200){
                    obj.text(result.data);
                } else {
                    obj.text("未知用户");
                }
            });
        }
        function menuActive(id) {
            $(".nav-link").removeClass("active");
            let selector = "#" + id;
            let obj = $(selector);
            obj.addClass("active");
            let objParent = obj.parents()[2];
            if ($.nodeName(objParent, "li")) {
                $(objParent).addClass("menu-open");
            }
        }
    }
    menuInit();
    uiCopyright();
});