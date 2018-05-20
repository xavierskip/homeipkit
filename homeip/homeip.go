package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "time"
    "net/http"
    "io/ioutil"
    "crypto/sha256"
    "strings"
    "flag"
    "encoding/json"
)
const Appname = "Home IP kit"
const Version = "0.2.6.1"


var Step       int64  = 30
var IPfile     string = "myip.txt"
var CMDfile    string = "cmd.txt"
var Urandom    string = "YouMustSetThisInConfigFile!"
var Addr       string = "127.0.0.1:9090"
var TIMEFORMAT string = "2006-01-02 15:04:05"
var Configfile string
var Proxypass bool
var Config map[string]interface{}
var Bot *TelegramApi
var NotificationId string  // telegram Id

func check(err error) {
    if err != nil {
        log.Panicln(err)
    }
}

func checkFileIsExist(filename string) bool {
    var exist = true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false
    }
    return exist
}

func getSha256(message []byte) []byte{
    h := sha256.New()
    h.Write(message)
    return h.Sum(nil)
}

func getTimestep() (time.Time, int64) {
    t := time.Now()
    return t, t.Unix() / Step
}

func getIP(r *http.Request) string{
    var ip string
    var err error
    // fmt.Printf("Proxypass:%v\n", Proxypass)
    if Proxypass {
        ip = r.Header.Get("X-Forwarded-For")
    } else {
        if ip, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
            log.Printf("SplitHostPort ERROR:%s\n", err)
        }
    }
    return ip
}

type indexHandler struct{
    name string
    value string
}

func (index *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        fmt.Fprintf(w, "%s %s\n", index.name, index.value)
        return
    }
    if r.Method == "POST" {
        var dat map[string]interface{}
        body, err := ioutil.ReadAll(r.Body)
        check(err)
        err = json.Unmarshal(body, &dat)
        check(err)
        // log.Println(&dat)
        message := dat["message"].(map[string]interface{})
        text := message["text"].(string)
        chat := message["chat"].(map[string]interface{})
        chatid := fmt.Sprintf("%.0f", chat["id"].(float64))
        // log.Println("chatid:",chatid)
        if text == "/start" {
            _, err := Bot.SendMessage(chatid, "别问我为什么，只有我告诉你什么。")
            if err != nil {
                log.Printf("BOT error:%s\n", err)
            }            
        }
        if text == "/help" {
            _, err := Bot.SendMessage(chatid, "一个私人机器人")
            if err != nil {
                log.Printf("BOT error:%s\n", err)
            }
        }
        return
    }
    http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
    return
}

func myip(w http.ResponseWriter, r *http.Request) {
    // fmt.Printf("r:\n%+v\n", r)
    if r.Method == "GET" {
        myip, err := ioutil.ReadFile(IPfile)
        check(err)
        fileinfo, err := os.Stat(IPfile)
        check(err)
        modtime := fileinfo.ModTime().Format(TIMEFORMAT)
        fmt.Fprintf(w, "%s\n%s", myip, modtime)
        return
    }
    if r.Method == "PUT" {
        t, t_step := getTimestep()
        t_format := t.Format(TIMEFORMAT)
        message := fmt.Sprintf("%s%d", Urandom, t_step)
        key := fmt.Sprintf("%x", getSha256([]byte(message)))
        // r.ParseForm()
        // homekey := strings.Join(r.Form["homekey"],",")
        // fmt.Printf("%v\n%v\n", r.Form, r.PostForm)
        homekey := r.FormValue("homekey")
        if homekey == key {
            ip := getIP(r)
            dat, err := ioutil.ReadFile(IPfile)
            check(err)
            // IP change check
            if ip == string(dat){
                log.Println("IP not change")
            }else{
                log.Println("Change IP!")
                // record to file
                if err := ioutil.WriteFile(IPfile, []byte(ip), 0640); err != nil {
                    log.Panicln(err)
                }
                // change ddns
                login_token := Config["dnspodtoken"].(string)
                format := Config["format"].(string)
                domain_id := Config["domain_id"].(string)
                record_id := Config["record_id"].(string)
                sub_domain := Config["sub_domain"].(string)
                record_line := Config["record_line"].(string)
                go RecordDdns(login_token,
                    format,
                    domain_id,
                    record_id,
                    sub_domain,
                    record_line,
                    ip)
                // telegram Bot
                text := fmt.Sprintf("%s\nYour IP change to %s\n",t_format, ip)
                go Bot.SendMessage(NotificationId, text)
                // resp, e := Bot.SendMessage(NotificationId, text)
                // if e != nil {
                //     log.Printf("BOT err: %s\n", e)
                //     return
                // }
                // log.Printf("BOT: %s\n",resp.Status)
            }
            fmt.Fprintf(w, ip)
            
        } else {
            text := fmt.Sprintf("%s\ntimestamp:%v\nstep:%v\nRequest:%s",t_format, t.Unix(), t_step, homekey)
            go Bot.SendMessage(NotificationId, text)
            log.Println(text)
            http.Error(w, http.StatusText(403), http.StatusForbidden)
        }
        return 
    }
    http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
    return
}

func main() {
    log.SetOutput(os.Stdout)
    // command-line flag parsing
    flag.StringVar(&Configfile, "c", "config.json", "json config file")
    version := flag.Bool("v", false, "version")
    flag.Parse()
    // show version
    if *version {
      fmt.Println(Version)
      os.Exit(0)
    }
    // read json and set Config
    dat, err := ioutil.ReadFile(Configfile)
    check(err)
    err = json.Unmarshal(dat, &Config)
    check(err)
    Addr = Config["addr"].(string)
    Step = int64(Config["step"].(float64))
    Urandom = Config["random"].(string)
    Proxypass = strings.HasPrefix(Addr, "127.0.0.1")
    NotificationId = Config["notification_id"].(string)
    // telegram Bot
    Bot = NewBot(Config["telegram_token"].(string))
    // check file
    filelist := []string{IPfile, CMDfile}
    for _, filename := range filelist {
        if checkFileIsExist(filename) == false {
            f, err1 := os.Create(filename)
            check(err1)
            defer f.Close()
            fmt.Printf("create file:%s\n", filename)
        }
    }
    // test
    // post_test()
    // server
    http.Handle("/tgbot", &indexHandler{Appname, Version})
    http.HandleFunc("/myip", myip)
    log.Fatalln(http.ListenAndServe(Addr, nil))
}