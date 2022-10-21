package tx

import (
	fm "file_manager"
)

type ConCurrencyManager struct {
	lock_table *LockTable
	lock_map   map[fm.BlockId]string
}

func NewConcurrencyManager() *ConCurrencyManager {
	concurrency_mgr := &ConCurrencyManager{
		lock_table: GetLockTableInstance(),
		lock_map:   make(map[fm.BlockId]string),
	}

	return concurrency_mgr
}

func (c *ConCurrencyManager) SLock(blk *fm.BlockId) error {
	_, ok := c.lock_map[*blk]
	if !ok {
		err := c.lock_table.SLock(blk)
		if err != nil {
			return err
		}
		c.lock_map[*blk] = "S"
	}
	return nil
}

func (c *ConCurrencyManager) XLock(blk *fm.BlockId) error {
	if !c.hasXLock(blk) {
		//c.SLock(blk) //判断区块是否已经被加上共享锁，如果别人已经获得共享锁那么就会挂起
		err := c.lock_table.XLock(blk)
		if err != nil {
			return err
		}
		c.lock_map[*blk] = "X"
	}

	return nil
}

func (c *ConCurrencyManager) Release() {
	for key, _ := range c.lock_map {
		c.lock_table.UnLock(&key)
	}
}

func (c *ConCurrencyManager) hasXLock(blk *fm.BlockId) bool {
	lock_type, ok := c.lock_map[*blk]
	return ok && lock_type == "X"
}
