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


func InitProject(tmpPath string) {

	// check if directory needed exist.
	paths := [...]string{
			tmpPath,
			tmpPath + "/metas",
			tmpPath + "/secrets",
			tmpPath + "/temp_videos",
			tmpPath + "/titles_videos",
			tmpPath + "/tokens",
	}

	for _,path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
		    os.Mkdir(path, os.FileMode(0777))
		}
	}


	// check if files needed exist.

}



