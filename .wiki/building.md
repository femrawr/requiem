1. Run `/.scripts/start-builder.ps1`.
2. Wait for the server to start.
3. Click on the URL it prints. It should look something like "listening on http // localhost:7305".
4. Once you are done with the config click the "BUILD" button at the bottom left of the screen.
5. Wait for the build to finish. The finished build will be in the `/.builds/` folder.

## Notes
- You can hold down any of the control keys to make the build button only update the config instead of building it too.
- To use "Obfuscate build" in the "BUILD SETTINGS" tab, you need to have [garble](https://github.com/burrowers/garble) installed.
- To use "Pack build" in the "BUILD SETTINGS" tab, you need to have [upx](https://github.com/upx/upx) installed and added to PATH.
