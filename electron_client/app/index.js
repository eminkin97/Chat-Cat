// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.
const {ipcRenderer} = require('electron')
var $ = require('jquery')
var localforage = require('localforage')

var online_users = []
var chattingWithCurrently = ""
var message_numb = 1		//id for each chat message. Lower ids are older messages
var store		//localforage instance

//add a chat channel receiver
ipcRenderer.on('chat-reply', (event, arg) => {
	//parse arg. Before first colon is username. After is message
	var msg = arg.msg.split(":", 2)
	console.log(msg)

	if (chattingWithCurrently == msg[0]) {
		$("#messages_display_box").append("<span class='other_message'>" + msg[1] + "</span>")
	} else {
		console.log("Updating message count circle")

		//show new_message_count circle
		$('#'+msg[0]).find('.new_message_count').css("visibility", "visible")
		$('#'+msg[0]).find('.new_message_count').html(parseInt($('#'+msg[0]).find('.new_message_count').html(), 10)+1)
	}

	//create object for message
	var msg_info = {
		name: msg[0],
		senderOrReceiver: 1,	//0 indicates user sent. 1 indicates user received
		data: msg[1]
	}

	//insert into localstorage
	store.setItem(message_numb.toString(), msg_info).catch(function(err) {
		console.log(err)
	});

	message_numb = message_numb + 1
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
				$("#contacts").append("<div id='" + user + "' class='contact_entry' onclick='contact_clicked(this)'><p>" + user + "</p><span class='new_message_count'>0</span></div>");
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
	$(".contact_entry").css("background-color", "white")
	$(obj).css("background-color", "yellow")	//set the background color to yellow
	$('#messages_display_box').html('');	//clear content of div
	$("#messages_display_box_title").css("visibility", "visible")
	$("#messages_display_box_title").text(chattingWithCurrently)

	//clear the new message count box
	$(obj).find('.new_message_count').html(0)
	$(obj).find('.new_message_count').css("visibility", "hidden")



	//iterate through localstorage and pull up previous messages
	store.iterate(function(value, key, iterationNumber) {
		if (value.name == chattingWithCurrently) {
			if (value.senderOrReceiver == 0) {
				//user sent
				$("#messages_display_box").append("<span class='my_message'>" + value.data + "</span>")
			} else if (value.senderOrReceiver == 1) {
				//user received
				$("#messages_display_box").append("<span class='other_message'>" + value.data + "</span>")
			}
		}
	}).catch(function(err) {
		console.log(err)
	});
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

			//create object for message
			var msg_info = {
				name: chattingWithCurrently,
				senderOrReceiver: 0,	//0 indicates user sent. 1 indicates user received
				data: user_text
			}

			//insert into localstorage
			store.setItem(message_numb.toString(), msg_info).catch(function(err) {
				console.log(err)
			});

			message_numb = message_numb + 1

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

	//clear local storage at beginning
	localforage.clear().then(function() {
		// Run this code once the database has been entirely deleted.
		console.log('Database is now empty.');
	}).catch(function(err) {
		// This code runs if there were any errors
		console.log(err);
	});

	//create new local storage instance
	var name = ipcRenderer.sendSync('get_name')
	store = localforage.createInstance({
		name: name
	});
});
