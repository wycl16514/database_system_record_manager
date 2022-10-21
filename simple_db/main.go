package main

import (
	bmg "buffer_manager"
	fm "file_manager"
	"fmt"
	lm "log_manager"
	"math/rand"
	record_mgr "record_manager"
	"tx"
)

func main() {
	file_manager, _ := fm.NewFileManager("recordtest", 400)
	log_manager, _ := lm.NewLogManager(file_manager, "logfile.log")
	buffer_manager := bmg.NewBufferManager(file_manager, log_manager, 3)

	tx := tx.NewTransation(file_manager, log_manager, buffer_manager)
	sch := record_mgr.NewSchema()

	sch.AddIntField("A")
	sch.AddStringField("B", 9)
	layout := record_mgr.NewLayoutWithSchema(sch)
	for _, field_name := range layout.Schema().Fields() {
		offset := layout.Offset(field_name)
		fmt.Printf("%s has offset %d\n", field_name, offset)
	}

	blk, err := tx.Append("testfile")
	if err != nil {
		err_str := fmt.Sprintf("err : %v\n", err)
		panic(err_str)
	}
	tx.Pin(blk)
	rp := record_mgr.NewRecordPage(tx, blk, record_mgr.LayoutInterface(layout))
	rp.Format()
	fmt.Println("Filling the page with random records")
	slot := rp.InsertAfter(-1) //找到第一条可用插槽
	for slot >= 0 {
		n := rand.Intn(50)
		rp.SetInt(slot, "A", n)                          //找到可用插槽后随机设定字段A的值
		rp.SetString(slot, "B", fmt.Sprintf("rec%d", n)) //设定字段B
		fmt.Printf("inserting into slot :%d :{ %d , rec%d}\n", slot, n, n)
		slot = rp.InsertAfter(slot) //查找当前插槽之后可用的插槽
	}

	fmt.Println("Deleted these records with A-values < 25.")
	count := 0
	slot = rp.NextAfter(-1)
	for slot >= 0 {
		a := rp.GetInt(slot, "A")
		b := rp.GetString(slot, "B")
		if a < 25 {
			count += 1
			fmt.Printf("slot %d: {%d, %s}\n", slot, a, b)
			rp.Delete(slot)
		}
		slot = rp.NextAfter(slot)
	}
	fmt.Printf("%d values under 25 were deleted.\n", count)
	fmt.Println("Here are the remaining records")
	slot = rp.NextAfter(-1)
	for slot >= 0 {
		a := rp.GetInt(slot, "A")
		b := rp.GetString(slot, "B")
		fmt.Printf("slot %d : {%d, %s}\n", slot, a, b)
		slot = rp.NextAfter(slot)
	}

	tx.UnPin(blk)
	tx.Commit()
}
