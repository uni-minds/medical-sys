$(function(){
    console.log("Copyright @ Uni-Minds 2019-2020");
    function uiCopyright() {
        $.get("/api/v1/raw?action=getversion", result => {
            $(".main-footer").html(result);
        })
    }

    class Menu {
        constructor(props) {
            this.root = $(".mt-2 .nav");

        }

        init() {
            $.get("/api/v1/user?action=getrealname", result => {
                if (result.code === 200) {
                    this.username = result.data
                } else {
                    this.username = null;
                }
            });
        }

        loadData(menuJson) {
            for (let i = 0; i < menuJson.length; i++) {
                let menu = menuJson[i];
                if (menu.child && menu.child.length>0) {
                    let parent = this.parent(menu.name, menu.icon);
                    let tree = $('<ul class="nav nav-treeview"></ul>');
                    parent.append(tree);
                    this.root.append(parent);

                    let childlen = menu.child.length;
                    for (let j = 0; j < childlen; j++) {
                        let child = menu.child[j];
                        let obj = this.child(child.id, child.name, child.icon, child.controller);
                        tree.append(obj);
                    }
                } else {
                    let obj = this.child(menu.id, menu.name, menu.icon, menu.controller);
                    this.root.append(obj);
                }
            }
        }

        parent(name, icon) {
            return $('<li class="nav-item has-treeview"><a href="#" class="nav-link"><i class="nav-icon ' +
                icon + '"></i>' +
                '<p>' + name + '<i class="right fas fa-angle-left"></p></i></a></li>');
        }

        child(id, name, icon, controller) {
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

        set username(u) {
            let obj = $('.sidebar .info .d-block');
            if (u) {
                obj.attr('href', '/ui/logout').text(u);
            } else {
                obj.attr("href", "#").text("????????????")
            }
        }

        set menudata(d) {
            this.loadData(d)
        }

        set active(id) {
            if (this.activeRef) {
                this.activeRef.removeClass("active")
            }
            let obj = $(`#${id}`);
            obj.addClass("active");
            let objParent = obj.parents()[2];
            if ($.nodeName(objParent, "li")) {
                $(objParent).addClass("menu-open");
            }
            this.activeRef = obj
        }
    }

    let menu =new Menu()
    menu.init()
    uiCopyright();
});