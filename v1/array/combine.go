package array

import (
	"errors"
	"strconv"
	"strings"
)

type CombineNum struct {
	Length int    //长度
	Sep    string //分割符号
}

func (c *CombineNum) Parse(text string) []string {
	var list []string
	if text == "" {
		list = make([]string, c.Length)
	} else {
		list = strings.Split(text, c.Sep)
	}

	return list
}

func (c *CombineNum) Get(text string, id int) (int, error) {
	list := c.Parse(text)
	if len(list) <= id {
		return 0, errors.New("下标不存在")
	}

	num, _ := strconv.Atoi(list[id])
	return num, nil
}

func (c *CombineNum) Set(text string, id int, num int) (string, error) {
	list := c.Parse(text)
	if len(list) <= id {
		return text, errors.New("下标不存在")
	}

	list[id] = strconv.Itoa(num)
	return strings.Join(list, "."), nil
}
