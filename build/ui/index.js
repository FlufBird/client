const { exit } = require("process");
const fs = require("fs");

const postcss = require("postcss");

const sass = require("sass");
const cssnano = require("cssnano");
const terser = require("terser");

const cssSourceFile = "../../src/ui/style.scss";
const jsSourceFile = "../../src/ui/index.js";

const cssDistributeFile = "../../src/ui/dist/style.min.css";
const jsDistributeFile = "../../src/ui/dist/index.min.js";

const cssProcessor = postcss([cssnano()]);
const jsMinifierOptions = JSON.parse(fs.readFileSync(".terserrc.config.json"));

const watchFileOptions = {
    persistent: true,
    interval: 1 * 1000,
};

const getTime = () => {
    let date = new Date();

    const addPrefix = (part) => { // its just too annoying i cant fucking take it
        return (part.toString().length !== 1) ? part : "0" + part;
    };

    return `${addPrefix(date.getHours())}:${addPrefix(date.getMinutes())}:${addPrefix(date.getSeconds())}`;
};

const errorMessage = (path, error) => console.log(`[${getTime()}] Couldn't process ${path}: ${error}\n`);
const modifiedMessage = (path) => console.log(`[${getTime()}] ${path} has been modified.`);
const processedMessage = (path) => console.log(`[${getTime()}] ${path} has been processed.\n`);

console.log(`[${getTime()}] Watching for changes.\n`);

const processCss = async () => {
    try {
        const compiled = (await sass.compileAsync(cssSourceFile)).css;
        const processed = (await cssProcessor.process(compiled)).css;

        fs.writeFile(cssDistributeFile, processed, (callback) => {
			if (callback !== null) {
                errorMessage(path, callback);

                return;
            }

            processedMessage(cssSourceFile);
        });
    } catch (error) {
		errorMessage(cssSourceFile, error);
	}
};
const minifyJs = () => {
    try {
        fs.readFile(jsSourceFile, async (error, data) => {
            if (error) {
                errorMessage(error);

                return;
            }

            const minified = (await terser.minify(data.toString(), jsMinifierOptions)).code;
    
            fs.writeFile(jsDistributeFile, minified, (callback) => {
				if (callback !== null) {
                    errorMessage(path, callback);

                    return;
                }

                processedMessage(jsSourceFile);
            });
        });
    } catch (error) {
		errorMessage(jsSourceFile, error);
	}
};

processCss();
minifyJs();

fs.watchFile(cssSourceFile, watchFileOptions, (_, __) => {
    modifiedMessage(cssSourceFile);
    processCss();
});
fs.watchFile(jsSourceFile, watchFileOptions, (_, __) => {
    modifiedMessage(jsSourceFile);
    minifyJs();
});