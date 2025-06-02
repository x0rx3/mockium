package logging

import (
	"encoding/json"
	"fmt"
	"mockium/internal/model"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

type ProcessLogger struct {
	log         *zap.Logger
	mu          sync.Mutex
	currentFile *os.File
	maxSize     int64
	dirPath     string
	baseName    string
	fileIndex   int
	currentSize int64
}

func NewProcessLogger(log *zap.Logger, dirPath, baseName string, maxSizeMB int64) (*ProcessLogger, error) {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	maxSize := maxSizeMB * 1024 * 1024

	rl := &ProcessLogger{
		maxSize:  maxSize,
		dirPath:  dirPath,
		baseName: baseName,
	}

	if err := rl.rotate(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (inst *ProcessLogger) Log(logFields *model.ProcessLoggingFileds) {
	p, err := json.MarshalIndent(logFields, "", "	")
	if err != nil {
		inst.log.Error("marshal", zap.Error(err))
		return
	}

	inst.mu.Lock()
	defer inst.mu.Unlock()

	if inst.currentSize+int64(len(p)) > inst.maxSize {
		if err := inst.rotate(); err != nil {
			inst.log.Error("rotate process log file", zap.Error(err))
			return
		}
	}

	n, err := inst.currentFile.Write(p)
	if err != nil {
		inst.log.Error("write process log", zap.Error(err))
		return
	}
	inst.currentSize += int64(n)
}

func (inst *ProcessLogger) rotate() error {
	if inst.currentFile != nil {
		if err := inst.currentFile.Close(); err != nil {
			return err
		}
	}

	newIndex := 0
	for {
		logPath := filepath.Join(inst.dirPath, fmt.Sprintf("%s.%d.log", inst.baseName, newIndex))
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			break
		}
		newIndex++
	}
	inst.fileIndex = newIndex

	logPath := filepath.Join(inst.dirPath, fmt.Sprintf("%s.%d.log", inst.baseName, inst.fileIndex))
	f, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	inst.currentFile = f
	inst.currentSize = 0
	return nil
}

func (inst *ProcessLogger) Close() error {
	inst.mu.Lock()
	defer inst.mu.Unlock()

	if inst.currentFile != nil {
		return inst.currentFile.Close()
	}
	return nil
}
