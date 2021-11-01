function infoSuccessHandler (response) {
    sessionStorage.setItem("currentUser", response.Username)
    sessionStorage.setItem("verbose", response.Verbose)

    document.title = response.Username + " - comuniGO"

    $("#me").text(response.Username)
    $("#tos").text(response.Tos)
    response.OtherMembers.forEach(m => 
        $("#memberList").append('<li class="list-group-item list-group-item-warning"><strong class="text-primary">' + 
            m.Username + '</strong> (addr: ' + m.Address + ')</li>')
    )  
    
    if (response.Verbose) {
        console.log("*** VERBOSE ENABLED ***")
        console.log("Current peer...................: " + response.Username)
        console.log("Current type of service offered: " + response.Tos)
        console.log("List of other peers connected..:")
        response.OtherMembers.forEach(m=>console.log("\t" + m.Username + "@" + m.Address))
    }
}

function sendSuccessHandler() {
    $("#message").val('')
    setTimeout(function(){
        $("#sendBtn").empty();
        $("#sendBtn").append('<span class="text-primary"><i class="fas fa-paper-plane"></i> Send Message</span>')
    }, 1500) 
}

function listSuccessHandler(response) {
    if (response.length > 0) {
        $("#emptyMessageListAlert").fadeOut(200)
        $("#messageList").empty()

        if (debug) {
            console.log("New message list retrieved from server:")
        }
        $.each(response, function(_, m) {
            if (debug) {
                if (Array.isArray(m.Timestamp)) {
                    console.log("\t[ID: [" + m.Timestamp + "], From: " + m.From + ", Body: " + m.Body + "]")
                } else if (m.Timestamp === undefined) {
                    console.log("\t[ID: 0, From: " + m.From + ", Body: " + m.Body + "]")
                } else {
                    console.log("\t[ID: " + m.Timestamp + ", From: " + m.From + ", Body: " + m.Body + "]")
                }
            }

            if (m.From == sessionStorage.getItem("currentUser")) {
                $("#messageList").append('<li class="list-group-item list-group-item-warning"><strong class="text-warning">(' + m.From + ')</strong> ' + m.Body + '</li>')                    
            } else {
                $("#messageList").append('<li class="list-group-item list-group-item-primary"><strong class="text-primary">(' + m.From + ')</strong> ' + m.Body + '</li>')
            }
        })
    }
}

function forbiddenHandler() {
    window.location = "login.html"
}


$( document ).ready(function() {
    $.ajax({
        type: "GET",
        url: "info",
        statusCode: {
            200: infoSuccessHandler,
            403: forbiddenHandler
        }
    });
});

$("#sendForm").submit(function (e) { 
    e.preventDefault(); 
    
    $("#sendBtn").empty();
    $("#sendBtn").append('<span class="text-success"><i class="fas fa-check"></i> Message succesfully sent</span>')

    $.ajax({
        type: "POST",
        url: "send",
        data: {
            "message" : $("#message").val()
        },
        statusCode: {
            200: sendSuccessHandler,
            403: forbiddenHandler
        }
    });
});

$("#reloadMessagesForm").submit(function (e) { 
    e.preventDefault();

    debug = sessionStorage.getItem("verbose")

    $.ajax({
        type: "GET",
        url: "list",
        statusCode: {
            200: listSuccessHandler,
            403: forbiddenHandler
        }
    });  
});    
    

