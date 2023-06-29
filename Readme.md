# kataGo-server
it's a kataGo analysis server by golang.

## Features
* just can analysis a go ban game

## How to use?
* check your machine has nvidia gpu
* your system has docker and nvidia-docker
* your nvidia-driver support opencl
* build the backend by dockerfile to image
* run the image by `nvidia-docker` or `docker with nvidia-docker runtime` and expose the 8080 port
* when you see the `Started, ready to begin handling requests` by the container logger ,the backend is running
* your need tauri runtime link(https://tauri.app/)
* cd demo-frontend
* do `npm install`
* run the frontend demo by `npm run tauri dev`

## Prod Practice
* one server can hold one request in atom time
* need a lb server or used task assignment server 