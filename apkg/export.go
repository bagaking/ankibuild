package apkg

import (
	"io"
	"os"
	"path/filepath"

	"archive/zip"
)

// ExportToAPKG exports the database and associated files as an .apkg file.
func (p *PkgInfo) ExportToAPKG(exportPath string) error {
	// Close the database to ensure no pending operations
	p.Close()

	apkgName := exportPath //filepath.Join(exportPath, "collection.apkg")
	// Create a new zip file
	zipFile, err := os.Create(apkgName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	collectionPth := filepath.Join(p.Path, "collection.anki2")
	// Add the database file to the zip
	if err = addFileToZip(archive, collectionPth); err != nil {
		return err
	}

	if err = os.Remove(collectionPth); err != nil {
		return err
	}

	// Add other necessary files (media files, etc) to zip here
	// ...

	// Change the file extension to .apkg if necessary
	// os.Rename(zipFileName, apkgName)

	return nil
}

// addFileToZip adds a single file to the specified zip archive.
func addFileToZip(archive *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file information
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate // The compression algorithm

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	if _, err := io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}
