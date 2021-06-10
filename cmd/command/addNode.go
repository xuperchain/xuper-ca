/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/service"
)

func NewAddNodeCommand() *cobra.Command {
	var address string
	var net string
	var adminAddress string
	addNodeCommand := &cobra.Command{
		Use:   "addNode",
		Short: "add a node for the net",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddNode(address, net, adminAddress)
		},
	}
	addNodeCommand.PersistentFlags().StringVar(&address, "Addr", "", "Address to add")
	addNodeCommand.PersistentFlags().StringVar(&adminAddress, "Admin", "", "Address for net admin")
	addNodeCommand.PersistentFlags().StringVar(&net, "Net", "", "the name of the net")

	return addNodeCommand
}

func runAddNode(address string, net, adminAddress string) error {
	if address == "" || net == "" || adminAddress == "" {
		fmt.Println("Please check params")
		return errors.New("params is not valid")
	}
	err := service.AddNode(net, adminAddress, address)
	if err != nil {
		fmt.Println("create node failed,", err)
	} else {
		fmt.Println("add node success")
	}
	return err
}
