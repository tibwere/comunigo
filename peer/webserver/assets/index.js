$( document ).ready(function() {
    
    $.ajax({
        type: "GET",
        url: "info",
        success: function (response) {
            sessionStorage.setItem("currentUser", response.Username)
            document.title = response.Username + " - comuniGO"
        
            $("#me").text(response.Username)
            $("#tos").text(response.Tos)
            response.OtherMembers.forEach(m => 
                $("#memberList").append('<li class="list-group-item list-group-item-warning"><strong class="text-primary">' + 
                    m.Username + '</strong> (addr: ' + m.Address + ')</li>')
            )          
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
        success: function () {
            $("#message").val('')
            setTimeout(function(){
                $("#sendBtn").empty();
                $("#sendBtn").append('<span class="text-primary"><i class="fas fa-paper-plane"></i> Send Message</span>')
            }, 1500)            
        }
    });
});

$("#reloadMessagesForm").submit(function (e) { 
    e.preventDefault();

    $.ajax({
        type: "GET",
        url: "list",
        success: function (response) {
            if (response != null) {
                $("#emptyMessageListAlert").fadeOut(200)
                $("#messageList").empty()

                $.each(response, function(_, m) {
                    if (Array.isArray(m.Timestamp)) {
                        console.log("New! [ID: [" + m.Timestamp + "], From: " + m.From + ", Body: " + m.Body + "]")
                    } else if (m.Timestamp === undefined) {
                        console.log("New! [ID: 0, From: " + m.From + ", Body: " + m.Body + "]")
                    } else {
                        console.log("New! [ID: " + m.Timestamp + ", From: " + m.From + ", Body: " + m.Body + "]")
                    }

                    if (m.From == sessionStorage.getItem("currentUser")) {
                        $("#messageList").append('<li class="list-group-item list-group-item-warning"><strong class="text-warning">(' + m.From + ')</strong> ' + m.Body + '</li>')                    
                    } else {
                        $("#messageList").append('<li class="list-group-item list-group-item-primary"><strong class="text-primary">(' + m.From + ')</strong> ' + m.Body + '</li>')
                    }
                })
            }                      
        }
    });  
});    
    

