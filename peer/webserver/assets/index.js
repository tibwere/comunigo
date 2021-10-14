$( document ).ready(function() {
    metadata = $.parseJSON(sessionStorage.getItem("comunigo-metadata"))
    sessionStorage.setItem("currentUser", metadata.Username)
    document.title = metadata.Username + " - comuniGO"

    $("#me").text(metadata.Username)
    $("#tos").text(metadata.Tos)
    metadata.Members.forEach(m => $("#memberList").append('<li class="list-group-item list-group-item-warning"><strong class="text-primary">' + m + '</strong></li>'))
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
        success: function (response) {

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

    var messages = []

    $.ajax({
        type: "POST",
        url: "list",
        success: function (response) {
            arrayOfJSONMessages = $.parseJSON(response)

            if (arrayOfJSONMessages != null) {

                $.each(arrayOfJSONMessages, function(_, elem) {
                    messages.push($.parseJSON(elem))
                })

                $("#emptyMessageListAlert").fadeOut(200)
                $("#messageList").empty()

                $.each(messages, function(i, m) {
                    if (Array.isArray(m.Timestamp)) {
                        console.log("New! [ID: [" + m.Timestamp + "], From: " + m.From + ", Body: " + m.Body + "]")
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
    

