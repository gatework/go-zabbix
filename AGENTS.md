# go-zabbix

Follow [the gatework Go engineering standard](https://github.com/gatework/go-standards/blob/v1.0.2/STANDARD.md) (go-standards/STANDARD.md); this file adds SDK-specific rules. Use golang-how-to for Go work.

- Keep this standard-library-only Zabbix 7.4 SDK in the root `zabbix` package. Namespace files expose API operations; shared request, authentication and ownership logic belongs in client/auth/options files.
- Read [CLAUDE.md](CLAUDE.md), [GLOSSARY.md](GLOSSARY.md) and the relevant versioned row of [the contract matrix](docs/api-7.4-contracts.md) before protocol changes.
- Preserve decimal-string IDs, atomic output decoding and omitted/empty/zero distinctions. Read models and mutation inputs remain separate. Credentials and payloads stay out of default errors and logs.
- Public API changes update MIGRATION.md, README examples and compile-checked examples together.
- Run `make check` before committing. This includes build, shuffled race tests and integration-test compilation. Protocol changes also run `make integration` against the disposable fixture; never point mutation tests at a development database.
- CI exercises Go 1.26.x and 1.27.x. The 1.27 job adds the disposable Zabbix round trips through CHECK_EXTRA.
