const { exit } = require("process");
const fs = require("fs");

const postcss = require("postcss");

const sass = require("sass");
const autoprefixer = require("autoprefixer");
const cssnano = require("cssnano");
const terser = require("terser");

const cssSourcePath = "../src/ui/css/src";
const cssDistributePath = "../src/ui/css/dist";

const jsSourcePath = "../src/ui/js/src";
const jsDistributePath = "../src/ui/js/dist";

const cssProcessor = postcss([autoprefixer(), cssnano()]);
const jsMinifierOptions = JSON.parse(fs.readFileSync(".terserrc.config.json"));

fs.readdir(cssSourcePath, (error, items) => {
    if (error) {
        exit(1);
    }

    items.forEach((item, _) => {
        const path = `${cssSourcePath}/${item}`;

        fs.stat(path, async (error, statistics) => {
            if (error) {
                return;
            }

            if (statistics.isFile()) {
                const compiled = (await sass.compileAsync(path)).css;
                const processed = (await cssProcessor.process(compiled)).css;

                fs.writeFile(`${cssDistributePath}/${item.slice(0, -5)}.min.css`, processed, (_) => {});
            }
        });
    });
});

fs.readdir(jsSourcePath, (error, items) => {
    if (error) {
        exit(1);
    }

    items.forEach((item, _) => {
        const path = `${jsSourcePath}/${item}`;

        fs.stat(path, async (error, statistics) => {
            if (error) {
                return;
            }

            if (statistics.isFile()) {
                fs.readFile(path, async (error, data) => {
                    if (error) {
                        return;
                    }

                    const minified = (await terser.minify(data.toString(), jsMinifierOptions)).code;

                    fs.writeFile(`${jsDistributePath}/${item.slice(0, -3)}.min.js`, minified, (_) => {});
                });
            }
        });
    });
});