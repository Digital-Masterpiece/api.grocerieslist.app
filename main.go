package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Example url to shorten.
// http://localhost:8080?import=eyJpZCI6Ijk5ZDg4YWExLTQxZWItNDM4MS1hZGYwLTUzY2VmYzY2NTkxNyIsIm5hbWUiOiJTaG9ydCBMaXN0IiwiaXRlbXMiOlt7ImlkIjoiZmRkMjAyOWEtNWY4YS00NWIxLWEwYjQtNmMwNTFhOWNmODU5IiwibmFtZSI6IkJyZWFkIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOnRydWUsImNyZWF0ZWQiOjE2MzY0NzYxNDQwNTcsInVwZGF0ZWQiOjE2MzY0ODQxNjIxMDcsImRlbGV0ZWQiOm51bGx9LHsiaWQiOiI3YTM0NDQ3ZS0xZWM3LTRjOTYtYjM5OS05NDM1M2ZjNWU0NmIiLCJuYW1lIjoiRWdncyIsInF1YW50aXR5IjoxLCJjaGVja2VkIjpmYWxzZSwiY3JlYXRlZCI6MTYzNjQ3NjE0NDg1NywidXBkYXRlZCI6MTYzNjQ3NjE0NDg1NywiZGVsZXRlZCI6bnVsbH0seyJpZCI6IjZiYmY4ZGM5LTExN2UtNDAyNC1hZTQ1LTVhMDA4Y2NkMmRiOCIsIm5hbWUiOiJNaWxrIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOmZhbHNlLCJjcmVhdGVkIjoxNjM2NDgzNDYyMDcxLCJ1cGRhdGVkIjoxNjM2NDgzNDYyMDcxLCJkZWxldGVkIjpudWxsfV19

func main() {
	http.HandleFunc("/", endpoint)
	fmt.Println("Starting endpoint for grocerieslist.app short urls.")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
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
			if _, err := fmt.Fprintf(w, "target query found: \"%s\"\n", target); err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}

			kuttResp, kErr := createKuttLink(target)
			if kErr != nil {
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			}

			// Return a JSON response.
			jsonResp, jErr := json.Marshal(kuttResp)
			if jErr != nil {
				log.Fatal(jErr)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			if _, wErr := w.Write(jsonResp); wErr != nil {
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
	jsonStr := []byte(`{"domain": "s.grocerieslist.app", "target": "` + t + `}`)

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
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.New("failed to read status code from kutt.it")
		}
		return string(body), nil
	} else {
		return "", errors.New("failed response from kutt.it")
	}
}
