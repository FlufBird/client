const { app, BrowserWindow, dialog, } = require("electron");
const path = require("path");
const { readFileSync, } = require("fs");
const axios = require("axios").default;

let window;

if (require("electron-squirrel-startup")) {
	app.quit();
}

const createWindow = () => {
	window = new BrowserWindow({
		show: false,
		width: 800,
		height: 600,
		autoHideMenuBar: true,
		webPreferences: {
			nodeIntegration: true,
		},
	});
	window.loadFile(path.join(__dirname, "index.html"));

	window.once("ready-to-show", () => {
		window.show();
	});
};

app.once("ready", () => {
	const error = () => {
		const response = dialog.showMessageBoxSync(null, {
			type: "warning",
			buttons: [
				"No",
				"Yes",
			],
			defaultId: 0,
			cancelId: 0,
			noLink: true,
			title: "Mozuli",
			message: "Couldn't check for updates, continue?",
		});

		if (response === 0) {
			app.quit();
		}
	};

	createWindow();

	try {
		const applicationData = JSON.parse(readFileSync("resources/data/application.json"));
		const server = (applicationData.development) ? applicationData.servers.development : applicationData.servers.production;

		axios.get(`${server}/latest_version`)
			.then((response) => {
				console.log(response);
			})
			.catch((_) => {
				error();
			})
			.then(() => {});
	} catch(_) {
		error();
	}
});

app.on("window-all-closed", () => {
	if (process.platform !== "darwin") {
		app.quit();
	}
});

app.on("activate", () => {
	if (BrowserWindow.getAllWindows().length === 0) {
		createWindow();
	}
});

app.once("closed", () => {
	window = null;

	app.quit();
});