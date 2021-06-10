/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/service"
)

func NewRevokeNodeCommand() *cobra.Command {
	var address string
	var net string

	revokeNodeCommand := &cobra.Command{
		Use:   "revokeNode",
		Short: "revoke the node for the net",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevokeNode(address, net)
		},
	}
	revokeNodeCommand.PersistentFlags().StringVar(&address, "Addr", "", "Address for the node")
	revokeNodeCommand.PersistentFlags().StringVar(&net, "Net", "", "the name of the net")

	return revokeNodeCommand
}

func runRevokeNode(address string, net string) error {

	_, err := service.RevokeNode(net, address)
	if err != nil {
		fmt.Println("revoke cert failed, ", err)
	} else {
		fmt.Println("revoke cert success")
	}
	return err
}
