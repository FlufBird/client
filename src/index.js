const parseJSON = require("json-templates");
const { app, BrowserWindow, dialog, } = require("electron");
const path = require("path");
const { existsSync, readFileSync, } = require("fs");
const axios = require("axios").default;

if (require("electron-squirrel-startup")) {
	app.quit();
}

const developmentMode = existsSync("../development");

const applicationData = JSON.parse(readFileSync("resources/data/application.json"));
const userData = JSON.parse(readFileSync("resources/data/user.json"));

const server = ((developmentMode) ? "http://localhost:5000" : "https://mozuli.deta.dev") + "/api/v1";

const language = JSON.parse(readFileSync(`resources/languages/${userData.language}.json`));

const createWindow = () => {
	const window = new BrowserWindow({
		show: false,
		width: 800,
		height: 600,
		frame: false,
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
			noLink: true,
			defaultId: 0,
			cancelId: 0,
			title: "Mozuli",
			message: language.updater.checkingError,
		});

		if (response === 0) {
			app.quit();
		}
	};

	createWindow();

	try {
		const currentVersion = applicationData.version;

		axios.get(`${server}/latest_version`)
			.then((response) => {
				const data = response.data;
				const latestVersion = data.latestVersion;

				if (data.maintenance) {
					dialog.showErrorBox(
						"Mozuli",
						language.general.serverOnMaintenance,
					);

					app.quit();
				}

				if (!("latestVersion" in data)) {
					error();
				}

				if (currentVersion !== latestVersion) {
					const response = dialog.showMessageBoxSync(null, {
						type: "question",
						buttons: [
							"No",
							"Yes",
						],
						noLink: true,
						defaultId: 0,
						cancelId: 0,
						title: "Mozuli",
						message: parseJSON(language.updater.updateAvailable)({
							currentVersion: currentVersion,
							latestVersion: latestVersion,
						}),
					});

					if (response === 0) {
						app.quit();
					}
				}
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