package test

import (
	"fmt"
	"testing"
	"time"
	"yueja/go-redisson/lock"
)

// TestDemo1 加锁简单测试
func TestDemo1(t *testing.T) {
	var (
		locker *lock.Locker
		err    error
	)
	if locker, err = lock.Lock("TestDemo1_", &lock.Options{RetryCount: 5, LockTimeout: 20 * time.Second}); err != nil {
		return
	}
	defer locker.Unlock()
	time.Sleep(10 * time.Second)
}

// TestDemo2 看门狗简单测试
func TestDemo2(t *testing.T) {
	var (
		locker *lock.Locker
		err    error
	)
	if locker, err = lock.Lock("TestDemo1_", &lock.Options{RetryCount: 5, LockTimeout: 20 * time.Second}); err != nil {
		return
	}
	time.Sleep(1 * time.Minute)
	if err = locker.Unlock(); err != nil {
		return
	}
}

// TestDemo3 并发加锁测试，同时测试加锁重试
func TestDemo3(t *testing.T) {
	var (
		err error
	)
	for i := 0; i < 5; i++ {
		go func(index int) {
			locker := new(lock.Locker)
			if locker, err = lock.Lock("TestDemo1_", &lock.Options{
				RetryCount:  50,
				RetryDelay:  200 * time.Millisecond,
				LockTimeout: 5 * time.Second,
				Index:       index}); err != nil {
				return
			}
			// Sleep 模拟业务处理耗时
			time.Sleep(2 * time.Second)
			defer locker.Unlock()
		}(i)
	}
	time.Sleep(20 * time.Second)
}

// TestDemo4 并发加锁测试，同时测试加锁重试和锁续命
func TestDemo4(t *testing.T) {
	var (
		err error
	)
	for i := 0; i < 5; i++ {
		go func(index int) {
			locker := new(lock.Locker)
			if locker, err = lock.Lock("TestDemo1_", &lock.Options{
				RetryCount:  50,
				RetryDelay:  300 * time.Millisecond,
				LockTimeout: 2 * time.Second,
				Index:       index}); err != nil {
				return
			}
			// Sleep 模拟业务处理耗时
			time.Sleep(5 * time.Second)
			defer locker.Unlock()
		}(i)
	}
	time.Sleep(20 * time.Second)
}

// TestDemo5 并发加锁测试，同时测试加锁重试，锁重入和锁续命
func TestDemo5(t *testing.T) {
	var (
		err error
	)
	for i := 0; i < 5; i++ {
		go func(index int) {
			locker := new(lock.Locker)
			if locker, err = lock.Lock("TestDemo1_", &lock.Options{
				RetryCount:  60,
				RetryDelay:  300 * time.Millisecond,
				LockTimeout: 2 * time.Second,
				Index:       index}); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
			var ok bool
			if ok, err = locker.Lock(); err != nil {
				panic(err)
			}
			if !ok {
				panic("锁重入失败：")
			}
			// Sleep 模拟业务处理耗时
			time.Sleep(2 * time.Second)
			locker.Unlock()
			time.Sleep(2 * time.Second)
			locker.Unlock()
		}(i)
	}
	time.Sleep(20 * time.Second)
}

// TestDemo6 锁重入(2次重入，即加锁3次，2次解锁)
func TestDemo6(t *testing.T) {
	var (
		locker *lock.Locker
		ok     bool
		err    error
	)
	if locker, err = lock.Lock("TestDemo1_", &lock.Options{RetryCount: 5, LockTimeout: 10 * time.Second}); err != nil {
		panic(err)
		return
	}
	time.Sleep(25 * time.Second)
	if ok, err = locker.Lock(); err != nil {
		panic(err)
		return
	}
	if !ok {
		panic("锁重入失败")
	}
	fmt.Println("锁重入成功 001")

	time.Sleep(5 * time.Second)
	if ok, err = locker.Lock(); err != nil {
		panic(err)
		return
	}
	if !ok {
		panic("锁重入失败")
	}
	fmt.Println("锁重入成功 002")

	time.Sleep(25 * time.Second)
	if err = locker.Unlock(); err != nil {
		panic(err)
		return
	}
	time.Sleep(10 * time.Second)
	if err = locker.Unlock(); err != nil {
		panic(err)
		return
	}
	fmt.Println("锁完全释放")
	time.Sleep(15 * time.Second)
}

// TestDemo7 锁重入(2次重入，即加锁3次，3次解锁)
func TestDemo7(t *testing.T) {
	var (
		locker *lock.Locker
		ok     bool
		err    error
	)
	if locker, err = lock.Lock("TestDemo1_", &lock.Options{RetryCount: 5, LockTimeout: 10 * time.Second}); err != nil {
		panic(err)
		return
	}
	time.Sleep(25 * time.Second)
	if ok, err = locker.Lock(); err != nil {
		panic(err)
		return
	}
	if !ok {
		panic("锁重入失败")
	}
	fmt.Println("锁重入成功 001")

	time.Sleep(5 * time.Second)
	if ok, err = locker.Lock(); err != nil {
		panic(err)
		return
	}
	if !ok {
		panic("锁重入失败")
	}
	fmt.Println("锁重入成功 002")

	time.Sleep(25 * time.Second)
	if err = locker.Unlock(); err != nil {
		panic(err)
		return
	}
	time.Sleep(10 * time.Second)
	if err = locker.Unlock(); err != nil {
		panic(err)
		return
	}
	if err = locker.Unlock(); err != nil {
		panic(err)
		return
	}
	fmt.Println("锁完全释放")
	time.Sleep(15 * time.Second)
}
