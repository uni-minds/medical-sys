function OpSearchHis(patientId) {
    let u = `/api/v1/his/${patientId}`
    $.get(u, result => {
        console.log(result)
        if (result.code === 200) {
            let disp = ""
            if (result.data.length === 0) {
                windowError(`未检索到关于病例号< ${patientId} >的诊断信息。`,2000)
            } else {
                windowResult("查询结果", hisResultAnalysis(result.data), 20000)
            }
        } else {
            windowError(result.msg)
        }
    })
}

function OpSeriesMemo(studiesId,seriesId) {
    let u = `/api/v1/screen/studies/${studiesId}/series/${seriesId}/memo`
    $.get(u, (resp) => {
        if (resp.code === 200) {
            let memoText = resp.data
            Swal.fire({
                title: "备注",
                input: 'textarea',
                inputValue: memoText,
                showCancelButton: true,
                confirmButtonText: '提交',
                cancelButtonText: '取消',
                showLoaderOnConfirm: true,
                preConfirm: (memoNew) => {
                    let data = {}
                    data["value"] = memoNew

                    return fetch(u, {
                        body: JSON.stringify(data), // must match 'Content-Type' header
                        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                        credentials: 'same-origin', // include, same-origin, *omit
                        headers: {
                            'user-agent': 'Mozilla/4.0 MDN Example',
                            'content-type': 'application/json'
                        },
                        method: 'POST', // *GET, POST, PUT, DELETE, etc.
                        mode: 'cors', // no-cors, cors, *same-origin
                        redirect: 'follow', // manual, *follow, error
                        referrer: 'no-referrer', // *client, no-referrer
                    })
                        .then(response => {
                            if (!response.ok) {
                                throw new Error(response.statusText)
                            }
                            return response.json()
                        }) // parses response to JSON
                        .catch(error => {
                            Swal.showValidationMessage(
                                `反馈错误: ${error}`
                            )
                        })
                },
                allowOutsideClick: () => !Swal.isLoading()
            }).then((result) => {
                let data = result.value
                if (data.code === 200) {
                    windowMessage("", "提交完成", 1000)
                }
            })
        } else {
            alert("memo error")
        }
    })
}

function OpReviewApprove(studiesId,seriesId) {
    console.log("approve")
    let u = `/api/v1/screen/studies/${studiesId}/series/${seriesId}/submit?action=review_approve`
    $.post(u, (resp) => {
        if (resp.code === 200) {
            alert("完成提交")
        } else {
            alert("提交失败，请重试")
        }
    })
}

function OpReviewReject(studiesId,seriesId) {
    console.log("reject")
    let u = `/api/v1/screen/studies/${studiesId}/series/${seriesId}/submit?action=review_reject`
    $.post(u, (resp) => {
        if (resp.code === 200) {
            alert("完成提交")
        } else {
            alert("提交失败，请重试")
        }
    })
}

function OpAuthorSubmit(studiesId,seriesId) {
    console.log("author_submit")
    let u = `/api/v1/screen/studies/${studiesId}/series/${seriesId}/submit?action=author`
    $.post(u, (resp) => {
        if (resp.code === 200) {
            alert("完成提交")
        } else {
            alert("提交失败，请重试")
        }
    })
}


function hisResultAnalysis(data) {
    let obj = $('<table class="table table-bordered">')
    data.forEach((ele, i) => {
        let objTr = $('<tr>').attr("align", "center")
        objTr.append($('<td colspan="2" class="bg-info">').text(`== 关联诊断 ${i + 1} ==`))
        obj.append(objTr)
        for (const eleKey in ele) {
            if (eleKey === "md5") {
                continue
            }

            let value = ele[eleKey]
            switch (value.toLowerCase()) {
                case "[]":
                case "":
                case "na":
                    continue
                default:
                    objTr = $('<tr>')
                    objTr.append($('<td>').text(`${eleKey}`))
                    objTr.append($('<td>').text(`${value}`))
                    obj.append(objTr)
            }
        }
    })
    return obj
}



function postData(url, data) {
    return fetch(url, {
        body: JSON.stringify(data), // must match 'Content-Type' header
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, same-origin, *omit
        headers: {
            'user-agent': 'Mozilla/4.0 MDN Example',
            'content-type': 'application/json'
        },
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, cors, *same-origin
        redirect: 'follow', // manual, *follow, error
        referrer: 'no-referrer', // *client, no-referrer
    })
        .then(response => response.json()) // parses response to JSON
}