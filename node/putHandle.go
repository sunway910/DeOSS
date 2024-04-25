/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"encoding/json"

	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CESSProject/DeOSS/pkg/utils"
	"github.com/CESSProject/cess-go-sdk/core/pattern"
	"github.com/CESSProject/cess-go-sdk/core/process"
	sutils "github.com/CESSProject/cess-go-sdk/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ChunksInfo struct {
	BlockNum      int    `json:"block_num"`
	SavedBlockId  int    `json:"saved_block_id"`
	SavedFileSize int64  `json:"saved_file_size"`
	FileName      string `json:"file_name"`
	TotalSize     int64  `json:"total_size"`
}

const max_concurrent_req = 10

var max_concurrent_req_ch chan bool

func init() {
	max_concurrent_req_ch = make(chan bool, 10)
	for i := 0; i < max_concurrent_req; i++ {
		max_concurrent_req_ch <- true
	}
}

// putHandle
func (n *Node) putHandle(c *gin.Context) {
	if _, ok := <-max_concurrent_req_ch; !ok {
		c.JSON(http.StatusTooManyRequests, "The server is busy, please try again later.")
		return
	}
	defer func() { max_concurrent_req_ch <- true }()

	var (
		ok       bool
		err      error
		clientIp string
		fpath    string
		roothash string
		savedir  string
		filename string
		pkey     []byte
	)

	// record client ip
	clientIp = c.Request.Header.Get("X-Forwarded-For")
	n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, INFO_PutRequest))

	// verify the authorization
	account := c.Request.Header.Get(HTTPHeader_Account)
	ethAccount := c.Request.Header.Get(HTTPHeader_EthAccount)
	message := c.Request.Header.Get(HTTPHeader_Message)
	signature := c.Request.Header.Get(HTTPHeader_Signature)
	bucketName := c.Request.Header.Get(HTTPHeader_BucketName)
	cipher := c.Request.Header.Get(HTTPHeader_Cipher)
	contentLength := c.Request.ContentLength
	n.Upfile("info", fmt.Sprintf("[%v] Acc: %s", clientIp, account))
	n.Upfile("info", fmt.Sprintf("[%v] EthAcc: %s", clientIp, ethAccount))
	n.Upfile("info", fmt.Sprintf("[%v] Message: %s", clientIp, message))
	n.Upfile("info", fmt.Sprintf("[%v] Signature: %s", clientIp, signature))
	n.Upfile("info", fmt.Sprintf("[%v] BucketName: %s", clientIp, bucketName))
	n.Upfile("info", fmt.Sprintf("[%v] ContentLength: %d", clientIp, contentLength))

	if err = n.AccessControl(account); err != nil {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	if ethAccount != "" {
		ethAccInSian, err := VerifyEthSign(message, signature)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %s", clientIp, err.Error()))
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if ethAccInSian != ethAccount {
			n.Upfile("err", fmt.Sprintf("[%v] %s", clientIp, "ETH signature verification failed"))
			c.JSON(http.StatusBadRequest, "ETH signature verification failed")
			return
		}
		pkey, err = sutils.ParsingPublickey(account)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %s", clientIp, err.Error()))
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid cess account: %s", account))
			return
		}
	} else {
		pkey, err = n.VerifyAccountSignature(account, message, signature)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %s", clientIp, err.Error()))
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if !sutils.CheckBucketName(bucketName) {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_HeaderFieldBucketName))
		c.JSON(http.StatusBadRequest, ERR_HeaderFieldBucketName)
		return
	}

	// verify mem availability
	if len(cipher) > 32 {
		n.Upfile("err", fmt.Sprintf("[%v] The length of cipher cannot exceed 32", clientIp))
		c.JSON(http.StatusBadRequest, "The length of cipher cannot exceed 32")
		return
	}

	// verify the bucket name

	if strings.Contains(bucketName, " ") {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_HeaderFieldBucketName))
		c.JSON(http.StatusBadRequest, ERR_HeaderFieldBucketName)
		return
	}

	// verify the space is authorized
	var flag bool
	authAccs, err := n.QueryAuthorizedAccounts(pkey)
	if err != nil {
		if err.Error() == pattern.ERR_Empty {
			n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_SpaceNotAuth))
			c.JSON(http.StatusForbidden, ERR_SpaceNotAuth)
			return
		}
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, v := range authAccs {
		if n.GetSignatureAcc() == v {
			flag = true
			break
		}
	}
	if !flag {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_SpaceNotAuth))
		c.JSON(http.StatusForbidden, fmt.Sprintf("please authorize your space usage to %s", n.GetSignatureAcc()))
		return
	}

	if contentLength == 0 {
		txHash, err := n.CreateBucket(pkey, bucketName)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		n.Upfile("info", fmt.Sprintf("[%v] create bucket [%v] successfully: %v", clientIp, bucketName, txHash))
		if len(txHash) != (pattern.FileHashLen + 2) {
			c.JSON(http.StatusOK, "bucket already exists")
		} else {
			c.JSON(http.StatusOK, map[string]string{"Block hash:": txHash})
		}
		return
	}

	userInfo, err := n.QueryUserSpaceSt(pkey)
	if err != nil {
		if err.Error() == pattern.ERR_Empty {
			n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_NoSpace))
			c.JSON(http.StatusForbidden, ERR_NoSpace)
			return
		}
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusForbidden, ERR_RpcFailed)
		return
	}

	blockheight, err := n.QueryBlockHeight("")
	if err != nil {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusForbidden, ERR_RpcFailed)
		return
	}

	if userInfo.Deadline < (blockheight + 30) {
		n.Upfile("info", fmt.Sprintf("[%v] %v [%d] [%d]", clientIp, ERR_SpaceExpiresSoon, userInfo.Deadline, blockheight))
		c.JSON(http.StatusForbidden, ERR_SpaceExpiresSoon)
		return
	}

	mem, err := utils.GetSysMemAvailable()
	if err == nil {
		if uint64(contentLength) > uint64(mem*90/100) {
			if uint64(contentLength) < MaxMemUsed {
				n.Upfile("err", fmt.Sprintf("[%v] %v, size: [%d] mem: [%d]", clientIp, ERR_SysMemNoLeft, contentLength, mem))
				c.JSON(http.StatusForbidden, ERR_SysMemNoLeft)
				return
			}
		}
	}

	// verify disk space availability
	if contentLength > MaxMemUsed {
		freeSpace, err := utils.GetDirFreeSpace("/tmp")
		if err == nil {
			if uint64(contentLength+pattern.SIZE_1MiB*16) > freeSpace {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, ERR_DeviceSpaceNoLeft))
				c.JSON(http.StatusForbidden, ERR_DeviceSpaceNoLeft)
				return
			}
		}
	}

	usedSpace := contentLength * 15 / 10
	remainingSpace, err := strconv.ParseUint(userInfo.RemainingSpace, 10, 64)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	if usedSpace > int64(remainingSpace) {
		n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, ERR_NotEnoughSpace))
		c.JSON(http.StatusForbidden, ERR_NotEnoughSpace)
		return
	}

	blockIdx, _ := strconv.Atoi(c.Request.Header.Get(HTTPHeader_BIdx))
	blockNum, _ := strconv.Atoi(c.Request.Header.Get(HTTPHeader_BNum))
	totalSize, _ := strconv.Atoi(c.Request.Header.Get(HTTPHeader_TSize))
	filename = c.Request.Header.Get(HTTPHeader_Fname)
	var chunksInfo ChunksInfo
	for blockNum <= 0 {
		savedir = filepath.Join(n.GetDirs().FileDir, account, fmt.Sprintf("%s-%s", uuid.New().String(), uuid.New().String()))
		_, err = os.Stat(savedir)
		if err != nil {
			err = os.MkdirAll(savedir, 0755)
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
				c.JSON(http.StatusInternalServerError, ERR_InternalServer)
				return
			}
		} else {
			continue
		}
		fpath = filepath.Join(savedir, fmt.Sprintf("%v", time.Now().Unix()))
		defer os.Remove(savedir)
		break
	}

	if blockNum > 0 {
		fdir, err := sutils.CalcSHA256(append([]byte(bucketName+filename), []byte(account)...))
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
		savedir = filepath.Join(n.GetDirs().FileDir, account, fdir)
		_, err = os.Stat(savedir)
		if err != nil {
			err = os.MkdirAll(savedir, 0755)
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
				c.JSON(http.StatusInternalServerError, ERR_InternalServer)
				return
			}
		}

		fpath = filepath.Join(savedir, filename)
		_, err = os.Stat(filepath.Join(savedir, "chunk-info"))
		if err != nil {
			chunksInfo = ChunksInfo{
				BlockNum:     blockNum,
				SavedBlockId: -1,
				FileName:     filename,
				TotalSize:    int64(totalSize),
			}
		} else {
			jbytes, err := os.ReadFile(filepath.Join(savedir, "chunk-info"))
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
				c.JSON(http.StatusInternalServerError, ERR_InternalServer)
				return
			}
			err = json.Unmarshal(jbytes, &chunksInfo)
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
				c.JSON(http.StatusInternalServerError, ERR_InternalServer)
				return
			}
		}
	}

	if blockNum > 0 && (blockNum != chunksInfo.BlockNum || blockIdx != chunksInfo.SavedBlockId+1) {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, "bad chunk index"))
		c.JSON(http.StatusBadRequest, "bad chunk index")
		return
	}

	formfile, fileHeder, err := c.Request.FormFile("file")
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
		if strings.Contains(err.Error(), "no space left on device") {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
			c.JSON(http.StatusForbidden, ERR_DeviceSpaceNoLeft)
			return
		}
		if err.Error() != http.ErrNotMultipart.ErrorString {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		buf, _ := io.ReadAll(c.Request.Body)
		if len(buf) == 0 {
			txHash, err := n.CreateBucket(pkey, bucketName)
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			n.Upfile("info", fmt.Sprintf("[%v] create bucket [%v] successfully: %v", clientIp, bucketName, txHash))
			c.JSON(http.StatusOK, map[string]string{"Block hash:": txHash})
			return
		}
		// save body content
		err = sutils.WriteBufToFile(buf, fpath)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
	} else {
		if blockNum <= 0 {
			filename = fileHeder.Filename
		}
		if len(filename) > pattern.MaxBucketNameLength {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, ERR_FileNameTooLang))
			c.JSON(http.StatusBadRequest, ERR_FileNameTooLang)
			return
		}
		var f *os.File
		if blockNum > 0 && blockIdx > 0 {
			f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND, 0666)
		} else {
			f, err = os.Create(fpath)
		}

		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}

		if blockNum > 0 && fileHeder.Size+chunksInfo.SavedFileSize > chunksInfo.TotalSize {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, "bad chunk size"))
			c.JSON(http.StatusBadRequest, "bad chunk size")
			return

		}

		_, err = io.Copy(f, formfile)
		if err != nil {
			f.Close()
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
		f.Close()

		if blockNum > 0 {
			chunksInfo.SavedFileSize += fileHeder.Size
		}
	}

	if blockNum > 0 && blockIdx < blockNum-1 {
		chunksInfo.SavedBlockId += 1
		jbytes, err := json.Marshal(chunksInfo)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
		err = sutils.WriteBufToFile(jbytes, filepath.Join(savedir, "chunk-info"))
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
		n.Upfile("info", fmt.Sprintf("[%v] uploaded file block %d/%d successfully", clientIp, blockIdx, blockNum))
		c.JSON(http.StatusOK, blockIdx)
		return
	}

	if blockNum > 0 && blockIdx == blockNum-1 {
		defer os.Remove(savedir)
		stat, err := os.Stat(fpath)
		if err != nil {
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err.Error()))
			c.JSON(http.StatusInternalServerError, ERR_InternalServer)
			return
		}
		if stat.Size() != int64(totalSize) && stat.Size() != chunksInfo.TotalSize {
			os.Remove(fpath)
			os.Remove(filepath.Join(savedir, "chunk-info"))
			n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, "file is not the same size as expected"))
			c.JSON(http.StatusBadRequest, fmt.Sprintf("file size mismatch,expected %d, actual %d", totalSize, stat.Size()))
			return
		}
	}

	if filename == "" {
		filename = "null"
	}

	if len(filename) < 3 {
		filename += ".ces"
	}

	fstat, err := os.Stat(fpath)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	if fstat.Size() == 0 {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, ERR_BodyEmptyFile))
		c.JSON(http.StatusBadRequest, ERR_BodyEmptyFile)
		return
	}

	segmentInfo, roothash, err := process.ShardedEncryptionProcessing(fpath, cipher)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	n.Upfile("info", fmt.Sprintf("[%v] segmentInfo: %v", clientIp, segmentInfo))

	ok, err = n.deduplication(pkey, segmentInfo, roothash, filename, bucketName, uint64(fstat.Size()))
	if err != nil {
		if strings.Contains(err.Error(), "[GenerateStorageOrder]") {
			n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, err))
			c.JSON(http.StatusForbidden, ERR_RpcFailed)
			return
		}
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}
	if ok {
		n.Upfile("info", fmt.Sprintf("[%v] [%s] uploaded successfully", clientIp, roothash))
		c.JSON(http.StatusOK, roothash)
		return
	}

	roothashDir := filepath.Join(n.GetDirs().FileDir, account, roothash)
	_, err = os.Stat(roothashDir)
	if err == nil {
		if ok := n.HasTrackFile(roothash); ok {
			err = n.MoveFileToCache(roothash, filepath.Join(roothashDir, roothash)) //move file to cache
			if err != nil {
				n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
			}
			c.JSON(http.StatusOK, roothash)
			n.Upfile("info", fmt.Sprintf("[%v] [%s] uploaded successfully", clientIp, roothash))
			return
		}
	}

	err = os.MkdirAll(roothashDir, 0755)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	err = utils.RenameDir(filepath.Dir(fpath), roothashDir)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	err = os.Rename(filepath.Join(roothashDir, filepath.Base(fpath)), filepath.Join(roothashDir, roothash))
	if err != nil {
		os.RemoveAll(roothashDir)
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	err = n.MoveFileToCache(roothash, filepath.Join(roothashDir, roothash)) //move file to cache
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
	}

	for i := 0; i < len(segmentInfo); i++ {
		segmentInfo[i].SegmentHash = filepath.Join(roothashDir, filepath.Base(segmentInfo[i].SegmentHash))
		for j := 0; j < len(segmentInfo[i].FragmentHash); j++ {
			segmentInfo[i].FragmentHash[j] = filepath.Join(roothashDir, filepath.Base(segmentInfo[i].FragmentHash[j]))
		}
	}

	var recordInfo = &RecordInfo{
		SegmentInfo: segmentInfo,
		Owner:       pkey,
		Roothash:    roothash,
		Filename:    filename,
		Buckname:    bucketName,
		Filesize:    uint64(fstat.Size()),
		Putflag:     false,
		Count:       0,
		Duplicate:   false,
	}

	b, err := json.Marshal(recordInfo)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	err = n.WriteTrackFile(roothash, b)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(http.StatusInternalServerError, ERR_InternalServer)
		return
	}

	txhash, err := n.GenerateStorageOrder(
		roothash,
		recordInfo.SegmentInfo,
		recordInfo.Owner,
		recordInfo.Filename,
		recordInfo.Buckname,
		recordInfo.Filesize,
	)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%s] GenerateStorageOrder failed, tx: %s err: %v", roothash, txhash, err))
	} else {
		n.Upfile("info", fmt.Sprintf("[%s] GenerateStorageOrder suc: %s", roothash, txhash))
	}

	n.Upfile("info", fmt.Sprintf("[%v] [%s] uploaded successfully", clientIp, roothash))
	c.JSON(http.StatusOK, roothash)
}

func (n *Node) deduplication(pkey []byte, segmentInfo []pattern.SegmentDataInfo, roothash, filename, bucketname string, filesize uint64) (bool, error) {
	fmeta, err := n.QueryFileMetadata(roothash)
	if err == nil {
		for _, v := range fmeta.Owner {
			if sutils.CompareSlice(v.User[:], pkey) {
				return true, nil
			}
		}
		_, err = n.GenerateStorageOrder(roothash, nil, pkey, filename, bucketname, filesize)
		if err != nil {
			return false, errors.Wrapf(err, "[GenerateStorageOrder]")
		}

		return true, nil
	}

	order, err := n.QueryStorageOrder(roothash)
	if err == nil {
		if sutils.CompareSlice(order.User.User[:], pkey) {
			return true, nil
		}

		_, err = os.Stat(filepath.Join(n.trackDir, roothash))
		if err == nil {
			return false, errors.New(ERR_DuplicateOrder)
		}

		var record RecordInfo
		record.SegmentInfo = segmentInfo
		record.Owner = pkey
		record.Roothash = roothash
		record.Filename = filename
		record.Buckname = bucketname
		record.Putflag = false
		record.Count = 0
		record.Duplicate = true

		f, err := os.Create(filepath.Join(n.trackDir, roothash))
		if err != nil {
			return false, errors.Wrapf(err, "[create file]")
		}
		defer f.Close()

		b, err := json.Marshal(&record)
		if err != nil {
			return false, errors.Wrapf(err, "[marshal data]")
		}
		_, err = f.Write(b)
		if err != nil {
			return false, errors.Wrapf(err, "[write file]")
		}
		err = f.Sync()
		if err != nil {
			return false, errors.Wrapf(err, "[sync file]")
		}
		return true, nil
	}

	return false, nil
}
