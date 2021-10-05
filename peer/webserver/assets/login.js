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
            
            decoded = $.parseJSON(response)
            if (decoded.IsError == true) {
                $("#errorMsg").html(decoded.ErrorMessage)
                $("#errorMsg").fadeIn()
                setTimeout(function(){
                    $("#errorMsg").fadeOut()
                }, 1500)
            } else {
                sessionStorage.setItem("comunigo-metadata", response)
                window.location = "index.html"
            }            
        }
    });
});

$("#username").click(function (e) { 
    $("#errorDiv").hide()    
});