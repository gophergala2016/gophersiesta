package server
import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"bytes"
	"strings"
)

type BoltDb struct {
	Path string
	Name string
	PropsBucketName []byte
		 *bolt.DB
}

func (s *BoltDb) Init() {

	s.Name = "gophersiesta.DB"
	s.PropsBucketName = []byte("props")
	db, err := bolt.Open(s.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	s.DB = db

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.PropsBucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (s *BoltDb) SetOption(appName string, label string, variable string, value string) {



	s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.PropsBucketName)

		k := getKey(appName, label, variable)

		err := b.Put(k, []byte(value))
		return err
	})

}

func (s *BoltDb) GetOption(appName string, label string, variable string) string {

	var value string
	s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.PropsBucketName)

		k := getKey(appName, label, variable)

		v := b.Get(k)

		value = parseValue(v)

		return nil
	})

	return value

}

func parseValue(v []byte) string {

	if v != nil{

		return string(v[:len(v)])
	}

	return ""
}


func parseKey(v []byte) (appName string, label string, variable string, err error) {

	if v != nil{
		k := string(v[:len(v)])

		parts := strings.Split(k, "-")

		if (len(parts) != 3){
			return appName, label, variable, fmt.Errorf("The key is not structured like appName-labels-variable")
		}

		appName = parts[0]
		label = parts[1]
		variable = parts[2]

		return appName, label, variable, nil
	}

	return appName, label, variable, fmt.Errorf("value is nil")
}

func (s *BoltDb) GetOptions(appName, label string) map[string]string {

	props := make(map[string]string)

	prefix := []byte(fmt.Sprintf("%s-%s", appName, label))

	s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(s.PropsBucketName)

		c := b.Cursor()

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {

			_, _, variable, err := parseKey(k)

			if (err != nil){
				log.Print(err)
			}

			props[variable] = parseValue(v)

		}

		return nil
	})

	return props
}

func getKey(appName string, label string, variable string) []byte {

	if label == "" {
		label = "default"
	}

	return []byte(fmt.Sprintf("%s-%s-%s", appName, label, variable))
}


func (s *BoltDb) Close(){
	s.Close()
}
