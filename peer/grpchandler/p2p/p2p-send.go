package p2p

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// "Metodo della classe P2PHandler" che permette di smistare i messaggi ricevuti
// dal frontend verso le varie goroutine in precedenza spawnate per l'inoltro dei messaggi
func (h *P2PHandler) MultiplexMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessageBody := <-h.peerStatus.GetFromFrontendBackendChannel():
			log.Printf("Received from frontend: %v\n", newMessageBody)

			if h.modality == P2P_SCALAR {
				newMessage := h.sData.GenerateNewMessage(h.peerStatus.GetCurrentUsername(), newMessageBody)
				log.Printf("Created new message with scalar clock %v\n", newMessage.GetTimestamp())

				h.sData.SendToAll(newMessage)
			} else {
				newMessage := h.vData.GenerateNewMessage(h.peerStatus.GetCurrentUsername(), newMessageBody)
				log.Printf("Created new message with vectorial clock %v\n", newMessage.GetTimestamp())

				h.vData.SendToAll(newMessage)

				// Questo messaggio può essere direttamente consegnato perché di sicuro
				// rispetta la causalità
				if err := h.peerStatus.RPUSHMessage(newMessage); err != nil {
					return err
				}
			}

		}
	}
}

// "Metodo della classe P2PHandler" che non fa altro di spawnare
// una goroutine per ciascun membro del gruppo di multicast
// dedicata all'inoltro dei messaggi provenienti dal frontend
func (h *P2PHandler) ConnectToPeers(ctx context.Context) error {
	errCh := make(chan error)

	for i := range h.peerStatus.GetOtherMembers() {
		index := i
		go func() {
			err := h.sendToOther(ctx, index)
			if err != nil {
				errCh <- err
			}
		}()
	}

	errMsg := ""
	for _, m := range h.peerStatus.GetOtherMembers() {
		errMsg += fmt.Sprintf("Handler for: %v->%v, ", m.GetUsername(), <-errCh)
	}
	// rimuove l'ulitmo ", "
	errMsg = errMsg[:len(errMsg)-2]

	return fmt.Errorf(errMsg)
}

// Entrypoint delle goroutine dedicate all'inoltro dei messaggi ai peer
// connessi al gruppo di multicast
func (h *P2PHandler) sendToOther(ctx context.Context, index int) error {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.peerStatus.GetSpecificMember(index).GetAddress(), h.comunicationPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf(
		"Succesfully linked to %v@%v\n",
		h.peerStatus.GetSpecificMember(index).GetUsername(),
		h.peerStatus.GetSpecificMember(index).GetAddress(),
	)

	c := proto.NewComunigoClient(conn)

	if h.modality == P2P_SCALAR {
		return h.sendLoopSC(ctx, c, index)
	} else {
		return h.sendLoopVC(ctx, c, index)
	}
}

// "Metodo della classe P2PHandler" per l'invio dei messaggi agli altri peer nel caso in cui
// la modalità scelta è 'scalar'
func (h *P2PHandler) sendLoopSC(ctx context.Context, c proto.ComunigoClient, index int) error {
	for {
		var newMessage *proto.ScalarClockMessage
		var newAck *proto.ScalarClockAck

		select {
		case <-ctx.Done():
			log.Printf("Message sender %v shutdown\n", index)
			return fmt.Errorf("signal caught")

		case newMessage = <-h.sData.GetIncomingMsgToBeSentCh(index):
			log.Printf(
				"Sending [%v] to %v@%v\n",
				newMessage,
				h.peerStatus.GetSpecificMember(index).GetUsername(),
				h.peerStatus.GetSpecificMember(index).GetAddress(),
			)
			_, err := c.SendUpdateP2PScalar(context.Background(), newMessage)
			if err != nil {
				return err
			}

		case newAck = <-h.sData.GetIncomingAckToBeSentCh(index):
			log.Printf(
				"Sending ack for [%v:%v] to %v@%v\n",
				newAck.GetFrom(),
				newAck.GetTimestamp(),
				h.peerStatus.GetSpecificMember(index).GetUsername(),
				h.peerStatus.GetSpecificMember(index).GetAddress(),
			)
			_, err := c.SendAckP2PScalar(context.Background(), newAck)
			if err != nil {
				return err
			}
		}
	}
}

// "Metodo della classe P2PHandler" per l'invio dei messaggi agli altri peer nel caso in cui
// la modalità scelta è 'vectorial'
func (h *P2PHandler) sendLoopVC(ctx context.Context, c proto.ComunigoClient, index int) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Message sender %v shutdown\n", index)
			return fmt.Errorf("signal caught")
		case newMessage := <-h.vData.GetIncomingMsgToBeSentCh(index):
			log.Printf(
				"Sending [%v] to %v@%v\n",
				newMessage, h.peerStatus.GetSpecificMember(index).GetUsername(),
				h.peerStatus.GetSpecificMember(index).GetAddress(),
			)
			_, err := c.SendUpdateP2PVectorial(context.Background(), newMessage)
			if err != nil {
				return err
			}
		}
	}
}
