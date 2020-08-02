# fast-r
A simple internet speed testing cli/wasm framework

### How To Get Started with the WASM framework

1. Copy The main.wasm binary and wasm_exec.js to your local folder
2. View the sample index.html to understand how to tie things together
3. Features are implemented with promises... The promise resolver global variable has to be set
    and should be accessible from the global dom. Also names should not change, if not the wasm binary might be unable to pick it up
 4. Run GetServer first to return the list of possible servers (This value would be cached for 10 minutes)
 5. Then Run StartDownload Or StartUpload
 
 
 ### Gotchas
 1. This is highly experimental and might not be suitable for production environments
 2. Cors would impact the process of picking up optimal serves, and might impact heavily on latency and accuracy of the test  



### CLI Version

The CLI Version is mostly stable, you can clone the repo and `go run main.go run` in the cli folder.

Offical binary release has not be done yet, as this project is highly experimental.