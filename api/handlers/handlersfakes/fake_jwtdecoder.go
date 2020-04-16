// Code generated by counterfeiter. DO NOT EDIT.
package handlersfakes

import (
	"sync"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
)

type FakeJWTDecoder struct {
	ClaimSetStub        func(string) (map[string]interface{}, error)
	claimSetMutex       sync.RWMutex
	claimSetArgsForCall []struct {
		arg1 string
	}
	claimSetReturns struct {
		result1 map[string]interface{}
		result2 error
	}
	claimSetReturnsOnCall map[int]struct {
		result1 map[string]interface{}
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeJWTDecoder) ClaimSet(arg1 string) (map[string]interface{}, error) {
	fake.claimSetMutex.Lock()
	ret, specificReturn := fake.claimSetReturnsOnCall[len(fake.claimSetArgsForCall)]
	fake.claimSetArgsForCall = append(fake.claimSetArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ClaimSet", []interface{}{arg1})
	fake.claimSetMutex.Unlock()
	if fake.ClaimSetStub != nil {
		return fake.ClaimSetStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.claimSetReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeJWTDecoder) ClaimSetCallCount() int {
	fake.claimSetMutex.RLock()
	defer fake.claimSetMutex.RUnlock()
	return len(fake.claimSetArgsForCall)
}

func (fake *FakeJWTDecoder) ClaimSetCalls(stub func(string) (map[string]interface{}, error)) {
	fake.claimSetMutex.Lock()
	defer fake.claimSetMutex.Unlock()
	fake.ClaimSetStub = stub
}

func (fake *FakeJWTDecoder) ClaimSetArgsForCall(i int) string {
	fake.claimSetMutex.RLock()
	defer fake.claimSetMutex.RUnlock()
	argsForCall := fake.claimSetArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeJWTDecoder) ClaimSetReturns(result1 map[string]interface{}, result2 error) {
	fake.claimSetMutex.Lock()
	defer fake.claimSetMutex.Unlock()
	fake.ClaimSetStub = nil
	fake.claimSetReturns = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeJWTDecoder) ClaimSetReturnsOnCall(i int, result1 map[string]interface{}, result2 error) {
	fake.claimSetMutex.Lock()
	defer fake.claimSetMutex.Unlock()
	fake.ClaimSetStub = nil
	if fake.claimSetReturnsOnCall == nil {
		fake.claimSetReturnsOnCall = make(map[int]struct {
			result1 map[string]interface{}
			result2 error
		})
	}
	fake.claimSetReturnsOnCall[i] = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeJWTDecoder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.claimSetMutex.RLock()
	defer fake.claimSetMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeJWTDecoder) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ handlers.JWTDecoder = new(FakeJWTDecoder)