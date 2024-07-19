package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

const ServerPort = 8080
const BaseUrl = "/"
const ShortenBaseUrl = "/shorten"
const ShortBaseUrl = "/short"

type UrlShortener struct {
	urls map[string]string
}

var formsHtml = map[string]string{
	"main": `<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>URL Shortener</title>
	<style>
	body {
		font-family: Arial, sans-serif;
		background-color: #f0f0f0;
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100vh;
		margin: 0;
	}
	.container {
		background-color: #fff;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		text-align: center;
		width: 500px;
	}
	h2 {
		font-size: 24px;
		color: #333;
	}
	form {
		display: flex;
		flex-direction: column;
	}
	input[type="url"] {
		padding: 10px;
		margin-bottom: 10px;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 16px;
	}
	input[type="submit"] {
		padding: 10px;
		background-color: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 16px;
		cursor: pointer;
		transition: background-color 0.3s;
	}
	input[type="submit"]:hover {
		background-color: #0056b3;
	}
	</style>
	</head>
	<body>
	<div class="container">
	<h2>URL Shortener</h2>
	<form method="post" action="/shorten">
		<input type="url" name="url" placeholder="Enter a URL" required>
		<input type="submit" value="Shorten">
	</form>
	</div>
	</body>
	</html>`,
	"shorten": `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>URL Shortener</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f0f0f0;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
			}
			.container {
				background-color: #fff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				text-align: center;
				width: 500px;
			}
			h2 {
				font-size: 24px;
				color: #333;
			}
			p {
				font-size: 16px;
				color: #666;
				text-align: left;
				text-decoration: left;
			}
			a {
				color: #007bff;
				text-decoration: none;
			}
			a:hover {
				text-decoration: underline;
			}
			form {
				display: flex;
				flex-direction: column;
				margin-top: 20px;
			}
			input[type="text"] {
				padding: 10px;
				margin-bottom: 10px;
				border: 1px solid #ccc;
				border-radius: 4px;
				font-size: 16px;
			}
			input[type="submit"] {
				padding: 10px;
				background-color: #007bff;
				color: white;
				border: none;
				border-radius: 4px;
				font-size: 16px;
				cursor: pointer;
				transition: background-color 0.3s;
			}
			input[type="submit"]:hover {
				background-color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>URL Shortener</h2>
			<p>Original URL: %s</span></p>
			<p>Shortened URL: <a href="%s">%s</a></p>
			<form method="post" action="/shorten">
				<input type="text" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</div>
	</body>
	</html>`,
}

func generateShortKey(url string) string {
	h := md5.New()
	h.Write([]byte(url))
	hashValue := h.Sum(nil)

	return hex.EncodeToString(hashValue)
}

func (us *UrlShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey(originalURL)
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:%v%s/%s", ServerPort, ShortBaseUrl, shortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(formsHtml["shorten"], originalURL, shortenedURL, shortenedURL)
	_, err := fmt.Fprintf(w, responseHTML)
	if err != nil {
		return
	}
}

func (us *UrlShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len(ShortBaseUrl):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (us *UrlShortener) mainForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, ShortenBaseUrl, http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, formsHtml["main"])
	if err != nil {
		return
	}
}

func main() {
	shortener := &UrlShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc(BaseUrl, shortener.mainForm)
	http.HandleFunc(ShortenBaseUrl, shortener.HandleShorten)
	http.HandleFunc(ShortBaseUrl, shortener.HandleRedirect)

	fmt.Println(fmt.Sprintf("URL Shortener is running on :%v", ServerPort))
	err := http.ListenAndServe(fmt.Sprintf(":%v", ServerPort), nil)
	if err != nil {
		return
	}

}
