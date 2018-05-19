package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"golang.org/x/net/html"
)

const basePath = "/api"

type Configuration struct {
	Port       string `json:"Port,omitempty"`
	FromEmail  string
	Password   string
	SmtpServer string
	SmtpPort   string
	ToEmails   []string
}

type Service struct {
	Auth       smtp.Auth
	SmtpServer string
	SmtpPort   string
	FromEmail  string
	ToEmails   []string
}

func main() {
	mux := http.NewServeMux()

	var c Configuration
	if err := loadConfiguration(&c); err != nil {
		log.Fatal(fmt.Sprintf("Could not load configuration: %v", err))
	}

	port := "80"
	if c.Port != "" {
		port = c.Port
	}

	service := NewService(&c)

	mux.HandleFunc(basePath+"/heartbeat", service.HeartbeatHandler)
	mux.HandleFunc(basePath+"/compliment", service.ComplimentHandler)
	log.Println("Starting Server on Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func NewService(c *Configuration) *Service {
	auth := smtp.PlainAuth("", c.FromEmail, c.Password, c.SmtpServer)
	return &Service{
		Auth:       auth,
		SmtpServer: c.SmtpServer,
		SmtpPort:   c.SmtpPort,
		FromEmail:  c.FromEmail,
		ToEmails:   c.ToEmails,
	}
}

func loadConfiguration(c *Configuration) error {
	configFile := "config.json"
	file, _ := os.Open(configFile)
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	if err != nil {
		return err
	}

	log.Print("Finished Loading Configuration from " + configFile)
	return nil
}

func (s *Service) HeartbeatHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func (s *Service) ComplimentHandler(w http.ResponseWriter, req *http.Request) {
	// Compliment API found at https://github.com/srobertson421/compliment-api
	complimentUrl := "https://compliment-api.herokuapp.com/"

	resp, err := http.Get(complimentUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Unable to retrieve compliment: %v", err))
	}
	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)

	compliment := doc.FirstChild.FirstChild.NextSibling.FirstChild.Data
	err = s.sendComplimentInEmail(compliment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Unable to send compliment: %v", err))
	}
	io.WriteString(w, "Wrote Compliment: "+compliment)
}

func (s *Service) sendComplimentInEmail(compliment string) error {
	if err := smtp.SendMail(s.SmtpServer+":"+s.SmtpPort, s.Auth, s.FromEmail, s.ToEmails, []byte(compliment)); err != nil {
		return err
	}

	return nil
}
