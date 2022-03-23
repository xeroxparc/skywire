// Code generated by mockery v1.0.0. DO NOT EDIT.

package setupclient

import (
	context "context"

	dmsg "github.com/skycoin/dmsg"
	logging "github.com/skycoin/skycoin/src/util/logging"
	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skywire-utilities/pkg/cipher"
	routing "github.com/skycoin/skywire/pkg/routing"
)

// MockRouteGroupDialer is an autogenerated mock type for the RouteGroupDialer type
type MockRouteGroupDialer struct {
	mock.Mock
}

// Dial provides a mock function with given fields: ctx, log, dmsgC, setupNodes, req
func (_m *MockRouteGroupDialer) Dial(ctx context.Context, log *logging.Logger, dmsgC *dmsg.Client, setupNodes []cipher.PubKey, req routing.BidirectionalRoute) (routing.EdgeRules, error) {
	ret := _m.Called(ctx, log, dmsgC, setupNodes, req)

	var r0 routing.EdgeRules
	if rf, ok := ret.Get(0).(func(context.Context, *logging.Logger, *dmsg.Client, []cipher.PubKey, routing.BidirectionalRoute) routing.EdgeRules); ok {
		r0 = rf(ctx, log, dmsgC, setupNodes, req)
	} else {
		r0 = ret.Get(0).(routing.EdgeRules)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *logging.Logger, *dmsg.Client, []cipher.PubKey, routing.BidirectionalRoute) error); ok {
		r1 = rf(ctx, log, dmsgC, setupNodes, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
