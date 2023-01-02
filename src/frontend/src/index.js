import {
    GetLanguageData_,
} from "../wailsjs/go/main/Application";

const $ = (query, _array) => {
	let results = document.querySelectorAll(query);

    if (_array) return results;
    if (!results) return null;

    return (results.length === 1) ? results[0] : results;
};

window.addEventListener("load", () => $(".container").classList.remove("hidden"));

const languageAttribute = "data-language";
const languageElements = $(`[${languageAttribute}]`, true);

for (let element = 0; element < languageElements.length; element++) {
    languageElements[element].innerText = GetLanguageData_(languageAttribute);
}