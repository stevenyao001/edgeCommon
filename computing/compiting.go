package computing

import (
	"encoding/json"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"io/ioutil"
	"regexp"
	"strconv"
)

type ruleData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Rule string `json:"rule"`
}

type InputTmp struct {
	Ts         int64                  `json:"ts"`
	Properties map[string]interface{} `json:"properties"`
}

type Commputing struct{}

var middle map[string]interface{}

func (c *Commputing) InitMiddle(path string) {
	middle = make(map[string]interface{}, 0)

	ruleD := make([]ruleData, 5)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &ruleD)
	if err != nil {
		panic(err)
	}
	for _, value := range ruleD {
		if value.Type == "int" {
			middle[value.Name] = 0
		} else if value.Type == "bool" {
			middle[value.Name] = true
		}
	}
}

func (c *Commputing) Computing(buf []byte, path string) ([]byte, error) {
	input, err := GetData(buf)
	ruleDatas, err := GetRule(input, path)
	if err != nil {
		return nil, err
	}

	framework(ruleDatas, input)

	pushMiddleData(input, ruleDatas)
	return json.Marshal(input)
}

func framework(ruleDatas []ruleData, input map[string]interface{}) {
	var ruleTtmp string
	for _, value := range ruleDatas {
		ruleTtmp = ruleTtmp + value.Rule + "\n"
	}
	//fmt.Println("---------", ruleTtmp)
	dataContext := context.NewDataContext()
	propertiesAdd(ruleDatas, input, middle, dataContext)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	err := ruleBuilder.BuildRuleFromString(ruleTtmp)
	if err != nil {
		panic(err)
	}
	gengine := engine.NewGengine()
	err = gengine.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}
	resultMap, _ := gengine.GetRulesResultMap()
	for key, value := range ruleDatas {
		if value.Rule != "" {
			input[value.Name] = resultMap[strconv.Itoa(key)]
		}
	}
}

func propertiesAdd(ruleDatas []ruleData, input map[string]interface{}, middle map[string]interface{}, dataContext *context.DataContext) {
	for _, value := range ruleDatas {
		if input[value.Name] == nil {
			input[value.Name] = middle[value.Name]
		}
		if middle[value.Name] == nil {
			return
		}
		switch input[value.Name].(type) {
		case int:
			dataContext.Add(value.Name, input[value.Name].(int))
		case bool:
			dataContext.Add(value.Name, input[value.Name].(bool))
		case float64:
			dataContext.Add(value.Name, input[value.Name].(float64))
		case int64:
			dataContext.Add(value.Name, input[value.Name].(int64))
		case int32:
			dataContext.Add(value.Name, input[value.Name].(int32))
		}
		switch middle[value.Name].(type) {
		case int:
			dataContext.Add(value.Name+"_last", middle[value.Name].(int))
		case bool:
			dataContext.Add(value.Name+"_last", middle[value.Name].(bool))
		case float64:
			dataContext.Add(value.Name+"_last", middle[value.Name].(float64))
		case int64:
			dataContext.Add(value.Name+"_last", middle[value.Name].(int64))
		case int32:
			dataContext.Add(value.Name+"_last", middle[value.Name].(int32))
		}
	}
}

func pushMiddleData(input map[string]interface{}, ruleDatas []ruleData) {
	for _, value := range ruleDatas {
		middle[value.Name] = input[value.Name]
	}
}

func GetRule(input map[string]interface{}, path string) ([]ruleData, error) {
	ruleDatas := make([]ruleData, 0)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &ruleDatas)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, value := range ruleDatas {
		names = append(names, value.Name)
	}

	for key, value := range ruleDatas {
		if value.Rule != "" {
			rule := `
rule "` + strconv.Itoa(key) + `" "test"
begin
`
			rule = rule + ruleTmpComputer(input, value.Rule, names)
			rule = rule + "\nend"
			ruleDatas[key].Rule = rule
			//rules = append(rules, rule)
		}
	}
	return ruleDatas, nil
}

func GetData(buf []byte) (map[string]interface{}, error) {
	tmp := new(InputTmp)
	err := json.Unmarshal(buf, tmp)
	if err != nil {
		return nil, err
	}
	return tmp.Properties, nil
}

func ruleTmpComputer(input map[string]interface{}, rule string, names []string) string {
	for _, name := range names {
		/*nameTmp2 := "[" + name + "_last]" + "{" + strconv.Itoa(len(name)+5) + ",}"
		replace2 := " middle.Properties[" + name + "]"
		regIf2 := regexp.MustCompile(nameTmp2)
		if regIf2 != nil {
			rule = regIf2.ReplaceAllString(rule, replace2)
		}

		nameTmp := "[ ]+[" + name + "]" + "{" + strconv.Itoa(len(name)+5) + ",}"
		replace := " input.Properties[" + name + "]"
		regIf := regexp.MustCompile(nameTmp)
		if regIf != nil {
			rule = regIf.ReplaceAllString(rule, replace)
		}*/
		if name == "threadhold" {
			nameTmp3 := "[ ]+[" + name + "]" + "{" + strconv.Itoa(len(name)+5) + ",}"
			replace := " 40"
			regIf := regexp.MustCompile(nameTmp3)
			if regIf != nil {
				rule = regIf.ReplaceAllString(rule, replace)
			}
		}
	}

	return rule
}
