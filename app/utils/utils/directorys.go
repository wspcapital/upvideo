package utils


import(
    "os"
    _"path/filepath"
)


func file_is_exists(f string) bool {
    _, err := os.Stat(f)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil
}


func InitProject(){

	// check if directory needed exist.
	paths := [...]string{
			"tmp",
			"tmp/metas",
			"tmp/secrets",
			"tmp/temp_videos",
			"tmp/titles_videos",
			"tmp/tokens",
	}

	for _,path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
		    os.Mkdir(path, os.FileMode(0777))
		}
	}


	// check if files needed exist.

}



