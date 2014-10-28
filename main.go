package main

import (
    "fmt"
    "flag"
    "time"
    "sync"
    "net/http"
    "log"
    "crypto/tls"
    "os"
    twilio "github.com/carlosdp/twiliogo"
)

// Command-line flags
var (
    httpAddr = flag.String("http", ":8080", "Listen address")
    pollPeriod = flag.Duration("poll", 10*time.Second, "Poll period")
    target = flag.String("target", "http://google.com", "Target to monitor")
)

var badCount int = 0

var TWILIO_SID string = os.Getenv("TWILIO_SID")
var TWILIO_TOKEN string = os.Getenv("TWILIO_TOKEN")
var SMS_FROM string = os.Getenv("SMS_FROM")
var SMS_TO string = os.Getenv("SMS_TO")
var twilio_client = twilio.NewClient(TWILIO_SID, TWILIO_TOKEN)

type Server struct {
    url string
    period time.Duration
    up bool

    mu sync.RWMutex // protects the up variable
}

func NewServer(url string, period time.Duration) *Server {
    s := &Server{url: url, period: period}
    go s.poll()
    return s
}

func (s *Server) poll() {
    for !isActive(s.url) {
        pollSleep(s.period)
    }
    s.mu.Lock()
    s.up = true
    s.mu.Unlock()
    pollDone()
}

type SMSMessage string

func (s SMSMessage) GetParam() (string, string){
    return "Body", string(s)
}

var (
    pollSleep = time.Sleep
    pollDone = func() {}
)

func isActive(url string) bool {
    // Ignoring Certificate Authorities for now -- we're mostly concerned with availability
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    r, err := client.Head(url)
    if err != nil {
        log.Print(err)
        badCount++
        if badCount >= 5 {
            log.Print("This is too much! Im sending a message!")
            s := fmt.Sprintf("%s is down", url)
            var m SMSMessage = SMSMessage(s)
            message, err := twilio.NewMessage(
                twilio_client,
                SMS_FROM,
                SMS_TO,
                m,
            )
            if err != nil {
                fmt.Println(err)
            } else {
                fmt.Println(message.Status)
            }
            badCount = 0
        }
        return false
    }
    badCount = 0
    return r.StatusCode == http.StatusOK
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mu.RLock()
    data := struct {
        URL string
        Up bool
    }{
        s.url,
        s.up,
    }
    s.mu.RUnlock()
    if data.Up{
        w.Write([]byte(fmt.Sprintf("Looks like %s is still up", data.URL)))
    } else {
        w.Write([]byte(fmt.Sprintf("%s appears to be down!", data.URL)))
    }

    log.Println(data.Up)

}

func main(){
    flag.Parse()
    fmt.Printf("Jim is now monitoring %s...\n", *target)
    http.Handle("/", NewServer(*target, time.Second*5))
    log.Print(http.ListenAndServe(*httpAddr, nil))
}
