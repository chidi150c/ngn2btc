package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"user-apiv2/apiapp"
	"user-apiv2/apichart"
	"user-apiv2/apichat"
	"user-apiv2/apiexch"
	"user-apiv2/apiuser"
	"user-apiv2/bolt"
	"user-apiv2/memory"

	boltdb "github.com/coreos/bbolt"
)

func main() {
	//Get the port number from the environmental variable
	addr := os.Getenv("PORT")
	//Because the handler is served in a goroutine, so the main function process will complete its process
	//and close while the goroutine is stil working. we must hang the main function process after the call
	//to Open() by waiting on a signal using the combination of 'sigs' os.signal channel and 'done' bool channel.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	//Database initialization Starts...
	//Initializing bolt DB
	db, err := boltdb.Open("bolt/data.db", 0600, &boltdb.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initializing memory DB.
	var memoryDB = make(memory.MDBType)
	_ = memoryDB
	//Database initialization Ends...

	//Model Sessioner Implementor Initiallization
	//Select between boltdb and memoryDB by commenting and uncomemting appropriately
	ms := bolt.NewSession(db)
	//ms := memory.NewSession(memoryDB)

	//Cache initialization
	var dbUser = make(apiuser.UDBType)
	var dbExchData = make(apiexch.ExDBType)
	var dbExchUser = make(apiexch.UDBType)
	var dbChatUser = make(apichat.UDBType)
	var dbChannel = make(apichat.ChDBType)

	//Service initialization
	as := apichart.NewSession(dbUser, ms)
	//initialize the handlers exch, chat and user handlers
	eh := apiexch.NewExchHandler(dbExchUser, dbExchData, ms)
	uh := apiuser.NewUserHandler(dbUser, ms)
	ch := apichat.NewChatHandler(dbChatUser, dbChannel, ms)
	//ah := amchart.NewChartHandler(as)
	//_ = uh.Session.Userservice.AddUser(&data.User{Username: "chidi", Password: "cc", Level: "Admin"})
	go as.Graphservice.GraphPointFromCoinmkt(done, sigs)
	go as.Graphservice.PopulateChartData()
	go ch.ChatManager()
	//initialize the outer handler wrapper for exch, chat and user handlers , ah
	h := apiapp.NewHandler(uh, eh, ch)
	//open apiuser server
	server := apiapp.NewServer(addr, h)
	//Start the webserver
	if err := server.Open(done, sigs); err != nil {
		log.Fatalf("Unable to Open Server for listen and serve: %v", err)
	}
	fmt.Println("Listening on: ", server.Port())
	<-done
	<-done
	<-done
	fmt.Println("exiting")
}
