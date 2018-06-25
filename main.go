package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
)

const basePath = "/api"

type Configuration struct {
	Port       string `json:"Port,omitempty"`
	FromEmail  string
	SmtpServer string
	SmtpPort   string
}

type Recipients struct {
	Emails []string
}

type Service struct {
	Auth       smtp.Auth
	SmtpServer string
	SmtpPort   string
	FromEmail  string
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

	service, err := NewService(&c)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not create service: %v", err))
	}

	mux.HandleFunc(basePath+"/heartbeat", service.HeartbeatHandler)
	mux.HandleFunc(basePath+"/compliment", service.ComplimentHandler)
	log.Println("Starting Server on Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func NewService(c *Configuration) (*Service, error) {
	fmt.Print("Enter Email Account Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	// Pring newline because ReadPassword doesn't add a newline
	fmt.Println("")
	if err != nil {
		return nil, err
	}
	password := string(bytePassword)
	auth := smtp.PlainAuth("", c.FromEmail, password, c.SmtpServer)
	return &Service{
		Auth:       auth,
		SmtpServer: c.SmtpServer,
		SmtpPort:   c.SmtpPort,
		FromEmail:  c.FromEmail,
	}, nil
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

func loadRecipients(r *Recipients) error {
	configFile := "recipients.json"
	file, _ := os.Open(configFile)
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(r)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) HeartbeatHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func retrieveComplimentFromDB() {
	query := "SELECT rowid, compliment from compliments where used is NULL order by type desc, rowid asc limit 1;"
}

func retrieveComplimentFromApi() {
	// Compliment API found at https://github.com/srobertson421/compliment-api
	complimentUrl := "https://compliment-api.herokuapp.com/"

	resp, err := http.Get(complimentUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Unable to retrieve compliment: %v", err))
	}
	defer resp.Body.Close()

	doc, _ := html.Parse(resp.Body)

	return doc.FirstChild.FirstChild.NextSibling.FirstChild.Data
}

func (s *Service) ComplimentHandler(w http.ResponseWriter, req *http.Request) {
	compliment = retrieveComplimentFromDB()
	err = s.sendComplimentInEmail(compliment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("Unable to send compliment: %v", err))
	}

	io.WriteString(w, "Wrote Compliment: "+compliment)
}

func (s *Service) sendComplimentInEmail(compliment string) error {
	// Load the recipients at real time so we can change them without having to reboot the service
	var r Recipients
	if err := loadRecipients(&r); err != nil {
		return err
	}

	if err := smtp.SendMail(s.SmtpServer+":"+s.SmtpPort, s.Auth, s.FromEmail, r.Emails, []byte(compliment)); err != nil {
		return err
	}

	return nil
}
