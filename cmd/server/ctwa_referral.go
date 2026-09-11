// Criado: 2026-09-11 04:43 BRT
package main

import (
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

// ctwaReferral carrega os dados de atribuição de um clique em anúncio Meta
// (Click To WhatsApp) extraídos do ExternalAdReply do ContextInfo da mensagem
// inbound. Nomes em camelCase de propósito: casam com o CtwaSourceFields que
// o Nexus já usa pra normalizar Baileys e Meta Cloud no mesmo vocabulário
// (src/core/ctwaAttribution.ts) e com o resto do envelope do webhook do Astra
// (id, chat, sender, pushName...), que também é camelCase.
type ctwaReferral struct {
	CtwaClid   string `json:"ctwaClid"`
	SourceID   string `json:"sourceId,omitempty"`
	SourceURL  string `json:"sourceUrl,omitempty"`
	SourceType string `json:"sourceType,omitempty"`
	SourceApp  string `json:"sourceApp,omitempty"`
	Ref        string `json:"ref,omitempty"`
	MediaType  string `json:"mediaType,omitempty"`
}

// extractCTWAReferral devolve o bloco de atribuição CTWA da mensagem, ou nil
// quando não há ExternalAdReply ou quando ctwaClid vem vazio.
//
// ExternalAdReply também aparece em forward e em quote que só referencia um
// anúncio (sem ser um clique real); nesses casos ctwaClid vem vazio e esta
// função devolve nil de propósito, para não gerar atribuição falsa no CRM.
// messageContextInfo (webhook.go) já sabe achar o ContextInfo nos tipos de
// mensagem que podem carregá-lo (extendedText, image, video, audio, document,
// sticker, contact, location); texto puro (conversation) nunca carrega
// ContextInfo.
func extractCTWAReferral(message *waE2E.Message) *ctwaReferral {
	contextInfo := messageContextInfo(message)
	adReply := contextInfo.GetExternalAdReply()
	clid := strings.TrimSpace(adReply.GetCtwaClid())
	if clid == "" {
		return nil
	}
	return &ctwaReferral{
		CtwaClid:   clid,
		SourceID:   strings.TrimSpace(adReply.GetSourceID()),
		SourceURL:  strings.TrimSpace(adReply.GetSourceURL()),
		SourceType: strings.TrimSpace(adReply.GetSourceType()),
		SourceApp:  strings.TrimSpace(adReply.GetSourceApp()),
		Ref:        strings.TrimSpace(adReply.GetRef()),
		MediaType:  ctwaMediaTypeLabel(adReply.GetMediaType()),
	}
}

func ctwaMediaTypeLabel(mediaType waE2E.ContextInfo_ExternalAdReplyInfo_MediaType) string {
	switch mediaType {
	case waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE:
		return "image"
	case waE2E.ContextInfo_ExternalAdReplyInfo_VIDEO:
		return "video"
	default:
		return ""
	}
}
