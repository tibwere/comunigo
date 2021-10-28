// Package per la gestione della logica applicativa
// basata sullo scambio dei messaggi gRPC tra i peer
// nel caso in cui è stata scelta la modalità 'scalar'
// oppure 'vectorial'
package p2p

import (
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

// Ridefinizione di un tipo di dato intero da
// poter utilizzare a mo' di enumerato per
// discriminare la scelta fra multicast totalmente
// o causalmente ordinato
type P2PModality uint8

const (
	P2P_SCALAR    P2PModality = 0
	P2P_VECTORIAL P2PModality = 1
)

// Costante moltiplicativa per la bufferizzazione dei canali
//
// Eventualmente potrebbe diventare un parametro se si prevede
// di dover gestire forti jitter del carico di lavoro
const (
	BUFFSIZE_FOR_PEER = 20
)

// In ottica OO, oggetto che rappresenta l'handler della comunicazione p2p
type P2PHandler struct {
	proto.UnimplementedComunigoServer
	comunicationPort uint16
	peerStatus       *peer.Status
	modality         P2PModality
	sData            *ScalarMetadata
	vData            *VectorialMetadata
}

// "Costruttore" dell'oggetto P2PHandler
func NewP2PHandler(port uint16, status *peer.Status, modality P2PModality) *P2PHandler {
	h := &P2PHandler{
		comunicationPort: port,
		peerStatus:       status,
		modality:         modality,
		sData:            nil,
	}

	if modality == P2P_SCALAR {
		h.sData = InitScalarMetadata(status.GetOtherMembers())
	} else {
		h.vData = InitVectorialMetadata(status.GetCurrentUsername(), status.GetOtherMembers())
	}

	return h
}
