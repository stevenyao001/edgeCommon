package td

import (
	"bytes"
	"database/sql"
	"fmt"
)

type DataAction interface {
	Create()
	CreateSuper()
	Query()
	Insert()
	Delete()
	DeleteSuper()
}

type TDGroup struct {
	DB *sql.DB
}

type Tags struct {
	Name     string
	TagsType string
}

type TagsSuper struct {
	Name     string
	TagsType string
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

func (group *TDGroup) Query(sql string, scan func(rows *sql.Rows)) error {
	rows, err := group.DB.Query(sql)
	if err != nil {
		return err
	}
	for rows.Next() {
		scan(rows)
	}
	rows.Close()
	return nil
}

func (group *TDGroup) Insert(sql string) error {
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
