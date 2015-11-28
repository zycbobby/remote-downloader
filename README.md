# remote-downloader
Call aria2c remotely

# API

## GET /tasks


## GET /tasks?taskUrl=xxxxx


## POST /tasks 

with json, which "url" is specified


# Run

```bash
go run main.go
```


# TODO

 - Introduce web framework to handle routes [v]
 - Better handling errors
 - Better log
 - MVC
 - Use RWLock or concurrent map to protect data
 - Use chan
 - Find a way for package Management (Postponed, dont try godep...dependency problem cannot be solve if you just move the dependencies into the subfolder, because the import path still using the relative path, so the project is still hard to deploy)
