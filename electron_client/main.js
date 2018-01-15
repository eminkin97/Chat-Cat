const electron = require('electron')
// Module to control application life.
const app = electron.app
// Module to create native browser window.
const BrowserWindow = electron.BrowserWindow
const ipcMain = electron.ipcMain

const path = require('path')
const url = require('url')
const WebSocket = require('ws')

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow
let promptWindow

//socket to chat server
const socket = new WebSocket("127.0.0.1:8001")

function createPrompt() {
  promptWindow = new BrowserWindow({width: 200, height: 150})

  // and load the index.html of the prompt.
  promptWindow.loadURL(url.format({
    pathname: path.join(__dirname, 'prompt/prompt.html'),
    protocol: 'file:',
    slashes: true
  }))

  // Open the DevTools.
  promptWindow.webContents.openDevTools()

  // Emitted when the window is closed.
  promptWindow.on('closed', function () {
    promptWindow = null
  })
}

function createWindow () {
  // Create the browser window.
  mainWindow = new BrowserWindow({width: 800, height: 600})

  // and load the index.html of the app.
  mainWindow.loadURL(url.format({
    pathname: path.join(__dirname, 'app/index.html'),
    protocol: 'file:',
    slashes: true
  }))

  // Open the DevTools.
  mainWindow.webContents.openDevTools()

  // Emitted when the window is closed.
  mainWindow.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null
  })
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', createPrompt)

// Quit when all windows are closed.
app.on('window-all-closed', function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

/*app.on('activate', function () {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow()
  }
})
*/


// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.

//get list of available clients from chat server
ipcMain.on('list', (event) => {
	ws.send('list')
});

ipcMain.on('name', (event, arg) => {
	ws.send(arg)
	createWindow()
});

socket.on('message', function incoming(data) {
	//parse message received and send to renderer process correctly
	let code = data.substring(0, 1)

	if (code == "l") {
		mainWindow.webContents.send('list-reply' , {msg: data.substring(2)});	
	} else if (code == "c") {
		mainWindow.webContents.send('chat-reply' , {msg: data.substring(2)});	
	}
});
