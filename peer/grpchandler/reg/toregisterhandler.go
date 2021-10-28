// Package per la gestione della comunicazione
// fra peer e nodo di registrazione al fine di registrarsi
// al gruppo di multicast
package reg

import "gitlab.com/tibwere/comunigo/peer"

// In ottica OO, oggetto che rappresenta l'handler della
// comunicazione verso il nodo di registrazione
type ToRegisterGRPCHandler struct {
	registerAddr string
	registerPort uint16
	peerStatus   *peer.Status
}

// "Costruttore" dell'oggetto ToRegisterGRPCHandler
func NewToRegisterGRPCHandler(addr string, port uint16, status *peer.Status) *ToRegisterGRPCHandler {
	return &ToRegisterGRPCHandler{
		registerAddr: addr,
		registerPort: port,
		peerStatus:   status,
	}
}
