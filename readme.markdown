# go-gcd

API interface implemented in GoLang to make requests to the `comics.org` API endpoints.

## Usage examples

### Series

```go
package main

import (
  "context"
  "fmt"

  gcd "github.com/ipkgs/go-gcd"
)

func main() {
    api := gcd.API{}
    
    resp, err := api.Series(context.Background(), gcd.SeriesReq{})
    if err != nil {
        panic(err)
    }
    
    for i, result := range resp.Results {
        fmt.Printf("%d: %s %q %s\n", i, result.Name, result.APIURL, result.Language)
    }
}
```

The return values can be narrowed down by providing a series name on the request, for example:

```go
gcd.SeriesReq{
    Name: "Superman"
}
```

And further narrowed down by the year the series started:
```go
gcd.SeriesReq{
    Name: "Superman", 
    Year: "2023",
}
```

## Authentication

If you have an account, the cookie value of `gcdsessionid` can be provided to the API to unlock more frequent requests,
by passing an extra parameter when initializing the API object:

```go
api := gcd.API{
    SessionID: "cookie-value-here"
}
```


## Author

Sergio Moura [https://sergio.moura.ca](https://sergio.moura.ca)
