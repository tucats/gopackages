package defs

// The subdirectory in "APP_PATH" where the .ego runtime files and assets
// are found.
const LibPathName = "lib"

// The environment variable that contains the name(s) of the loggers that
// are to be enabled by default at startup (before command line processing).
const DefaultLogging = "APP_DEFAULT_LOGGING"

// The environment variable that contains the name to use for writing log
// file messages. If not specified, defaults to writing to stdout.
const DefaultLogFileName = "APP_LOG_FILE"

// This is the name of a column automatically added to tables created using
// the 'tables' REST API.
const RowIDName = "_row_id_"

// This is the name for objects that otherwise have no name.
const Anon = "<anon>"

// This section describes the profile keys used by Ego.
const (
	// The prefix for all configuration keys reserved to Ego.
	PrivilegedKeyPrefix = "app."

	// File system location used to locate services, lib,
	// and test directories.
	PathSetting = PrivilegedKeyPrefix + "runtime.path"

	// Do we normalize the case of all symbols to a common
	// (lower) case string. If not true, symbol names are
	// case-sensitive.
	CaseNormalizedSetting = PrivilegedKeyPrefix + "compiler.normalized"

	// What is the output format that should be used by
	// default for operations that could return either
	// "text" , "indented", or "json" output.
	OutputFormatSetting = PrivilegedKeyPrefix + "console.format"

	// Set to true if the full stack should be listed during
	// tracing.
	FullStackListingSetting = PrivilegedKeyPrefix + "compiler.full.stack"

	// Should the Ego program(s) be run with "strict" or
	// "dynamic" typing? The default is "dynamic".
	StaticTypesSetting = PrivilegedKeyPrefix + "compiler.types"

	// The base URL of the Ego server providing logon services.
	LogonServerSetting = PrivilegedKeyPrefix + "logon.server"

	// The last token created by a ego logon command, which
	// is used by default for server admin commands as well
	// as rest calls.
	LogonTokenSetting = PrivilegedKeyPrefix + "logon.token"

	// Stores the expiration date from the last login. This can be
	// used to detect an expired token and provide a better message
	// to the client user than "not authorized".
	LogonTokenExpirationSetting = PrivilegedKeyPrefix + "logon.token.expiration"

	// Default allocation factor to set on symbol table create/expand
	// operations. Larger numbers are more efficient for larger symbol
	// tables, but too large a number wastes time and memory.
	SymbolTableAllocationSetting = PrivilegedKeyPrefix + "runtime.symbol.allocation"

	// If true, functions that return multiple values including an
	// error that do not assign that error to a value will result in
	// the error being thrown.
	ThrowUncheckedErrorsSetting = PrivilegedKeyPrefix + "runtime.unchecked.errors"

	// If true, the TRACE operation will print the full stack instead of
	// a shorter single-line version.
	FullStackTraceSetting = PrivilegedKeyPrefix + "runtime.stack.trace"
)

// Agent identifiers for REST calls, which indicate the role of the client.
const (
	LogonAgent  = "logon"
	StatusAgent = "status"
	TableAgent  = "tables"
)

const (
	True    = "true"
	False   = "false"
	Any     = "any"
	Strict  = "strict"
	Loose   = "relaxed"
	Dynamic = "dynamic"
	Main    = "main"
)

const (
	ByteCodeReflectionTypeString = "<*bytecode.ByteCode Value>"

	TypeCheckingVariable   = InvisiblePrefix + "type_checking"
	StrictTypeEnforcement  = 0
	RelaxedTypeEnforcement = 1
	NoTypeEnforcement      = 2

	InvisiblePrefix          = "__"
	ThisVariable             = InvisiblePrefix + "this"
	MainVariable             = InvisiblePrefix + "main"
	ErrorVariable            = InvisiblePrefix + "error"
	ArgumentListVariable     = InvisiblePrefix + "args"
	CLIArgumentListVariable  = InvisiblePrefix + "cli_args"
	ModeVariable             = InvisiblePrefix + "exec_mode"
	DebugServicePathVariable = InvisiblePrefix + "debug_service_path"
	PathsVariable            = InvisiblePrefix + "paths"
	LocalizationVariable     = InvisiblePrefix + "localization"
	LineVariable             = InvisiblePrefix + "line"
	ExtensionsVariable       = InvisiblePrefix + "extensions"
	ModuleVariable           = InvisiblePrefix + "module"
	RestStatusVariable       = InvisiblePrefix + "rest_status"
	DiscardedVariable        = "_"
	ReadonlyVariablePrefix   = "_"
	VersionName              = ReadonlyVariablePrefix + "version"
	CopyrightVariable        = ReadonlyVariablePrefix + "copyright"
	InstanceUUIDVariable     = ReadonlyVariablePrefix + "server_instance"
	BuildTimeVariable        = ReadonlyVariablePrefix + "buildtime"
	PlatformVariable         = ReadonlyVariablePrefix + "platform"
)

// ValidSettings describes the list of valid settings, and whether they can be set by the
// command line.
var ValidSettings map[string]bool = map[string]bool{
	PathSetting:                  true,
	CaseNormalizedSetting:        true,
	OutputFormatSetting:          true,
	FullStackListingSetting:      true,
	StaticTypesSetting:           true,
	LogonServerSetting:           true,
	LogonTokenSetting:            false,
	LogonTokenExpirationSetting:  false,
	FullStackTraceSetting:        true,
	SymbolTableAllocationSetting: true,
}
