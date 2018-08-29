/*----------------------------------------------------------------------------------
queue package contains a backup module which assures all orders being saved.
------------------------------------------------------------------------------------
CONTENT:
	runBackup: 		BacksUp the queues to a file
	saveToDisk:		saves the queue to a file on disk.
	loadFromDisk:	checks if the file is avaliable and returns nil if everything ok
----------------------------------------------------------------------------------*/
package queue

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"utilities"
)

//runBackup backs up the queues to a file
func runBackup(outgoingMsg chan<- utilities.Message) {
	const offlineFileName = "OfflineBackup"
	const onlineFileName = "OnlineBackup"

	var backup queue
	backup.loadFromDisk(offlineFileName)

	// Resend all orders found on loaded backup file:
	if !backup.isEmpty() {
		for f := 0; f < utilities.NFloors; f++ {
			for b := 0; b < utilities.NButtons; b++ {
				if backup.isOrder(f, b) {
					if b == utilities.BtnInside {
						AddLocalOrder(f, b)
					} else {
						outgoingMsg <- utilities.Message{Category: utilities.NewOrder, Floor: f, Button: b}
					}
				}
			}
		}
	}
	//go routine that waits for takeBackup channel to write to file.
	go func() {
		for {
			<-takeBackup
			if err := local.saveToDisk(offlineFileName); err != nil {
				log.Println(err)
			}
			if err := remote.saveToDisk(onlineFileName); err != nil {
				log.Println(err)
			}
		}
	}()
}

//saveToDisk saves the queue to a file on disk.
func (q *queue) saveToDisk(filename string) error {

	data, err := json.Marshal(&q)
	if err != nil {
		log.Println("json.Marshal() error: Failed to Marshal.")
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println("ioutil.WriteFile() error: Failed to backup.")
		return err
	}
	return nil
}

//loadFromDisk checks if the file is avaliable and returns nil if everything ok
func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println("Loading from existing backup file...")

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("loadFromDisk() error: failed to read file.")
		}
		if err := json.Unmarshal(data, q); err != nil {
			log.Println("loadFromDisk() error: failed to Unmarshal.")
		}
	}
	return nil
}
