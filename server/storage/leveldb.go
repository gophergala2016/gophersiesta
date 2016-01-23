package server
import (
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/jmhodges/levigo"
	"log"
	"fmt"
)

type LevelDB struct {
	Path 		string
	Db			*levigo.DB
	Ro 			*levigo.ReadOptions
	Wo 			*levigo.WriteOptions
}


func (this *LevelDB) Init() {
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	tmpDb, err := levigo.Open(this.Path, opts)
	if err != nil {
		log.Fatal(err)
	}

	this.Db = tmpDb
	this.Ro = levigo.NewReadOptions()
	this.Wo = levigo.NewWriteOptions()

}

func (this *LevelDB) SetOption(appname, label, variable, value string) {

	if label=="" {
		label = "default"
	}

	this.Db.Put(this.Wo, []byte(fmt.Sprintf("%s-%s-%s", appname, label, variable)), []byte(value))
}

func (this *LevelDB) GetOption(appname, label, variable string) string {

	if label=="" {
		label = "default"
	}

	data, err := this.Db.Get(this.Ro, []byte(fmt.Sprintf("%s-%s-%s", appname, label, variable)))
	if err != nil {
		return ""
	}

	return string(data)
}

func (this *LevelDB) GetOptions(appname, label string) map[string]string {

	if label=="" {
		label = "default"
	}

	options := make(map[string]string)
	keystr := fmt.Sprintf("%s-%s", appname, label)
	l := len(keystr)

	it := this.Db.NewIterator(this.Ro)
	defer it.Close()
	for it.Seek([]byte(keystr)); it.Valid(); it.Next() {
		if string(it.Key())[0:l] == keystr {
			options[string(it.Key())] = string(it.Value())
		} else {
			break;
		}
	}

	return options

}
