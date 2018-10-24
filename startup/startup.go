package startup

import (
	"github.com/brharrelldev/jag/startup/cache"
	"log"
	"os"
	"path/filepath"
)

// This seem pretty confusing, need to rethink this approach.
func Check(dir, fileName string) (bool, error) {

	fullPath := filepath.Join(dir, fileName)

	var f *os.File

	//Has jag run before?  Check if db file exists and metadata
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println("Could not create dir ", err)
			return true, nil
		}
		if f, err = os.Create(fullPath); err != nil {
			log.Fatal("could not create file ", err)
			return true, err

		}

		defer f.Close()

		fr := cache.FTC{}

		ftc, err := fr.New(fullPath)
		if err != nil {
			log.Println("Could not instantiate due to error ", err)
			return true, nil
		}

		if err = ftc.Write(); err != nil {
			log.Println("could not write to database ", err)
			return true, nil
		}

		return true, err

	}

	return false, nil
}
