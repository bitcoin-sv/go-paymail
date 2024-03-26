package server

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/spv"
)

type PaymailServiceLocator struct {
	paymailService PaymailServiceProvider
	pikeService    PikeServiceProvider
}

func (l *PaymailServiceLocator) RegisterPaymailService(s PaymailServiceProvider) {
	l.paymailService = s
}

func (l *PaymailServiceLocator) GetPaymailService() PaymailServiceProvider {
	if l.paymailService == nil {
		panic("PaymailServiceProvider was not registered")
	}

	return l.paymailService
}

func (l *PaymailServiceLocator) RegisterPikeService(s PikeServiceProvider) {
	l.pikeService = s
}

func (l *PaymailServiceLocator) GetPikeService() PikeServiceProvider {
	if l.pikeService == nil {
		panic("PikeServiceProvider was not registered")
	}

	return l.pikeService
}

// PaymailServiceProvider the paymail server interface that needs to be implemented
type PaymailServiceProvider interface {
	CreateAddressResolutionResponse(
		ctx context.Context,
		alias, domain string,
		senderValidation bool,
		metaData *RequestMetadata,
	) (*paymail.ResolutionPayload, error)

	CreateP2PDestinationResponse(
		ctx context.Context,
		alias, domain string,
		satoshis uint64,
		metaData *RequestMetadata,
	) (*paymail.PaymentDestinationPayload, error)

	GetPaymailByAlias(
		ctx context.Context,
		alias, domain string,
		metaData *RequestMetadata,
	) (*paymail.AddressInformation, error)

	RecordTransaction(
		ctx context.Context,
		p2pTx *paymail.P2PTransaction,
		metaData *RequestMetadata,
	) (*paymail.P2PTransactionPayload, error)

	VerifyMerkleRoots(
		ctx context.Context,
		merkleProofs []*spv.MerkleRootConfirmationRequestItem,
	) error
}

type PikeServiceProvider interface {
	AddContact(
		ctx context.Context,
		requesterPaymail string,
		contact *paymail.PikeContactRequestPayload,
	) error
}
