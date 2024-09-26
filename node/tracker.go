/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CESSProject/DeOSS/common/coordinate"
	"github.com/CESSProject/DeOSS/common/record"
	"github.com/CESSProject/DeOSS/common/utils"
	"github.com/CESSProject/cess-go-sdk/chain"
	schain "github.com/CESSProject/cess-go-sdk/chain"
	"github.com/CESSProject/cess-go-sdk/core/process"
	sutils "github.com/CESSProject/cess-go-sdk/utils"
	"github.com/pkg/errors"
)

type StorageType uint8

type StorageDataType struct {
	Fid         string
	Complete    []string
	Data        [][]string
	StorageType StorageType
	Assignments []string
	Range       coordinate.Range
}

const (
	Storage_NoAssignment StorageType = iota
	Storage_PartAssignment
	Storage_FullAssignment
	Storage_RangeAssignment
)

const maxConcurrentStorages = 20

var concurrentStoragesCh chan bool

func init() {
	concurrentStoragesCh = make(chan bool, maxConcurrentStorages)
	for i := 0; i < maxConcurrentStorages; i++ {
		concurrentStoragesCh <- true
	}
}

// tracker
func (n *Node) TrackerV2() {
	n.Logtrack("info", ">>>>> start trackerv2 <<<<<")
	var tNow time.Time
	for {
		tNow = time.Now()
		n.processTrackFiles()
		if time.Since(tNow).Minutes() < 2.0 {
			time.Sleep(time.Minute * 2)
		}
	}
}

func (n *Node) processTrackFiles() {
	defer func() {
		if err := recover(); err != nil {
			n.Pnc(utils.RecoverError(err))
		}
	}()

	var err error
	var trackFiles []string
	trackFiles, err = n.ListTraceFiles()
	if err != nil {
		n.Logtrack("err", err.Error())
		return
	}
	if len(trackFiles) <= 0 {
		return
	}

	n.Logtrack("info", fmt.Sprintf("number of track files: %d", len(trackFiles)))

	//count := 0
	fid := ""
	var dealFiles = make([]StorageDataType, 0)
	for i := 0; i < len(trackFiles); i++ {
		fid = filepath.Base(trackFiles[i])
		storageDataType, ok, err := n.checkFileState(fid)
		if err != nil {
			n.Logtrack("err", fmt.Sprintf("checkFileState: %v", err))
			continue
		}
		if ok {
			n.Logtrack("info", fmt.Sprintf(" %s storage suc", fid))
			continue
		}

		switch storageDataType.StorageType {
		case Storage_NoAssignment:
			dealFiles = append(dealFiles, storageDataType)
		case Storage_PartAssignment, Storage_FullAssignment:
			if len(concurrentStoragesCh) > 0 {
				<-concurrentStoragesCh
				go n.StoragePartAssignment(concurrentStoragesCh, storageDataType, storageDataType.Assignments)
			}
		// case Storage_FullAssignment:
		// 	if len(concurrentStoragesCh) > 0 {
		// 		<-concurrentStoragesCh
		// 		go n.StorageFullAssignment(concurrentStoragesCh, storageDataType)
		// 	}
		case Storage_RangeAssignment:
			if len(concurrentStoragesCh) > 0 {
				<-concurrentStoragesCh
				go n.StorageRangeAssignment(concurrentStoragesCh, storageDataType)
			}
		}

		// count++
		// if count >= 10 {
		// 	n.Logtrack("info", fmt.Sprintf(" will storage %d files: %v", len(dealFiles), dealFiles))
		// 	err = n.storageFiles(dealFiles)
		// 	if err != nil {
		// 		n.Logtrack("err", err.Error())
		// 		return
		// 	}
		// 	count = 0
		// 	dealFiles = make([]StorageDataType, 0)
		// }
	}
	if len(dealFiles) > 0 {
		n.Logtrack("info", fmt.Sprintf(" will storage no assignment %d files", len(dealFiles)))
		err = n.storageFiles(dealFiles)
		if err != nil {
			n.Logtrack("err", err.Error())
		}
	}
}

func (n *Node) StoragePartAssignment(ch chan<- bool, data StorageDataType, assignments []string) {
	defer func() {
		ch <- true
		if err := recover(); err != nil {
			n.Pnc(utils.RecoverError(err))
		}
	}()
	var ok bool
	var accountInfo record.Minerinfo
	n.Logpart("info", " will storage file: "+data.Fid)
	n.Logpart("info", fmt.Sprintf(" file have %d batch fragments", len(data.Data)))
	for i := 0; i < len(data.Data); {
		n.Logpart("info", fmt.Sprintf(" will storage %d batch fragments", i))
		for j := 0; j < len(assignments); j++ {
			n.Logpart("info", " will storage to "+assignments[j])
			if IsStoraged(assignments[j], data.Complete) {
				n.Logpart("info", " the miner already storaged")
				continue
			}
			accountInfo, ok = n.GetMinerinfo(assignments[j])
			if !ok {
				n.Logpart("err", " not a miner")
				continue
			}
			if accountInfo.State != schain.MINER_STATE_POSITIVE {
				n.Logpart("err", " miner status is not "+schain.MINER_STATE_POSITIVE)
				continue
			}
			if accountInfo.Idlespace < chain.FragmentSize*(chain.ParShards+chain.DataShards) {
				n.Logpart("err", " miner space < 96M ")
				continue
			}
			err := n.storageBatchFragment(accountInfo, data)
			if err != nil {
				n.Logpart("err", " storage failed: "+err.Error())
				continue
			}
			n.Logpart("err", " transfer suc")
			if len(data.Data) == 1 {
				return
			}
			if len(data.Data) > 1 {
				data.Data = data.Data[1:]
			}
		}
	}
}

func (n *Node) StorageRangeAssignment(ch chan<- bool, data StorageDataType) {
	defer func() {
		ch <- true
		if err := recover(); err != nil {
			n.Pnc(utils.RecoverError(err))
		}
	}()
	minerinfolist := n.GetAllWhitelistInfos()
	minerinfolist = append(minerinfolist, n.GetAllMinerinfos()...)
	length := len(minerinfolist)
	var ok bool
	var accountInfo record.Minerinfo
	for i := 0; i < length; i++ {
		if IsStoraged(minerinfolist[i].Account, data.Complete) {
			n.Logrange("info", " the miner already storaged "+minerinfolist[i].Account)
			continue
		}
		if n.IsInBlacklist(minerinfolist[i].Account) {
			continue
		}
		coordinateInfo, err := GetCoordinate(minerinfolist[i].Addr)
		if err != nil {
			n.Logrange("err", fmt.Sprintf("[%s] getAddrCoordinate: %v", minerinfolist[i].Account, err))
			continue
		}
		if !coordinate.PointInRange(coordinateInfo, data.Range) {
			n.Logrange("err", fmt.Sprintf("[%s] %v not in range: %v", minerinfolist[i].Account, coordinateInfo, data.Range))
			continue
		}
		accountInfo, ok = n.GetMinerinfo(minerinfolist[i].Account)
		if !ok {
			n.Logrange("err", " not a miner")
			continue
		}
		if accountInfo.State != schain.MINER_STATE_POSITIVE {
			n.Logrange("err", " miner status is not "+schain.MINER_STATE_POSITIVE)
			continue
		}
		if accountInfo.Idlespace < chain.FragmentSize*(chain.ParShards+chain.DataShards) {
			n.Logrange("err", " miner space < 96M ")
			continue
		}

		n.Logrange("info", " will storage file: "+data.Fid)
		n.Logrange("info", fmt.Sprintf(" file have %d batch fragments", len(data.Data)))
		for i := 0; i < len(data.Data); {
			n.Logrange("info", " will storage to "+minerinfolist[i].Account)
			err := n.storageBatchFragment(accountInfo, data)
			if err != nil {
				n.Logrange("err", " storage failed: "+err.Error())
				continue
			}
			n.Logrange("err", " transfer suc")
			if len(data.Data) == 1 {
				return
			}
			if len(data.Data) > 1 {
				data.Data = data.Data[1:]
			}
		}
	}
}

func (n *Node) checkFileState(fid string) (StorageDataType, bool, error) {
	recordFile, err := n.ParsingTraceFile(fid)
	if err != nil {
		return StorageDataType{}, false, fmt.Errorf("[ParseTrackFromFile(%s)] %v", fid, err)
	}

	_, err = n.QueryFile(fid, -1)
	if err != nil {
		if err.Error() != chain.ERR_Empty {
			return StorageDataType{}, false, err
		}
	} else {
		for i := 0; i < len(recordFile.Segment); i++ {
			for j := 0; j < len(recordFile.Segment[i].FragmentHash); j++ {
				os.Remove(filepath.Join(recordFile.CacheDir, recordFile.Segment[i].FragmentHash[j]))
			}
		}
		n.DeleteTraceFile(fid)
		return StorageDataType{}, true, nil
	}

	flag := false
	if recordFile.Segment == nil {
		flag = true
	}

	for i := 0; i < len(recordFile.Segment); i++ {
		for j := 0; j < len(recordFile.Segment[i].FragmentHash); j++ {
			_, err = os.Stat(recordFile.Segment[i].FragmentHash[j])
			if err != nil {
				flag = true
				break
			}
		}
		if flag {
			break
		}
	}

	if flag {
		segment, hash, err := n.reFullProcessing(fid, recordFile.Cipher, recordFile.CacheDir)
		if err != nil {
			return StorageDataType{}, false, errors.Wrapf(err, "reFullProcessing")
		}
		if recordFile.Fid != hash {
			return StorageDataType{}, false, fmt.Errorf("The fid after reprocessing is inconsistent [%s != %s] %v", recordFile.Fid, hash, err)
		}
		recordFile.Segment = segment
		err = n.AddToTraceFile(fid, recordFile)
		if err != nil {
			return StorageDataType{}, false, errors.Wrapf(err, "[%s] [WriteTrackFile]", fid)
		}
	}

	var storageDataType = StorageDataType{
		Fid:      fid,
		Complete: make([]string, 0),
		Data:     make([][]string, 0),
	}

	dealmap, err := n.QueryDealMap(fid, -1)
	if err != nil {
		if err.Error() != chain.ERR_Empty {
			return StorageDataType{}, false, err
		}
	} else {
		for index := 0; index < (chain.DataShards + chain.ParShards); index++ {
			acc, ok := IsComplete(index+1, dealmap.CompleteList)
			if ok {
				storageDataType.Complete = append(storageDataType.Complete, acc)
				continue
			}
			var value = make([]string, 0)
			for i := 0; i < len(recordFile.Segment); i++ {
				value = append(value, string(recordFile.Segment[i].FragmentHash[index]))
			}
			storageDataType.Data = append(storageDataType.Data, value)
		}
		if len(recordFile.ShuntMiner) >= (chain.DataShards + chain.ParShards) {
			storageDataType.StorageType = Storage_FullAssignment
			storageDataType.Assignments = recordFile.ShuntMiner
		} else if len(recordFile.ShuntMiner) > 0 {
			suc := 0
			for i := 0; i < len(recordFile.ShuntMiner); i++ {
				for j := 0; j < len(storageDataType.Complete); j++ {
					if recordFile.ShuntMiner[i] == storageDataType.Complete[j] {
						suc++
						break
					}
				}
			}
			if suc == len(recordFile.ShuntMiner) {
				storageDataType.StorageType = Storage_NoAssignment
			} else {
				storageDataType.StorageType = Storage_PartAssignment
				storageDataType.Assignments = recordFile.ShuntMiner
			}
		} else if len(recordFile.Points.Coordinate) > 3 {
			storageDataType.StorageType = Storage_RangeAssignment
			storageDataType.Range = recordFile.Points
		} else {
			storageDataType.StorageType = Storage_NoAssignment
		}
		return storageDataType, false, nil
	}

	// verify the space is authorized
	authAccs, err := n.QueryAuthorityList(recordFile.Owner, -1)
	if err != nil {
		if err.Error() != chain.ERR_Empty {
			return StorageDataType{}, false, err
		}
	}

	flag = false
	for _, v := range authAccs {
		if sutils.CompareSlice(n.GetSignatureAccPulickey(), v[:]) {
			flag = true
			break
		}
	}

	if !flag {
		// os.RemoveAll(recordFile.CacheDir)
		// n.DeleteTrackFile(roothash)
		user, _ := sutils.EncodePublicKeyAsCessAccount(recordFile.Owner)
		return StorageDataType{}, true, errors.Errorf("[%s] user [%s] has revoked authorization", fid, user)
	}

	txhash, err := n.PlaceStorageOrder(
		fid,
		recordFile.FileName,
		recordFile.BucketName,
		recordFile.TerritoryName,
		recordFile.Segment,
		recordFile.Owner,
		recordFile.FileSize,
	)
	if err != nil {
		return StorageDataType{}, false, fmt.Errorf(" %s [UploadDeclaration] %v", fid, err)
	}
	n.Logtrack("info", fmt.Sprintf(" %s [UploadDeclaration] suc: %s", fid, txhash))

	for index := 0; index < (chain.DataShards + chain.ParShards); index++ {
		var value = make([]string, 0)
		for i := 0; i < len(recordFile.Segment); i++ {
			value = append(value, string(recordFile.Segment[i].FragmentHash[index]))
		}
		storageDataType.Data = append(storageDataType.Data, value)
	}
	if len(recordFile.ShuntMiner) >= (chain.DataShards + chain.ParShards) {
		storageDataType.StorageType = Storage_FullAssignment
	} else if len(recordFile.ShuntMiner) > 0 {
		storageDataType.StorageType = Storage_PartAssignment
	} else if len(recordFile.Points.Coordinate) > 3 {
		storageDataType.StorageType = Storage_RangeAssignment
	} else {
		storageDataType.StorageType = Storage_NoAssignment
	}
	return storageDataType, false, nil
}

func (n *Node) storageFiles(tracks []StorageDataType) error {
	if len(n.Config.Shunt.Account) >= (chain.DataShards + chain.ParShards) {
		for i := 0; i < len(tracks); i++ {
			n.StoragePartAssignment(make(chan<- bool, 1), tracks[i], n.Config.Shunt.Account)
		}
		return nil
	}
	var continueStorage = make([]StorageDataType, 0)
	if len(n.Config.Shunt.Account) > 0 {
		for i := 0; i < len(tracks); i++ {
			value, err := n.StoragePartFixed(tracks[i], n.Config.Shunt.Account)
			if err == nil {
				continueStorage = append(continueStorage, value)
			}
		}
		tracks = continueStorage
	}

	minerinfolist := n.GetAllWhitelistInfos()
	minerinfolist = append(minerinfolist, n.GetAllMinerinfos()...)
	length := len(minerinfolist)
	n.Logtrack("info", fmt.Sprintf("miner length: %d", length))
	for i := 0; i < length; i++ {
		n.Logtrack("info", fmt.Sprintf(" use miner: %s", minerinfolist[i].Account))
		// if n.IsInBlacklist(minerinfolist[i].Account) {
		// 	n.Logtrack("info", " miner in blacklist")
		// 	continue
		// }
		err := n.storageToMiner(minerinfolist[i].Account, tracks)
		if err != nil {
			n.Logtrack("err", err.Error())
		}
	}
	return nil
}

func (n *Node) StoragePartFixed(data StorageDataType, assignments []string) (StorageDataType, error) {
	var ok bool
	var accountInfo record.Minerinfo
	n.Logpart("info", " will storage file: "+data.Fid)
	n.Logpart("info", fmt.Sprintf(" file have %d batch fragments", len(data.Data)))
	var allsuc int
	for i := 0; i < len(data.Data); {
		n.Logpart("info", fmt.Sprintf(" will storage %d batch fragments", i))
		for j := 0; j < len(assignments); j++ {
			n.Logpart("info", " will storage to "+assignments[j])
			if IsStoraged(assignments[j], data.Complete) {
				allsuc++
				n.Logpart("info", " the miner already storaged")
				continue
			}
			accountInfo, ok = n.GetMinerinfo(assignments[j])
			if !ok {
				n.Logpart("err", " not a miner")
				continue
			}
			if accountInfo.State != schain.MINER_STATE_POSITIVE {
				n.Logpart("err", " miner status is not "+schain.MINER_STATE_POSITIVE)
				continue
			}
			if accountInfo.Idlespace < chain.FragmentSize*(chain.ParShards+chain.DataShards) {
				n.Logpart("err", " miner space < 96M ")
				continue
			}
			err := n.storageBatchFragment(accountInfo, data)
			if err != nil {
				n.Logpart("err", " storage failed: "+err.Error())
				continue
			}
			allsuc++
			data.Complete = append(data.Complete, assignments[j])
			n.Logpart("err", " transfer suc")
			if len(data.Data) == 1 {
				if allsuc == len(assignments) {
					return data, nil
				}
				return StorageDataType{}, errors.New("StoragePartFixed failed")
			}
			if len(data.Data) > 1 {
				data.Data = data.Data[1:]
			}
		}
		if allsuc == len(assignments) {
			return data, nil
		}
	}
	return StorageDataType{}, errors.New("StoragePartFixed failed")
}

func (n *Node) storageToMiner(account string, tracks []StorageDataType) error {
	accountInfo, ok := n.GetMinerinfo(account)
	if !ok {
		n.Logtrack("err", " not a miner")
		return nil
	}
	if accountInfo.State != schain.MINER_STATE_POSITIVE {
		n.Logtrack("err", fmt.Sprintf(" miner status is not %s", schain.MINER_STATE_POSITIVE))
		return fmt.Errorf(" %s status is not %s", account, schain.MINER_STATE_POSITIVE)
	}
	if accountInfo.Idlespace < chain.FragmentSize*(chain.ParShards+chain.DataShards) {
		n.Logtrack("err", " miner space < 96M")
		return fmt.Errorf(" %s space < 96M", account)
	}
	length := len(tracks)
	for i := 0; i < length; i++ {
		n.Logtrack("info", fmt.Sprintf(" miner will storage file %s", tracks[i].Fid))
		if IsStoraged(account, tracks[i].Complete) {
			n.Logtrack("info", " miner already storaged this file")
			continue
		}
		err := n.storageBatchFragment(accountInfo, tracks[i])
		if err != nil {
			return err
		}
		if len(tracks[i].Data) > 1 {
			tracks[i].Data = tracks[i].Data[1:]
		} else {
			tracks[i].Data = make([][]string, 0)
		}
		accountInfo.Idlespace -= chain.FragmentSize * (chain.ParShards + chain.DataShards)
		if accountInfo.Idlespace < chain.FragmentSize*(chain.ParShards+chain.DataShards) {
			n.Logtrack("info", " miner space < 96M, stop storage")
			return nil
		}
	}
	n.Logtrack("info", " all files transferred")
	return nil
}

func (n *Node) storageBatchFragment(minerinfo record.Minerinfo, tracks StorageDataType) error {
	var err error
	if len(tracks.Data) <= 0 {
		n.Logtrack("info", " miner transferred this batch of fragments")
		return nil
	}
	if len(tracks.Data[0]) <= 0 {
		n.Logtrack("info", " miner transferred all fragments of the file")
		return nil
	}
	for j := 0; j < len(tracks.Data[0]); j++ {
		err = n.UploadFragmentToMiner(minerinfo.Addr, tracks.Fid, filepath.Base(tracks.Data[0][j]), tracks.Data[0][j])
		if err != nil {
			n.Logtrack("info", fmt.Sprintf(" miner transfer %d fragment failed: %v", j, err))
			//if strings.Contains(err.Error(), "refused") || strings.Contains(err.Error(), "timeout") {
			n.AddToBlacklist(minerinfo.Account, minerinfo.Addr, err.Error())
			//}
			return err
		}
		n.Logtrack("info", fmt.Sprintf(" miner transfer %d fragment suc", j))
	}
	n.Logtrack("info", " miner transfer all fragment suc")
	minerinfo.Idlespace -= chain.FragmentSize * (chain.ParShards + chain.DataShards)
	n.AddToWhitelist(minerinfo.Account, minerinfo)
	return nil
}

func (n *Node) UploadFragmentToMiner(addr, fid, fragmentHash, fragmentPath string) error {
	message := sutils.GetRandomcode(16)
	sig, err := sutils.SignedSR25519WithMnemonic(n.GetURI(), message)
	if err != nil {
		return fmt.Errorf("[SignedSR25519WithMnemonic] %v", err)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", filepath.Base(fragmentPath))
	if err != nil {
		return err
	}

	fd, err := os.Open(fragmentPath)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(formFile, fd)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	tmp := strings.Split(addr, "\x00")
	if len(tmp) < 1 {
		return errors.New("invalid addr")
	}
	url := tmp[0]
	if strings.HasSuffix(url, "/") {
		url = url + "fragment"
	} else {
		url = url + "/fragment"
	}
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Fid", fid)
	req.Header.Set("Fragment", fragmentHash)
	req.Header.Set("Account", n.GetSignatureAcc())
	req.Header.Set("Message", message)
	req.Header.Set("Signature", hex.EncodeToString(sig))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	client.Transport = globalTransport
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed code: %d", resp.StatusCode)
	}
	var respinfo RespType
	err = json.Unmarshal(respbody, &respinfo)
	if err != nil {
		return errors.New("server returns invalid data")
	}
	if respinfo.Code != 200 {
		return fmt.Errorf("server returns code: %d", respinfo.Code)
	}
	return nil
}

func IsStoraged(account string, complete []string) bool {
	length := len(complete)
	for i := 0; i < length; i++ {
		if account == complete[i] {
			return true
		}
	}
	return false
}

func IsComplete(index int, completeInfo []schain.CompleteInfo) (string, bool) {
	length := len(completeInfo)
	for i := 0; i < length; i++ {
		if int(completeInfo[i].Index) == index {
			acc, _ := sutils.EncodePublicKeyAsCessAccount(completeInfo[i].Miner[:])
			return acc, true
		}
	}
	return "", false
}

func GetCoordinate(addr string) (coordinate.Coordinate, error) {
	longitude, latitude, ok := ParseCity(addr)
	if !ok {
		return coordinate.Coordinate{}, errors.New("parsing addr failed")
	}
	return coordinate.Coordinate{Longitude: longitude, Latitude: latitude}, nil
}

func (n *Node) reFullProcessing(fid, cipher, cacheDir string) ([]chain.SegmentDataInfo, string, error) {
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, "", err
	}
	segmentDataInfo, hash, err := process.FullProcessing(filepath.Join(n.GetFileDir(), fid), cipher, cacheDir)
	if err != nil {
		return process.FullProcessing(filepath.Join(n.GetStoringDir(), fid), cipher, cacheDir)
	}
	return segmentDataInfo, hash, err
}
