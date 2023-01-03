// FIXME: ensure EventsOnce is executed before EventsEmit

import {
    GetLanguageData_,
} from "../wailsjs/go/main/Application"; // this path is only for the ide features, this path gets replaced in the assets server handler

const $ = (query, _array) => {
	let results = document.querySelectorAll(query);

    if (!_array) return results[0];

    return results;
};

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

const updateCheckerMessage = $(".update-checker .message");
const setUpdateCheckerMessage = (key) => {
    GetLanguageData_(key)
        .then(
            (value) => updateCheckerMessage.innerText = value,
            (_) => {},
        );
};

setUpdateCheckerMessage("update.message");

window.runtime.EventsEmit("contentLoaded");

window.runtime.EventsOnce("startupUpdateCheckerError", () => {
    updateCheckerMessage.style.color = "rgb(255, 100, 100)";
    setUpdateCheckerMessage("update.error");

    // sleep and start app
});
window.runtime.EventsOnce("startupUpdateCheckerUpdateAvailable", () => {
    console.log("update available");
});