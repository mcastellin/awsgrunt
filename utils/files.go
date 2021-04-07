package utils

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

func ReadTemplateBodyFromFile(path string) (*string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := string(data)
	return &template, nil
}

func CreateZipFile(source string, fileName string, destDir string) (destination *string, err error) {
	if destDir == "" {
		destDir = os.TempDir()
	} else {
		//make sure dest directory exists
		_ = os.Mkdir(destDir, fs.ModePerm)
	}
	destPath := fmt.Sprintf("%s/%s.zip", destDir, fileName)

	os.Remove(destPath) // make sure we are replacing the file
	outFile, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	addFiles(w, source, "")

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return &destPath, nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}
