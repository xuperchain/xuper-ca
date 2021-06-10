/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/config"
	"github.com/xuperchain/xuper-ca/service"
)

func NewInitCommand() *cobra.Command {
	var tlsPath string
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "init ca, regenerate the ca cert",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(tlsPath)
		},
	}

	initCommand.PersistentFlags().StringVar(&tlsPath, "tlsPath", "", "the cert of path")

	return initCommand
}

func runInit(tlsPath string) error {
	if tlsPath == "" {
		tlsPath = config.GetCertPath()
	}

	// 创建根证书
	cert, err := service.GenerateCert(nil, "root", false)
	if err != nil {
		fmt.Println("cant init root cert", err)
		return err
	}

	// 写入文件
	err = service.WriteCert(tlsPath, cert)
	if err != nil {
		fmt.Println("write root cert failed, ", err)
		return err
	}
	fmt.Println("init root cert success")
	return nil
}
