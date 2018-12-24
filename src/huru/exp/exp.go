package main

import (
	"github.com/sirupsen/logrus"
	"huru/h"
	"huru/utils"
)

func main()  {
	list := h.CreateXList();

	v := list.Filter(func(v interface{}, i int) interface{} {
		return v == 3;
	}).Map(func(v interface{}, i int) interface{} {
         return i;
	})

	logrus.Info("vvv:",v);
}

func mainx() {

	list := utils.FlattenDeep(
		[]interface{}{3, 4, 5, []interface{}{3, 4, 5}}, 1, 2, []interface{}{3, 4, 5, []interface{}{3, 4, 5}, []interface{}{3, 4, 5}})

	logrus.Info(list)

	x := h.List{}
	x.Add(3,5,6,7)
	logrus.Info("value:", x.GetValue())
	logrus.Info("len:", x.GetLength())
	logrus.Info(x.Pop())
	logrus.Info("len:", x.GetLength())

	mapped := x.Map(func(v interface{}, i int) interface{} {
		return v.(int)*3;
	})

	z := h.MakeList(mapped);

	logrus.Info(mapped)

	filtered := z.Filter(func(v interface{}, i int) interface{} {
		if v.(int) % 2 != 0 {
			return true;
		}
		return nil;
	})

	logrus.Info(filtered)
}
