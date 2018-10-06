package msgParser

import (
	"log"

	"github.com/DOSNetwork/core/examples/ethereum_vss_p2p/internalMsg"
	"github.com/DOSNetwork/core/share/vss/pedersen"
	"github.com/dedis/kyber"
	"github.com/golang/protobuf/proto"
)

func PackEncryptedDeal(encryptedDeal *vss.EncryptedDeal) proto.Message {
	p := &internalMsg.EncryptedDeal{
		DHKey:     encryptedDeal.DHKey,
		Signature: encryptedDeal.Signature,
		Nonce:     encryptedDeal.Nonce,
		Cipher:    encryptedDeal.Cipher,
	}
	return proto.Message(p)
}

func UnpackEncryptedDeal(encryptedDeal *internalMsg.EncryptedDeal) *vss.EncryptedDeal {
	return &vss.EncryptedDeal{
		DHKey:     encryptedDeal.GetDHKey(),
		Signature: encryptedDeal.GetSignature(),
		Nonce:     encryptedDeal.GetNonce(),
		Cipher:    encryptedDeal.GetCipher(),
	}
}

func PackEncryptedDeals(encryptedDeals []*vss.EncryptedDeal) proto.Message {
	size := len(encryptedDeals)
	deals := make([]*internalMsg.EncryptedDeal, size)
	for i := 0; i < size; i++ {
		encryptedDeal := encryptedDeals[i]
		p := &internalMsg.EncryptedDeal{
			DHKey:     encryptedDeal.DHKey,
			Signature: encryptedDeal.Signature,
			Nonce:     encryptedDeal.Nonce,
			Cipher:    encryptedDeal.Cipher,
		}
		deals[i] = p
	}
	msg := &internalMsg.EncryptedDeals{
		Deals: deals,
	}
	return proto.Message(msg)
}

func UnpackEncryptedDeals(encryptedDeals *internalMsg.EncryptedDeals) []*vss.EncryptedDeal {
	size := len(encryptedDeals.Deals)
	deals := make([]*vss.EncryptedDeal, size)
	for i := 0; i < size; i++ {
		encryptedDeal := encryptedDeals.Deals[i]
		p := &vss.EncryptedDeal{
			DHKey:     encryptedDeal.GetDHKey(),
			Signature: encryptedDeal.GetSignature(),
			Nonce:     encryptedDeal.GetNonce(),
			Cipher:    encryptedDeal.GetCipher(),
		}
		deals[i] = p
	}

	return deals
}

func PackResponse(response *vss.Response) proto.Message {
	p := &internalMsg.Response{
		SessionID: response.SessionID,
		Index:     response.Index,
		Status:    response.Status,
		Signature: response.Signature,
	}
	return proto.Message(p)
}

func UnpackResponse(response *internalMsg.Response) *vss.Response {
	return &vss.Response{
		SessionID: response.GetSessionID(),
		Index:     response.GetIndex(),
		Status:    response.GetStatus(),
		Signature: response.GetSignature(),
	}
}

func PackResponses(responses []*vss.Response) proto.Message {
	size := len(responses)
	r := make([]*internalMsg.Response, size)
	for i := 0; i < size; i++ {
		response := responses[i]
		p := &internalMsg.Response{
			SessionID: response.SessionID,
			Index:     response.Index,
			Status:    response.Status,
			Signature: response.Signature,
		}
		r[i] = p
	}
	msg := &internalMsg.Responses{
		Responses: r,
	}
	return proto.Message(msg)
}

func UnpackResponses(responses *internalMsg.Responses) []*vss.Response {
	size := len(responses.Responses)
	r := make([]*vss.Response, size)
	for i := 0; i < size; i++ {
		response := responses.Responses[i]
		p := &vss.Response{
			SessionID: response.GetSessionID(),
			Index:     response.GetIndex(),
			Status:    response.GetStatus(),
			Signature: response.GetSignature(),
		}
		r[i] = p
	}

	return r
}

func PackJustification(justification *vss.Justification) proto.Message {
	dealBytes, err := justification.Deal.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	p := &internalMsg.Justification{
		SessionID: justification.SessionID,
		Index:     justification.Index,
		Deal:      dealBytes,
		Signature: justification.Signature,
	}
	return proto.Message(p)
}

func UnpackJustification(s vss.Suite, justification *internalMsg.Justification) *vss.Justification {
	deal := &vss.Deal{}
	err := deal.UnmarshalBinary(s, justification.GetDeal())
	if err != nil {
		log.Fatal(err)
	}

	return &vss.Justification{
		SessionID: justification.GetSessionID(),
		Index:     justification.GetIndex(),
		Deal:      deal,
		Signature: justification.GetSignature(),
	}
}

func PackSignature(index uint32, queryID string, content, sig []byte) proto.Message {
	p := &internalMsg.Signature{
		Index:     index,
		QueryId:   queryID,
		Content:   content,
		Signature: sig,
	}
	return proto.Message(p)
}

func UnpackSignature(signature *internalMsg.Signature) (index uint32, queryID string, content, sig []byte) {
	return signature.GetIndex(), signature.GetQueryId(), signature.GetContent(), signature.GetSignature()
}

func PackPublicKey(publicKey kyber.Point) proto.Message {
	publicKeyBytes, err := publicKey.MarshalBinary()

	if err != nil {
		log.Fatal(err)
	}

	p := &internalMsg.PublicKey{
		PublicKey: publicKeyBytes,
	}
	return proto.Message(p)
}

func UnpackPublicKey(s vss.Suite, publicKeyBytes *internalMsg.PublicKey) *kyber.Point {
	publicKey := s.G2().Point()
	err := publicKey.UnmarshalBinary(publicKeyBytes.GetPublicKey())
	if err != nil {
		log.Fatal(err)
	}

	return &publicKey
}

func PackPublicKeys(publicKeys []kyber.Point) proto.Message {
	publicKeysBytes := make([]*internalMsg.PublicKey, 0)
	for _, publicKey := range publicKeys {
		publicKeyBytes, err := publicKey.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}

		publicKeysBytes = append(publicKeysBytes, &internalMsg.PublicKey{
			PublicKey: publicKeyBytes,
		})
	}

	p := &internalMsg.PublicKeys{
		PublicKey: publicKeysBytes,
	}
	return proto.Message(p)
}

func UnpackPublicKeys(s vss.Suite, publicKeysBytes *internalMsg.PublicKeys) *[]kyber.Point {
	publicKeys := make([]kyber.Point, 0)
	for _, publicKeyBytes := range publicKeysBytes.GetPublicKey() {
		publicKey := s.G2().Point()
		err := publicKey.UnmarshalBinary(publicKeyBytes.GetPublicKey())
		if err != nil {
			log.Fatal(err)
		}

		publicKeys = append(publicKeys, publicKey)
	}

	return &publicKeys
}