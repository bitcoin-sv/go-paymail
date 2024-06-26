package main

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
)

// Example demo implementation of a service provider
type demoServiceProvider struct {
	// Extend your dependencies or custom values
}

// GetPaymailByAlias is a demo implementation of this interface
func (d *demoServiceProvider) GetPaymailByAlias(_ context.Context, alias, domain string,
	_ *server.RequestMetadata,
) (*paymail.AddressInformation, error) {
	// Get the data from the demo database
	return DemoGetPaymailByAlias(alias, domain)
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (d *demoServiceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string,
	senderValidation bool, _ *server.RequestMetadata,
) (*paymail.ResolutionPayload, error) {
	// Generate a new destination / output for the basic address resolution
	return DemoCreateAddressResolutionResponse(ctx, alias, domain, senderValidation)
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (d *demoServiceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string,
	satoshis uint64, _ *server.RequestMetadata,
) (*paymail.PaymentDestinationPayload, error) {
	// Generate a new destination for the p2p request
	return DemoCreateP2PDestinationResponse(ctx, alias, domain, satoshis)
}

// RecordTransaction is a demo implementation of this interface
func (d *demoServiceProvider) RecordTransaction(ctx context.Context,
	p2pTx *paymail.P2PTransaction, _ *server.RequestMetadata,
) (*paymail.P2PTransactionPayload, error) {
	// Record the tx into your datastore layer
	return DemoRecordTransaction(ctx, p2pTx)
}

// VerifyMerkleRoots is a demo implementation of this interface
func (d *demoServiceProvider) VerifyMerkleRoots(ctx context.Context, merkleProofs []*spv.MerkleRootConfirmationRequestItem) error {
	// Verify the Merkle roots
	return nil
}

func (d *demoServiceProvider) AddContact(
	ctx context.Context,
	requesterPaymail string,
	contact *paymail.PikeContactRequestPayload,
) error {
	return nil
}

func (d *demoServiceProvider) CreatePikeOutputResponse(
	ctx context.Context,
	alias, domain, senderPubKey string,
	satoshis uint64,
	metaData *server.RequestMetadata,
) (*paymail.PikePaymentOutputsResponse, error) {
	return nil, nil
}
