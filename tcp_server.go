package main

import (
	"encoding/csv"
	"fmt"
	"github.com/B4KO/is105sem03/mycrypt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}

					switch msg := string(mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)); msg {

					case "ping":
						log.Println("Dekrypter melding: ", msg)
						msg = string(mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, 4))
						_, err = c.Write([]byte("pong"))

					default:
						log.Println("Dekrypter melding: ", msg)

						if strings.HasPrefix(msg, "Kjevik;") {
							msg = formatString(msg)
							log.Println("Formaterer melding: ", msg)
						}

						msg = string(mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, 4))

						_, err = c.Write([]byte(msg))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}

func celsiusToFahrenheit(celsius int) int {
	return int(float64(celsius)*9/5 + 32)
}

func formatString(input string) string {

	r := csv.NewReader(strings.NewReader(input))
	r.Comma = ';'

	record, err := r.Read()
	if err != nil {
		fmt.Println("Error parsing the CSV string:", err)
		return "Invalid format"
	}

	// Get the last value and convert it from Celsius to Fahrenheit
	lastValue, err := strconv.Atoi(record[len(record)-1])
	if err != nil {
		fmt.Println("Error converting the last value to an integer:", err)
		return "Invalid format"
	}

	fahrenheit := celsiusToFahrenheit(lastValue)

	record[len(record)-1] = strconv.Itoa(fahrenheit)

	// Reconstruct the CSV string with the new value
	output := strings.Join(record, ";")

	return output
}
