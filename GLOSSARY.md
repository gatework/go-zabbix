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

**Service access tag**: A role permission condition selecting services by tag.
An empty tag clears the condition; unlike an entity tag, it is not attached to
the service itself.

**HTTP agent item**: An item whose value is collected by an HTTP request. Its
ordered header and query-field lists retain repeated names and empty values.

**Preprocessing step**: A typed transformation or validation of a collected item
value, with optional failure handling. A replacement step supplies its complete
configuration.

**Mutation receipt**: The server's affected IDs, updated field names or boolean
success result. Correspondence to the request is a protocol check; persisted
state requires a separate read-back.

**History**: Individual values collected for a monitored item.

**Trend**: An hourly aggregate of a numeric item's values, including the count,
average, minimum and maximum.

**Monitored host**: A validated host identity returned by the bounded monitoring queries.

**Monitored item**: A validated item identity and latest observation. A network
interface is identified by its host and traffic item IDs, not `hostinterfaceid`,
which identifies an SNMP/agent connection endpoint.

**Traffic series**: Rate observations normalized to bits per second, preserving
missing samples, observed peaks and whether the requested window is incomplete.

**Frontend session**: The signed cookie issued by the Zabbix web login. It is
separate from a JSON-RPC API token or login-session bearer credential.

**Chart image**: The PNG rendered by the Zabbix frontend for item IDs or a graph
ID and an explicit time range. A graph API object describes configuration only.
