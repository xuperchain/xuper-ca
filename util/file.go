/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package util

import (
	"io/ioutil"
)

/**
 * 生成文件
 */
func WriteFileUsingFilename(filename string, content []byte) error {
	//	//判断文件是否存在
	//	if checkFileIsExist(filename) {
	//		//打开文件
	//		f, err = os.OpenFile(filename, os.O_TRUNC, 0666)
	//		log.Printf("File [%v] exist", filename)
	//	} else {
	//		//创建文件
	//		f, err = os.Create(filename)
	//		log.Printf("File [%v] does not exist", filename)
	//	}
	//
	//	if err != nil {
	//		return err
	//	}
	//	var data = []byte(content)
	//函数向filename指定的文件中写入数据(字节数组)。如果文件不存在将按给出的权限创建文件，否则在写入数据之前清空文件。
	err := ioutil.WriteFile(filename, content, 0666)
	return err
}
