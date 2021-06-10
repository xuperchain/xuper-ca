/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/service"
)

func NewAddNetCommand() *cobra.Command {
	var address string
	var net string

	addNetCommand := &cobra.Command{
		Use:   "addNet",
		Short: "add a net with a net admin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddNet(address, net)
		},
	}
	addNetCommand.PersistentFlags().StringVar(&address, "Addr", "", "Address for net admin")
	addNetCommand.PersistentFlags().StringVar(&net, "Net", "", "the name of the net")

	return addNetCommand
}

func runAddNet(address string, net string) error {
	err := service.AddNetAdmin(net, address)
	if err != nil {
		fmt.Println("create net admin failed,", err)
	} else {
		fmt.Println("add net admin success")
	}
	return err
}
