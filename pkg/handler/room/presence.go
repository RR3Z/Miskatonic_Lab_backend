package room

import (
	"context"
	"log/slog"
	"sync"
	"time"

	roomModel "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/room"
	roomService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/room"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	DefaultRoomDisconnectGrace = 30 * time.Second
	presenceRetryDelay         = 5 * time.Second
)

type PresenceConfig struct {
	DisconnectGrace time.Duration
}

type roomPresenceKey struct {
	roomID string
	userID string
}

type roomPresenceState struct {
	roomID      pgtype.UUID
	connections int
	removing    bool
	timer       *time.Timer
}

type roomPresenceCoordinator struct {
	service roomService.IRoom
	grace   time.Duration
	onLeave func(roomModel.RoomMutationResult[roomModel.LeaveRoomResult])

	mu     sync.Mutex
	states map[roomPresenceKey]*roomPresenceState
}

func newRoomPresenceCoordinator(
	service roomService.IRoom,
	config PresenceConfig,
	onLeave func(roomModel.RoomMutationResult[roomModel.LeaveRoomResult]),
) *roomPresenceCoordinator {
	grace := config.DisconnectGrace
	if grace <= 0 {
		grace = DefaultRoomDisconnectGrace
	}

	return &roomPresenceCoordinator{
		service: service,
		grace:   grace,
		onLeave: onLeave,
		states:  make(map[roomPresenceKey]*roomPresenceState),
	}
}

func (p *roomPresenceCoordinator) AwaitFirstConnection(roomID pgtype.UUID, userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := roomPresenceKey{roomID: roomID.String(), userID: userID}
	state := p.states[key]
	if state == nil {
		state = &roomPresenceState{roomID: roomID}
		p.states[key] = state
	}
	p.scheduleLocked(key, state, p.grace)
}

func (p *roomPresenceCoordinator) Connected(roomID pgtype.UUID, userID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := roomPresenceKey{roomID: roomID.String(), userID: userID}
	state := p.states[key]
	if state == nil {
		state = &roomPresenceState{roomID: roomID}
		p.states[key] = state
	}
	if state.removing {
		return false
	}
	if state.timer != nil {
		state.timer.Stop()
		state.timer = nil
	}
	state.connections++
	return true
}

func (p *roomPresenceCoordinator) Disconnected(roomID pgtype.UUID, userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := roomPresenceKey{roomID: roomID.String(), userID: userID}
	state := p.states[key]
	if state == nil || state.removing || state.connections == 0 {
		return
	}

	state.connections--
	if state.connections == 0 {
		p.scheduleLocked(key, state, p.grace)
	}
}

func (p *roomPresenceCoordinator) Forget(roomID pgtype.UUID, userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.forgetLocked(roomPresenceKey{roomID: roomID.String(), userID: userID})
}

func (p *roomPresenceCoordinator) ForgetRoom(roomID pgtype.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for key := range p.states {
		if key.roomID == roomID.String() {
			p.forgetLocked(key)
		}
	}
}

func (p *roomPresenceCoordinator) scheduleLocked(key roomPresenceKey, state *roomPresenceState, delay time.Duration) {
	if state.removing || state.connections > 0 || state.timer != nil {
		return
	}
	state.timer = time.AfterFunc(delay, func() {
		p.expire(key)
	})
}

func (p *roomPresenceCoordinator) expire(key roomPresenceKey) {
	p.mu.Lock()
	state := p.states[key]
	if state == nil || state.removing || state.connections > 0 {
		p.mu.Unlock()
		return
	}
	state.timer = nil
	state.removing = true
	roomID := state.roomID
	p.mu.Unlock()

	result, err := p.service.LeaveRoom(context.Background(), roomModel.LeaveRoomInput{
		RoomID: roomID,
		UserID: key.userID,
	})
	if err != nil {
		slog.Warn("room auto-leave failed", "component", "room_presence", "room_id", key.roomID, "user_id", key.userID, "error", err)
		p.mu.Lock()
		state = p.states[key]
		if state != nil && state.removing {
			state.removing = false
			p.scheduleLocked(key, state, presenceRetryDelay)
		}
		p.mu.Unlock()
		return
	}

	p.mu.Lock()
	p.forgetLocked(key)
	p.mu.Unlock()

	if p.onLeave != nil {
		p.onLeave(result)
	}
}

func (p *roomPresenceCoordinator) forgetLocked(key roomPresenceKey) {
	state := p.states[key]
	if state != nil && state.timer != nil {
		state.timer.Stop()
	}
	delete(p.states, key)
}
