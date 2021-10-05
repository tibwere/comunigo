$( document ).ready(function() {
    metadata = $.parseJSON(sessionStorage.getItem("comunigo-metadata"))
    sessionStorage.setItem("currentUser", metadata.Username)
    document.title = metadata.Username + " - comuniGO"

    $("#me").text(metadata.Username)
    metadata.Members.forEach(m => $("#memberList").append('<li class="list-group-item">' + m + '</li>'))
    sessionStorage.setItem("index", -1)
});

$("#sendForm").submit(function (e) { 
    e.preventDefault();  

    $.ajax({
        type: "POST",
        url: "send",
        data: {
            "message" : $("#message").val()
        },
        success: function (response) {

            $("#message").text('')

            $("#alertMsg").fadeIn()
            setTimeout(function(){
                $("#alertMsg").fadeOut()
            }, 1500)            
        }
    });
});

setInterval(function(){
    var index = sessionStorage.getItem("index")
    var messages = []

    $.ajax({
        type: "POST",
        url: "list",
        data: {
            "next" : parseInt(index) + 1
        },
        success: function (response) {
            arrayOfJSONMessages = $.parseJSON(response)

            if (arrayOfJSONMessages != null) {

                $.each(arrayOfJSONMessages, function(_, elem) {
                    messages.push($.parseJSON(elem))
                })
                highestReceivedID = messages.at(-1).ID

                if (index < highestReceivedID) {

                    if (index == -1) {
                        $("#emptyMessageListAlert").fadeOut(200)
                    }

                    sessionStorage.setItem("index", highestReceivedID)
                    $.each(messages, function(i, m) {
                        console.log("New! [ID: " + m.ID + ", From: " + m.From + ", Body: " + m.Body + "]")
                        if (m.From == sessionStorage.getItem("currentUser")) {
                            $("#messageList").append('<li class="list-group-item list-group-item-warning"><strong class="text-warning">(' + m.From + ')</strong> ' + m.Body + '</li>')                    
                        } else {
                            $("#messageList").append('<li class="list-group-item list-group-item-dark"><strong class="text-secondary">(' + m.From + ')</strong> ' + m.Body + '</li>')
                        }
                    })
                }                      
            }
        }
    });
}, 2000)