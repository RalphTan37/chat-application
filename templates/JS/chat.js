$(function(){

    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $("#messages");
  
    $("#chatbox").submit(function(){
      if (!msgBox.val()) return false;
      if (!socket) {
        alert("Error: No socket connection.");
        return false;
      }
      socket.send(msgBox.val());
      msgBox.val("");
      return false;
    });

    if (!window["WebSocket"]) {
      alert("Error: Browser does not support web sockets.")
    } else {
      socket = new WebSocket("ws://{{.Host}}/room");
      socket.onclose = function() {
        alert("Connection has been closed.");
      }
      socket.onmessage = function(e) {
        messages.append($("<li>").text(e.data));
      }
    }
  });