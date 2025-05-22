package durablelogs

import (
	"bufio"
	"durablelogs/durablelogs/pb"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	segmentPrefix = "dl-"
)

type DurableLogger struct {
	directory      string
	maxPerFile     int
	bufWriter      *bufio.Writer
	currentFile    *os.File
	mu             sync.Mutex
	currentFileNum int
}

func NewDLServer(directory string, maxPerFile int) *DurableLogger {

	if err := os.MkdirAll(directory, 0755); err != nil {
		panic(err)
	}

	files, err := os.ReadDir(directory)
	var file *os.File
	if err != nil {
		panic(err)
	}
	if files != nil {
		//TODO : IF wal files are present then read the most recent one and start from there
	} else {
		file, err = os.Create(directory + "/" + segmentPrefix + "0")
		if err != nil {
			panic(err)
		}
	}

	return &DurableLogger{
		directory:      directory,
		maxPerFile:     maxPerFile,
		bufWriter:      bufio.NewWriter(file),
		currentFile:    file,
		currentFileNum: 0,
	}
}

func (dl *DurableLogger) Log(message string) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	logEntry := &pb.Log{
		Log:       message,
		Timestamp: time.Now().String(),
	}

	marsheledLog := MustMarshal(logEntry)

	if len(dl.buffer) == dl.maxPerFile {
		dl.Flush()
	}
	fmt.Println(marsheledLog)
	// dl.buffer = append(dl.buffer, marsheledLog)

}

func (dl *DurableLogger) GetBufferedLogs() []string {
	return []string{}
}

func (dl *DurableLogger) GetCurrentFile() int {
	return dl.currentFile
}

func (dl *DurableLogger) Flush() {

}
