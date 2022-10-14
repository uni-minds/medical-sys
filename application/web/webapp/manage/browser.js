$("#search").click(()=> {
    $.get("/api/v1/blockchain/height?height=" + $("#SearchHeight").val(), (resp) => {
        if (resp.code !== 200) {
            console.log(resp)
        } else {
            console.log(resp.data)
            $("#block-content").css("word-break", "break-all").css("word-wrap", "break-word").text(JSON.stringify(resp.data.Block))
            $("#votes-content").css("word-break", "break-all").css("word-wrap", "break-word").text(JSON.stringify(resp.data.Votes))
        }
    })
})