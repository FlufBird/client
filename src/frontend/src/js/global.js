import {
    GetLanguageData_,
} from "../../wailsjs/go/main/Application"; // this path is only for the ide features, itll get replaced in the assets server handler

const $ = (query, _array) => {
	let results = document.querySelectorAll(query);

    if (!_array) return results[0];

    return results;
};

const setLanguageData = (element, key) => GetLanguageData_(key).then((value) => element.innerText = value);

const languageAttribute = "data-language";
const languageElements = $(`[${languageAttribute}]`, true);

for (let index = 0; index < languageElements.length; index++) {
    const element = languageElements[index];

    setLanguageData(element, element.getAttribute(languageAttribute));
}