import {
    GetLanguageData_,
} from "../wailsjs/go/main/Application"; // this path is only for the ide features, this path gets replaced in the assets server handler

const $ = (query, _array) => {
	let results = document.querySelectorAll(query);

    if (!_array) return results[0];

    return results;
};

const sleep = duration => new Promise(resolve => setTimeout(resolve, duration * 1000));

window.runtime.EventsOnce("domReady", () => $(".container", false).style.display = "block");

const languageAttribute = "data-language";
const languageElements = $(`[${languageAttribute}]`, true);

for (let element = 0; element < languageElements.length; element++) {
    GetLanguageData_(languageElements[element].getAttribute(languageAttribute))
        .then(
            (value) => languageElements[element].innerText = value,
            (_) => {},
        );
}

const checkedUpdates = false;
const updateCheckerMessage = $(".update-checker .message");
const updateCheckerSleepDuration = 2;

const setUpdateCheckerMessage = (key) => {
    GetLanguageData_(key)
        .then(
            (value) => updateCheckerMessage.innerText = value,
            (_) => {},
        );
};

setUpdateCheckerMessage("update.messages.checking");

// TODO: unlisten these events after one triggers
window.runtime.EventsOnce("startupUpdateCheckerError", async () => {
    updateCheckerMessage.style.color = "rgb(255, 100, 100)";
    setUpdateCheckerMessage("update.messages.error");

    await sleep(updateCheckerSleepDuration);
});
window.runtime.EventsOnce("startupUpdateCheckerUpToDate", async () => {
    updateCheckerMessage.style.color = "rgb(100, 255, 100)";
    setUpdateCheckerMessage("update.messages.upToDate");

    await sleep(updateCheckerSleepDuration);
});
window.runtime.EventsOnce("startupUpdateCheckerUpdateAvailable", () => { // TODO: latestversion argument
    console.log("update available");

    // TODO: hide app & ask user, if accepted, redirect; else, show app when started
});

window.runtime.EventsEmit("contentLoaded");

// TODO: wait until checkedUpdates is true

console.log("start app");

// TODO: start app