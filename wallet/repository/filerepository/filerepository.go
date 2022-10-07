package filerepository

import (
	"fmt"
	"os"

	"github.com/mauroalderete/weasel/wallet"
	"github.com/mauroalderete/weasel/wallet/repository"
)

type fileWriter struct {
	file *os.File
}

func (f *fileWriter) Write(data []byte) (int, error) {

	n, err := fmt.Fprintf(f.file, "%s,\n", data)
	if err != nil {
		return 0, fmt.Errorf("failed to save in file: %v", err)
	}
	f.file.Sync()

	return n, err
}

type FileRepository struct {
	base    *repository.Base
	writter *fileWriter
}

func (r *FileRepository) Save(w wallet.Wallet) error {

	err := r.base.Save(w)
	if err != nil {
		return fmt.Errorf("fileRepository fail to save the wallet: %v", err)
	}

	return nil
}

func (r *FileRepository) Close() {
	if r.writter != nil && r.writter.file != nil {
		r.writter.file.Write([]byte("]"))

		if r.writter.file != nil {
			r.writter.file.Close()
		}
	}

	if r.base != nil {
		r.base.Close()
	}
}

func New(filepath string) (*FileRepository, error) {

	if len(filepath) == 0 {
		return nil, fmt.Errorf("filepath is required")
	}

	repo := FileRepository{}

	file, err := openFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("fail open file: %v", err)
	}

	fw := fileWriter{}
	fw.file = file

	base, err := repository.New(&fw)
	if err != nil {
		return nil, fmt.Errorf("failed to instance a base repository: %v", err)
	}

	repo.base = base
	repo.writter = &fw

	return &repo, nil
}

func openFile(filepath string) (*os.File, error) {
	var err error
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed open file '%s' to store the data: %v", filepath, err)
	}

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get info of the file '%s': %v", filepath, err)
	}

	if info.Size() == 0 {
		return file, nil
	}

	lastByte := make([]byte, 1)
	_, err = file.ReadAt(lastByte, info.Size()-1)
	if err != nil {
		return nil, fmt.Errorf("failed to get last byte of the file '%s': %v", filepath, err)
	}

	if lastByte[0] == ']' {
		file.Truncate(info.Size() - 1)
		file.Sync()
	}

	return file, nil
}
