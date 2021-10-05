$( document ).ready(function() {
    metadata = $.parseJSON(sessionStorage.getItem("comunigo-metadata"))
    sessionStorage.setItem("currentUser", metadata.Username)
    console.log(metadata)
    document.title = metadata.Username + " - comuniGO"

    $("#me").text(metadata.Username)
    metadata.Members.forEach(m => $("#memberList").append('<li class="list-group-item">' + m + '</li>'))
    sessionStorage.setItem("index", 0)
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
            $("#alertMsg").fadeIn()
            setTimeout(function(){
                $("#alertMsg").fadeOut()
            }, 1500)            
        }
    });
});

setInterval(function(){
    var index = sessionStorage.getItem("index")

    $.ajax({
        type: "POST",
        url: "list",
        data: {
            "next" : index
        },
        success: function (response) {
            console.log(response)
            messages = $.parseJSON(response).MessageList
            if (messages != null) {
                sessionStorage.setItem("index", index + messages.length)
                $.each(messages, function(i, m) {
                    if (m.From == sessionStorage.getItem("currentUser")) {
                        $("#messageList").append('<li class="list-group-item list-group-item-success"><strong class="text-success">(' + m.From + ')</strong> ' + m.Body + '</li>')                    
                    } else {
                        $("#messageList").append('<li class="list-group-item list-group-item-secondary"><strong class="text-secondary">(' + m.From + ')</strong> ' + m.Body + '</li>')
                    }
                })                      
            }
        }
    });
}, 5000)