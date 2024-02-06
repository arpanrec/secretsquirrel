package keyvalue

import (
	"encoding/json"
	"errors"
	"github.com/arpanrec/secretsquirrel/internal/physical"
	"log"
	"net/http"
)

func ReadWriteFilesFromURL(data *[]byte, operation string, key *string) (*physical.KVData, error) {

	switch operation {

	case http.MethodGet:
		log.Println("http.MethodGet for KeyValue called " + *key)
		d, err := physical.Get(key, nil)
		if err != nil {
			log.Println("Error while getting data: ", err)
			return nil, err
		}
		return d, nil
	case http.MethodPost:
		var kvData physical.KVData
		errUnmarshal := json.Unmarshal(*data, &kvData)
		if errUnmarshal != nil {
			log.Println("Error while unmarshalling data from request: ", errUnmarshal)
			return nil, errUnmarshal
		}
		err := physical.Save(key, &kvData)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("unsupported Method")
	}
}
