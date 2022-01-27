package database

const (
	COND_GRATE = 10
	COND_EQUAL = 15
	COND_LESS  = 20
)

type TYPE_COND int
type WhereCond struct {
	field string
	value interface{}
	cond  TYPE_COND
}

func NewWhereCond(field string, value interface{}, cond TYPE_COND) *WhereCond {
	return &WhereCond{
		field: field,
		value: value,
		cond:  cond,
	}
}

type IDatabase interface {
	Insert(data *InsertPayLoad) (int, error)
	InsertArray(data []*InsertPayLoad) (int, error)
	Update(data *UpdatePayLoad) (int, error)
	Delete(data *DeletePayLoad) (int, error)
	Select(data *SelectPayLoad) ([]interface{}, error)
}

type InsertPayLoad struct {
	tid  string
	data map[string]interface{}
}

func NewInsertPayLoad(tid string, data map[string]interface{}) *InsertPayLoad {
	return &InsertPayLoad{
		tid:  tid,
		data: data,
	}
}

type UpdatePayLoad struct {
	tid        string
	conditions []WhereCond
	data       map[string]interface{}
}

func NewUpdatePayLoad(tid string, cond []WhereCond, data map[string]interface{}) *UpdatePayLoad {
	return &UpdatePayLoad{
		tid:        tid,
		conditions: cond,
		data:       data,
	}
}

type DeletePayLoad UpdatePayLoad

type SelectPayLoad struct {
	tid        string
	field      []string
	conditions []WhereCond
}

func NewSelectPayLoad(tid string, field []string, cond []WhereCond) *SelectPayLoad {
	return &SelectPayLoad{
		tid:        tid,
		field:      field,
		conditions: cond,
	}
}
