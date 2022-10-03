const $ = (query) => {
	return document.querySelector(query);
};

window.onload = () => {
    $(".container").style.display = "block";
};