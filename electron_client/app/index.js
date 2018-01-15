// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.
const {ipcRenderer} = require('electron')
var $ = require('jquery')

//add a chat channel receiver
ipcRenderer.on('chat-reply', (event, arg) => {
	console.log(arg)
})

//add a list channel receiver
ipcRenderer.on('list-reply', (event, arg) => {
	console.log(arg)
	//replace online clients with list received
})

$(document).ready(function() {
	//send request for list
	ipcRenderer.send('list')
	
});
