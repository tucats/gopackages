package i18n

// messages contains a map of internalized strings. The map is organized
// by message id passed into the i18n.T() function as the first key, and
// then the language derived from the environment.  If a string for a key
// in a given language is not found, it reverts to using the "en" key.
//
// If the text isn't found in English either, the key is returned
// as the unlocalizable result.
var messages = map[string]map[string]string{
	"app.config": {
		"en": "View or set application configuration",
	},
	"app.config.delete": {
		"en": "Delete an application configuration item",
	},
	"app.config.list": {
		"en": "List the application configuration profiles",
	},
	"app.config.remove": {
		"en": "Remove an application configuration profile",
	},
	"app.config.set": {
		"en": "Set an application configuration item",
	},
	"app.config.set.description": {
		"en": "Set the application configuration profile description",
	},
	"app.config.set.output": {
		"en": "Set the default output format",
	},
	"app.config.show": {
		"en": "Show the application configuration items in the current profile",
	},
	"app.logon": {
		"en": "Log on to a remote server",
	},
	"global.version": {
		"en": "Show application version",
	},
	"error.arg.count": {
		"en": "incorrect function argument count",
	},
	"error.arg.type": {
		"en": "incorrect function argument type",
	},
	"error.bit.shift": {
		"en": "invalid bit shift specification",
	},
	"error.boolean.option": {
		"en": "invalid boolean option value",
	},
	"error.bytecode.address": {
		"en": "invalid bytecode address",
	},
	"error.bytecode.not.found": {
		"en": "unimplemented bytecode instruction",
	},
	"error.cannot.delete.profile": {
		"en": "cannot delete active profile",
	},
	"cert.parse.err": {
		"en": "error parsing certficate file",
	},
	"error.cli.command.not.found": {
		"en": "unrecognized command",
	},
	"error.cli.extra": {
		"en": "unexpected text after command",
	},
	"error.cli.option": {
		"en": "unknown command line option",
	},
	"error.cli.parms": {
		"en": "too many parameters on command line",
	},
	"error.cli.subcommand": {
		"en": "unexpected parameters or invalid subcommand",
	},
	"error.colon": {
		"en": "missing ':'",
	},
	"error.column.count": {
		"en": "incorrect number of columns",
	},
	"error.column.name": {
		"en": "invalid column name",
	},
	"error.column.number": {
		"en": "invalid column number",
	},
	"error.column.width": {
		"en": "invalid column width",
	},
	"error.compiler": {
		"en": "internal compiler error",
	},
	"error.conditional.bool": {
		"en": "invalid conditional expression type",
	},
	"error.constant": {
		"en": "invalid constant expression",
	},
	"error.credentials": {
		"en": "invalid credentials",
	},
	"error.credentials.missing": {
		"en": "no credentials provided",
	},
	"error.div.zero": {
		"en": "division by zero",
		"fr": "division par zéro",
	},
	"error.dup.column": {
		"en": "duplicate column name",
	},
	"error.dup.type": {
		"en": "duplicate type name",
	},
	"error.empty.column": {
		"en": "empty column list",
	},
	"error.equals": {
		"en": "missing '='",
	},
	"error.expired": {
		"en": "expired token",
	},
	"error.expression": {
		"en": "missing expression",
	},
	"error.expression.term": {
		"en": "missing term",
	},
	"error.extension": {
		"en": "unsupported language extension",
	},
	"error.filter.term.invalid": {
		"en": "Unrecognized operator {{term}}",
	},
	"error.filter.term.missing": {
		"en": "Missing filter term",
	},
	"error.for.assignment": {
		"en": "missing ':='",
	},
	"error.for.body": {
		"en": "for{} body empty",
	},
	"error.for.exit": {
		"en": "for{} has no exit",
	},
	"error.for.init": {
		"en": "missing for-loop initializer",
	},
	"error.format.spec": {
		"en": "invalid or unsupported format specification",
	},
	"error.format.type": {
		"en": "invalid output format type",
	},
	"error.fucntion.values": {
		"en": "missing return values",
	},
	"error.func.arg": {
		"en": "invalid function argument",
	},
	"error.func.call": {
		"en": "invalid function invocation",
	},
	"error.func.exists": {
		"en": "function already defined",
	},
	"error.func.name": {
		"en": "invalid function name",
	},
	"error.func.return.count": {
		"en": "incorrect number of return values",
	},
	"error.func.unused": {
		"en": "function call used as parameter has unused error return value",
	},
	"error.func.void": {
		"en": "function did not return a value",
	},
	"error.function": {
		"en": "missing function",
	},
	"error.function.body": {
		"en": "missing function body",
	},
	"error.function.list": {
		"en": "missing function parameter list",
	},
	"error.function.name": {
		"en": "missing function name",
	},
	"error.function.ptr": {
		"en": "unable to convert {{ptr}} to function pointer",
	},
	"error.function.receiver": {
		"en": "no function receiver",
	},
	"error.function.return": {
		"en": "missing function return type",
	},
	"error.general": {
		"en": "general error",
	},
	"error.go.error": {
		"en": "Go routine {{name}} failed, {{err}}",
	},
	"error.http": {
		"en": "received HTTP",
	},
	"error.identifier": {
		"en": "invalid identifier",
	},
	"error.identifier.not.found": {
		"en": "unknown identifier",
	},
	"error.immutable.array": {
		"en": "cannot change an immutable array",
	},
	"error.immutable.map": {
		"en": "cannot change an immutable map",
	},
	"error.import": {
		"en": "import not permitted inside a block or loop",
	},
	"error.import.not.found": {
		"en": "attempt to use imported package not in package cache",
	},
	"error.instruction": {
		"en": "invalid instruction",
	},
	"error.integer.option": {
		"en": "invalid integer option value",
	},
	"error.interface.imp": {
		"en": "missing interface implementation",
	},
	"error.invalid.alignment.spec": {
		"en": "invalid alignment specification",
	},
	"error.invalid.blockquote": {
		"en": "invalid block quote",
	},
	"error.keyword.option": {
		"en": "invalid option keyword",
	},
	"error.list": {
		"en": "invalid list",
	},
	"error.logger.confict": {
		"en": "conflicting logger state",
	},
	"error.logger.name": {
		"en": "invalid logger name",
	},
	"error.logon.endpoint": {
		"en": "logon endpoint not found",
	},
	"error.logon.server": {
		"en": "no --logon-server specified",
	},
	"error.media.type": {
		"en": "invalid media type",
	},
	"error.nil": {
		"en": "nil pointer reference",
	},
	"error.no.info": {
		"en": "no information for item",
	},
	"error.no.symbol.table": {
		"en": "no symbol table available",
	},
	"error.not.assignment.list": {
		"en": "not an assignment list",
	},
	"error.not.found": {
		"en": "not found",
	},
	"error.not.type": {
		"en": "not a type",
	},
	"error.opcode.defined": {
		"en": "opcode already defined",
	},
	"error.option.required": {
		"en": "required option not found",
	},
	"error.option.value": {
		"en": "missing option value",
	},
	"error.panic": {
		"en": "Panic",
	},
	"error.parens": {
		"en": "missing parenthesis",
	},
	"error.parm.count": {
		"en": "incorrect number of parameters",
	},
	"error.parm.value.count": {
		"en": "wrong number of parameter values",
	},
	"error.privilege": {
		"en": "no privilege for operation",
	},
	"error.profile.key": {
		"en": "no such profile key",
	},
	"error.profile.name": {
		"en": "invalid configuration name",
	},
	"error.profile.not.found": {
		"en": "no such profile",
	},
	"error.range": {
		"en": "invalid range",
	},
	"error.readonly": {
		"en": "invalid write to read-only item",
	},
	"error.readonly.write": {
		"en": "invalid write to read-only value",
	},
	"error.reserved.name": {
		"en": "reserved profile setting name",
	},
	"error.row.number": {
		"en": "invalid row number",
	},
	"error.sandbox.path": {
		"en": "invalid sandbox path",
	},
	"error.semicolon": {
		"en": "missing ';'",
	},
	"error.slice.index": {
		"en": "invalid slice index",
	},
	"error.spacing": {
		"en": "invalid spacing value",
	},
	"error.stack.underflow": {
		"en": "stack underflow",
	},
	"error.statement": {
		"en": "missing statement",
	},
	"error.statement.not.found": {
		"en": "unexpected token",
	},
	"error.step.type": {
		"en": "invalid step type",
	},
	"error.symbol.exists": {
		"en": "symbol already exists",
	},
	"error.symbol.name": {
		"en": "invalid symbol name",
	},
	"error.symbol.not.found": {
		"en": "unknown symbol",
	},
	"error.symbol.overflow": {
		"en": "too many local symbols defined",
	},
	"error.table.closed": {
		"en": "table closed",
	},
	"error.table.processing": {
		"en": "table processing",
	},
	"error.terminated": {
		"en": "terminated with errors",
	},
	"error.token.encryption": {
		"en": "invalid token encryption",
	},
	"error.token.extra": {
		"en": "unexpected token",
	},
	"error.type": {
		"en": "invalid or unsupported data type for this operation",
	},
	"error.type.check": {
		"en": "invalid @type keyword",
	},
	"error.type.def": {
		"en": "missing type definition",
	},
	"error.type.mismatch": {
		"en": "type mismatch",
	},
	"error.type.name": {
		"en": "invalid type name",
	},
	"error.type.not.found": {
		"en": "no such type",
	},
	"error.type.spec": {
		"en": "invalid type specification",
	},
	"error.unsupported.on.os": {
		"en": "command not implemented for this operating system",
	},
	"error.user.defined": {
		"en": "user-supplied error",
	},
	"error.user.not.found": {
		"en": "no such user",
	},
	"error.value": {
		"en": "invalid value",
	},
	"error.value.extra": {
		"en": "unexpected value",
	},
	"error.var.type": {
		"en": "invalid type for this variable",
	},
	"error.version.parse": {
		"en": "Unable to process version number {{v}; count={{c}}, err={{err}\n",
	},
	"errors.terminated": {
		"en": "terminated due to errors",
	},
	"label.Active": {
		"en": "Active",
	},
	"label.active.loggers": {
		"en": "Active loggers: ",
	},
	"label.Columns": {
		"en": "Columns",
	},
	"label.Command": {
		"en": "Command",
	},
	"label.Commands": {
		"en": "Commands",
	},
	"label.Default.configuration": {
		"en": "Default configuration",
	},
	"label.Description": {
		"en": "Description",
	},
	"label.Error": {
		"en": "Error",
	},
	"label.Field": {
		"en": "Field",
	},
	"label.had.default.verb": {
		"en": "(*) indicates the default subcommand if none given",
	},
	"label.ID": {
		"en": "ID",
	},
	"label.Key": {
		"en": "Key",
	},
	"label.Logger": {
		"en": "Logger",
	},
	"label.Member": {
		"en": "Member",
	},
	"label.Name": {
		"en": "Name",
	},
	"label.Nullable": {
		"en": "Nullable",
	},
	"label.Parameters": {
		"en": "Parameters",
	},
	"label.Permissions": {
		"en": "Permissions",
	},
	"label.Row": {
		"en": "Row",
	},
	"label.Rows": {
		"en": "Rows",
	},
	"label.Schema": {
		"en": "Schema",
	},
	"label.Size": {
		"en": "Size",
	},
	"label.Table": {
		"en": "Table",
	},
	"label.Type": {
		"en": "Type",
	},
	"label.Unique": {
		"en": "Unique",
	},
	"label.Usage": {
		"en": "Usage",
	},
	"label.User": {
		"en": "User",
	},
	"label.Value": {
		"en": "Value",
	},
	"label.break.at": {
		"en": "Break at",
	},
	"label.command": {
		"en": "command",
	},
	"label.configuration": {
		"en": "configuration",
	},
	"label.debug.commands": {
		"en": "Debugger commands:",
	},
	"label.options": {
		"en": "options",
	},
	"label.parameter": {
		"en": "parameter",
	},
	"label.parameters": {
		"en": "parameters",
	},
	"label.password.prompt": {
		"en": "Password: ",
	},
	"label.since": {
		"en": "since",
	},
	"label.stepped.to": {
		"en": "Step to",
	},
	"label.symbols": {
		"en": "symbols",
	},
	"label.username.prompt": {
		"en": "Username: ",
	},
	"label.version": {
		"en": "version",
	},
	"msg.config.deleted": {
		"en": "Configuration {{name}} deleted",
	},
	"msg.config.written": {
		"en": "Configuration key {{key}} written",
	},
	"msg.debug.break.added": {
		"en": "Added break {{break}}",
	},
	"msg.debug.break.exists": {
		"en": "Breakpoint already set",
	},
	"msg.debug.error": {
		"en": "Debugger error, {{err}}",
	},
	"msg.debug.load.count": {
		"en": "Loaded {{count}} breakpoints",
	},
	"msg.debug.no.breakpoints": {
		"en": "No breakpoints defined",
	},
	"msg.debug.no.source": {
		"en": "No source available for debugging",
	},
	"msg.debug.return": {
		"en": "Return from entrypoint",
	},
	"msg.debug.save.count": {
		"en": "Saving {{count}} breakpoints",
	},
	"msg.debug.scope": {
		"en": "Symbol table scope:",
	},
	"msg.debug.start": {
		"en": "Start program with call to entrypoint {{name}}()",
	},
	"msg.enter.blank.line": {
		"en": "Enter a blank line to terminate command input",
	},
	"msg.logged.in": {
		"en": "Successfully logged in as {{user}}, valid until {{expires}}",
	},
	"msg.server.cache": {
		"en": "Server Cache, hostname {{host}}, ID {{id}}",
	},
	"msg.server.cache.assets": {
		"en": "There are {{count}} HTML assets in cache, for a total size of {{size}} bytes.",
	},
	"msg.server.cache.emptied": {
		"en": "Server cache emptied",
	},
	"msg.server.cache.no.assets": {
		"en": "There are no HTML assets cached.",
	},
	"msg.server.cache.no.services": {
		"en": "There are no service items in cache. The maximum cache size is {{limit}} items.",
	},
	"msg.server.cache.one.asset": {
		"en": "There is 1 HTML asset in cache, for a total size of {{size}} bytes.",
	},
	"msg.server.cache.one.service": {
		"en": "There is 1 service item in cache. The maximum cache size is {{limit}} items.",
	},
	"msg.server.cache.services": {
		"en": "There are {{count}} service items in cache. The maximum cache size is {{limit}} items.",
	},
	"msg.server.cache.updated": {
		"en": "Server cache size updated",
	},
	"msg.server.logs.file": {
		"en": "Server log file is {{name}}",
	},
	"msg.server.logs.no.retain": {
		"en": "Server does not retain previous log files",
	},
	"msg.server.logs.purged": {
		"en": "Purged {{count}} old log files",
	},
	"msg.server.logs.retains": {
		"en": "Server also retains last {{count}} previous log files",
	},
	"msg.server.logs.status": {
		"en": "Logging status, hostname {{host}}, ID {{id}}",
	},
	"msg.server.not.running": {
		"en": "Server not running",
	},
	"msg.server.started": {
		"en": "Server started as process {{pid}}",
	},
	"msg.server.status": {
		"en": "Ego {{version}}, pid {{pid}}, host {{host}}, session {{id}}",
	},
	"msg.server.stopped": {
		"en": "Server (pid {{pid}}) stopped",
	},
	"msg.table.created": {
		"en": "Created table {{name}} with {{count}} columns",
	},
	"msg.table.delete.count": {
		"en": "Deleted {{count}} tables",
	},
	"msg.table.deleted": {
		"en": "Table {{name}} deleted",
	},
	"msg.table.deleted.no.rows": {
		"en": "No rows deleted",
	},
	"msg.table.deleted.rows": {
		"en": "{{count}} rows deleted",
	},
	"msg.table.empty.rowset": {
		"en": "No rows in result",
	},
	"msg.table.insert.count": {
		"en": "Added {{count}} rows to table {{name}}",
	},
	"msg.table.no.insert": {
		"en": "Nothing to insert into table",
	},
	"msg.table.sql.no.rows": {
		"en": "No rows modified",
	},
	"msg.table.sql.one.row": {
		"en": "1 row modified",
	},
	"msg.table.sql.rows": {
		"en": "{{count}} rows modified",
	},
	"msg.table.update.count": {
		"en": "Updated {{count}} rows in table {{name}}",
	},
	"msg.table.user.permissions": {
		"en": "User {{user}} permissions for {{schema}}.{{table}} {{verb}}: {{perms}}",
	},
	"msg.user.added": {
		"en": "User {{user}} added",
	},
	"msg.user.deleted": {
		"en": "User {{user}} deleted",
	},
	"opt.address.port": {
		"en": "Specify address (and optionally port) of server",
	},
	"opt.config.force": {
		"en": "Do not signal error if option not found",
	},
	"opt.filter": {
		"en": "List of optional filter clauses",
	},
	"opt.global.format": {
		"en": "Specify text, json or indented output format",
	},
	"opt.global.log": {
		"en": "Loggers to enable",
	},
	"opt.global.log.file": {
		"en": "Name of file where log messages are written",
	},
	"opt.global.profile": {
		"en": "Name of profile to use",
	},
	"opt.global.quiet": {
		"en": "If specified, suppress extra messaging",
	},
	"opt.global.version": {
		"en": "Show version number of command line tool",
	},
	"opt.help.text": {
		"en": "Show this help text",
	},
	"opt.insecure": {
		"en": "Do not require X509 server certificate verification",
	},
	"opt.limit": {
		"en": "If specified, limit the result set to this many rows",
	},
	"opt.local": {
		"en": "Show local server status info",
	},
	"opt.logon.server": {
		"en": "URL of server to authenticate with",
	},
	"opt.password": {
		"en": "Password for logon",
	},
	"opt.port": {
		"en": "Specify port number of server",
	},
	"opt.run.auto.import": {
		"en": "Override auto-import configuration setting",
	},
	"opt.run.debug": {
		"en": "Run with interactive debugger",
	},
	"opt.run.disasm": {
		"en": "Display a disassembly of the bytecode before execution",
	},
	"opt.run.entry.point": {
		"en": "Name of entrypoint function (defaults to main)",
	},
	"opt.run.log": {
		"en": "Direct log output to this file instead of stdout",
	},
	"opt.run.optimize": {
		"en": "Enable bytecode optimizer",
	},
	"opt.run.project": {
		"en": "Source is a directory instead of a file",
	},
	"opt.run.static": {
		"en": "Specify value typing during program execution",
	},
	"opt.run.symbols": {
		"en": "Display symbol table",
	},
	"opt.scope": {
		"en": "Blocks can access any symbol in call stack",
	},
	"opt.server.delete.user": {
		"en": "Username to delete",
	},
	"opt.server.logging.disable": {
		"en": "List of loggers to disable",
	},
	"opt.server.logging.enable": {
		"en": "List of loggers to enable",
	},
	"opt.server.logging.file": {
		"en": "Show only the active log file name",
	},
	"opt.server.logging.keep": {
		"en": "Specify how many log files to keep",
	},
	"opt.server.logging.session": {
		"en": "Limit display to log entries for this session number",
	},
	"opt.server.logging.status": {
		"en": "Display the state of each logger",
	},
	"opt.server.run.cache": {
		"en": "Number of service programs to cache in memory",
	},
	"opt.server.run.code": {
		"en": "Enable /code endpoint",
	},
	"opt.server.run.debug": {
		"en": "Service endpoint to debug",
	},
	"opt.server.run.force": {
		"en": "If set, override existing PID file",
	},
	"opt.server.run.is.detached": {
		"en": "If set, server assumes it is already detached",
	},
	"opt.server.run.keep": {
		"en": "The number of log files to keep",
	},
	"opt.server.run.log": {
		"en": "File path of server log",
	},
	"opt.server.run.no.log": {
		"en": "Suppress server log",
	},
	"opt.server.run.not.secure": {
		"en": "If set, use HTTP instead of HTTPS",
	},
	"opt.server.run.realm": {
		"en": "Name of authentication realm",
	},
	"opt.server.run.sandbox": {
		"en": "File path of sandboxed area for file I/O",
	},
	"opt.server.run.static": {
		"en": "Specify value typing during service execution",
	},
	"opt.server.run.superuser": {
		"en": "Designate this user as a super-user with ROOT privileges",
	},
	"opt.server.run.users": {
		"en": "File with authentication JSON data",
	},
	"opt.server.run.uuid": {
		"en": "Sets the optional session UUID value",
	},
	"opt.server.show.id": {
		"en": "Display the UUID of each user",
	},
	"opt.server.user.pass": {
		"en": "Password to assign to user",
	},
	"opt.server.user.perms": {
		"en": "Permissions to grant to user",
	},
	"opt.server.user.user": {
		"en": "Username to create or update",
	},
	"opt.sql.file": {
		"en": "Filename of SQL command text",
	},
	"opt.sql.row.ids": {
		"en": "Include the row UUID in the output",
	},
	"opt.sql.row.numbers": {
		"en": "Include the row number in the output",
	},
	"opt.start": {
		"en": "If specified, start result set at this row",
	},
	"opt.symbol.allocation": {
		"en": "Allocation size (in symbols) when expanding storage for a symbol table ",
	},
	"opt.table.create.file": {
		"en": "File name containing JSON column info",
	},
	"opt.table.delete.filter": {
		"en": "Filter for rows to delete. If not specified, all rows are deleted",
	},
	"opt.table.grant.permission": {
		"en": "Permissions to set for this table updated",
	},
	"opt.table.grant.user": {
		"en": "User (if other than current user) to update",
	},
	"opt.table.insert.file": {
		"en": "File name containing JSON row info",
	},
	"opt.table.list.no.row.counts": {
		"en": "If specified, listing does not include row counts",
	},
	"opt.table.permission.user": {
		"en": "User (if other than current user) to list)",
	},
	"opt.table.permissions.user": {
		"en": "If specified, list only this user",
	},
	"opt.table.read.columns": {
		"en": "List of columns to display; default is all columns",
	},
	"opt.table.read.order.by": {
		"en": "List of optional columns use to sort output",
	},
	"opt.table.read.row.ids": {
		"en": "Include the row UUID column in the output",
	},
	"opt.table.read.row.numbers": {
		"en": "Include the row number in the output",
	},
	"opt.table.update.filter": {
		"en": "Filter for rows to update. If not specified, all rows are updated",
	},
	"opt.trace": {
		"en": "Display trace of bytecode execution",
	},
	"opt.username": {
		"en": "Username for logon",
	},
	"parm.address.port": {
		"en": "address:port",
	},
	"parm.config.key.value": {
		"en": "key=value",
	},
	"parm.file": {
		"en": "file",
	},
	"parm.file.or.path": {
		"en": "file or path",
	},
	"parm.key": {
		"en": "key",
	},
	"parm.name": {
		"en": "name",
	},
}
