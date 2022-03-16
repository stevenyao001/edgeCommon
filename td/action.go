package td

import (
	"bytes"
	"database/sql"
	"fmt"
)

type TDAction interface {
	AddLabelData()
	Create()
	CreateSuper()
	Query()
	Insert()
	Delete()
	DeleteSuper()
}

type TDGroup struct {
	DB     *sql.DB
	Labels []SQLLabel
}

type Tags struct {
	Name     string
	TagsType string
}

type TagsSuper struct {
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

func (group *TDGroup) Create(table string, data []Tags) error {
	buffer := bytes.Buffer{}
	buffer.WriteString("create table ")
	buffer.WriteString(table)
	buffer.WriteString("(ts timestamp")
	for _, v := range data {
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

func (group *TDGroup) CreateSuper(table string, tags []Tags, tagsSuper []TagsSuper) error {
	bufferTmp := bytes.Buffer{}
	for _, v := range tags {
		bufferTmp.WriteString(v.Name)
		bufferTmp.WriteString(" ")
		bufferTmp.WriteString(v.TagsType)
		bufferTmp.WriteString(" ,")
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("create stable ")
	buffer.WriteString(table)
	buffer.WriteString(" (ts timestamp")
	for _, v := range tagsSuper {
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

func (group *TDGroup) Query(table string) ([][]interface{}, error) {
	var (
		num        int
		scanReplay []interface{}
		replay     [][]interface{}
		err        error
	)
	sql := setSelectSQL(group, table)
	rows, err := group.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataTmp := make([]interface{}, len(group.Labels))

	for rows.Next() {
		for rows.Next() {
			num = len(dataTmp)
			if num == 3 {
				scanReplay, err = scanThreeTmp(rows, dataTmp)
			} else if num == 5 {
				scanReplay, err = scanFiveTmp(rows, dataTmp)
			}
			if err != nil {
				return nil, err
			}
			replay = append(replay, scanReplay)
		}
	}
	return replay, nil
}

func (group *TDGroup) Insert(table string, ts int64) error {
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

func setSelectSQL(group *TDGroup, table string) string {
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
	buffer.WriteString(" order by ts desc limit 1")
	return buffer.String()
}

func scanThreeTmp(rows *sql.Rows, data []interface{}) (replay []interface{}, err error) {
	err = rows.Scan(&data[0], &data[1], &data[2])
	if err != nil {
		return nil, err
	}
	replay = append(replay, data[0], data[1], data[2])
	return replay, nil
}

func scanFiveTmp(rows *sql.Rows, data []interface{}) (replay []interface{}, err error) {
	err = rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4])
	if err != nil {
		return nil, err
	}
	replay = append(replay, data[0], data[1], data[2], data[3], data[4])
	return replay, nil
}
