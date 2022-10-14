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

function windowError(html,autoCounter) {
    Swal.fire({
        icon: 'error',
        title: "错误",
        html: html,
        showConfirmButton: true,
        timer: autoCounter,
        timerProgressBar: autoCounter
    });
}

function windowMessage(title,html,autoCounter,func) {
    Swal.fire({
        title: title,
        html: html,
        icon: 'success',
        showConfirmButton: true,
        timer: autoCounter,
        timerProgressBar: autoCounter,
    }).then(func)
}

function windowResult(title,html,autoCounter) {
    Swal.fire({
        title: title,
        html: html,
        showConfirmButton: true,
        timer: autoCounter,
        timerProgressBar: autoCounter
    })
}

function windowInput(title, text) {
    Swal.fire({
        title: title,
        input: 'textarea',
        inputValue: text,
        showCancelButton: true,
        confirmButtonText: '提交',
        cancelButtonText: '取消',
        showLoaderOnConfirm: true,
        preConfirm: (login) => {
            return fetch(`//api.github.com/users/${login}`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(response.statusText)
                    }
                    return response.json()
                })
                .catch(error => {
                    Swal.showValidationMessage(
                        `Request failed: ${error}`
                    )
                })
        },
        allowOutsideClick: () => !Swal.isLoading()
    }).then((result) => {
        console.log(result)
        windowMessage("resp", result.data, 2000)
        // if (result.value) {
        //     Swal.fire({
        //         title: `${result.value.login}'s avatar`,
        //         imageUrl: result.value.avatar_url
        //     })
        // }
    })
}