/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package server

import "errors"

var (
	ErrAuth           = errors.New("address has no authority")
	ErrSign           = errors.New("sign is error")
	ErrAddNetAndAdmin = errors.New("add net and netAdmin failed")
	ErrAddNode        = errors.New("add node failed")
)
