/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package crypto

import (
	"github.com/xuperchain/crypto/client/service/base"
	"github.com/xuperchain/crypto/client/service/gm"
	"github.com/xuperchain/crypto/client/service/xchain"
)

// var CryptoTypeConfig = crypto_client.CryptoTypeDefault
func getInstance() interface{} {
	return &xchain.XchainCryptoClient{}
}

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {
	cryptoClient := getInstance().(base.CryptoClient)
	return cryptoClient
}

// Get HdCrypto Client
func GetHdCryptoClient() *xchain.XchainCryptoClient {
	return new(xchain.XchainCryptoClient)
}

func GetGMHdCryptoClient() *gm.GmCryptoClient {
	return new(gm.GmCryptoClient)
}
