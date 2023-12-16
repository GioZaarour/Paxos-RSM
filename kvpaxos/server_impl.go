package kvpaxos

import (
	"fmt"
	// "net/rpc"
	"time"

	"usc.edu/csci499/proj3/paxos"
)

// Define what goes into "value" that Paxos is used to agree upon.
// Field names must start with capital letters,
// otherwise RPC will break.
type Op struct {
	Type  string
	Key   string
	Value string
	ID    int64
	Seq int
}

// additions to KVPaxos state
type KVPaxosImpl struct {
	kvMap             map[string]string
	operationsApplied map[int64]bool
<<<<<<< HEAD
	px                *paxos.Paxos
=======
	seq int
	seqCache map[int64]int 
	// px *paxos.Paxos 
>>>>>>> 7f3342fe3b401929da129a4f8dc5bd5826624e8e
}

// initialize kv.impl.*
func (kv *KVPaxos) InitImpl() {
	kv.impl = KVPaxosImpl{
		operationsApplied: make(map[int64]bool),
		seqCache: make(map[int64]int),
		seq: 0,
	}
	// fmt.Println("peers: ", kv.impl.peers)

}

// Handler for Get RPCs
func (kv *KVPaxos) Get(args *GetArgs, reply *GetReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Check if duplicate
	_, doneBefore := kv.impl.operationsApplied[args.Impl.ID]
	if doneBefore {
		reply.Value = kv.impl.kvMap[args.Key]
		reply.Err = OK
		return nil
	}

	// Not a duplicate
	op := Op{
		Type: "Get",
		Key:  args.Key,
		ID:   args.Impl.ID,
	}

	// add seq num to cache
	seq := kv.incrementSeq()
	op.Seq = seq 
	kv.addSeqToCache(args.Impl.ID, args.Impl.Seq)

	kv.rsm.AddOp(op.Seq)

	// Wait for Paxos consensus
	for {
<<<<<<< HEAD
		status, _ := kv.impl.px.Status(index)
		if kv.impl.px.StatusString(status) == "Decided" {
=======
		status, _ := kv.rsm.Status(op.Seq)
		if status == paxos.Decided {
>>>>>>> 7f3342fe3b401929da129a4f8dc5bd5826624e8e
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	decidedOp := kv.waitForDecision(int(op.ID))

	kv.ApplyOp(decidedOp)

	return nil
}

// Handler for Put and Append RPCs
func (kv *KVPaxos) PutAppend(args *PutAppendArgs, reply *PutAppendReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	op := Op{
		Type: args.Op,
		Key:  args.Key,
		Value: args.Value,
		ID:   args.Impl.ID,
		Seq: args.Impl.Seq,
	}

	fmt.Println("op: ", op.Type)
	fmt.Println("key: ", op.Key)
	fmt.Println("value: ", op.Value)

	kv.rsm.AddOp(op.Seq)

	// check for duplicate, dont process duplicates
	if kv.impl.operationsApplied[args.Impl.ID] {
		reply.Err = paxos.OK
		return nil
	}

	kv.addSeqToCache(args.Impl.ID, args.Impl.Seq)

	// wait for paxos consensus
	for {
<<<<<<< HEAD
		status, _ := kv.px.Status(index)
		if kv.impl.px.StatusString(status) == "Decided" {
=======
		// fmt.Println("peers: ", kv.impl.peers)

		status, _ := kv.rsm.Status(args.Impl.Seq)
		if status == paxos.Decided {
>>>>>>> 7f3342fe3b401929da129a4f8dc5bd5826624e8e
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	decidedOp := kv.waitForDecision(int(op.ID))
	kv.rsm.AddOp(decidedOp)

// check if val exists
	if _, exists := kv.impl.kvMap[args.Key]; exists {
		reply.Err = paxos.OK
	} else {
		reply.Err = "error"
	}
	return nil 

}

func (kv *KVPaxos) waitForDecision(index int) Op {
	for {
<<<<<<< HEAD
		status, decidedValue := kv.impl.px.Status(index)
		if kv.impl.px.StatusString(status) == "Decided" {
=======
		status, decidedValue := kv.rsm.Status(index)
		if status == paxos.Decided {
>>>>>>> 7f3342fe3b401929da129a4f8dc5bd5826624e8e
			return decidedValue.(Op)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Execute operation encoded in decided value v and update local state
func (kv *KVPaxos) ApplyOp(v interface{}) {
	operation, ok := v.(Op)

	if !ok {
		return
	}

	// check again for duplicate
<<<<<<< HEAD
	if kv.impl.operationsApplied(operation.ID) {
=======
	if _, exists := kv.impl.operationsApplied[operation.ID]; exists {
>>>>>>> 7f3342fe3b401929da129a4f8dc5bd5826624e8e
		return
	}

	// apply operation
	curr, exists := kv.impl.kvMap[operation.Key]
	switch operation.Type {
	case "Put":
		kv.impl.kvMap[operation.Key] = operation.Value
	case "Append":
		if exists {
			kv.impl.kvMap[operation.Key] = curr + operation.Value
		} else {
			kv.impl.kvMap[operation.Key] = operation.Value
		}

		// op has been applied
		kv.impl.operationsApplied[operation.ID] = true

	}

}

func (kv *KVPaxos) incrementSeq() int {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.impl.seq++
	return kv.impl.seq
}


func (kv *KVPaxos) addSeqToCache(id int64, seq int) {
	kv.impl.seqCache[id] = seq
}