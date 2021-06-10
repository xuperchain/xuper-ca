/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/cmd/command"
	"github.com/xuperchain/xuper-ca/config"
	"github.com/xuperchain/xuper-ca/dao"
	"github.com/xuperchain/xuper-ca/server"
	"github.com/xuperchain/xuper-ca/util"
)

const defaultConfigFile = "./conf/caserver.yaml"

func newCaServerCommand() (*cobra.Command, error) {
	var configFile string

	caServerCmd := &cobra.Command{
		Use:   "ca-server",
		Short: "ca server, for node management of the different net",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetOutput(util.GetLogOut())
			level, _ := log.ParseLevel(config.GetLog().Level)
			log.SetLevel(level)

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
			defer signal.Stop(sigc)
			quit := make(chan int)
			server.Start(quit)
			for {
				select {
				case <-sigc:
					pprof.StopCPUProfile()
					return nil
				case <-quit:
					pprof.StopCPUProfile()
					return nil
				}
			}

			return nil
		},
	}

	caServerCmd.PersistentFlags().StringVar(&configFile, "config-file", defaultConfigFile, "CA Server configuration file")

	// 从配置文件中加载配置
	config.InstallCaConfig(configFile)

	dao.InitTables()

	return caServerCmd, nil
}

func runCaServer() error {
	rootCmd, err := newCaServerCommand()
	if err != nil {
		return err
	}

	rootCmd.AddCommand(command.NewAddNetCommand())
	rootCmd.AddCommand(command.NewAddNodeCommand())
	rootCmd.AddCommand(command.NewRevokeNodeCommand())
	rootCmd.AddCommand(command.NewInitCommand())
	rootCmd.AddCommand(command.NewDecryptByHdKeyCommand())

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
		FullTimestamp:   true,
	})

	return rootCmd.Execute()
}

func main() {
	// init config

	onError := func(err error) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// run
	if err := runCaServer(); err != nil {
		onError(err)
	}
}
