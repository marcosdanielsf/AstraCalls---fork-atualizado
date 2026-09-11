// Criado: 2026-09-11 04:43 BRT
package main

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestExtractCTWAReferralGenuineClick(t *testing.T) {
	message := &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text: proto.String("quero saber mais"),
		ContextInfo: &waE2E.ContextInfo{
			ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
				CtwaClid:   proto.String("ctwa-clid-real-123"),
				SourceID:   proto.String("120210000000000"),
				SourceURL:  proto.String("https://fb.me/anuncio"),
				SourceType: proto.String("ad"),
				SourceApp:  proto.String("instagram"),
				Ref:        proto.String("ref-campanha-x"),
				MediaType:  waE2E.ContextInfo_ExternalAdReplyInfo_VIDEO.Enum(),
			},
		},
	}}

	referral := extractCTWAReferral(message)
	if referral == nil {
		t.Fatal("esperava bloco de referral, veio nil")
	}
	want := ctwaReferral{
		CtwaClid:   "ctwa-clid-real-123",
		SourceID:   "120210000000000",
		SourceURL:  "https://fb.me/anuncio",
		SourceType: "ad",
		SourceApp:  "instagram",
		Ref:        "ref-campanha-x",
		MediaType:  "video",
	}
	if *referral != want {
		t.Fatalf("referral = %+v, want %+v", *referral, want)
	}
}

func TestExtractCTWAReferralExternalAdReplySemCtwaClidNaoEmite(t *testing.T) {
	// ExternalAdReply também aparece em forward e em quote que só referencia um
	// anúncio (sem ser um clique real). Nesses casos ctwaClid vem vazio e não
	// deve gerar atribuição falsa no CRM.
	message := &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text: proto.String("encaminhado"),
		ContextInfo: &waE2E.ContextInfo{
			ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
				SourceID:  proto.String("120210000000000"),
				SourceURL: proto.String("https://fb.me/anuncio"),
			},
		},
	}}

	if referral := extractCTWAReferral(message); referral != nil {
		t.Fatalf("esperava nil sem ctwaClid, veio %+v", *referral)
	}
}

func TestExtractCTWAReferralMensagemComumSemContextInfo(t *testing.T) {
	message := &waE2E.Message{Conversation: proto.String("oi, tudo bem?")}

	if referral := extractCTWAReferral(message); referral != nil {
		t.Fatalf("esperava nil sem contextInfo, veio %+v", *referral)
	}
}
