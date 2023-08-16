package main

import (
	"backup-db/mysql"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func main() {
	log.Println("Beginning DB backup script")
	ch := make(chan string)

	// Non-blocking goroutine. Blocking would be ok here but this is for practice
	go func() {
		wg := sync.WaitGroup{}

		for _, val := range mysql.C {
			wg.Add(1)

			go mysql.DumpDB(val, ch, &wg)
		}

		wg.Wait()

		close(ch)
	}()

	for {
		result, ok := <-ch

		if !ok {
			log.Println("All DBs backup processes complete")
			break
		}

		log.Println("Output: ", result)
	}
}

func init() {
	viper.AddConfigPath("./setup")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.UnmarshalKey("databases", &mysql.C)
	if err != nil {
		log.Fatal(err)
	}
}
