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
        success: function (response) {
            console.log(response)
            result = $.parseJSON(response)
            if (result.Status == "ERROR") {
                $("#errorMsg").html(result.Message)
                $("#errorMsg").fadeIn()

                $("#signBtn").empty()
                $("#signBtn").append('<i class="fas fa-sign-in-alt"></i> Login')
                $("#signBtn").prop("disabled", false);

                setTimeout(function(){
                    $("#errorMsg").fadeOut()
                }, 1500)
            } else {
                window.location = "index.html"
            }            
        }
    });
});

$("#username").click(function (e) { 
    $("#errorDiv").hide()    
});