package module

import (
	"github.com/04Akaps/block-event/repository"
)

type Writer struct {
	repository *repository.Repository
	writerChan <-chan *WriterChan
}

func NewWriter(
	repository *repository.Repository,
	writerChan <-chan *WriterChan,
) *Writer {
	w := &Writer{repository: repository, writerChan: writerChan}

	return w
}

func (w *Writer) LookingEvent() {
	for {
		event := w.writerChan

	}
}
