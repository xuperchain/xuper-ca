/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package service

import (
	"log"

	"github.com/xuperchain/xuper-ca/crypto"
)

// The range for hardened child keys is [2^31, 2^32 - 1].
const (
	HardenedKeyStart = 0x80000000
	Language         = 1
	StrengthHard     = 3
)

// 生成网络根私钥
func GenerateNetHdPriKey() (string, error) {
	// 生成助记词
	cryptoClient := crypto.GetHdCryptoClient()
	ecdsaAccount, err := cryptoClient.CreateNewAccountWithMnemonic(Language, StrengthHard)
	if err != nil {
		log.Printf("CreateNewAccountWithMnemonic failed and err is: %v", err)
		return "", err
	}
	log.Printf("mnemonic is %v, Address: %v , jsonPrivateKey: %v, jsonPublicKey: %v", ecdsaAccount.Mnemonic, ecdsaAccount.Address, ecdsaAccount.JsonPrivateKey, ecdsaAccount.JsonPublicKey)

	jsonMasterKey, err := cryptoClient.GenerateMasterKeyByMnemonic(ecdsaAccount.Mnemonic, Language)
	if err != nil {
		log.Printf("GenerateMasterKeyByMnemonic failed and err is: %v", err)
		return "", err
	}
	return jsonMasterKey, err
}

// 生成全节点的一级私钥
func GenerateNodeHdPriKey(total uint32, netHdPriKey string) (string, error) {
	// 兼容旧网络
	if netHdPriKey == "" {
		return "", nil
	}
	nodeHdKeyStart := HardenedKeyStart + total

	cryptoClient := crypto.GetHdCryptoClient()
	childHdKey, err := cryptoClient.GenerateChildKey(netHdPriKey, nodeHdKeyStart)
	if err != nil {
		log.Printf("GenerateChildKey failed and err is: %v", err)
		return "", err
	}
	return childHdKey, err
}

// 交易解密
func DecryptByNetHdPriKey(netHdPriKey, childHdPubKey, cypherText string) (string, error) {
	if netHdPriKey == "" || childHdPubKey == "" || cypherText == "" {
		return "", ErrParam
	}
	// hd客户端
	cryptoClient := crypto.GetHdCryptoClient()
	// test
	/*
		parentPublicKey, _ := cryptoClient.ConvertPrvKeyToPubKey(netHdPriKey)
		// hdMsg := "Hello hd msg!"
		newChildPublicKey, err := cryptoClient.GenerateChildKey(parentPublicKey, 18)
		log.Printf("newChildPublicKey is %v and err is %v", newChildPublicKey, err)
		cryptoMsg, err := cryptoClient.EncryptByHdKey(newChildPublicKey, cypherText)
		log.Printf("cryptoMsg is %c", []byte(cryptoMsg))
	*/
	realMsg, err := cryptoClient.DecryptByHdKey(childHdPubKey, netHdPriKey, cypherText)
	if err != nil {
		log.Printf("DecryptByNetHdPriKey failed and err is: %v", err)
		return "", err
	}
	return realMsg, err
}
