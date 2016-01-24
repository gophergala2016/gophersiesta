package storage

import (
	"bytes"
	"fmt"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/boltdb/bolt"
	"log"
	"strings"
)

type BoltDb struct {
	Path             string
	Name             string
	PropsBucketName  []byte
	LabelsBucketName []byte
	*bolt.DB
}

func (s *BoltDb) Init() {

	s.Name = "gophersiesta.DB"
	s.PropsBucketName = []byte("props")
	s.LabelsBucketName = []byte("labels")

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

		_, err = tx.CreateBucketIfNotExists(s.LabelsBucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (s *BoltDb) GetLabels(appName string) []string {
	lbls := make([]string, 0)

	prefix := []byte(appName)

	s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(s.LabelsBucketName)

		c := b.Cursor()

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {

			_, _, err := parseLabelKey(k)

			if err != nil {
				log.Print(err)
			} else {
				lbls = append(lbls, parseValue(v))
			}

		}

		return nil
	})

	return lbls
}

func (s *BoltDb) SetOption(appName string, label string, variable string, value string) {

	s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.PropsBucketName)
		bl := tx.Bucket(s.LabelsBucketName)

		k := getPropertyKey(appName, label, variable)
		lk := getLabelKey(appName, label)

		err := b.Put(k, []byte(value))
		err = bl.Put(lk, []byte(getLabel(label)))

		return err
	})

}

func (s *BoltDb) GetOption(appName string, label string, variable string) string {

	var value string
	s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.PropsBucketName)

		k := getPropertyKey(appName, label, variable)

		v := b.Get(k)

		value = parseValue(v)

		return nil
	})

	return value

}

func parseValue(v []byte) string {

	if v != nil {

		return string(v[:len(v)])
	}

	return ""
}

func parsePropertyKey(v []byte) (appName string, label string, variable string, err error) {

	if v != nil {
		k := string(v[:len(v)])

		parts := strings.Split(k, "-")

		if len(parts) != 3 {
			return appName, label, variable, fmt.Errorf("The key is not structured like appName-labels-variable\n")
		}

		appName = parts[0]
		label = parts[1]
		variable = parts[2]

		return appName, label, variable, nil
	}

	return appName, label, variable, fmt.Errorf("value is nil")
}

func parseLabelKey(v []byte) (appName string, label string, err error) {

	if v != nil {
		k := string(v[:len(v)])

		parts := strings.Split(k, "-")

		if len(parts) != 2 {
			return appName, label, fmt.Errorf("The key %s is not structured like appName-labels\n", k)
		}

		appName = parts[0]
		label = parts[1]

		return appName, label, nil
	}

	return appName, label, fmt.Errorf("value is nil")
}

func (s *BoltDb) GetOptions(appName, label string) map[string]string {

	props := make(map[string]string)

	prefix := getLabelKey(appName, label)

	s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(s.PropsBucketName)

		c := b.Cursor()

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {

			_, _, variable, err := parsePropertyKey(k)

			if err != nil {
				log.Print(err)
			} else {
				props[variable] = parseValue(v)
			}

		}

		return nil
	})

	return props
}

func getLabel(label string) string {
	if label == "" {
		label = "default"
	}

	return label
}

func getPropertyKey(appName string, label string, variable string) []byte {

	return []byte(fmt.Sprintf("%s-%s-%s", appName, getLabel(label), variable))
}

func getLabelKey(appName string, label string) []byte {

	return []byte(fmt.Sprintf("%s-%s", appName, getLabel(label)))
}

func (s *BoltDb) Close() {
	s.Close()
}
