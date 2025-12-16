package zilog

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

const (
	Error = "error" //错误日志
	Debug = "debug" //调试日志
	Info  = "info"  //信息日志
	Warn  = "warn"  //警告日志
	Gift  = "gift"  //礼物日志
)

type (
	// 单日志结构体
	logFileStruct struct {
		Lock     *sync.RWMutex
		File     *os.File
		OpenTime int64
	}

	// 日志结构体
	LogStruct struct {
		debugFile   logFileStruct
		infoFile    logFileStruct
		errorFile   logFileStruct
		warnFile    logFileStruct
		giftFile    logFileStruct
		logInfoFile logFileStruct
		level       string
		maxSize     int64
		maxBackups  int
		maxAge      int
		rotateTime  int64
	}
	LogInfoStruct struct {
		LogInfoContext map[string][]string
		Lock           *sync.RWMutex
	}
)

var (
	logInfo = LogInfoStruct{
		LogInfoContext: make(map[string][]string),
		Lock:           &sync.RWMutex{},
	}
)

// 日志结构体初始化
func (logS *LogStruct) Init(dataDir, level string, maxSize int64, maxBacks, maxAge int, rotateTime int64) {
	var err error
	if dataDir == "" {
		dataDir, err = os.Executable()
		if err != nil {
			panic("get log dir err: " + err.Error())

		}
		dataDir = filepath.Join(filepath.Dir(dataDir), "logs")
		if !pathExists(dataDir) {
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				panic("create log dir err: " + err.Error())
			}
		}
	}
	logS.rotateTime = rotateTime
	logS.level = level
	logS.maxSize = maxSize
	logS.maxBackups = maxBacks
	logS.maxAge = maxAge
	logS.debugFile.Lock = &sync.RWMutex{}
	logS.infoFile.Lock = &sync.RWMutex{}
	logS.errorFile.Lock = &sync.RWMutex{}
	logS.warnFile.Lock = &sync.RWMutex{}
	logS.giftFile.Lock = &sync.RWMutex{}
	logS.logInfoFile.Lock = &sync.RWMutex{}

	for _, label := range []string{"debug", "info", "error", "warn", "gift", "loginfo"} {
		if label == "debug" {
			if logS.level != "debug" {
				continue
			}
		}
		if label == "loginfo" {
			if logS.maxBackups < 1 {
				continue
			}
			if err := logS.open(label, filepath.Join(dataDir, "."+label)); err != nil {
				panic("log info init err: " + err.Error())

			}
			logS.getLogInfoContext()
			continue
		}
		if err := logS.open(label, filepath.Join(dataDir, label+".log")); err != nil {
			panic("log init err: " + err.Error())
		}
	}
	go logS.checkLogRotate()
}

// 获取上次日志内容
func (logS *LogStruct) getLogInfoContext() {
	logS.logInfoFile.Lock.Lock()
	defer logS.logInfoFile.Lock.Unlock()
	dataByte, err := io.ReadAll(logS.logInfoFile.File)
	if err != nil {
		logS.Error("read log info file err: "+err.Error(), true)
	}
	if len(dataByte) > 0 {
		if err := json.Unmarshal(dataByte, &logInfo.LogInfoContext); err != nil {
			logS.Error("read log info file err: "+err.Error(), true)
		}
	} else {
		logInfo.LogInfoContext = make(map[string][]string)
	}
}

// open 日志文件
func (logS *LogStruct) open(label, fileName string) error {
	var err error
	switch label {
	case "debug":
		if logS.debugFile.File == nil {
			logS.debugFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			logS.debugFile.OpenTime = time.Now().Unix()
		}
	case "info":
		if logS.infoFile.File == nil {
			logS.infoFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			logS.infoFile.OpenTime = time.Now().Unix()
		}
	case "error":
		if logS.errorFile.File == nil {
			logS.errorFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			logS.errorFile.OpenTime = time.Now().Unix()
		}
	case "warn":
		if logS.warnFile.File == nil {
			logS.warnFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			logS.warnFile.OpenTime = time.Now().Unix()
		}
	case "gift":
		if logS.giftFile.File == nil {
			logS.giftFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			logS.giftFile.OpenTime = time.Now().Unix()
		}
	case "loginfo":
		if logS.logInfoFile.File == nil {
			logS.logInfoFile.File, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			logS.logInfoFile.OpenTime = time.Now().Unix()
		}
	default:
		return nil
	}
	return nil
}

// Close 日志文件
func (logS *LogStruct) Closes() {
	logS.debugFile.Lock.Lock()
	if logS.debugFile.File != nil {
		logS.debugFile.File.Close()
	}
	logS.debugFile.Lock.Unlock()

	logS.infoFile.Lock.Lock()
	if logS.infoFile.File != nil {
		logS.infoFile.File.Close()
	}
	logS.infoFile.Lock.Unlock()

	logS.errorFile.Lock.Lock()
	if logS.errorFile.File != nil {
		logS.errorFile.File.Close()
	}
	logS.errorFile.Lock.Unlock()

	logS.warnFile.Lock.Lock()
	if logS.warnFile.File != nil {
		logS.warnFile.File.Close()
	}
	logS.warnFile.Lock.Unlock()
	logS.giftFile.Lock.Lock()
	if logS.giftFile.File != nil {
		logS.giftFile.File.Close()
	}
	logS.giftFile.Lock.Unlock()
}

// close 日志文件
func (logS *LogStruct) close(label string) {
	switch label {
	case "debug":
		if logS.debugFile.File != nil {
			logS.debugFile.File.Close()
			logS.debugFile.File = nil // 关闭后置空，避免重复关闭
		}
	case "info":
		if logS.infoFile.File != nil {
			logS.infoFile.File.Close()
			logS.infoFile.File = nil // 关闭后置空，避免重复关闭
		}
	case "error":
		if logS.errorFile.File != nil {
			logS.errorFile.File.Close()
			logS.errorFile.File = nil // 关闭后置空，避免重复关闭
		}
	case "warn":
		if logS.warnFile.File != nil {
			logS.warnFile.File.Close()
			logS.warnFile.File = nil // 关闭后置空，避免重复关闭
		}
	case "gift":
		if logS.giftFile.File != nil {
			logS.giftFile.File.Close()
			logS.giftFile.File = nil // 关闭后置空，避免重复关闭
		}
	}
}

// Write To debug
func (logS *LogStruct) debugWrite(data []byte) {
	logS.debugFile.Lock.Lock()
	if logS.debugFile.File != nil {
		logS.debugFile.File.Write(data)
	}
	logS.debugFile.Lock.Unlock()
}

// write写日志
func (logS *LogStruct) Write(label string, data string) {
	dataByte := []byte(data)
	if logS.level == "debug" {
		logS.debugWrite(dataByte)
	}
	switch label {
	case "debug":
		if logS.level != "debug" {
			logS.debugWrite(dataByte)
		}
	case "info":
		logS.infoFile.Lock.Lock()
		if logS.infoFile.File != nil {
			logS.infoFile.File.Write(dataByte)
		}
		logS.infoFile.Lock.Unlock()
	case "error":
		logS.errorFile.Lock.Lock()
		if logS.errorFile.File != nil {
			logS.errorFile.File.Write(dataByte)
		}
		logS.errorFile.Lock.Unlock()
	case "warn":
		logS.warnFile.Lock.Lock()
		if logS.warnFile.File != nil {
			logS.warnFile.File.Write(dataByte)
		}
		logS.warnFile.Lock.Unlock()
	case "gift":
		logS.giftFile.Lock.Lock()
		if logS.giftFile.File != nil {
			logS.giftFile.File.Write(dataByte)
		}
		logS.giftFile.Lock.Unlock()
	default:
		return
	}
}

// info日志写入
func (logS *LogStruct) Info(data string, debug bool) {
	timeFormat := time.Now().Format("2006-01-02 15:04:05")
	data = timeFormat + " " + data + "\n"
	if debug {
		fmt.Print(data)
	}
	go logS.Write(Info, data)
}

// error日志写入
func (logS *LogStruct) Error(data string, debug bool) {
	timeFormat := time.Now().Format("2006-01-02 15:04:05")
	data = timeFormat + " " + data + "\n"
	if debug {
		fmt.Print(data)
	}
	go logS.Write(Error, data)
}

// warn日志写入
func (logS *LogStruct) Warn(data string, debug bool) {
	timeFormat := time.Now().Format("2006-01-02 15:04:05")
	data = timeFormat + " " + data + "\n"
	if debug {
		fmt.Print(data)
	}
	go logS.Write(Warn, data)
}

// gift日志写入
func (logS *LogStruct) Gift(data string, debug bool) {
	timeFormat := time.Now().Format("2006-01-02 15:04:05")
	data = timeFormat + " " + data + "\n"
	if debug {
		fmt.Print(data)
	}
	go logS.Write(Gift, data)
}

// debug日志写入
func (logS *LogStruct) Debug(data string, debug bool) {
	timeFormat := time.Now().Format("2006-01-02 15:04:05")
	data = timeFormat + " " + data + "\n"
	if debug {
		fmt.Print(data)
	}
	go logS.Write(Debug, data)
}

// logtotate
func (logS *LogStruct) logRotate(label string) error {
	var (
		file *os.File
		lock *sync.RWMutex
	)
	switch label {
	case "debug":
		file = logS.debugFile.File
		lock = logS.debugFile.Lock
	case "info":
		file = logS.infoFile.File
		lock = logS.infoFile.Lock
	case "error":
		file = logS.errorFile.File
		lock = logS.errorFile.Lock
	case "warn":
		file = logS.warnFile.File
		lock = logS.warnFile.Lock
	case "gift":
		file = logS.giftFile.File
		lock = logS.giftFile.Lock
	default:
		return nil
	}
	fileName := file.Name()
	lock.Lock()
	defer lock.Unlock()
	logS.close(label)
	bakFileName := fileName + "." + time.Now().Format("2006_01_02_15_04_05")
	os.Rename(fileName, bakFileName)
	go logS.compressLogFile(bakFileName, label)
	return logS.open(label, fileName)
}

// 获取日志大小
func (logS *LogStruct) getLogSize(label string) (int64, error) {
	var file *os.File
	switch label {
	case "debug":
		file = logS.debugFile.File
	case "info":
		file = logS.infoFile.File
	case "error":
		file = logS.errorFile.File
	case "warn":
		file = logS.warnFile.File
	case "gift":
		file = logS.giftFile.File
	default:
		return 0, nil
	}
	if file == nil {
		return 0, nil
	}
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// 获取日志打开时间
func (logS *LogStruct) getLogOpenTime(label string) int64 {
	switch label {
	case "debug":
		return logS.debugFile.OpenTime
	case "info":
		return logS.infoFile.OpenTime
	case "error":
		return logS.errorFile.OpenTime
	case "warn":
		return logS.warnFile.OpenTime
	case "gift":
		return logS.giftFile.OpenTime
	default:
		return 0
	}
}

// 检查日志文件是否需要轮转
func (logS *LogStruct) checkLogRotate() {
	if logS.maxSize == 0 && logS.maxAge == 0 {
		return
	}
	t := time.NewTicker(30 * time.Minute)
	for {
		<-t.C
		for _, label := range []string{"debug", "info", "error", "warn", "gift"} {
			if label == "debug" {
				if logS.level != "debug" {
					continue
				}
			}
			logS.timeLogRorate(label)
		}
	}
}

// size轮转
func (logS *LogStruct) sizeLogRorate(label string) {
	size, err := logS.getLogSize(label)
	if err != nil {
		logS.Write("error", "get log size err: "+err.Error())
		return
	}
	if size <= 0 {
		return
	}
	if size >= int64(logS.maxSize) {
		if err := logS.logRotate(label); err != nil {
			logS.Write("error", label+" log rotate err: "+err.Error()+"\n")
		} else {
			logS.Write(label, label+" log file rotated"+"\n")
		}
	}
}

// 时间轮转
func (logS *LogStruct) timeLogRorate(label string) {
	// 测试更改时间设置
	if time.Now().Unix()-logS.getLogOpenTime(label) >= int64(logS.maxAge)*logS.rotateTime {
		if err := logS.logRotate(label); err != nil {
			logS.Write("error", label+" log rotate err: "+err.Error()+"\n")
		} else {
			logS.Write(label, label+" log file rotated"+"\n")
		}
		return
	}
	logS.sizeLogRorate(label)
}

// 压缩日志文件
func (logS *LogStruct) compressLogFile(filename, label string) error {
	zipFileName := filename + ".zip"
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 打开日志文件
	logFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// 获取文件信息
	info, err := logFile.Stat()
	if err != nil {
		return err
	}

	// 创建一个新的 ZIP 文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// 确保文件路径正确
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	// 创建一个新的 ZIP 文件条目
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 将日志文件内容复制到 ZIP 文件中
	if _, err = io.Copy(writer, logFile); err != nil {
		return err
	}
	logFile.Close()
	os.Remove(filename)
	logS.logInfoFile.Lock.Lock()
	defer logS.logInfoFile.Lock.Unlock()
	logInfo.LogInfoContext[label] = append(logInfo.LogInfoContext[label], zipFileName)
	lenCount := len(logInfo.LogInfoContext[label])
	for {
		if lenCount > logS.maxBackups && lenCount >= 1 {
			if err := os.Remove(logInfo.LogInfoContext[label][0]); err != nil {
				break
			}
			logInfo.LogInfoContext[label] = slices.Delete(logInfo.LogInfoContext[label], 0, 1)
			lenCount--
		} else {
			break
		}
	}
	// 更新日志信息文件
	if err := logS.logInfoFile.File.Truncate(0); err != nil {
		return fmt.Errorf("truncate log info file err: %v", err)
	}
	dataByte, err := json.Marshal(logInfo.LogInfoContext)
	if err != nil {
		return fmt.Errorf("marshal log info err: %v", err)
	}
	_, err = logS.logInfoFile.File.Write(dataByte)
	if err != nil {
		return fmt.Errorf("write log info file err: %v", err)
	}
	return nil
}

// 判断文件是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
