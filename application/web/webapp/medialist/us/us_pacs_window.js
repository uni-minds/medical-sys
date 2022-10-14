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

function windowMessage(title,html,autoCounter) {
    Swal.fire({
        title: title,
        html: html,
        icon: 'success',
        showConfirmButton: true,
        timer: autoCounter,
        timerProgressBar: autoCounter
    })
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