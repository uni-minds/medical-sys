function main() {
    $("#s-exec").click(function () {
        let mask = $('<div id="main-mask" class="overlay dark" ><i class="fas fa-3x fa-sync-alt fa-spin"></i></div>')
        $("#main-body").prepend(mask)

        let dev = $("#s-dev").val()
        let algo = $("#s-algo").val()
        let url = `/mobi/exec?dev=${dev}&algoid=${algo}`
        $.get(url, function (resp) {
            if (resp.code !== 200) {
                alert(`异常反馈：${resp}`)
            } else {
                $("#main-mask").remove()
                window.location.href = `/mobi/result/${resp.data}`
            }
        })
    })
}

main()