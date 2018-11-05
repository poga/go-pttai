// Copyright 2018 The go-pttai Authors
// This file is part of the go-pttai library.
//
// The go-pttai library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pttai library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pttai library. If not, see <http://www.gnu.org/licenses/>.

package me

import (
	"bytes"
	"crypto/ecdsa"
	"os"

	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/content"
	"github.com/ailabstw/go-pttai/log"
)

func (spm *ServiceProtocolManager) CreateMe(myID *types.PttID, myKey *ecdsa.PrivateKey, contentBackend *content.Backend) error {
	ptt := spm.myPtt

	// new my info
	myInfo, err := NewMyInfo(myID, myKey, ptt, spm.Service())
	if err != nil {
		return err
	}

	// save
	err = myInfo.Save()
	if err != nil {
		return err
	}

	// add to entities
	err = spm.RegisterEntity(myInfo.ID, myInfo)
	if err != nil {
		return err
	}

	// SPM MyInfo
	spm.MyInfo = myInfo

	return nil
}

func (pm *ProtocolManager) CreateFullMe(oplog *MasterOplog) error {
	myInfo := pm.Entity().(*MyInfo)
	ptt := pm.myPtt

	// create-me-oplog

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	myHostname := bytes.TrimSuffix([]byte(hostname), []byte(".local"))

	myID := myInfo.ID
	myNodeID := ptt.MyNodeID()
	myNodeType := ptt.MyNodeType()
	opData := &MeOpCreateMe{
		NodeID:   myNodeID,
		NodeType: myNodeType,
		NodeName: myHostname,
	}

	meOplog, err := pm.CreateMeOplog(myID, oplog.UpdateTS, MeOpTypeCreateMe, opData)
	if err != nil {
		return err
	}

	// my-info
	myInfo.Status = types.StatusAlive
	myInfo.CreateTS = meOplog.UpdateTS
	myInfo.UpdateTS = meOplog.UpdateTS
	myInfo.LogID = meOplog.ID

	err = myInfo.Save()
	if err != nil {
		return err
	}

	// my-node
	myRaftID := ptt.MyRaftID()
	myNode, ok := pm.MyNodes[myRaftID]
	if !ok {
		return ErrInvalidNode
	}

	myNode.NodeName = myHostname
	myNode.Save()

	// meOplog save
	meOplog.Save(false)

	// op-key
	if len(pm.OpKeyInfos()) == 0 {
		err = pm.CreateOpKeyInfo()
		if err != nil {
			return err
		}
	}

	log.Debug("CreateFullMe: done")

	return nil
}