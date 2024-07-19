# Url Shortener
This is a simple basic url shortener mechanism in golang. 

This code include a `http server`, and `hash` for generating short keys (now we use MD5 as our hashing algorithm).

## Main Routes
- / : base url, show main form
- /shorten : generate short urls
- /short : short urls main route for redirection

This project is mainly based on this [link](https://dev.to/envitab/how-to-build-a-url-shortener-with-go-5hn5)