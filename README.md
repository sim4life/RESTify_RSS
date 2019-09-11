# RESTify_RSS

This project exposes some RSS feeds using a REST interface

## Usage

Clone the repo in any path other than `$GOPATH`

```
$ git clone https://github.com/sim4life/RESTify_RSS.git
$ cd RESTify_RSS
$ go install
$ RESTify_RSS
```

From another terminal:

```
$ curl -i "http://localhost:3333/articles?category=UK"
$ curl -i "http://localhost:3333/articles?provider=Reuters"
$ curl -i "http://localhost:3333/articles?category=UK&provider=Reuters"
```

## Notes

Only one REST endpoint was exposed, which can serve all the required client requirements so more endpoints were NOT developed.  
Unit testing is NOT very thorough as mocking is NOT used in this exercise testing.  
Very basic single machine in-memory caching is used.

## Requirements

Go v1.12  
Go Modules enabled
