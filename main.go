package main

import (
	"bytes"
	"challenge_9/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/generate", Generate_data)

	server := new(http.Server)
	server.Addr = ":9000"

	fmt.Println("Server Listering in the port", server.Addr)
	server.ListenAndServe()
}

func Generate_data(w http.ResponseWriter, r *http.Request) {

	var dataUnmarshal entity.Data
	iteration := 1

	go func() {
		for true {
			body, err := generate_wind_and_water()
			if err != nil {
				log.Fatalln(err)
				return
			}

			err = json.Unmarshal([]byte(body), &dataUnmarshal)
			if err != nil {
				log.Fatalln(err)
				return
			}

			responseWater := ""
			responseWind := ""

			switch {
			case dataUnmarshal.Water <= 5:
				responseWater = "Air Aman"
			case dataUnmarshal.Water < 8:
				responseWater = "Air Siaga"
			case dataUnmarshal.Water >= 9:
				responseWater = "Air Bahaya"
			}

			switch {
			case dataUnmarshal.Water <= 6:
				responseWind = "Angin Aman"
			case dataUnmarshal.Water < 15:
				responseWind = "Angin Siaga"
			case dataUnmarshal.Water >= 16:
				responseWind = "Angin Bahaya"
			}
			fmt.Println("number of iteration:", iteration)
			fmt.Println("time: ", time.Now().Local())

			// fmt.Println(dataUnmarshal.Wind)
			fmt.Println(string(body))
			fmt.Println(responseWater)
			fmt.Println(responseWind)
			fmt.Println()
			time.Sleep(15 * time.Second)

			iteration++
		}
	}()

	time.Sleep(15 * time.Second)
}

func generate_wind_and_water() ([]byte, error) {
	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 100
	waterValue := rand.Intn(max-min+1) + min
	windValue := rand.Intn(max-min+1) + min

	data := map[string]interface{}{
		"water": waterValue,
		"wind":  windValue,
	}

	requestJSON, err := json.Marshal(data)
	client := &http.Client{}

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://jsonplaceholder.typicode.com/posts", bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return body, nil
}
