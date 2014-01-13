package main

import (
  "net/http"
  "log"
  "io"
  "io/ioutil"
  "encoding/json"
  "github.com/sfreiberg/gotwilio"
)

type Options struct {
  Path string
  Port string
}

var messages map[string] chan string = make(map[string] chan string)
var queues = make(map[string]chan []byte)
//set twilio details
var accountSid = "AC5c233662524ac52727857f6a02024d4b"
var authToken = "89e66e03ed6af99ff58c578ea57c54f8"
var twilio = gotwilio.NewTwilioClient(accountSid, authToken)

func PushHandler(w http.ResponseWriter, req *http.Request) {
  rcpt := req.FormValue("rcpt")
  var length = len(rcpt)
  var isPhone = ""
  //verify valid phone number, prepend +1
  if (length == 10) {
    t := 0
    for _, value := range rcpt {
      switch {
        case value >= '0' && value <= '9':
          t++
      }
    }
    if (t == length) {
       //add leading +1
      isPhone = "+1" + rcpt
      log.Println("isPhone now equal: " + isPhone)
    }
  }
  
  //verify valid phone number, prepend +
  if (length == 11) {
    t := 0
    for _, value := range rcpt {
      switch {
      case value >= '0' && value <= '9':
        t++
      }
    }
    if (t == length) {
      //add leading +
      isPhone = "+" + rcpt
      log.Println("isPhone now equal: " + isPhone)
    }
  }
  
  //get data from push request
  body, err := ioutil.ReadAll(req.Body)
  //check for error, write success header
  if err != nil || rcpt == "" || req.Method != "POST" {
		w.WriteHeader(400)
		return
	}
  
  ch := messages[rcpt]
	
	// new user?
	if ch == nil {
		ch = make (chan string)
		messages[rcpt] = ch
	}
	
	// store message
  ch <- string(body)
  
  //set message send details
  from := "+16364338457"
  message := string(body)
  twilio.SendSMS(from, isPhone, message, "", "") 
  log.Println("\nFrom: " + from + "\nPhone: " + isPhone + "\nMessage: " + message)
}

func PollResponse(w http.ResponseWriter, req *http.Request) {
  rcpt := req.FormValue("rcpt")
  
  if req.Method != "GET" || rcpt == "" {
		w.WriteHeader(400)
		return		
	}
  ch := messages[rcpt]
	
	// new user?
	if ch == nil {
		ch = make (chan string)
		messages[rcpt] = ch
	}
	
	io.WriteString(w, <-ch)
  
}

func main() {
   // set default options
  op := &Options{Path: "./", Port: "8008"}

  // read config file into memory
  data, jsonErr := ioutil.ReadFile("./config.json")
  if jsonErr != nil {
    log.Println("JSONReadFileError: ", jsonErr) 
  }
  
  // parse config file, store results in "op"
  json.Unmarshal(data, op)
  log.Println("Parsed options from config file: ", op) 

  //set handlers
  http.Handle("/", http.FileServer(http.Dir("./")))
  http.HandleFunc("/poll", PollResponse)
  http.HandleFunc("/push", PushHandler)
  //serve http
  err := http.ListenAndServe(":" + op.Port, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}