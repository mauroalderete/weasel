package stdoutrepository

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mauroalderete/weasel/wallet"
	"github.com/mauroalderete/weasel/wallet/repository"
)

type StdoutRepository struct {
	base *repository.Base
}

func (s *StdoutRepository) Save(w wallet.Wallet) error {

	err := s.base.Save(w)
	if err != nil {
		return fmt.Errorf("stdout repository fail to save the wallet: %v", err)
	}

	return nil
}

func (s *StdoutRepository) Close() {

}

func New(foregroundcolor int) (*StdoutRepository, error) {

	mw := &maskWriter{
		writer:  os.Stdout,
		suffix:  []byte("\n"),
		preffix: []byte(""),
		tint:    color.New(color.Attribute(foregroundcolor)),
	}

	base, err := repository.New(mw)
	if err != nil {
		return nil, fmt.Errorf("failed to instance a base repository: %v", err)
	}

	repo := StdoutRepository{}
	repo.base = base

	return &repo, nil
}

type maskWriter struct {
	writer  io.Writer
	suffix  []byte
	preffix []byte
	tint    *color.Color
}

func (m *maskWriter) Write(data []byte) (int, error) {

	n, err := m.tint.Fprintf(m.writer, "%s%s%s", m.preffix, data, m.suffix)
	if err != nil {
		return 0, fmt.Errorf("failed to save in file: %v", err)
	}

	return n, err
}
