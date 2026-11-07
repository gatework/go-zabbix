# Zabbix API language

These terms distinguish Zabbix concepts with different ownership or protocol
behavior.

## Language

**API token**: A credential issued for API access without a user login session.
It has its own permissions and expiration; session logout does not revoke it.

**Login session**: Authentication state created by a username/password login.
Its session ID authenticates requests until logout or expiration.
_Avoid_: API token when referring to session lifecycle.

**Network discovery rule**: A rule that discovers hosts and network services
within a network range.

**Low-level discovery rule**: A rule that discovers resources belonging to a
host and creates entities from prototypes.
_Avoid_: Network discovery rule for low-level discovery.

**Host macro**: A user macro belonging to a specific host or template.

**Global macro**: A user macro defined independently of a host or template.

**Entity tag**: A name/value pair attached to a monitored entity.

**Tag filter**: A tag name, value and comparison operator used to select entities.
_Avoid_: Entity tag for a filtering condition.

**History**: Individual values collected for a monitored item.

**Trend**: An hourly aggregate of a numeric item's values, including the count,
average, minimum and maximum.
