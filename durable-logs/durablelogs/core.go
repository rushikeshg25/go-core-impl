package durablelogs

import (
	"bufio"
	"durablelogs/durablelogs/pb"
	"encoding/binary"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	segmentPrefix = "dl-"
)

type DurableLogger struct {
	directory         string
	maxPerFile        int
	bufWriter         *bufio.Writer
	currentFile       *os.File
	mu                sync.Mutex
	currentFileNum    int
	logsInCurrentFile int
}

func NewDLServer(directory string, maxPerFile int) *DurableLogger {

	dl := &DurableLogger{
		directory:         directory,
		maxPerFile:        maxPerFile,
		logsInCurrentFile: 0,
		bufWriter:         nil,
		currentFile:       nil,
		currentFileNum:    0,
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		panic(err)
	}

	files, err := os.ReadDir(directory)
	var file *os.File
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		file, err = os.Create(directory + "/" + segmentPrefix + "0")
		if err != nil {
			panic(err)
		}
	} else {
		latestFile := files[len(files)-1]
		file, err = os.OpenFile(directory+"/"+latestFile.Name(), os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		lastChar := latestFile.Name()[len(latestFile.Name())-1:]
		num, err := strconv.Atoi(lastChar)
		if err != nil {
			panic(err)
		}
		dl.logsInCurrentFile = dl.maxPerFile
		dl.currentFileNum = num
		dl.currentFile = file
	}

	dl.bufWriter = bufio.NewWriter(file)

	return dl
}

func (dl *DurableLogger) Log(message string) {
	var err error
	dl.mu.Lock()

	logEntry := &pb.Log{
		Log:       message,
		Timestamp: time.Now().String(),
	}

	marsheledLog := MustMarshal(logEntry)
	marsheledLogSize := len(marsheledLog)

	err = binary.Write(dl.bufWriter, binary.LittleEndian, uint32(marsheledLogSize))
	if err != nil {
		panic(err)
	}

	err = binary.Write(dl.bufWriter, binary.LittleEndian, marsheledLog)
	if err != nil {
		panic(err)
	}

	dl.logsInCurrentFile++
	dl.mu.Unlock()

	if dl.logsInCurrentFile == dl.maxPerFile {
		dl.Flush()
		dl.NewFile()
	}

}

func (dl *DurableLogger) GetBufferedLogs() []string {
	return []string{}
}

func (dl *DurableLogger) GetCurrentFile() int {
	return dl.currentFileNum
}

func (dl *DurableLogger) Flush() {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	dl.bufWriter.Flush()
}

func (dl *DurableLogger) Close() {
	dl.Flush()
	dl.currentFile.Close()
}

func (dl *DurableLogger) NewFile() {
	dl.currentFile.Close()
	dl.logsInCurrentFile = 0
	dl.currentFileNum++
	file, err := os.Create(dl.directory + "/" + segmentPrefix + strconv.Itoa(dl.currentFileNum))
	if err != nil {
		panic(err)
	}
	dl.currentFile = file
	dl.bufWriter = bufio.NewWriter(file)
}

func MustMarshal(v interface{})
