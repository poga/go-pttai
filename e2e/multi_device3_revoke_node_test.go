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

package e2e

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ailabstw/go-pttai/account"
	"github.com/ailabstw/go-pttai/common/types"
	"github.com/ailabstw/go-pttai/log"
	"github.com/ailabstw/go-pttai/me"
	pkgservice "github.com/ailabstw/go-pttai/service"
	"github.com/stretchr/testify/assert"
	baloo "gopkg.in/h2non/baloo.v3"
)

func TestMultiDevice3RevokeNode(t *testing.T) {
	NNodes = 3
	isDebug := true

	var err error
	var bodyString string
	var marshaled []byte
	var dummyBool bool
	assert := assert.New(t)

	setupTest(t)
	defer teardownTest(t)

	t0 := baloo.New("http://127.0.0.1:9450")
	t1 := baloo.New("http://127.0.0.1:9451")
	t2 := baloo.New("http://127.0.0.1:9452")

	// 0 test-error
	err = testError("http://127.0.0.1:9450")
	assert.Equal(nil, err)

	err = testError("http://127.0.0.1:9451")
	assert.Equal(nil, err)

	err = testError("http://127.0.0.1:9452")
	assert.Equal(nil, err)

	// 1. get
	bodyString = `{"id": "testID", "method": "me_get", "params": []}`

	me0_1 := &me.BackendMyInfo{}
	testCore(t0, bodyString, me0_1, t, isDebug)
	assert.Equal(types.StatusAlive, me0_1.Status)

	me1_1 := &me.BackendMyInfo{}
	testCore(t1, bodyString, me1_1, t, isDebug)
	assert.Equal(types.StatusAlive, me1_1.Status)

	me2_1 := &me.BackendMyInfo{}
	testCore(t2, bodyString, me2_1, t, isDebug)
	assert.Equal(types.StatusAlive, me2_1.Status)

	// 3. getRawMe
	bodyString = `{"id": "testID", "method": "me_getRawMe", "params": [""]}`

	me0_3 := &me.MyInfo{}
	testCore(t0, bodyString, me0_3, t, isDebug)
	assert.Equal(types.StatusAlive, me0_3.Status)
	assert.Equal(me0_1.ID, me0_3.ID)
	assert.Equal(1, len(me0_3.OwnerIDs))
	assert.Equal(me0_3.ID, me0_3.OwnerIDs[0])
	assert.Equal(true, me0_3.IsOwner(me0_3.ID))
	//profileID0_3 := me0_3.ProfileID

	me1_3 := &me.MyInfo{}
	testCore(t1, bodyString, me1_3, t, isDebug)
	assert.Equal(types.StatusAlive, me1_3.Status)
	assert.Equal(me1_1.ID, me1_3.ID)
	assert.Equal(1, len(me1_3.OwnerIDs))
	assert.Equal(me1_3.ID, me1_3.OwnerIDs[0])
	assert.Equal(true, me1_3.IsOwner(me1_3.ID))
	profileID1_3 := me1_3.ProfileID

	me2_3 := &me.MyInfo{}
	testCore(t2, bodyString, me2_3, t, isDebug)
	assert.Equal(types.StatusAlive, me2_3.Status)
	assert.Equal(me2_1.ID, me2_3.ID)
	assert.Equal(1, len(me2_3.OwnerIDs))
	assert.Equal(me2_3.ID, me2_3.OwnerIDs[0])
	assert.Equal(true, me2_3.IsOwner(me2_3.ID))
	//profileID2_3 := me2_3.ProfileID

	// 4. show-my-key
	bodyString = `{"id": "testID", "method": "me_showMyKey", "params": []}`

	var myKey0_4 string

	testCore(t0, bodyString, &myKey0_4, t, isDebug)
	if isDebug {
		t.Logf("myKey0_4: %v\n", myKey0_4)
	}

	var myKey2_4 string

	testCore(t2, bodyString, &myKey2_4, t, isDebug)
	if isDebug {
		t.Logf("myKey2_4: %v\n", myKey2_4)
	}

	// 5. show-me-url
	bodyString = `{"id": "testID", "method": "me_showMeURL", "params": []}`

	dataShowMeURL1_5 := &pkgservice.BackendJoinURL{}
	testCore(t1, bodyString, dataShowMeURL1_5, t, isDebug)
	meURL1_5 := dataShowMeURL1_5.URL

	// 6. me_GetMyNodes
	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes0_6 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetMyNodes0_6, t, isDebug)
	assert.Equal(1, len(dataGetMyNodes0_6.Result))

	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes1_6 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t1, bodyString, dataGetMyNodes1_6, t, isDebug)
	assert.Equal(1, len(dataGetMyNodes1_6.Result))

	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes2_6 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t2, bodyString, dataGetMyNodes2_6, t, isDebug)
	assert.Equal(1, len(dataGetMyNodes2_6.Result))

	// 7. join-me
	log.Debug("7. join-me")

	bodyString = fmt.Sprintf(`{"id": "testID", "method": "me_joinMe", "params": ["%v", "%v", false]}`, meURL1_5, myKey0_4)

	dataJoinMe0_7 := &pkgservice.BackendJoinRequest{}
	testCore(t0, bodyString, dataJoinMe0_7, t, true)

	assert.Equal(me1_3.ID, dataJoinMe0_7.CreatorID)
	assert.Equal(me1_1.NodeID, dataJoinMe0_7.NodeID)

	// wait 15
	t.Logf("wait 10 seconds for hand-shaking")
	time.Sleep(TimeSleepRestart)

	// 7.1
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "me_joinMe", "params": ["%v", "%v", false]}`, meURL1_5, myKey2_4)

	dataJoinMe2_7 := &pkgservice.BackendJoinRequest{}
	testCore(t2, bodyString, dataJoinMe2_7, t, true)

	assert.Equal(me1_3.ID, dataJoinMe2_7.CreatorID)
	assert.Equal(me1_1.NodeID, dataJoinMe2_7.NodeID)

	// wait 15
	t.Logf("wait 10 seconds for hand-shaking")
	time.Sleep(TimeSleepRestart)

	// 8. me_GetMyNodes

	t.Logf("8. me_GetMyNodes t0")
	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes0_8 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetMyNodes0_8, t, isDebug)
	assert.Equal(3, len(dataGetMyNodes0_8.Result))
	myNode0_8_0 := dataGetMyNodes0_8.Result[0]
	myNode0_8_1 := dataGetMyNodes0_8.Result[1]
	myNode0_8_2 := dataGetMyNodes0_8.Result[2]

	assert.Equal(types.StatusAlive, myNode0_8_0.Status)
	assert.Equal(types.StatusAlive, myNode0_8_1.Status)
	assert.Equal(types.StatusAlive, myNode0_8_2.Status)

	t.Logf("8. me_GetMyNodes t1")
	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes1_8 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t1, bodyString, dataGetMyNodes1_8, t, isDebug)
	assert.Equal(3, len(dataGetMyNodes1_8.Result))
	myNode1_8_0 := dataGetMyNodes1_8.Result[0]
	myNode1_8_1 := dataGetMyNodes1_8.Result[1]
	myNode1_8_2 := dataGetMyNodes1_8.Result[2]

	assert.Equal(types.StatusAlive, myNode1_8_0.Status)
	assert.Equal(types.StatusAlive, myNode1_8_1.Status)
	assert.Equal(types.StatusAlive, myNode1_8_2.Status)

	t.Logf("8. me_GetMyNodes t2")
	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes2_8 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t2, bodyString, dataGetMyNodes2_8, t, isDebug)
	assert.Equal(3, len(dataGetMyNodes2_8.Result))
	myNode2_8_0 := dataGetMyNodes2_8.Result[0]
	myNode2_8_1 := dataGetMyNodes2_8.Result[1]
	myNode2_8_2 := dataGetMyNodes2_8.Result[2]

	assert.Equal(types.StatusAlive, myNode2_8_0.Status)
	assert.Equal(types.StatusAlive, myNode2_8_1.Status)
	assert.Equal(types.StatusAlive, myNode2_8_2.Status)

	// 9. me_getPeers

	t.Logf("9. me_getPeers")
	bodyString = `{"id": "testID", "method": "me_getPeers", "params": [""]}`
	dataGetPeers0_9 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetPeers0_9, t, isDebug)
	assert.Equal(2, len(dataGetPeers0_9.Result))

	dataGetPeers1_9 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t1, bodyString, dataGetPeers1_9, t, isDebug)
	assert.Equal(2, len(dataGetPeers1_9.Result))

	dataGetPeers2_9 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t2, bodyString, dataGetPeers2_9, t, isDebug)
	assert.Equal(2, len(dataGetPeers2_9.Result))

	// 10. revoke-node
	marshaled, _ = me1_1.NodeID.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "me_removeNode", "params": ["%v"]}`, string(marshaled))

	testCore(t0, bodyString, &dummyBool, t, isDebug)
	assert.Equal(true, dummyBool)

	// wait 10 seconds

	time.Sleep(TimeSleepRestart)

	// 11.0 test-error
	err = testError("http://127.0.0.1:9450")
	assert.Equal(nil, err)

	err = testError("http://127.0.0.1:9451")
	assert.NotEqual(nil, err)

	err = testError("http://127.0.0.1:9452")
	assert.Equal(nil, err)

	// 11. get my nodes
	bodyString = `{"id": "testID", "method": "me_getMyNodes", "params": []}`
	dataGetMyNodes0_11 := &struct {
		Result []*me.MyNode `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetMyNodes0_11, t, isDebug)
	assert.Equal(2, len(dataGetMyNodes0_11.Result))
	myNode0_11_0 := dataGetMyNodes0_11.Result[0]
	myNode0_11_1 := dataGetMyNodes0_11.Result[1]
	myNode0_11_me := myNode0_11_0
	myNode0_11_other := myNode0_11_1

	if !reflect.DeepEqual(me0_1.NodeID, myNode0_11_0.NodeID) {
		myNode0_11_me = myNode0_11_1
		myNode0_11_other = myNode0_11_0
	}

	assert.Equal(types.StatusAlive, myNode0_11_0.Status)
	assert.Equal(me0_1.NodeID, myNode0_11_me.NodeID)
	assert.Equal(me2_1.NodeID, myNode0_11_other.NodeID)

	// 12. getPeers
	bodyString = `{"id": "testID", "method": "me_getPeers", "params": [""]}`

	dataPeers0_12 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t0, bodyString, dataPeers0_12, t, isDebug)
	assert.Equal(1, len(dataPeers0_12.Result))

	dataPeers2_12 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t2, bodyString, dataPeers2_12, t, isDebug)
	assert.Equal(1, len(dataPeers2_12.Result))

	// 12.1 account_getPeers
	marshaled, _ = profileID1_3.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_getPeers", "params": ["%v"]}`, string(marshaled))

	dataPeers0_12_1 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t0, bodyString, dataPeers0_12_1, t, isDebug)
	assert.Equal(1, len(dataPeers0_12_1.Result))
	peer0_12_1_0 := dataPeers0_12_1.Result[0]
	assert.Equal(me2_1.NodeID, peer0_12_1_0.NodeID)

	dataPeers2_12_1 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t2, bodyString, dataPeers2_12_1, t, isDebug)
	assert.Equal(1, len(dataPeers2_12_1.Result))
	peer2_12_1_0 := dataPeers2_12_1.Result[0]
	assert.Equal(me0_1.NodeID, peer2_12_1_0.NodeID)

	// 13. getPeers
	bodyString = `{"id": "testID", "method": "ptt_getPeers", "params": []}`

	dataPeers0_13 := &struct {
		Result []*pkgservice.BackendPeer `json:"result"`
	}{}
	testListCore(t0, bodyString, dataPeers0_13, t, isDebug)
	assert.Equal(1, len(dataPeers0_13.Result))

	// 10.1. force sync.
	marshaled, _ = profileID1_3.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_forceSync", "params": ["%v"]}`, string(marshaled))

	bool0_10_1 := false
	testCore(t0, bodyString, &bool0_10_1, t, isDebug)
	assert.Equal(true, bool0_10_1)
	time.Sleep(TimeSleepRestart)

	bool2_10_1 := false
	testCore(t2, bodyString, &bool2_10_1, t, isDebug)
	assert.Equal(true, bool2_10_1)
	time.Sleep(TimeSleepRestart)

	// 14. getUserNode
	t.Logf("10.6 GetUserNodeList")
	marshaled, _ = profileID1_3.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_getUserNodeList", "params": ["%v", "", 0, 2]}`, string(marshaled))
	dataGetUserNodeList0_14 := &struct {
		Result []*account.UserNode `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetUserNodeList0_14, t, isDebug)
	assert.Equal(3, len(dataGetUserNodeList0_14.Result))

	for _, each := range dataGetUserNodeList0_14.Result {
		if reflect.DeepEqual(each.NodeID, me1_1.NodeID) {
			assert.Equal(types.StatusDeleted, each.Status)
		} else {
			assert.Equal(types.StatusAlive, each.Status)
		}
	}

	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_getUserNodeList", "params": ["%v", "", 0, 2]}`, string(marshaled))
	dataGetUserNodeList2_14 := &struct {
		Result []*account.UserNode `json:"result"`
	}{}
	testListCore(t2, bodyString, dataGetUserNodeList2_14, t, isDebug)
	assert.Equal(3, len(dataGetUserNodeList2_14.Result))

	for _, each := range dataGetUserNodeList2_14.Result {
		if reflect.DeepEqual(each.NodeID, me1_1.NodeID) {
			assert.Equal(types.StatusDeleted, each.Status)
		} else {
			assert.Equal(types.StatusAlive, each.Status)
		}
	}

	// 15. new t0
	t.Logf("15 GetUserOplogList: new t0")
	marshaled, _ = profileID1_3.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_getUserOplogList", "params": ["%v", "", 0, 2]}`, string(marshaled))
	dataGetUserOplogList0_15 := &struct {
		Result []*account.UserOplog `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetUserOplogList0_15, t, isDebug)
	assert.Equal(9, len(dataGetUserOplogList0_15.Result))
	//assert.Equal(dataGetUserOplogList0_9_4_1.Result, dataGetUserOplogList0_15.Result[:5])
	userOplog0_15 := dataGetUserOplogList0_15.Result[7]
	assert.Equal(types.StatusAlive, userOplog0_15.ToStatus())
	assert.Equal(account.UserOpTypeRemoveUserNode, userOplog0_15.Op)

	marshaled, _ = profileID1_3.MarshalText()
	bodyString = fmt.Sprintf(`{"id": "testID", "method": "account_getUserOplogList", "params": ["%v", "", 0, 2]}`, string(marshaled))
	dataGetUserOplogList2_15 := &struct {
		Result []*account.UserOplog `json:"result"`
	}{}
	testListCore(t0, bodyString, dataGetUserOplogList2_15, t, isDebug)
	assert.Equal(9, len(dataGetUserOplogList2_15.Result))
	//assert.Equal(dataGetUserOplogList0_9_4_1.Result, dataGetUserOplogList2_15.Result[:5])
	userOplog2_15 := dataGetUserOplogList2_15.Result[7]
	assert.Equal(types.StatusAlive, userOplog2_15.ToStatus())
	assert.Equal(account.UserOpTypeRemoveUserNode, userOplog2_15.Op)
}
