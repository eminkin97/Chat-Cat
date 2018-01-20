// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.
const {ipcRenderer} = require('electron')
var $ = require('jquery')

var online_users = []
var chattingWithCurrently = ""

//add a chat channel receiver
ipcRenderer.on('chat-reply', (event, arg) => {
	//parse arg. Before first colon is username. After is message
	var msg = arg.msg.split(":", 2)
	console.log(msg)

	if (chattingWithCurrently == msg[0]) {
		$("#messages_display_box").append("<span class='other_message'>" + msg[1] + "</span>")
	}
})

//add a list channel receiver
ipcRenderer.on('list-reply', (event, arg) => {
	console.log(arg)
	console.log(arg.msg)

	user_list = arg.msg.split(',')

	user_list.forEach(function(user) {
		if (user) {
			var index = online_users.indexOf(user)

			//check to see if user is already in online contacts list
			if (index == -1) {
				//need to add user to list and append entry to online contacts div
				online_users.push(user)
				$("#contacts").append("<div id='" + user + "' class='contact_entry' onclick='contact_clicked(this)'><p>" + user + "</p></div>");
			}
		}
	});

	//removing contacts that have gone offline from online contact list
	i = 0
	while (i < online_users.length) {
		var user = online_users[i]
		var index = user_list.indexOf(user)

		if (index == -1) {
			online_users.splice(i, 1)
			$("#" + user).remove()	//remove div
		} else {
			i++;
		}
	}
})

function contact_clicked(obj) {
	console.log("clicked")
	chattingWithCurrently = $(obj).find("p").text()
	$(obj).css("background-color", "yellow")	//set the background color to yellow
	$("#messages_display_box_title").text(chattingWithCurrently)
}

$("#text_input").keyup(function(e) {
	var code = e.which

	if (code == 13) {
		//enter key was pressed
		if (chattingWithCurrently) {
			var user_text = $(this).val()
			var textToSend = "chat:"+chattingWithCurrently+":"+user_text
			ipcRenderer.send('chat', textToSend)
			$(this).val('')	//clear textarea

			//append text to messages display text
			$("#messages_display_box").append("<span class='my_message'>" + user_text + "</span>")
		}
	}
});

$(document).ready(function() {
	//send request for list
	ipcRenderer.send('list')

	//update the list every 5 seconds
	setInterval(function() {
		ipcRenderer.send('list')
	}, 5000);
	
});
