const { exit } = require("process");
const fs = require("fs");

const postcss = require("postcss");

const sass = require("sass");
const autoprefixer = require("autoprefixer");
const cssnano = require("cssnano");
const terser = require("terser");

const cssSourceFile = "../../src/ui/style.scss";
const cssDistributeFile = "../../src/ui/style.min.css";

const jsSourceFile = "../../src/ui/index.js";
const jsDistributeFile = "../../src/ui/index.min.js";

const cssProcessor = postcss([autoprefixer(), cssnano()]);
const jsMinifierOptions = JSON.parse(fs.readFileSync(".terserrc.config.json"));

const watchFileOptions = {
    persistent: true,

    interval: 1000,
};

const getTime = () => {
    let date = new Date();

    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
};

const errorMessage = (path, error) => {
    console.log(`[${getTime()}] Couldn't process ${path}: ${error}\n`);
};
const modifiedMessage = (path) => {
    console.log(`[${getTime()}] ${path} is modified.`);
}; 
const processedMessage = (path) => {
    console.log(`[${getTime()}] ${path} has been processed.\n`);
};

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