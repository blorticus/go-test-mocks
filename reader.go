package mock

import "io"

type readerResult struct {
	bytesToCopyToReadBuffer []byte
	numberOfBytesRead       int
	err                     error
}

type Reader struct {
	readResults        []*readerResult
	nextResultToReturn int
	atEOF              bool
}

func NewReader() *Reader {
	return &Reader{
		readResults:        make([]*readerResult, 0, 20),
		nextResultToReturn: 0,
		atEOF:              false,
	}
}

func (reader *Reader) AddGoodRead(bytesToCopyToCallerBuffer []byte) *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: bytesToCopyToCallerBuffer,
		numberOfBytesRead:       len(bytesToCopyToCallerBuffer),
		err:                     nil,
	})
	return reader
}

func (reader *Reader) AddEmptyRead() *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     nil,
	})
	return reader
}

func (reader *Reader) AddError(err error) *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     err,
	})
	return reader
}

func (reader *Reader) AddEOF() *Reader {
	reader.readResults = append(reader.readResults, &readerResult{
		bytesToCopyToReadBuffer: []byte{},
		numberOfBytesRead:       0,
		err:                     io.EOF,
	})

	return reader
}

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
