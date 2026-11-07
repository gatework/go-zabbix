package zabbix

import "errors"

func (p EventAcknowledgeParams) validateResult(ids []string) error {
	// PHP deduplicates the original strings before converting array keys to
	// numbers. "001" and "1" therefore remain distinct requested entries.
	seen := make(map[string]struct{}, len(p.EventIDs))
	expected := make([]string, 0, len(p.EventIDs))
	for _, id := range p.EventIDs {
		if _, ok := seen[id]; !ok {
			expected = append(expected, id)
			seen[id] = struct{}{}
		}
	}
	return correspondingIDs(ids, expected)
}

func (p HostInterfaceMassAddParams) validateResult(ids []string) error {
	// Division avoids overflowing the size of the Cartesian product.
	if len(ids)%len(p.Hosts) != 0 || len(ids)/len(p.Hosts) != len(p.Interfaces) {
		return errors.New("created interface count does not match request")
	}
	return uniqueIDs(ids)
}

func (p HostInterfaceMassRemoveParams) validateResult(ids []string) error {
	return uniqueIDs(ids)
}

func (p HostInterfaceReplaceParams) validateResult(ids []string) error {
	if len(ids) != len(p.Interfaces) {
		return errors.New("replacement interface count does not match request")
	}
	existing := make(map[string]struct{}, len(ids))
	for i, item := range p.Interfaces {
		if item.InterfaceID == nil {
			continue
		}
		id := canonicalID(*item.InterfaceID)
		if id != canonicalID(ids[i]) {
			return errors.New("replacement interface id does not match its input position")
		}
		existing[id] = struct{}{}
	}
	for i, item := range p.Interfaces {
		if item.InterfaceID != nil {
			continue
		}
		id := canonicalID(ids[i])
		if _, ok := existing[id]; ok {
			return errors.New("replacement generated a duplicate interface id")
		}
		existing[id] = struct{}{}
	}
	return nil
}
