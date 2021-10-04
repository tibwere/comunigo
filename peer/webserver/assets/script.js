$("#signform").submit(function (e) { 
    $("#signBtn").empty()
    $("#signBtn").append('<div class="spinner-grow spinner-grow-sm text-light" role="status"></div> Waiting for other partecipants')
    $("#signBtn").prop("disabled", true);
});

$("#username").click(function (e) { 
    $("#errorDiv").hide()    
});