// Package mock provides mock objects that can be used to make golang testing easier.
package mock

import "io"

type readerResult struct {
	bytesToCopyToReadBuffer []byte
	numberOfBytesRead       int
	err                     error
}

// Reader is a mock io.Reader.  It is primed with values that will be extracted with each
// successive Read().  It can be primed with a GoodRead (meaning no error is returned and
// at least some bytes are present), an EmptyRead (meaning no error is returned, but no
// bytes are also returned), an Error (meaning an error other than io.EOF is returned) or
// EOF (meaning io.EOF is returned).  Once EOF is reached, it will be returned on all
// subsequent calls to Read().
type Reader struct {
	readResults        []*readerResult
	nextResultToReturn int
	atEOF              bool
}

// NewReader create a new, empty mock Reader.
func NewReader() *Reader {
	return &Reader{
		readResults:        make([]*readerResult, 0, 20),
		nextResultToReturn: 0,
		atEOF:              false,
	}
}

// AddGoodRead adds a GoodRead to the queue of return values.
func (reader *Reader) AddGoodRead(bytesToCopyToCallerBuffer []byte) *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: bytesToCopyToCallerBuffer,
		numberOfBytesRead:       len(bytesToCopyToCallerBuffer),
		err:                     nil,
	})
	return reader
}

// AddEmptyRead adds an EmptyRead to the queue of return values.
func (reader *Reader) AddEmptyRead() *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     nil,
	})
	return reader
}

// AddError adds an Error to the queue of return values.
func (reader *Reader) AddError(err error) *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     err,
	})
	return reader
}

// AddEOF adds EOF to the queue of return values.
func (reader *Reader) AddEOF() *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     io.EOF,
	})

	return reader
}

// Read implements the io.Reader Read() method.
func (reader *Reader) Read(b []byte) (n int, err error) {
	if reader.atEOF || len(reader.readResults) == 0 {
		return 0, io.EOF
	}

	if reader.readResults[reader.nextResultToReturn].err == io.EOF {
		reader.atEOF = true
		return 0, io.EOF
	}

	nv := reader.readResults[reader.nextResultToReturn]
	reader.nextResultToReturn++

	if len(nv.bytesToCopyToReadBuffer) > 0 {
		copy(b[0:], nv.bytesToCopyToReadBuffer)
	}

	return nv.numberOfBytesRead, nv.err
}
