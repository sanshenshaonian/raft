package kvraft

import (
	"6.824/labrpc"
	"sync"
)
import "crypto/rand"
import "math/big"

type Clerk struct {
	servers []*labrpc.ClientEnd
	// You will have to modify this struct.
	lastLeader           int
	mu                   sync.Mutex
	clientId             int64
	lastAppliedCommandId int
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	// You'll have to add code here.
	ck.lastLeader = 0
	ck.lastAppliedCommandId = 0
	ck.clientId = nrand()
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {

	// You will have to modify this function.
	commandId := ck.lastAppliedCommandId + 1
	args := GetArgs{
		Key:       key,
		ClientId:  ck.clientId,
		CommandId: commandId,
	}
	O_DPrintf("client[%d]: 开始发送Get RPC;args=[%v]\n", ck.clientId, args)

	serverId := ck.lastLeader
	serverNum := len(ck.servers)
	for ; ; serverId = (serverId + 1) % serverNum {
		var reply GetReply
		O_DPrintf("client[%d]: 开始发送Get RPC;args=[%v]到server[%d]\n", ck.clientId, args, serverId)
		ok := ck.servers[serverId].Call("KVServer.Get", &args, &reply)
		if !ok || reply.Err == ErrTimeout || reply.Err == ErrWrongLeader {
			O_DPrintf("client[%d]: 发送Get RPC;args=[%v]到server[%d]失败,ok = %v,Reply=[%v]\n", ck.clientId, args, serverId, ok, reply)
			continue
		}
		O_DPrintf("client[%d]: 发送Get RPC;args=[%v]到server[%d]成功,Reply=[%v]\n", ck.clientId, args, serverId, reply)
		ck.lastLeader = serverId
		ck.lastAppliedCommandId = commandId
		if reply.Err == ErrNoKey {
			return ""
		}
		return reply.Value
	}
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	commandId := ck.lastAppliedCommandId + 1
	args := PutAppendArgs{
		Key:       key,
		Value:     value,
		Op:        op,
		ClientId:  ck.clientId,
		CommandId: commandId,
	}
	O_DPrintf("client[%d]: 开始发送PutAppend RPC;args=[%v]\n", ck.clientId, args)
	serverId := ck.lastLeader
	serverNum := len(ck.servers)
	for ; ; serverId = (serverId + 1) % serverNum {
		var reply PutAppendReply
		O_DPrintf("client[%d]: 开始发送PutAppend RPC;args=[%v]到server[%d]\n", ck.clientId, args, serverId)
		ok := ck.servers[serverId].Call("KVServer.PutAppend", &args, &reply)
		//当发送失败或者返回不是leader时,则继续到下一个server进行尝试
		if !ok || reply.Err == ErrTimeout || reply.Err == ErrWrongLeader {
			O_DPrintf("client[%d]: 发送PutAppend RPC;args=[%v]到server[%d]失败,ok = %v,Reply=[%v]\n", ck.clientId, args, serverId, ok, reply)
			continue
		}
		O_DPrintf("client[%d]: 发送PutAppend RPC;args=[%v]到server[%d]成功,Reply=[%v]\n", ck.clientId, args, serverId, reply)
		//若发送成功,则更新最近发现的leader以及commandId
		ck.lastLeader = serverId
		ck.lastAppliedCommandId = commandId
		return
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
