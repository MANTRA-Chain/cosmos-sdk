//go:build linux
// +build linux

package keyring

import (
	"fmt"
	"io"

	"github.com/99designs/keyring"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// Linux-only backend options.
const BackendKeyctl = "keyctl"

// KeyctlScopeUser sets the keyctl scope to "user".
func KeyctlScopeUser(options *Options)        { setKeyctlScope(options, "user") }
// KeyctlScopeUserSession sets the keyctl scope to "usersession".
func KeyctlScopeUserSession(options *Options) { setKeyctlScope(options, "usersession") }
// KeyctlScopeSession sets the keyctl scope to "session".
func KeyctlScopeSession(options *Options)     { setKeyctlScope(options, "session") }
// KeyctlScopeProcess sets the keyctl scope to "process".
func KeyctlScopeProcess(options *Options)     { setKeyctlScope(options, "process") }
// KeyctlScopeThread sets the keyctl scope to "thread".
func KeyctlScopeThread(options *Options)      { setKeyctlScope(options, "thread") }

// Options define the options of the Keyring.
type Options struct {
	// SupportedAlgos is the list of supported signing algorithms for the keyring.
	SupportedAlgos SigningAlgoList
	// SupportedAlgosLedger is the list of supported signing algorithms for the Ledger.
	SupportedAlgosLedger SigningAlgoList
	// LedgerDerivation defines the function for Ledger derivation.
	LedgerDerivation func() (ledger.SECP256K1, error)
	// LedgerCreateKey defines the function to create a Ledger key.
	LedgerCreateKey func([]byte) types.PubKey
	// LedgerAppName defines the Ledger application name.
	LedgerAppName string
	// LedgerSigSkipDERConv indicates whether Ledger should skip DER conversion on signatures,
	// depending on which format (DER or BER) the Ledger app returns signatures.
	LedgerSigSkipDERConv bool
	// KeyctlScope defines the scope of the keyctl keyring.
	KeyctlScope string
}

func newKeyctlBackendConfig(appName, _ string, _ io.Reader, opts ...Option) keyring.Config {
	options := Options{
		KeyctlScope: keyctlDefaultScope, // currently "process"
	}

	for _, optionFn := range opts {
		optionFn(&options)
	}

	return keyring.Config{
		AllowedBackends: []keyring.BackendType{keyring.KeyCtlBackend},
		ServiceName:     appName,
		KeyCtlScope:     options.KeyctlScope,
	}
}

// New creates a new instance of a keyring.
// Keyring options can be applied when generating the new instance.
// Available backends are "os", "file", "kwallet", "memory", "pass", "test", "keyctl".
func New(
	appName, backend, rootDir string, userInput io.Reader, cdc codec.Codec, opts ...Option,
) (Keyring, error) {
	if backend != BackendKeyctl {
		return newKeyringGeneric(appName, backend, rootDir, userInput, cdc, opts...)
	}

	db, err := keyring.Open(newKeyctlBackendConfig(appName, "", userInput, opts...))
	if err != nil {
		return nil, fmt.Errorf("couldn't open keyring for %q: %w", appName, err)
	}

	return newKeystore(db, cdc, backend, opts...), nil
}

func setKeyctlScope(options *Options, scope string) { options.KeyctlScope = scope }

// this is private as it is meant to be here for SDK devs convenience
// as the user does not need to pick any default when he wants to
// initialize keyctl with the default scope.
const keyctlDefaultScope = "process"
