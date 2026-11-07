package zabbix

import "errors"

// HostInterfaceMatchInput identifies interfaces by their exact address and port.
// An empty DNS or IP is meaningful and is always sent.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/massremove
type HostInterfaceMatchInput struct {
	IP   string `json:"ip"`
	DNS  string `json:"dns"`
	Port string `json:"port"`
}

func (p HostInterfaceMatchInput) validate() error {
	if p.Port == "" {
		return errors.New("zabbix: HostInterfaceMatchInput.Port is required")
	}
	return nil
}

// HostInterfaceReplaceParams replaces all interfaces on a host. Existing
// interfaces are identified by InterfaceID; entries without it create new ones.
// A non-nil empty Interfaces slice removes every interface, when Zabbix permits.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/replacehostinterfaces
type HostInterfaceReplaceParams struct {
	HostID     string                     `json:"hostid"`
	Interfaces []HostInterfaceUpdateInput `json:"interfaces"`
}

func (p HostInterfaceReplaceParams) validate() error {
	if err := requireWriteID(p.HostID, "HostInterfaceReplaceParams.HostID"); err != nil {
		return err
	}
	if p.Interfaces == nil {
		return errors.New("zabbix: HostInterfaceReplaceParams.Interfaces is required")
	}
	for _, value := range p.Interfaces {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}
