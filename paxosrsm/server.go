package paxosrsm

import (
	"usc.edu/csci499/proj3/paxos"
)

type PaxosRSM struct {
	me      int
	px      *paxos.Paxos
	applyOp func(interface{})
	impl    PaxosRSMImpl
}

func (rsm *PaxosRSM) Kill() {
	rsm.px.Kill()
}

// applyOp(v) is a callback which the RSM invokes to let the application
// know that it can apply v (a value decided for some Paxos instance) to
// its state
func MakeRSM(me int, px *paxos.Paxos) *PaxosRSM {
	rsm := new(PaxosRSM)

	rsm.me = me
	rsm.px = px

	rsm.InitRSMImpl()

	return rsm
}
