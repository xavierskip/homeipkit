package main

import "net/http"
import "net/url"
import "io/ioutil"
import "log"
import "strings"

var UserAgent string = "MyHomeIPkit/0_0"

func RecordDdns(login_token string,
                format string,
                domain_id string,
                record_id string,
                sub_domain string,
                record_line string,
                value string) (*http.Response, error) {
    var api string = "https://dnsapi.cn/Record.Ddns"
    form := url.Values{
        "login_token": {login_token},
        "format": {format},
        "domain_id": {domain_id},
        "record_id": {record_id},
        "sub_domain": {sub_domain},
        "record_line": {record_line},
        "value": {value}}
    client := &http.Client{}
    req, err := http.NewRequest("POST", api, strings.NewReader(form.Encode()))
    check(err)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", UserAgent)
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func post_test() {
    form := url.Values{}
    form.Add("region", "San Francisco")
    client := &http.Client{}
    req, err := http.NewRequest("POST", "http://httpbin.org/anything", strings.NewReader(form.Encode()))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "Golang_Spider_Bot/3")
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
}