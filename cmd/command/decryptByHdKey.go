/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xuperchain/xuper-ca/service"
)

func NewDecryptByHdKeyCommand() *cobra.Command {
	var address string
	var net string
	var hdPubKey string
	var cypherText string

	decryptByHdKeyCommand := &cobra.Command{
		Use:   "decryptByHdKey",
		Short: "decryptBy a transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDecryptByHdKey(net, address, hdPubKey, cypherText)
		},
	}
	decryptByHdKeyCommand.PersistentFlags().StringVar(&address, "Addr", "", "Address for net admin")
	decryptByHdKeyCommand.PersistentFlags().StringVar(&net, "Net", "", "the name of the net")
	decryptByHdKeyCommand.PersistentFlags().StringVar(&hdPubKey, "HdPubKey", "", "the child hdPubKey of transaction")
	decryptByHdKeyCommand.PersistentFlags().StringVar(&cypherText, "CypherText", "", "the cyphertext of transaction")

	return decryptByHdKeyCommand
}

func runDecryptByHdKey(net, address, hdPubKey, cypherText string) error {
	realMsg, err := service.DecryptByHdKey(net, address, hdPubKey, cypherText)
	if err != nil {
		fmt.Println("decryt by hdKey error,", err)
	} else {
		fmt.Println("decryt transaction success, realMsg: ", realMsg)
	}
	return err
}
