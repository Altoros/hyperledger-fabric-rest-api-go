package notifications

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func HandleChaincodeEvent(ccEvent *fab.CCEvent) {
	fmt.Printf("Received CC event: %v\n", ccEvent)
	fmt.Println("Payload: " + string(ccEvent.Payload))
	// TODO future research of zero payload
}

func HandleFilteredBlockEvent(bEvent *fab.FilteredBlockEvent)  {
	fmt.Println("Received FilteredBlockEvent event")
	fmt.Println(bEvent.FilteredBlock)
}

func HandleTxStatusEvent(txEvent *fab.TxStatusEvent){
	fmt.Println("Received TxStatus event")
	fmt.Println(txEvent)
}