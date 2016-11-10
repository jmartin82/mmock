package persist

import (
    "log"
    
    "github.com/jmartin82/mmock/definition"
)

//FileBodyPersister persists body in file
type FileBodyPersister struct {
}

//Persist the body of the response to fiel if needed
func (fbp FileBodyPersister) Persist(per *definition.Persist, res *definition.Response) {
	log.Printf("Body : %s\n", res.Body)
}