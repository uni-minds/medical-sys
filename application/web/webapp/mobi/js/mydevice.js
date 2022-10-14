function main() {
    let algolist = []

    $.get("/api/v1/algo").done((resp) => {
        algolist = resp.data
        algolist.forEach((val,i) => {
            $("#s-algo").append($("<option>").val(i).text(val["algo-name"]))
        })
    })

    $("#s-exec").click(function () {
        let mask = $('<div id="main-mask" class="overlay dark" ><i class="fas fa-3x fa-sync-alt fa-spin"></i></div>')
        $("#main-body").prepend(mask)

        let algo = algolist[$("#s-algo").val()]
        let data = {}
        data["dev"] = $("#s-dev").val()
        data["algo-ref"] = algo["algo-ref"]
        $.post("/mobi/exec", JSON.stringify(data), "json").done((resp) => {
            $("#main-mask").remove()
            window.location.href = `/mobi/result/${resp.data}`
        }).fail(() => {
            alert("设备无法连接")
            $("#main-mask").remove()
        })
    })
}

main()