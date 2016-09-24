### What is storm?
`storm` is a simple command line http/https stress testing or benchmarking utility written in `Go`.


### How to install storm?
You may either download the binary from [here](https://github.com/dlutxx/storm/releases).   
Or compile from source:
```
git clone https://github.com/dlutxx/storm.git
cd storm/ && go build -o storm main/main.go
```

### How to use storm?
Just issue `./storm -h` and it will show all available options and give self-explanatory help info.


### Example
Show config and data:   
`> ./storm -requestData post-localhost.txt -method POST -h`
```
Address: localhost:80
Concurrency: 64
RequestsPerThread: 0
data:
---------
POST / HTTP/1.1
Host: localhost
User-Agent: stormer/1.0
X-Requested-With: XMLHttpRequest
Accept: */*
Connection: keep-alive

name=abc&age=123
```


Request `localhost/` with 16 concurrent connections in 4 in seconds:  
`> ./storm -requestData get-localhost.txt -url http://localhost/ -ttr 4 -concurrency 16`
```
Storm the front!
Stopping...
Total: 1600, Failure: 0, Time: 4.925928652s
```
