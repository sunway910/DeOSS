/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package configs

// account
const (
	// CESS token precision
	CESSTokenPrecision = 1_000_000_000_000_000_000
	// MinimumBalance is the minimum balance required for the program to run
	// The unit is pico
	MinimumBalance = 2 * CESSTokenPrecision
	//
	MaxTrackThread = 10
	//
	DefaultTxTimeOut = 30
	//
	DefaultRpcAddress = "wss://testnet-rpc.cess.network/ws/"
)

const (
	Access_Public  = "public"
	Access_Private = "private"
)

const (
	App_Mode_Release = "release"
	App_Mode_Debug   = "debug"
)

const (
	SIZE_1KiB = 1024
	SIZE_1MiB = 1024 * SIZE_1KiB
	SIZE_1GiB = 1024 * SIZE_1MiB
	SIZE_1TiB = 1024 * SIZE_1GiB
)
