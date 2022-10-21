const { exit } = require("process");
const fs = require("fs");

const postcss = require("postcss");

const sass = require("sass");
const autoprefixer = require("autoprefixer");
const cssnano = require("cssnano");
const terser = require("terser");

const cssSourceFile = "../src/ui/style.scss";
const cssDistributeFile = "../src/ui/style.min.css";

const jsSourceFile = "../src/ui/index.js";
const jsDistributeFile = "../src/ui/index.min.js";

const cssProcessor = postcss([autoprefixer(), cssnano()]);
const jsMinifierOptions = JSON.parse(fs.readFileSync(".terserrc.config.json"));

const processCss = async () => {
    const compiled = (await sass.compileAsync(cssSourceFile)).css;
    const processed = (await cssProcessor.process(compiled)).css;

    fs.writeFile(cssDistributeFile, processed, (_) => {});
};
const minifyJs = async () => {
    fs.readFile(jsSourceFile, async (error, data) => {
        if (error) {
            return;
        }

        const minified = (await terser.minify(data.toString(), jsMinifierOptions)).code;

        fs.writeFile(jsDistributeFile, minified, (_) => {});
    });
};

(async () => {
    await processCss();
})();
(async () => {
    await minifyJs();
})();