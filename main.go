package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Example url to shorten.
// http://localhost:8080?import=eyJpZCI6Ijk5ZDg4YWExLTQxZWItNDM4MS1hZGYwLTUzY2VmYzY2NTkxNyIsIm5hbWUiOiJTaG9ydCBMaXN0IiwiaXRlbXMiOlt7ImlkIjoiZmRkMjAyOWEtNWY4YS00NWIxLWEwYjQtNmMwNTFhOWNmODU5IiwibmFtZSI6IkJyZWFkIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOnRydWUsImNyZWF0ZWQiOjE2MzY0NzYxNDQwNTcsInVwZGF0ZWQiOjE2MzY0ODQxNjIxMDcsImRlbGV0ZWQiOm51bGx9LHsiaWQiOiI3YTM0NDQ3ZS0xZWM3LTRjOTYtYjM5OS05NDM1M2ZjNWU0NmIiLCJuYW1lIjoiRWdncyIsInF1YW50aXR5IjoxLCJjaGVja2VkIjpmYWxzZSwiY3JlYXRlZCI6MTYzNjQ3NjE0NDg1NywidXBkYXRlZCI6MTYzNjQ3NjE0NDg1NywiZGVsZXRlZCI6bnVsbH0seyJpZCI6IjZiYmY4ZGM5LTExN2UtNDAyNC1hZTQ1LTVhMDA4Y2NkMmRiOCIsIm5hbWUiOiJNaWxrIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOmZhbHNlLCJjcmVhdGVkIjoxNjM2NDgzNDYyMDcxLCJ1cGRhdGVkIjoxNjM2NDgzNDYyMDcxLCJkZWxldGVkIjpudWxsfV19

func main() {
	jsonStr := []byte(`{"target": "http://localhost:8080?import=eyJpZCI6Ijk5ZDg4YWExLTQxZWItNDM4MS1hZGYwLTUzY2VmYzY2NTkxNyIsIm5hbWUiOiJTaG9ydCBMaXN0IiwiaXRlbXMiOlt7ImlkIjoiZmRkMjAyOWEtNWY4YS00NWIxLWEwYjQtNmMwNTFhOWNmODU5IiwibmFtZSI6IkJyZWFkIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOnRydWUsImNyZWF0ZWQiOjE2MzY0NzYxNDQwNTcsInVwZGF0ZWQiOjE2MzY0ODQxNjIxMDcsImRlbGV0ZWQiOm51bGx9LHsiaWQiOiI3YTM0NDQ3ZS0xZWM3LTRjOTYtYjM5OS05NDM1M2ZjNWU0NmIiLCJuYW1lIjoiRWdncyIsInF1YW50aXR5IjoxLCJjaGVja2VkIjpmYWxzZSwiY3JlYXRlZCI6MTYzNjQ3NjE0NDg1NywidXBkYXRlZCI6MTYzNjQ3NjE0NDg1NywiZGVsZXRlZCI6bnVsbH0seyJpZCI6IjZiYmY4ZGM5LTExN2UtNDAyNC1hZTQ1LTVhMDA4Y2NkMmRiOCIsIm5hbWUiOiJNaWxrIiwicXVhbnRpdHkiOjEsImNoZWNrZWQiOmZhbHNlLCJjcmVhdGVkIjoxNjM2NDgzNDYyMDcxLCJ1cGRhdGVkIjoxNjM2NDgzNDYyMDcxLCJkZWxldGVkIjpudWxsfV19"}`)

	client := http.Client{}
	req, err := http.NewRequest("POST", "https://kutt.it/api/v2/links", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "")
	if err != nil {
		log.Fatalln(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode == 201 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(body))
	} else {
		fmt.Println(res.StatusCode)
	}
}
