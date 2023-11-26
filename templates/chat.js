$(function () {
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $("#messages");
    var messageSentFlag = false;

    $("#chatbox").submit(function () {
        if (!msgBox.val()) return false;
        if (!socket) {
            alert("Error: There is no socket connection.");
            return false;
        }

        // User's message
        var userMessage = msgBox.val();
        messages.append($("<li class='user-message'>").text(userMessage));

        // Set the message sent flag
        messageSentFlag = true;

        socket.send(userMessage);

        msgBox.val("");
        return false;
    });

    if (!window["WebSocket"]) {
        alert("Error: Your browser does not support web sockets.")
    } else {
        socket = new WebSocket("ws://{{.Host}}/room");
        socket.onclose = function () {
            alert("Connection has been closed.");
        }
        socket.onmessage = function (e) {
            // Other person's message
            var otherPersonMessage = e.data;

            // Check the message sent flag
            if (!messageSentFlag) {
                messages.append($("<li class='other-person-message'>").text(otherPersonMessage));
            }

            // Set the message sent flag to false
            messageSentFlag = false;
        }
    }
});
