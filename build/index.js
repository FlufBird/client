const { exit } = require("process");
const fs = require("fs");

const sourcePath = "../src/ui/css/src";
const distributePath = "../src/ui/css/dist";

fs.readdir(sourcePath, (error, items) => {
    if (error) {
        exit(1);
    }

    items.forEach((item, _) => {
        const path = `${sourcePath}/${item}`;

        fs.stat(path, async (error, statistics) => {
            if (error) {
                return;
            }

            if (statistics.isFile()) {
                const compiled = (await sass.compileAsync(path)).css;

                fs.writeFile(`${distributePath}/${item.slice(0, -5)}.css`, compiled, (_) => {});
            }
        });
    });
});