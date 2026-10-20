package main

import (
	"context"
	"errors"
	"testing"

	"go.mau.fi/whatsmeow/types"
)

type fakeLIDResolver struct {
	result    types.JID
	err       error
	got       types.JID
	callCount int
}

func (f *fakeLIDResolver) GetPNForLID(_ context.Context, lid types.JID) (types.JID, error) {
	f.callCount++
	f.got = lid
	return f.result, f.err
}

func TestResolveRecipientJIDLeavesPhoneNumberJIDUnchanged(t *testing.T) {
	input := types.JID{User: "573135804424", Server: types.DefaultUserServer}
	resolver := &fakeLIDResolver{}

	got, err := resolveRecipientJID(context.Background(), resolver, input)
	if err != nil {
		t.Fatalf("resolveRecipientJID() error = %v", err)
	}
	if got != input {
		t.Fatalf("resolveRecipientJID() = %v, want %v", got, input)
	}
	if resolver.callCount != 0 {
		t.Fatalf("GetPNForLID called %d times for a phone-number JID, want 0", resolver.callCount)
	}
}

func TestResolveRecipientJIDMapsLIDToPhoneNumber(t *testing.T) {
	lid := types.JID{User: "88197354754048", Server: types.HiddenUserServer}
	pn := types.JID{User: "573135804424", Server: types.DefaultUserServer}
	resolver := &fakeLIDResolver{result: pn}

	got, err := resolveRecipientJID(context.Background(), resolver, lid)
	if err != nil {
		t.Fatalf("resolveRecipientJID() error = %v", err)
	}
	if got != pn {
		t.Fatalf("resolveRecipientJID() = %v, want %v", got, pn)
	}
	if resolver.callCount != 1 {
		t.Fatalf("GetPNForLID called %d times, want 1", resolver.callCount)
	}
	if resolver.got != lid {
		t.Fatalf("GetPNForLID argument = %v, want %v", resolver.got, lid)
	}
}

func TestResolveRecipientJIDRejectsUnmappedLID(t *testing.T) {
	lid := types.JID{User: "88197354754048", Server: types.HiddenUserServer}
	resolver := &fakeLIDResolver{}

	_, err := resolveRecipientJID(context.Background(), resolver, lid)
	if err == nil {
		t.Fatal("resolveRecipientJID() error = nil, want an unmapped-LID error")
	}
	if resolver.callCount != 1 {
		t.Fatalf("GetPNForLID called %d times, want 1", resolver.callCount)
	}
}

func TestResolveRecipientJIDPropagatesMappingError(t *testing.T) {
	lid := types.JID{User: "88197354754048", Server: types.HiddenUserServer}
	wantErr := errors.New("mapping store unavailable")
	resolver := &fakeLIDResolver{err: wantErr}

	_, err := resolveRecipientJID(context.Background(), resolver, lid)
	if !errors.Is(err, wantErr) {
		t.Fatalf("resolveRecipientJID() error = %v, want wrapped %v", err, wantErr)
	}
}

func TestResolveRecipientJIDRejectsInvalidMapping(t *testing.T) {
	lid := types.JID{User: "88197354754048", Server: types.HiddenUserServer}
	resolver := &fakeLIDResolver{result: types.JID{User: "573135804424", Server: "g.us"}}

	_, err := resolveRecipientJID(context.Background(), resolver, lid)
	if err == nil {
		t.Fatal("resolveRecipientJID() error = nil, want an invalid-mapping error")
	}
}
