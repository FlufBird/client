const $ = (query) => {
	let results = document.querySelectorAll(query);

    return (results.length === 1) ? results[0] : results;
};

window.onload = () => {
    $(".container").classList.remove("hidden");
};