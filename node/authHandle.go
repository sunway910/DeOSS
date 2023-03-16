/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"errors"
	"net/http"
	"time"

	"github.com/CESSProject/DeOSS/configs"
	"github.com/CESSProject/DeOSS/pkg/utils"
	"github.com/CESSProject/go-keyring"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	Account string `json:"account"`
	jwt.StandardClaims
}

type AuthType struct {
	Account   string `json:"account"`
	Message   string `json:"message"`
	Signature []byte `json:"signature"`
}

// It is used to authorize users
func (n *Node) authHandle(c *gin.Context) {
	var (
		err error
		req AuthType
	)
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(400, "Invalid.Body")
		return
	}

	// Check publickey
	pubkey, err := utils.DecodePublicKeyOfCessAccount(req.Account)
	if err != nil {
		c.JSON(400, "InvalidParameter.Account")
		return
	}

	if req.Message == "" {
		c.JSON(400, "InvalidParameter.Message")
		return
	}

	if len(req.Signature) < 64 {
		c.JSON(400, "InvalidParameter.Signature")
		return
	}

	ok, err := VerifySign(pubkey, []byte(req.Message), req.Signature)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	if !ok {
		c.JSON(403, "NoPermission")
		return
	}

	claims := CustomClaims{
		req.Account,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 60),
			ExpiresAt: int64(time.Now().Unix() + int64(configs.TokenDated)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signKey, err := utils.CalcMD5(n.Confile.GetCtrlPrk())
	if err != nil {
		c.JSON(500, "InvalidProfile")
		return
	}

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		c.JSON(500, "InternalError")
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	return
}

func VerifySign(pkey, signmsg, sign []byte) (bool, error) {
	if len(signmsg) == 0 || len(sign) < 64 {
		return false, errors.New("Invalid.Signature")
	}

	ss58, err := utils.EncodePublicKeyAsSubstrateAccount(pkey)
	if err != nil {
		return false, err
	}

	verkr, _ := keyring.FromURI(ss58, keyring.NetSubstrate{})

	var sign_array [64]byte
	for i := 0; i < 64; i++ {
		sign_array[i] = sign[i]
	}

	// Verify signature
	return verkr.Verify(verkr.SigningContext(signmsg), sign_array), nil
}
