package h

import (
	"sync"
	log "github.com/sirupsen/logrus"
)

type List struct {
   internalList []interface{}
	mux sync.Mutex
}


func MakeList(l []interface{}) List{
	return List{l, sync.Mutex{}} // return List{l, nil}
}


func (l *List) Map(mapper func(v interface{}, i int) interface{}) []interface{}{
   output:= []interface{}{}
   for i, v := range l.internalList {
   	output = append(output, mapper(v,i));
   }
   return output;
}

func (l *List) Filter(mapper func(v interface{}, i int) interface{}) []interface{}{
	output:= []interface{}{}
	for i, v := range l.internalList {
		if mapper(v,i) == true {
			output = append(output, v)
		}
	}
	return output;
}


func (l *List) Add(args ...interface{}) *List {
	l.mux.Lock()
	l.internalList = append(l.internalList, args...)
	log.Info(l.internalList)
	l.mux.Unlock()
	return l;
}


func (l *List) Push(v interface{}) *List {
	l.mux.Lock()
	l.internalList = append(l.internalList, v)
	log.Info("internal:",l.internalList)
	l.mux.Unlock()
	return l
}


func (l *List) Pop() interface{}{
	l.mux.Lock()
	length :=len(l.internalList);
	log.Info("the length is:", length)
	if length < 1 {
		return nil;
	}
	last := l.internalList[length-1]
	l.internalList = l.internalList[:length-1]
	l.mux.Unlock()
	return last;
}


func (l *List) GetLength() int {
	return len(l.internalList);
}


func (l *List) Shift() interface{} {
	l.mux.Lock()
	if len(l.internalList) < 1 {
		return nil;
	}
	first := l.internalList[0];
	l.internalList = l.internalList[1:]
	l.mux.Unlock()
	return first;
}


func (l *List) Unshift(v interface{}){
	l.mux.Lock()
	l.internalList = append([]interface{}{v}, l.internalList...)
	l.mux.Unlock()
}

func (l *List) GetValue() []interface{}{
	return l.internalList
  //return copy([]interface{}{},l.internalList)
}






