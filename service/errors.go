/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package service

import "errors"

var (
	ErrNetExisted    = errors.New("net has been existed")
	ErrParam         = errors.New("params is illegal")
	ErrDB            = errors.New("DB operate failed")
	ErrCACert        = errors.New("can not get ca cert")
	ErrCreateCert    = errors.New("create cert failed")
	ErrCertNoExisted = errors.New("cert is not existed")
)
