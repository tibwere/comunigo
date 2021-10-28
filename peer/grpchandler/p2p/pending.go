package p2p

import (
	"context"
	"fmt"
	"log"
)

// "Metodo della classe P2PHandler" wrapper che ha l'unico scopo di invocare gli effettivi
// handler delle code pendenti a seconda della modalità scelta
func (h *P2PHandler) MessageQueueHandler(ctx context.Context) error {
	if h.modality == P2P_SCALAR {
		return h.messageQueueHandlerSC(ctx)
	} else {
		return h.messageQueueHandlerVC(ctx)
	}
}

// "Metodo della classe P2PHandler" per la gestione della coda dei messaggi
// e degli ack da processare nel caso in cui la modalità scelta è 'scalar'
//
// n.b. utilizzare una procedura singola che riceve input dai canali permette
// di gestire la sincronizzazione evitanto l'uso esplicito di lock
func (h *P2PHandler) messageQueueHandlerSC(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-h.sData.GetMessageCh():
			log.Printf("Insert [%v] into pendant queue\n", newMessage)
			h.sData.PushIntoPendingList(newMessage)
			h.sData.SyncDatastore(h.peerStatus)
		case newAck := <-h.sData.GetAckCh():
			log.Printf("Increment ack counter of [%v:%v]\n", newAck.GetFrom(), newAck.GetTimestamp())
			h.sData.IncrementAckCounter(newAck)
			h.sData.SyncDatastore(h.peerStatus)
		}
	}
}

// "Metodo della classe P2PHandler" per la gestione della coda dei messaggi
// e degli ack da processare nel caso in cui la modalità scelta è 'vectorial'
//
// n.b. utilizzare una procedura singola che riceve input dai canali permette
// di gestire la sincronizzazione evitanto l'uso esplicito di lock
func (h *P2PHandler) messageQueueHandlerVC(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-h.vData.GetReceivedCh():
			log.Printf("Insert [%v] into pendant queue\n", newMessage)
			h.vData.PushIntoPendingList(newMessage)
			if err := h.vData.SyncDatastore(h.peerStatus); err != nil {
				return err
			}
		}

	}
}
