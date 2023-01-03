This is the source code of FlufBird's (offical) client, made in Go and uses Wails for the frontend.

To build the program:
1. Run `wails generate module` inside `src`.
2. Run `npm run build` inside `build/frontend` (Requires [Node.JS](https://nodejs.org/en/download)).
3. Run `build.py` inside `build/application` (Requires [Python](https://www.python.org/downloads)).
4. The binaries and ready to distribute `.zip` archives for supported OSes and architectures will be located inside `dist`.

For developers: 
- To build the frontend (HTML + SCSS + JS), run `npm run build` inside `build/frontend` (Requires [Node.JS](https://nodejs.org/en/download)), this script will watch for changes for you and process the files automatically. The built assets will be located in `src/frontend/dist`.
- To start the application, `wails dev` inside `src` (Requires [Wails](https://wails.io/docs/gettingstarted/installation)).