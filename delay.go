package main

import (
        "fmt"
        "time"
        "math"
)

type Data struct {
        round int
        index int
}

type Item struct {
        round int
        index int
        f FormatFunc
}

type FormatFunc func()

const NUM = 3600

var itemMap map[int64]*Item

func custom(ch <-chan FormatFunc) {
        for {
                select {
                case f := <-ch:
                        go func() {
                                f()
                        }()
                }
        }
}

func main() {
        data := Data{}
        item_func := map[string]FormatFunc{"test1":test1, "test2":test2}
        itemMap = make(map[int64]*Item)
        ch := make(chan FormatFunc) 

        go custom(ch)
        go run(&data, ch)

        var (
                tt int64
                ff string
        )
        for {
                fmt.Scanln(&tt, &ff)
                itemMap[tt] = &Item{
                        f: item_func[ff],
                        round:1,
                        index:1,
                }

                length := int(tt - time.Now().Unix())
                itemMap[tt].round = data.round + int(math.Floor(float64((data.index+length) / NUM)))
                itemMap[tt].index = (data.index+length) % NUM
                fmt.Println(*itemMap[tt])
        }
}

func test1() {
        fmt.Println(time.Now().Unix())
        fmt.Println("test1")
}

func test2() {
        fmt.Println(time.Now().Unix())
        fmt.Println("test2")
}

func run(data *Data, ch chan<- FormatFunc) {
        for {
                <-time.After(time.Second)
                if data.index%NUM == 0 {
                        data.round++
                        data.index = 1
                } else {
                        data.index++
                }

                //fmt.Printf("round is %d, index is %d\n", data.round, data.index)

                for k,v := range itemMap {
                        if data.round == v.round && data.index == v.index {
                                ch <- v.f
                                delete(itemMap, k)
                        }
                }       
        }

}