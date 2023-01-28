# CFCS - Concurrent File Client and Server Project
The aim of this project is to create a server that will
listen to file changes in a configured directory. 
It will accept incoming tcp connections and "broadcast" the files to the client.

The client will listen to the available files, check if that file had been downloaded,
and download it if not already downloaded.

## Docker Stuffs

### Build Server Image
```shell
docker build . -f ./dockerfiles/server.dockerfile -t rmarken5/test-cfs
```

### Build Client