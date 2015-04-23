package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type Settings struct {
	Benutzername string
	Passwort     string
	Server       string
}

/*
	Appends timestaps + Log MSGs into a file
*/
func writeLog(text string) {
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0600)
	// f.. microsoft
	f.WriteString(time.Now().Format("Mon Jan _2 15:04:05 2006") + " " + text + "\r\n")
	f.Sync()

	defer f.Close()
}

/*
	Attempts SMTP Login at given server
*/
func testSmtpPlain(cfg Settings) (success bool, err error) {
	auth := smtp.PlainAuth(
		"",
		cfg.Benutzername,
		cfg.Passwort,
		cfg.Server,
	)
	success = true
	// try to auth
	logTxt := "[ Info ] Testing SMTP Auth Plain: " + cfg.Server + ":25"
	log.Println(logTxt)
	writeLog(logTxt)

	conn, err := smtp.Dial(cfg.Server + ":25")
	if err != nil {
		success = false
		err = errors.New("[ Fatal ] SMTP Auth Plain: " + err.Error())
		writeLog(err.Error())

		return success, err
	}

	err = conn.Auth(auth)
	if err != nil {
		success = false
		err = errors.New("[ Fatal ] SMTP Auth Plain: " + err.Error())
		writeLog(err.Error())

		return success, err
	}
	defer conn.Close()

	return success, err
}

/*
	Wait for user interaction to finish execution
	and return the correct exit code
*/
func endProgram(code int) {
	logTxt := "[ Info ] Tests finished."
	log.Println(logTxt)
	writeLog(logTxt)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press Enter to exit.")
	exit, _ := reader.ReadString('\n')
	fmt.Print(exit)
	open.Run(logFile)
	os.Exit(code)
}

var (
	logFile = "scan.txt"
)

func init() {
	/*
		Create a logfile for later
	*/
	F, err := os.Create(logFile)
	if err != nil {
		log.Println("[ Fatal ] Failed to create logfile.")
		endProgram(1)
	}
	defer F.Close()
	writeLog("Bitte senden Sie die folgenden Daten an unseren Support:\r\n\r\n")

	log.Println("[ Info ] Mailserver check started.")
}

func main() {
	/*
		Initialize the program
		Read Settings from the cmdline
	*/
	var cfg Settings
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ihre E-Mail Adresse: ")
	cfg.Benutzername, _ = reader.ReadString('\n')
	cfg.Benutzername = strings.TrimSpace(cfg.Benutzername)
	fmt.Print("Ihr Kennwort: ")
	cfg.Passwort = string(gopass.GetPasswd())
	cfg.Passwort = strings.TrimSpace(cfg.Passwort)
	fmt.Print("Mailserver: ")
	cfg.Server, _ = reader.ReadString('\n')
	cfg.Server = strings.TrimSpace(cfg.Server)

	logTxt := "[ Info ] Beginning Tests.."
	log.Println(logTxt)
	writeLog(logTxt)

	/*
		perform Smtp Plain login check
		ToDo: Perform more checks
	*/
	_, err := testSmtpPlain(cfg)
	if err != nil {
		log.Println(err.Error())
		endProgram(1)
	}

	logTxt = "[ Ok ] Connection to " + cfg.Server + " looks fine"
	log.Println(logTxt)
	writeLog(logTxt)

	endProgram(0)
}
