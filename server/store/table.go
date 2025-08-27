package store

type TableItem struct {
	Key   string
	Value string
}

func (s *Store) GetTable() []TableItem {
	var table []TableItem

	for k := range s.data {
		table = append(table, TableItem{
			Key:   k,
			Value: string(*s.data[k].Value),
		})
	}

	return table
}
