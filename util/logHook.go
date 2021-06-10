/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package util

import (
	"io"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/xuperchain/xuper-ca/config"
)

type LogIdHook struct {
	LogId string
}

func NewlogIdHook(logId string) logrus.Hook {
	hook := LogIdHook{
		LogId: logId,
	}
	return &hook
}

func (hook *LogIdHook) Fire(entry *logrus.Entry) error {
	entry.Data["logId"] = hook.LogId
	return nil
}

func (hook *LogIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func GetLogId() string {
	logid, _ := uuid.NewV1()
	return logid.String()
}

func GetLogOut() io.Writer {
	path := config.GetLog().Path
	if path == "" {
		return os.Stderr
	}
	if strings.LastIndex(path, "/") != len([]rune(path))-1 {
		path = path + "/"
	}
	file, err := os.Create(path + "caserver.log")
	if err != nil {
		return os.Stderr
	}

	return file
}
