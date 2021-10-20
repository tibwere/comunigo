package p2p

import (
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type P2PModality uint8

const (
	P2P_SCALAR    P2PModality = 0
	P2P_VECTORIAL P2PModality = 1
)

type P2PHandler struct {
	proto.UnimplementedComunigoServer
	comunicationPort uint16
	peerStatus       *peer.Status
	modality         P2PModality
	sData            *ScalarMetadata
	vData            *VectorialMetadata
}

func NewP2PHandler(port uint16, status *peer.Status, modality P2PModality) *P2PHandler {
	h := &P2PHandler{
		comunicationPort: port,
		peerStatus:       status,
		modality:         modality,
		sData:            nil,
	}

	if modality == P2P_SCALAR {
		h.sData = InitScalarMetadata(status.OtherMembers)
	} else {
		h.vData = InitVectorialMetadata(status, len(status.OtherMembers))
	}

	return h
}
