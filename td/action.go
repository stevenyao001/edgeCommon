package td

import (
	"bytes"
	"database/sql"
	"fmt"
)

type TDAction interface {
	AddLabelData()
	AddTagData()
	AddSuperTagData()
	Create()
	CreateSuper()
	Query()
	Insert()
	Delete()
	DeleteSuper()
}

type TDGroup struct {
	DB        *sql.DB
	Labels    []SQLLabel
	Tags      []Tags
	SuperTags []SuperTags
}

type Tags struct {
	Name     string
	TagsType string
}

type SuperTags struct {
	Name     string
	TagsType string
}

type SQLLabel struct {
	Name string
	Data interface{}
}

func (group *TDGroup) AddLabelData(name string, data interface{}) {
	tmp := SQLLabel{
		Name: name,
		Data: data,
	}
	group.Labels = append(group.Labels, tmp)
}

func (group *TDGroup) AddTagData(name string, tagsType string) {
	tmp := Tags{
		Name:     name,
		TagsType: tagsType,
	}
	group.Tags = append(group.Tags, tmp)
}

func (group *TDGroup) AddSuperTagData(name string, tagsType string) {
	tmp := SuperTags{
		Name:     name,
		TagsType: tagsType,
	}
	group.SuperTags = append(group.SuperTags, tmp)
}

func (group *TDGroup) Create(table string) error {
	defer group.cleanTagData()
	buffer := bytes.Buffer{}
	buffer.WriteString("create table ")
	buffer.WriteString(table)
	buffer.WriteString("(ts timestamp")
	for _, v := range group.Tags {
		buffer.WriteString(", ")
		buffer.WriteString(v.Name)
		buffer.WriteString(" ")
		buffer.WriteString(v.TagsType)
	}
	buffer.WriteString(")")

	fmt.Println(buffer.String())

	_, err := group.DB.Exec(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

func (group *TDGroup) CreateSuper(table string) error {
	defer group.cleanTagData()
	defer group.cleanSuperTagData()
	bufferTmp := bytes.Buffer{}
	for _, v := range group.Tags {
		bufferTmp.WriteString(v.Name)
		bufferTmp.WriteString(" ")
		bufferTmp.WriteString(v.TagsType)
		bufferTmp.WriteString(" ,")
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("create stable ")
	buffer.WriteString(table)
	buffer.WriteString(" (ts timestamp")
	for _, v := range group.SuperTags {
		buffer.WriteString(", ")
		buffer.WriteString(v.Name)
		buffer.WriteString(" ")
		buffer.WriteString(v.TagsType)
	}
	buffer.WriteString(") tags (")
	buffer.WriteString(bufferTmp.String()[:bufferTmp.Len()-1])
	buffer.WriteString(")")

	fmt.Println(buffer.String())

	_, err := group.DB.Exec(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

func (group *TDGroup) Query(table string, order string, limit string) ([][]interface{}, error) {
	defer group.cleanLabelData()
	var (
		num    int
		replay [][]interface{}
		err    error
	)
	replay = make([][]interface{}, 0)
	sql := setSelectSQL(group, table, order, limit)
	rows, err := group.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataTmp := make([]interface{}, len(group.Labels))

	for rows.Next() {
		num = len(dataTmp)
		if num == 3 {
			err = scanThreeTmp(rows, dataTmp)
		} else if num == 5 {
			err = scanFiveTmp(rows, dataTmp)
		} else if num == 6 {
			err = scanSixTmp(rows, dataTmp)
		}
		if err != nil || dataTmp == nil {
			return nil, err
		}
		replay = append(replay, dataTmp)
	}
	return replay, nil
}

func (group *TDGroup) Insert(table string, ts int64) error {
	defer group.cleanLabelData()
	sql := setInsertSQL(group, table, ts)
	_, err := group.DB.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (group *TDGroup) Delete(table string) error {
	buffer := bytes.Buffer{}
	buffer.WriteString("drop table if exists ")
	buffer.WriteString(table)

	_, err := group.DB.Exec(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

func (group *TDGroup) DeleteSuper(table string) error {
	buffer := bytes.Buffer{}
	buffer.WriteString("drop stable if exists ")
	buffer.WriteString(table)

	_, err := group.DB.Exec(buffer.String())
	if err != nil {
		return err
	}
	return nil
}

func (group *TDGroup) cleanLabelData() {
	group.Labels = group.Labels[0:0]
}

func (group *TDGroup) cleanTagData() {
	group.Tags = group.Tags[0:0]
}

func (group *TDGroup) cleanSuperTagData() {
	group.SuperTags = group.SuperTags[0:0]
}

func setInsertSQL(group *TDGroup, table string, ts int64) string {
	buffer := bytes.Buffer{}
	buffer.WriteString("insert into ")
	buffer.WriteString(table)
	buffer.WriteString(" (ts")
	for _, v := range group.Labels {
		buffer.WriteString(" ,")
		buffer.WriteString(v.Name)
	}
	buffer.WriteString(") values (")
	buffer.WriteString(fmt.Sprintf("%d", ts))
	for _, v := range group.Labels {
		buffer.WriteString(", ")
		switch v.Data.(type) {
		case int:
			buffer.WriteString(fmt.Sprintf("%d", v.Data))
		case int32:
			buffer.WriteString(fmt.Sprintf("%d", v.Data))
		case int64:
			buffer.WriteString(fmt.Sprintf("%d", v.Data))
		case float64:
			buffer.WriteString(fmt.Sprintf("%f", v.Data))
		case float32:
			buffer.WriteString(fmt.Sprintf("%f", v.Data))
		case bool:
			buffer.WriteString(fmt.Sprintf("%t", v.Data))
		case string:
			buffer.WriteString(fmt.Sprintf("%s", v.Data))
		}
	}
	buffer.WriteString(")")
	return buffer.String()
}

func setSelectSQL(group *TDGroup, table string, order string, limit string) string {
	bufferTmp := bytes.Buffer{}
	bufferTmp.WriteString("select ")
	for _, v := range group.Labels {
		bufferTmp.WriteString(v.Name)
		bufferTmp.WriteString(",")
	}
	buffer := bytes.Buffer{}
	buffer.WriteString(bufferTmp.String()[:bufferTmp.Len()-1])
	buffer.WriteString(" from ")
	buffer.WriteString(table)
	buffer.WriteString(" ")
	buffer.WriteString(order)
	buffer.WriteString(" ")
	buffer.WriteString(limit)
	return buffer.String()
}

func scanThreeTmp(rows *sql.Rows, data []interface{}) error {
	err := rows.Scan(&data[0], &data[1], &data[2])
	if err != nil {
		return err
	}
	return nil
}

func scanFiveTmp(rows *sql.Rows, data []interface{}) error {
	err := rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4])
	if err != nil {
		return err
	}
	return nil
}

func scanSixTmp(rows *sql.Rows, data []interface{}) error {
	err := rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4], &data[5])
	if err != nil {
		return err
	}
	return nil
}
