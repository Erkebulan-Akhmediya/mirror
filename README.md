# Mirror

### About
This program mirrors a website by finding and downloading links using depth-first search.
It uses goroutines to download files concurrently. The files are stored in a directory named after the website host.

### Compile
Compile the program using command below.
```sh
go build mirror
```

### Usage
The program takes one command line argument, which is the starting point for the program.
```sh
mirror https://golang.org
```