package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", endpoint)
	fmt.Println("Listening on port 8080 for requests.")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	ao := os.Getenv("ALLOWED_ORIGIN")
	w.Header().Set("Access-Control-Allow-Origin", ao)
	if ao != "*" {
		if r.Header.Get("Origin") != ao {
			http.Error(w, "403 Forbidden: Invalid Origin", http.StatusForbidden)
			return
		}
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Println("GET request received.")
		if _, err := io.WriteString(w, "Send a POST request with the target URL you want shortened."); err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		fmt.Println("POST request received.")
		if err := r.ParseForm(); err != nil {
			if _, pErr := io.WriteString(w, "There was an error parsing your request."); pErr != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		target := r.FormValue("target")
		if target != "" {
			_, err := url.ParseRequestURI(target)
			if err != nil {
				http.Error(w, "400 Bad Request: Malformed Target URL", http.StatusBadRequest)
				return
			}

			kuttResp, kErr := createKuttLink(target)
			if kErr != nil {
				http.Error(w, "400 Bad Request: "+kErr.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if _, wErr := w.Write([]byte(kuttResp)); wErr != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		} else {
			if _, err := io.WriteString(w, "Your designated target was empty."); err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	default:
		if _, err := io.WriteString(w, "Only GET and POST methods are supported."); err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
	}
}

func createKuttLink(t string) (string, error) {
	// Currently kutt.it doesn't support SSL for custom domains. This is obviously an issue with the .app TLD.
	// https://github.com/thedevs-network/kutt/issues/18 is open.
	//jsonStr := []byte(`{"domain": "s.grocerieslist.app", "expire_in": "24 hours", "target": "` + t + `"}`)
	jsonStr := []byte(`{"expire_in": "24 hours", "target": "` + t + `"}`)

	client := http.Client{}
	req, nrErr := http.NewRequest("POST", "https://kutt.it/api/v2/links", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", os.Getenv("KUTT_KEY"))
	if nrErr != nil {
		return "", errors.New("failed to instantiate a new request to kutt.it")
	}

	res, cErr := client.Do(req)
	if cErr != nil {
		return "", errors.New("failed to execute request to kutt.it")
	}

	if res.StatusCode == 201 {
		body, ioErr := ioutil.ReadAll(res.Body)
		if ioErr != nil {
			return "", errors.New("failed to read status code from kutt.it")
		}
		return string(body), nil
	} else {
		return "", errors.New("failed response from kutt.it")
	}
}
