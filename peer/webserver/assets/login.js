function signSuccessHandler(response) {
    if (response.Status == "ERROR") {
        $("#errorMsg").html(response.Message)
        $("#errorMsg").fadeIn()

        $("#signBtn").empty()
        $("#signBtn").append('<span class="text-primary"><i class="fas fa-sign-in-alt"></i> Login</span>')
        $("#signBtn").prop("disabled", false);

        setTimeout(function(){
            $("#errorMsg").fadeOut()
        }, 1500)
    } else {
        window.location = "index.html"
    }
}

function badRequestHandler(response) {
    $("#errorMsg").html("You are already logged in")
    $("#errorMsg").fadeIn()

    setTimeout(function(){
        window.location = "index.html"
    }, 1500)
}

$("#signform").submit(function (e) { 

    e.preventDefault();

    $("#signBtn").empty()
    $("#signBtn").append('<div class="spinner-grow spinner-grow-sm text-light" role="status"></div> Waiting for other partecipants')
    $("#signBtn").prop("disabled", true);

    $.ajax({
        type: "POST",
        url: "sign",
        data: {
            "username" : $("#username").val()
        },
        statusCode: {
            200: signSuccessHandler,
            400: badRequestHandler
        }
    });
});

$("#username").click(function (e) { 
    $("#errorDiv").hide()    
});