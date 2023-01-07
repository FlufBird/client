const { exit } = require("process");
const fs = require("fs");

const postcss = require("postcss");

const htmlMinifier = require("html-minifier-terser");
const sass = require("sass");
const cssnano = require("cssnano");
const terser = require("terser");

const frontendDirectory = "../../src/frontend";

const frontendSourceDirectory = `${frontendDirectory}/src`;
const frontendDistributeDirectory = `${frontendDirectory}/dist`;

const htmlFile = "index.html";

const htmlSourceFile = `${frontendSourceDirectory}/${htmlFile}`;
const htmlDistributeFile = `${frontendDistributeDirectory}/${htmlFile}`;

const cssVariablesFile = "variables.scss";

const cssSourceDirectory = `${frontendSourceDirectory}/css`;
const cssDistributeDirectory = `${frontendDistributeDirectory}/css`;

const jsSourceDirectory = `${frontendSourceDirectory}/js`;
const jsDistributeDirectory = `${frontendDistributeDirectory}/js`;

const htmlMinifierOptions = JSON.parse(fs.readFileSync(".htmlminifierrc.config.json"));
const cssProcessor = postcss([cssnano()]);
const jsMinifierOptions = JSON.parse(fs.readFileSync(".terserrc.config.json"));

const watchOptions = {
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
const modifiedMessage = (file, path) => console.log(`[${getTime()}] ${path} has been modified${(file !== cssVariablesFile) ? "" : " (Variables file)"}. ${(file !== cssVariablesFile) ? "" : "\n"}`);
const processedMessage = (path) => console.log(`[${getTime()}] ${path} has been processed.\n`);

const minifyHtml = async () => {
    try {
        fs.readFile(htmlSourceFile, async (error, data) => {
            if (error) {
                errorMessage(error);

                return;
            }

            const minified = await htmlMinifier.minify(data.toString(), htmlMinifierOptions);

            fs.writeFile(htmlDistributeFile, minified, (callback) => {
				if (callback !== null) {
                    errorMessage(htmlSourceFile, callback);

                    return;
                }

                processedMessage(htmlSourceFile);
            });
        });
    } catch (error) {
		errorMessage(htmlSourceFile, error);
	}
}
const processCss = async (firstTime, file, path) => {
    if (file === cssVariablesFile) {
        if (firstTime) return;

        try {
            walkCssSourceDirectory(false, true, async (file, path) => processCss(false, file, path));
        } catch (error) {
            console.log(`[${getTime()}] Error occured while processing all files (Variables file changed): ${error}\n`);
        }

        return;
    }

    try {
        const compiled = (await sass.compileAsync(path)).css;
        const processed = (await cssProcessor.process(compiled)).css;

        fs.writeFile(`${cssDistributeDirectory}/${file.slice(0, -5)}.css`, processed, (callback) => {
            if (callback !== null) {
                errorMessage(path, callback);

                return;
            }

            processedMessage(path);
        });
    } catch (error) {
        errorMessage(path, error);
    }
};
const minifyJs = (file, path) => {
    try {
        fs.readFile(path, async (error, data) => {
            if (error) {
                errorMessage(path, error);

                return;
            }

            const minified = (await terser.minify(data.toString(), jsMinifierOptions)).code;

            fs.writeFile(`${jsDistributeDirectory}/${file}`, minified, (callback) => {
                if (callback !== null) {
                    errorMessage(path, callback);

                    return;
                }

                processedMessage(path);
            });
        });
    } catch (error) {
        errorMessage(path, error);
    }
};

const walkCssSourceDirectory = (readDirectoryErrorExit, ignoreCssVariablesFile, callback) => {
    fs.readdir(cssSourceDirectory, (error, items) => {
        if (error) {
            if (!readDirectoryErrorExit) throw error;

            exit(1);
        }

        items.forEach((item, _) => {
            const path = `${cssSourceDirectory}/${item}`;

            fs.stat(path, (error, statistics) => {
                if (error) return;

                if (statistics.isFile()) {
                    if (item === cssVariablesFile && ignoreCssVariablesFile) return;

                    callback(item, path);
                }
            });
        });
    });
};

const startWatchingHtmlFile = () => {
    minifyHtml();

    fs.watchFile(htmlSourceFile, watchOptions, (_, __) => {
        modifiedMessage(htmlSourceFile);
        minifyHtml();
    });
};
const startWatchingCssFile = (file, path) => {
    processCss(true, file, path);

    fs.watchFile(path, watchOptions, (_, __) => {
        modifiedMessage(file, path);
        processCss(false, file, path);
    });
};
const startWatchingJsFile = (file, path) => {
    minifyJs(file, path);

    fs.watchFile(path, watchOptions, (_, __) => {
        modifiedMessage(file, path);
        minifyJs(file, path);
    });
};

for (directory of [frontendDistributeDirectory, cssDistributeDirectory, jsDistributeDirectory]) {
    if (!fs.existsSync(directory)) {
        console.log(`${directory} doesn't exist, attemping to create directory...`);

        try {
            fs.mkdirSync(directory)
        } catch (error) {
            console.log(`Couldn't create ${directory}: ${error}`);
            console.log("Exiting...");

            exit(1);
        }
    }
}

console.log(`\n[${getTime()}] Watching for changes.\n`);

walkCssSourceDirectory(true, false, (file, path) => startWatchingCssFile(file, path));

fs.readdir(jsSourceDirectory, (error, items) => {
    if (error) exit(1);

    items.forEach((item, _) => {
        const path = `${jsSourceDirectory}/${item}`;

        fs.stat(path, async (error, statistics) => {
            if (error) return;

            if (statistics.isFile()) startWatchingJsFile(item, path);
        });
    });
});

startWatchingHtmlFile();
fs.watch(cssSourceDirectory, watchOptions, (event, file) => {
    let path = `${cssSourceDirectory}/${file}`;

    if (event != "rename" || !fs.existsSync(path) || !file) return;

    startWatchingCssFile(file, path);
});
fs.watch(jsSourceDirectory, watchOptions, (event, file) => {
    let path = `${jsSourceDirectory}/${file}`;

    if (event != "rename" || !fs.existsSync(path) || !file) return;

    startWatchingJsFile(file, path);
});