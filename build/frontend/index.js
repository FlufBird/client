const fs = require("fs");

const postcss = require("postcss");

const htmlMinifier = require("html-minifier-terser");
const sass = require("sass");
const cssnano = require("cssnano");
const terser = require("terser");

const frontendSourceDirectory = "../../src/frontend/src";
const frontendDistributeDirectory = "../../src/frontend/dist";

const htmlSourceFile = `${frontendSourceDirectory}/index.html`;
const cssSourceFile = `${frontendSourceDirectory}/index.scss`;
const jsSourceFile = `${frontendSourceDirectory}/index.js`;

const htmlDistributeFile = `${frontendDistributeDirectory}/index.html`;
const cssDistributeFile = `${frontendDistributeDirectory}/index.css`;
const jsDistributeFile = `${frontendDistributeDirectory}/index.js`;

const htmlMinifierOptions = JSON.parse(fs.readFileSync(".htmlminifierrc.config.json"));
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

if (!fs.existsSync(frontendDistributeDirectory)) {
    console.log(`${frontendDistributeDirectory} doesn't exist, attemping to create directory...`);

    try {
        fs.mkdirSync(frontendDistributeDirectory)
    } catch (error) {
        console.log(`Couldn't create ${frontendDistributeDirectory}: ${error}`);
        console.log("Exiting...");

        exit(1);
    }
}

console.log(`[${getTime()}] Watching for changes.\n`);

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
const processCss = async () => {
    try {
        const compiled = (await sass.compileAsync(cssSourceFile)).css;
        const processed = (await cssProcessor.process(compiled)).css;

        fs.writeFile(cssDistributeFile, processed, (callback) => {
			if (callback !== null) {
                errorMessage(cssSourceFile, callback);

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
                    errorMessage(jsSourceFile, callback);

                    return;
                }

                processedMessage(jsSourceFile);
            });
        });
    } catch (error) {
		errorMessage(jsSourceFile, error);
	}
};

minifyHtml();
processCss();
minifyJs();

fs.watchFile(htmlSourceFile, watchFileOptions, (_, __) => {
    modifiedMessage(htmlSourceFile);
    minifyHtml();
});
fs.watchFile(cssSourceFile, watchFileOptions, (_, __) => {
    modifiedMessage(cssSourceFile);
    processCss();
});
fs.watchFile(jsSourceFile, watchFileOptions, (_, __) => {
    modifiedMessage(jsSourceFile);
    minifyJs();
});