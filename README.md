This is the source code of FlufBird's (offical) client, made in Go.

To build: Run `npm run build` inside `build/frontend` (Requires [Node.JS](https://nodejs.org/en/download)) and run `build.py` inside `build/application` (Requires [Python](https://www.python.org/downloads)).

For developers, to build the frontend (HTML + SCSS + JS), run `npm run build` inside `build/frontend` (Requires [Node.JS](https://nodejs.org/en/download)), this script will watch for changes for you and process the files automatically. The built assets will be located in `src/frontend/dist`.

To start the application, run `wails dev` inside `src` (Requires [Wails](https://wails.io/docs/gettingstarted/installation)).