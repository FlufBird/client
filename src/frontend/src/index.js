import {
    GetLanguageData_,
} from "../wailsjs/go/main/Application"; // this path is only for the ide features, itll get replaced in the assets server handler

const $ = (query, _array) => {
	let results = document.querySelectorAll(query);

    if (!_array) return results[0];

    return results;
};

const setLanguageData = (element, key) => GetLanguageData_(key).then((value) => element.innerText = value);

window.runtime.EventsOnce("domReady", () => $(".container", false).style.display = "block");

const languageAttribute = "data-language";
const languageElements = $(`[${languageAttribute}]`, true);

for (let index = 0; index < languageElements.length; index++) {
    const element = languageElements[index];

    setLanguageData(element, element.getAttribute(languageAttribute));
}

const updateCheckerMessage = $(".update-checker .message");
const updateCheckerEventNames = ["startupUpdateCheckerError", "startupUpdateCheckerUpToDate", "startupUpdateCheckerUpdateAvailable"];

const setUpdateCheckerMessage = (key) => setLanguageData(updateCheckerMessage, key);
const unlistenStartUpUpdateCheckerEvents = () => window.runtime.EventsOff(updateCheckerEventNames[0], updateCheckerEventNames[1], updateCheckerEventNames[2]);

setUpdateCheckerMessage("update.messages.checking");

// TODO: sleep for 2 seconds after an event
window.runtime.EventsOnce(updateCheckerEventNames[0], async () => {
    unlistenStartUpUpdateCheckerEvents();

    updateCheckerMessage.style.color = "rgb(255, 100, 100)";
    setUpdateCheckerMessage("update.messages.error");
});
window.runtime.EventsOnce(updateCheckerEventNames[1], async () => {
    unlistenStartUpUpdateCheckerEvents();

    updateCheckerMessage.style.color = "rgb(100, 255, 100)";
    setUpdateCheckerMessage("update.messages.upToDate");
});
window.runtime.EventsOnce(updateCheckerEventNames[2], async (latestVersion) => {
    unlistenStartUpUpdateCheckerEvents();

    console.log("update available" + latestVersion);

    // TODO: hide app & ask user, if accepted, redirect; else, show app when started
});

window.runtime.EventsEmit("contentLoaded");

// TODO: block until updates are checked

$(".update-checker").remove();

console.log("start app");

// TODO: start app