package mock_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/blorticus/go-test-mocks"
)

func TestReader(t *testing.T) {
	reader := mock.NewReader().
		AddGoodRead([]byte("first good read")).
		AddEmptyRead().
		AddGoodRead([]byte("second good read\n now with a newline!")).
		AddGoodRead([]byte{1, 2, 3}).
		AddEOF()

	if err := expectOnRead(reader, 15, []byte("first good read"), false, false); err != nil {
		t.Errorf("(reader) (test 1: first read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, []byte{}, false, false); err != nil {
		t.Errorf("(reader) (test 1: second read) %s", err.Error())
	}

	if err := expectOnRead(reader, 37, []byte("second good read\n now with a newline!"), false, false); err != nil {
		t.Errorf("(reader) (test 1: third read) %s", err.Error())
	}

	if err := expectOnRead(reader, 3, []byte{1, 2, 3}, false, false); err != nil {
		t.Errorf("(reader) (test 1: fourth read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, nil, false, true); err != nil {
		t.Errorf("(reader) (test 1: fifth read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, nil, false, true); err != nil {
		t.Errorf("(reader) (test 1: sixth read) %s", err.Error())
	}

	reader = mock.NewReader().
		AddGoodRead([]byte("first good read")).
		AddEmptyRead().
		AddError(fmt.Errorf("error")).
		AddGoodRead([]byte("second good read\n now with a newline!")).
		AddEOF()

	if err := expectOnRead(reader, 15, []byte("first good read"), false, false); err != nil {
		t.Errorf("(reader) (test 2: first read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, []byte{}, false, false); err != nil {
		t.Errorf("(reader) (test 2: second read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, nil, true, false); err != nil {
		t.Errorf("(reader) (test 2: third read) %s", err.Error())
	}

	if err := expectOnRead(reader, 37, []byte("second good read\n now with a newline!"), false, false); err != nil {
		t.Errorf("(reader) (test 2: fourth read) %s", err.Error())
	}

	if err := expectOnRead(reader, 0, nil, false, true); err != nil {
		t.Errorf("(reader) (test 1: fifth read) %s", err.Error())
	}
}

var incomingBuffer []byte = make([]byte, 9000)

func expectOnRead(reader *mock.Reader, expectedNumberOfBytesRead int, expectedBytes []byte, expectAnError bool, expectEOF bool) error {
	numberOfBytesRead, err := reader.Read(incomingBuffer)

	if err != nil {
		if err == io.EOF {
			if !expectEOF {
				return fmt.Errorf("did not expect EOF, got EOF")
			}

			return nil
		} else if !expectAnError {
			return fmt.Errorf("expected no error, got error = (%s)", err.Error())
		}

		return nil
	} else if expectAnError {
		return fmt.Errorf("expected an error, got no error")
	}

	if numberOfBytesRead != expectedNumberOfBytesRead {
		return fmt.Errorf("expected number of bytes read = (%d), got (%d)", expectedNumberOfBytesRead, numberOfBytesRead)
	}

	if !bytes.Equal(incomingBuffer[:numberOfBytesRead], expectedBytes) {
		return fmt.Errorf("bytes do not match expected value")
	}

	return nil
}
