package friends

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iceiko/vibit/runtime/internal/app"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
	"github.com/iceiko/vibit/runtime/internal/platform/tx"
)

type Operation string

const (
	OperationNewService                  Operation = "NewService"
	OperationSendFriendRequest           Operation = "SendFriendRequest"
	OperationAcceptFriendRequest         Operation = "AcceptFriendRequest"
	OperationRejectFriendRequest         Operation = "RejectFriendRequest"
	OperationRemoveFriend                Operation = "RemoveFriend"
	OperationBlockPlayer                 Operation = "BlockPlayer"
	OperationUnblockPlayer               Operation = "UnblockPlayer"
	OperationListFriendRelationships     Operation = "ListFriendRelationships"
	OperationGetFriendRelationshipStatus Operation = "GetFriendRelationshipStatus"
)

type FailureClass string

const (
	FailureClassInvalidRequest        FailureClass = "invalid_request"
	FailureClassUnauthenticated       FailureClass = "unauthenticated"
	FailureClassForbidden             FailureClass = "forbidden"
	FailureClassTargetNotFound        FailureClass = "target_not_found"
	FailureClassRelationshipNotFound  FailureClass = "relationship_not_found"
	FailureClassDuplicateRequest      FailureClass = "duplicate_request"
	FailureClassAlreadyFriends        FailureClass = "already_friends"
	FailureClassBlockedRelationship   FailureClass = "blocked_relationship"
	FailureClassInvalidTransition     FailureClass = "invalid_transition"
	FailureClassVersionMismatch       FailureClass = "version_mismatch"
	FailureClassDependencyUnavailable FailureClass = "dependency_unavailable"
)

type PublicErrorCode string

const (
	PublicErrorFriendshipInvalidRequest       PublicErrorCode = "FRIENDSHIP_INVALID_REQUEST"
	PublicErrorFriendshipUnauthenticated      PublicErrorCode = "FRIENDSHIP_UNAUTHENTICATED"
	PublicErrorFriendshipForbidden            PublicErrorCode = "FRIENDSHIP_FORBIDDEN"
	PublicErrorFriendshipTargetNotFound       PublicErrorCode = "FRIENDSHIP_TARGET_NOT_FOUND"
	PublicErrorFriendshipRelationshipNotFound PublicErrorCode = "FRIENDSHIP_RELATIONSHIP_NOT_FOUND"
	PublicErrorFriendshipDuplicateRequest     PublicErrorCode = "FRIENDSHIP_DUPLICATE_REQUEST"
	PublicErrorFriendshipAlreadyFriends       PublicErrorCode = "FRIENDSHIP_ALREADY_FRIENDS"
	PublicErrorFriendshipBlockedRelationship  PublicErrorCode = "FRIENDSHIP_BLOCKED_RELATIONSHIP"
	PublicErrorFriendshipInvalidTransition    PublicErrorCode = "FRIENDSHIP_INVALID_TRANSITION"
	PublicErrorFriendshipVersionMismatch      PublicErrorCode = "FRIENDSHIP_VERSION_MISMATCH"
	PublicErrorFriendshipUnavailable          PublicErrorCode = "FRIENDSHIP_UNAVAILABLE"
)

type ServiceError struct {
	Operation  Operation
	Class      FailureClass
	PublicCode PublicErrorCode
	Err        error
}

func (e *ServiceError) Error() string {
	if e == nil {
		return ""
	}
	if e.Operation == "" {
		return fmt.Sprintf("friends relationship service: %s", e.PublicCode)
	}
	return fmt.Sprintf("friends relationship service: %s: %s", e.Operation, e.PublicCode)
}

func (e *ServiceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *ServiceError) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

type UnitOfWorkRunner interface {
	WithinUnitOfWork(context.Context, func(context.Context, tx.UnitOfWork) error) error
}

type RelationshipIDGenerator interface {
	GenerateFriendRelationshipID(context.Context) (string, error)
}

type ServiceDependencies struct {
	UnitOfWorkRunner        UnitOfWorkRunner
	RelationshipIDGenerator RelationshipIDGenerator
}

type Service struct {
	unitOfWorkRunner        UnitOfWorkRunner
	relationshipIDGenerator RelationshipIDGenerator
}

func NewService(dependencies ServiceDependencies) (Service, error) {
	if isNilInterface(dependencies.UnitOfWorkRunner) || isNilInterface(dependencies.RelationshipIDGenerator) {
		return Service{}, serviceFailure(OperationNewService, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, errMissingDependency)
	}
	return Service{
		unitOfWorkRunner:        dependencies.UnitOfWorkRunner,
		relationshipIDGenerator: dependencies.RelationshipIDGenerator,
	}, nil
}

type FriendRelationshipOperationStatus string

const (
	FriendRelationshipOperationStatusRejected        FriendRelationshipOperationStatus = "rejected"
	FriendRelationshipOperationStatusSent            FriendRelationshipOperationStatus = "sent"
	FriendRelationshipOperationStatusAccepted        FriendRelationshipOperationStatus = "accepted"
	FriendRelationshipOperationStatusRequestRejected FriendRelationshipOperationStatus = "request_rejected"
	FriendRelationshipOperationStatusRemoved         FriendRelationshipOperationStatus = "removed"
	FriendRelationshipOperationStatusBlocked         FriendRelationshipOperationStatus = "blocked"
	FriendRelationshipOperationStatusUnblocked       FriendRelationshipOperationStatus = "unblocked"
	FriendRelationshipOperationStatusListed          FriendRelationshipOperationStatus = "listed"
	FriendRelationshipOperationStatusFound           FriendRelationshipOperationStatus = "found"
)

type FriendRelationshipPublicStatus string

const (
	FriendRelationshipPublicStatusNone                   FriendRelationshipPublicStatus = "none"
	FriendRelationshipPublicStatusOutgoingRequestPending FriendRelationshipPublicStatus = "outgoing_request_pending"
	FriendRelationshipPublicStatusIncomingRequestPending FriendRelationshipPublicStatus = "incoming_request_pending"
	FriendRelationshipPublicStatusFriends                FriendRelationshipPublicStatus = "friends"
	FriendRelationshipPublicStatusBlockedByActor         FriendRelationshipPublicStatus = "blocked_by_actor"
	FriendRelationshipPublicStatusBlockedActor           FriendRelationshipPublicStatus = "blocked_actor"
	FriendRelationshipPublicStatusMutualBlocked          FriendRelationshipPublicStatus = "mutual_blocked"
	FriendRelationshipPublicStatusRemoved                FriendRelationshipPublicStatus = "removed"
	FriendRelationshipPublicStatusRejected               FriendRelationshipPublicStatus = "rejected"
)

type SendFriendRequestRequest struct {
	Identity       app.RequestIdentity
	TargetPlayerID string
}

type AcceptFriendRequestRequest struct {
	Identity        app.RequestIdentity
	TargetPlayerID  string
	ExpectedVersion *friendsmodule.FriendRelationshipVersion
}

type RejectFriendRequestRequest struct {
	Identity        app.RequestIdentity
	TargetPlayerID  string
	ExpectedVersion *friendsmodule.FriendRelationshipVersion
}

type RemoveFriendRequest struct {
	Identity        app.RequestIdentity
	TargetPlayerID  string
	ExpectedVersion *friendsmodule.FriendRelationshipVersion
}

type BlockPlayerRequest struct {
	Identity        app.RequestIdentity
	TargetPlayerID  string
	ExpectedVersion *friendsmodule.FriendRelationshipVersion
}

type UnblockPlayerRequest struct {
	Identity        app.RequestIdentity
	TargetPlayerID  string
	ExpectedVersion *friendsmodule.FriendRelationshipVersion
}

type ListFriendRelationshipsRequest struct {
	Identity       app.RequestIdentity
	Status         friendsmodule.FriendRelationshipStatus
	Limit          int
	AfterPairToken string
}

type GetFriendRelationshipStatusRequest struct {
	Identity       app.RequestIdentity
	TargetPlayerID string
}

type FriendRelationshipView struct {
	Relationship friendsmodule.FriendRelationship
	PublicStatus FriendRelationshipPublicStatus
	Version      friendsmodule.FriendRelationshipVersion
}

type FriendRelationshipResult struct {
	Status          FriendRelationshipOperationStatus
	PublicErrorCode PublicErrorCode
	FailureClass    FailureClass
	Relationship    friendsmodule.FriendRelationship
	Relationships   []FriendRelationshipView
	PublicStatus    FriendRelationshipPublicStatus
	NextPairToken   string
	Version         friendsmodule.FriendRelationshipVersion
}

func (s Service) SendFriendRequest(ctx context.Context, request SendFriendRequestRequest) (FriendRelationshipResult, error) {
	actor, err := actorFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorFriendshipUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(OperationSendFriendRequest, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated, err)
	}
	targetPlayerID := strings.TrimSpace(request.TargetPlayerID)
	input, err := friendsmodule.NormalizeSendFriendRequestInput(friendsmodule.SendFriendRequestInput{
		RelationshipID: "placeholder",
		Actor:          actor,
		TargetPlayerID: targetPlayerID,
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationSendFriendRequest, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}
	relationshipID, err := s.generatedRelationshipID(ctx)
	if err != nil {
		return rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(OperationSendFriendRequest, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
	}
	input.RelationshipID = relationshipID

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationSendFriendRequest, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		record, err := repository.CreateOrUpdateFriendRequest(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationSendFriendRequest, err)
			return err
		}
		committedResult, err = relationshipResult(OperationSendFriendRequest, FriendRelationshipOperationStatusSent, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationSendFriendRequest, committedResult, failedResult, err)
}

func (s Service) AcceptFriendRequest(ctx context.Context, request AcceptFriendRequestRequest) (FriendRelationshipResult, error) {
	return s.incomingPendingTransition(ctx, OperationAcceptFriendRequest, FriendRelationshipOperationStatusAccepted, request.Identity, request.TargetPlayerID, request.ExpectedVersion, func(repository friendsmodule.Repository, runCtx context.Context, input friendsmodule.AcceptFriendRequestInput) (friendsmodule.FriendRelationship, error) {
		return repository.AcceptFriendRequest(runCtx, input)
	})
}

func (s Service) RejectFriendRequest(ctx context.Context, request RejectFriendRequestRequest) (FriendRelationshipResult, error) {
	return s.incomingPendingTransition(ctx, OperationRejectFriendRequest, FriendRelationshipOperationStatusRequestRejected, request.Identity, request.TargetPlayerID, request.ExpectedVersion, func(repository friendsmodule.Repository, runCtx context.Context, input friendsmodule.AcceptFriendRequestInput) (friendsmodule.FriendRelationship, error) {
		return repository.RejectFriendRequest(runCtx, friendsmodule.RejectFriendRequestInput{
			Actor:           input.Actor,
			Pair:            input.Pair,
			ExpectedVersion: input.ExpectedVersion,
		})
	})
}

func (s Service) RemoveFriend(ctx context.Context, request RemoveFriendRequest) (FriendRelationshipResult, error) {
	actor, pair, err := actorAndPair(request.Identity, request.TargetPlayerID)
	if err != nil {
		return identityOrInvalidRequestFailure(OperationRemoveFriend, err)
	}
	input, err := friendsmodule.NormalizeRemoveFriendInput(friendsmodule.RemoveFriendInput{
		Actor:           actor,
		Pair:            pair,
		ExpectedVersion: request.ExpectedVersion,
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationRemoveFriend, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationRemoveFriend, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		record, err := repository.RemoveFriend(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationRemoveFriend, err)
			return err
		}
		committedResult, err = relationshipResult(OperationRemoveFriend, FriendRelationshipOperationStatusRemoved, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationRemoveFriend, committedResult, failedResult, err)
}

func (s Service) BlockPlayer(ctx context.Context, request BlockPlayerRequest) (FriendRelationshipResult, error) {
	actor, err := actorFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorFriendshipUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(OperationBlockPlayer, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated, err)
	}
	input, err := friendsmodule.NormalizeBlockPlayerInput(friendsmodule.BlockPlayerInput{
		Actor:           actor,
		TargetPlayerID:  strings.TrimSpace(request.TargetPlayerID),
		ExpectedVersion: request.ExpectedVersion,
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationBlockPlayer, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationBlockPlayer, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		record, err := repository.SetPlayerBlock(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationBlockPlayer, err)
			return err
		}
		committedResult, err = relationshipResult(OperationBlockPlayer, FriendRelationshipOperationStatusBlocked, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationBlockPlayer, committedResult, failedResult, err)
}

func (s Service) UnblockPlayer(ctx context.Context, request UnblockPlayerRequest) (FriendRelationshipResult, error) {
	actor, err := actorFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorFriendshipUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(OperationUnblockPlayer, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated, err)
	}
	input, err := friendsmodule.NormalizeUnblockPlayerInput(friendsmodule.UnblockPlayerInput{
		Actor:           actor,
		TargetPlayerID:  strings.TrimSpace(request.TargetPlayerID),
		ExpectedVersion: request.ExpectedVersion,
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationUnblockPlayer, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationUnblockPlayer, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		record, err := repository.ClearPlayerBlock(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationUnblockPlayer, err)
			return err
		}
		committedResult, err = relationshipResult(OperationUnblockPlayer, FriendRelationshipOperationStatusUnblocked, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationUnblockPlayer, committedResult, failedResult, err)
}

func (s Service) ListFriendRelationships(ctx context.Context, request ListFriendRelationshipsRequest) (FriendRelationshipResult, error) {
	actor, err := actorFromValidatedIdentity(request.Identity)
	if err != nil {
		return rejectedResult(PublicErrorFriendshipUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(OperationListFriendRelationships, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated, err)
	}
	input, err := friendsmodule.NormalizeListFriendRelationshipsInput(friendsmodule.ListFriendRelationshipsInput{
		PlayerID:       actor.PlayerID,
		Status:         request.Status,
		Limit:          request.Limit,
		AfterPairToken: strings.TrimSpace(request.AfterPairToken),
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationListFriendRelationships, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationListFriendRelationships, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		result, err := repository.ListRelationshipsForPlayer(runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(OperationListFriendRelationships, err)
			return err
		}
		result, err = friendsmodule.NormalizeListFriendRelationshipsResult(result)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationListFriendRelationships, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, errRelationshipRecordInvalid)
		}
		views := make([]FriendRelationshipView, 0, len(result.Relationships))
		for _, relationship := range result.Relationships {
			view, err := relationshipView(actor.PlayerID, relationship)
			if err != nil {
				failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
				return err
			}
			views = append(views, view)
		}
		committedResult = FriendRelationshipResult{
			Status:        FriendRelationshipOperationStatusListed,
			Relationships: views,
			NextPairToken: strings.TrimSpace(result.NextPairToken),
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationListFriendRelationships, committedResult, failedResult, err)
}

func (s Service) GetFriendRelationshipStatus(ctx context.Context, request GetFriendRelationshipStatusRequest) (FriendRelationshipResult, error) {
	actor, pair, err := actorAndPair(request.Identity, request.TargetPlayerID)
	if err != nil {
		return identityOrInvalidRequestFailure(OperationGetFriendRelationshipStatus, err)
	}
	input, err := friendsmodule.NormalizeGetFriendRelationshipInput(friendsmodule.GetFriendRelationshipInput{Pair: pair})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(OperationGetFriendRelationshipStatus, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(OperationGetFriendRelationshipStatus, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		record, err := repository.GetRelationshipByPair(runCtx, input)
		if err != nil {
			if repositoryConflictClass(err) == friendsmodule.FriendRelationshipConflictNotFound {
				committedResult = FriendRelationshipResult{
					Status:       FriendRelationshipOperationStatusFound,
					PublicStatus: FriendRelationshipPublicStatusNone,
				}
				return nil
			}
			failedResult, err = mapRepositoryFailure(OperationGetFriendRelationshipStatus, err)
			return err
		}
		committedResult, err = relationshipResult(OperationGetFriendRelationshipStatus, FriendRelationshipOperationStatusFound, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(OperationGetFriendRelationshipStatus, committedResult, failedResult, err)
}

type acceptTransition func(friendsmodule.Repository, context.Context, friendsmodule.AcceptFriendRequestInput) (friendsmodule.FriendRelationship, error)

func (s Service) incomingPendingTransition(ctx context.Context, operation Operation, successStatus FriendRelationshipOperationStatus, identity app.RequestIdentity, targetPlayerID string, expectedVersion *friendsmodule.FriendRelationshipVersion, transition acceptTransition) (FriendRelationshipResult, error) {
	actor, pair, err := actorAndPair(identity, targetPlayerID)
	if err != nil {
		return identityOrInvalidRequestFailure(operation, err)
	}
	input, err := friendsmodule.NormalizeAcceptFriendRequestInput(friendsmodule.AcceptFriendRequestInput{
		Actor:           actor,
		Pair:            pair,
		ExpectedVersion: expectedVersion,
	})
	if err != nil {
		return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
			serviceFailure(operation, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
	}

	var committedResult FriendRelationshipResult
	var failedResult FriendRelationshipResult
	err = s.unitOfWorkRunner.WithinUnitOfWork(ctx, func(runCtx context.Context, unit tx.UnitOfWork) error {
		repository, err := friendRelationshipRepositoryFromUnitOfWork(unit)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
		}
		current, err := repository.GetRelationshipByPair(runCtx, friendsmodule.GetFriendRelationshipInput{Pair: input.Pair})
		if err != nil {
			failedResult, err = mapRepositoryFailure(operation, err)
			return err
		}
		current, err = friendsmodule.NormalizeFriendRelationshipRecord(current)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, errRelationshipRecordInvalid)
		}
		if current.LifecycleState != friendsmodule.FriendRelationshipLifecyclePending || current.RequestedByPlayerID == actor.PlayerID {
			failedResult = rejectedResult(PublicErrorFriendshipInvalidTransition, FailureClassInvalidTransition)
			return serviceFailure(operation, FailureClassInvalidTransition, PublicErrorFriendshipInvalidTransition, errInvalidTransition)
		}
		record, err := transition(repository, runCtx, input)
		if err != nil {
			failedResult, err = mapRepositoryFailure(operation, err)
			return err
		}
		committedResult, err = relationshipResult(operation, successStatus, actor.PlayerID, record)
		if err != nil {
			failedResult = rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable)
			return err
		}
		return nil
	})
	return resultAfterUnitOfWork(operation, committedResult, failedResult, err)
}

func actorAndPair(identity app.RequestIdentity, targetPlayerID string) (friendsmodule.FriendRelationshipActor, friendsmodule.FriendRelationshipPair, error) {
	actor, err := actorFromValidatedIdentity(identity)
	if err != nil {
		return friendsmodule.FriendRelationshipActor{}, friendsmodule.FriendRelationshipPair{}, err
	}
	pair, err := friendsmodule.NormalizeFriendRelationshipPair(friendsmodule.FriendRelationshipPair{
		PlayerLowID:  actor.PlayerID,
		PlayerHighID: strings.TrimSpace(targetPlayerID),
	})
	if err != nil {
		return friendsmodule.FriendRelationshipActor{}, friendsmodule.FriendRelationshipPair{}, errInvalidRequest
	}
	return actor, pair, nil
}

func actorFromValidatedIdentity(identity app.RequestIdentity) (friendsmodule.FriendRelationshipActor, error) {
	playerID := strings.TrimSpace(identity.PlayerID)
	actorID := strings.TrimSpace(identity.ActorID)
	if identity.Status != app.IdentityValidationValidated ||
		identity.ActorKind != app.ActorKindPlayer ||
		!identity.PlayerIDValidated ||
		playerID == "" ||
		actorID == "" ||
		playerID != actorID {
		return friendsmodule.FriendRelationshipActor{}, errValidatedPlayerIdentityRequired
	}
	return friendsmodule.NormalizeFriendRelationshipActor(friendsmodule.FriendRelationshipActor{PlayerID: playerID})
}

type friendRelationshipUnitOfWork interface {
	NewFriendRelationshipRepository() (friendsmodule.Repository, error)
}

func friendRelationshipRepositoryFromUnitOfWork(unit tx.UnitOfWork) (friendsmodule.Repository, error) {
	repositories, ok := unit.(friendRelationshipUnitOfWork)
	if !ok {
		return nil, errMissingFriendRelationshipUnitOfWork
	}
	repository, err := repositories.NewFriendRelationshipRepository()
	if err != nil {
		return nil, err
	}
	if isNilInterface(repository) {
		return nil, errMissingRepository
	}
	return repository, nil
}

func relationshipResult(operation Operation, status FriendRelationshipOperationStatus, actorPlayerID string, record friendsmodule.FriendRelationship) (FriendRelationshipResult, error) {
	view, err := relationshipView(actorPlayerID, record)
	if err != nil {
		return FriendRelationshipResult{}, serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
	}
	return FriendRelationshipResult{
		Status:       status,
		Relationship: view.Relationship,
		PublicStatus: view.PublicStatus,
		Version:      view.Version,
	}, nil
}

func relationshipView(actorPlayerID string, record friendsmodule.FriendRelationship) (FriendRelationshipView, error) {
	normalized, err := friendsmodule.NormalizeFriendRelationshipRecord(record)
	if err != nil {
		return FriendRelationshipView{}, errRelationshipRecordInvalid
	}
	publicStatus, err := actorRelativePublicStatus(actorPlayerID, normalized)
	if err != nil {
		return FriendRelationshipView{}, errRelationshipRecordInvalid
	}
	return FriendRelationshipView{
		Relationship: normalized,
		PublicStatus: publicStatus,
		Version:      normalized.Version,
	}, nil
}

func actorRelativePublicStatus(actorPlayerID string, record friendsmodule.FriendRelationship) (FriendRelationshipPublicStatus, error) {
	actorPlayerID = strings.TrimSpace(actorPlayerID)
	if actorPlayerID == "" || (actorPlayerID != record.Pair.PlayerLowID && actorPlayerID != record.Pair.PlayerHighID) {
		return "", errRelationshipRecordInvalid
	}
	actorIsLow := actorPlayerID == record.Pair.PlayerLowID
	actorBlockedTarget := record.BlockState.BlockedByLowAt != nil
	targetBlockedActor := record.BlockState.BlockedByHighAt != nil
	if !actorIsLow {
		actorBlockedTarget = record.BlockState.BlockedByHighAt != nil
		targetBlockedActor = record.BlockState.BlockedByLowAt != nil
	}
	switch {
	case actorBlockedTarget && targetBlockedActor:
		return FriendRelationshipPublicStatusMutualBlocked, nil
	case actorBlockedTarget:
		return FriendRelationshipPublicStatusBlockedByActor, nil
	case targetBlockedActor:
		return FriendRelationshipPublicStatusBlockedActor, nil
	}

	switch record.LifecycleState {
	case friendsmodule.FriendRelationshipLifecyclePending:
		if record.RequestedByPlayerID == actorPlayerID {
			return FriendRelationshipPublicStatusOutgoingRequestPending, nil
		}
		return FriendRelationshipPublicStatusIncomingRequestPending, nil
	case friendsmodule.FriendRelationshipLifecycleFriends:
		return FriendRelationshipPublicStatusFriends, nil
	case friendsmodule.FriendRelationshipLifecycleRemoved:
		return FriendRelationshipPublicStatusRemoved, nil
	case friendsmodule.FriendRelationshipLifecycleRejected:
		return FriendRelationshipPublicStatusRejected, nil
	default:
		return "", errRelationshipRecordInvalid
	}
}

func mapRepositoryFailure(operation Operation, err error) (FriendRelationshipResult, error) {
	publicCode, class := publicFailureForRepositoryError(err)
	return rejectedResult(publicCode, class), serviceFailure(operation, class, publicCode, nil)
}

func publicFailureForRepositoryError(err error) (PublicErrorCode, FailureClass) {
	switch repositoryConflictClass(err) {
	case friendsmodule.FriendRelationshipConflictSelfRelationshipForbidden:
		return PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest
	case friendsmodule.FriendRelationshipConflictTargetPlayerNotFound:
		return PublicErrorFriendshipTargetNotFound, FailureClassTargetNotFound
	case friendsmodule.FriendRelationshipConflictNotFound:
		return PublicErrorFriendshipRelationshipNotFound, FailureClassRelationshipNotFound
	case friendsmodule.FriendRelationshipConflictDuplicatePendingRequest:
		return PublicErrorFriendshipDuplicateRequest, FailureClassDuplicateRequest
	case friendsmodule.FriendRelationshipConflictAlreadyFriends:
		return PublicErrorFriendshipAlreadyFriends, FailureClassAlreadyFriends
	case friendsmodule.FriendRelationshipConflictBlockedRelationship:
		return PublicErrorFriendshipBlockedRelationship, FailureClassBlockedRelationship
	case friendsmodule.FriendRelationshipConflictInvalidTransition,
		friendsmodule.FriendRelationshipConflictPairIdentity:
		return PublicErrorFriendshipInvalidTransition, FailureClassInvalidTransition
	case friendsmodule.FriendRelationshipConflictVersionMismatch,
		friendsmodule.FriendRelationshipConflictStaleVersion:
		return PublicErrorFriendshipVersionMismatch, FailureClassVersionMismatch
	case friendsmodule.FriendRelationshipConflictStorageUnavailable:
		return PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable
	default:
		if errors.Is(err, friendsmodule.ErrFriendRelationshipInvalidInput) {
			return PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest
		}
		if errors.Is(err, friendsmodule.ErrFriendRelationshipUnavailable) {
			return PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable
		}
		return PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable
	}
}

func repositoryConflictClass(err error) friendsmodule.FriendRelationshipConflictClass {
	var repositoryErr *friendsmodule.FriendRelationshipRepositoryError
	if errors.As(err, &repositoryErr) {
		return repositoryErr.Conflict.Class
	}
	var conflict friendsmodule.FriendRelationshipConflict
	if errors.As(err, &conflict) {
		return conflict.Class
	}
	return ""
}

func (s Service) generatedRelationshipID(ctx context.Context) (friendsmodule.FriendRelationshipID, error) {
	relationshipID, err := s.relationshipIDGenerator.GenerateFriendRelationshipID(ctx)
	if err != nil {
		return "", errMissingDependency
	}
	normalized, err := friendsmodule.NormalizeFriendRelationshipID(friendsmodule.FriendRelationshipID(relationshipID))
	if err != nil {
		return "", errMalformedGeneratedID
	}
	return normalized, nil
}

func resultAfterUnitOfWork(operation Operation, committedResult FriendRelationshipResult, failedResult FriendRelationshipResult, err error) (FriendRelationshipResult, error) {
	if err != nil {
		if failedResult.Status != "" {
			return failedResult, err
		}
		return rejectedResult(PublicErrorFriendshipUnavailable, FailureClassDependencyUnavailable),
			serviceFailure(operation, FailureClassDependencyUnavailable, PublicErrorFriendshipUnavailable, err)
	}
	return committedResult, nil
}

func identityOrInvalidRequestFailure(operation Operation, err error) (FriendRelationshipResult, error) {
	if errors.Is(err, errValidatedPlayerIdentityRequired) {
		return rejectedResult(PublicErrorFriendshipUnauthenticated, FailureClassUnauthenticated),
			serviceFailure(operation, FailureClassUnauthenticated, PublicErrorFriendshipUnauthenticated, err)
	}
	return rejectedResult(PublicErrorFriendshipInvalidRequest, FailureClassInvalidRequest),
		serviceFailure(operation, FailureClassInvalidRequest, PublicErrorFriendshipInvalidRequest, errInvalidRequest)
}

func rejectedResult(publicCode PublicErrorCode, class FailureClass) FriendRelationshipResult {
	return FriendRelationshipResult{
		Status:          FriendRelationshipOperationStatusRejected,
		PublicErrorCode: publicCode,
		FailureClass:    class,
	}
}

func serviceFailure(operation Operation, class FailureClass, publicCode PublicErrorCode, err error) error {
	return &ServiceError{
		Operation:  operation,
		Class:      class,
		PublicCode: publicCode,
		Err:        err,
	}
}

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

var (
	errMissingDependency                   = errors.New("friends relationship service: dependency is required")
	errValidatedPlayerIdentityRequired     = errors.New("friends relationship service: validated player identity is required")
	errInvalidRequest                      = errors.New("friends relationship service: invalid request")
	errInvalidTransition                   = errors.New("friends relationship service: invalid relationship transition")
	errMissingFriendRelationshipUnitOfWork = errors.New("friends relationship service: friend relationship unit-of-work capability is required")
	errMissingRepository                   = errors.New("friends relationship service: friend relationship repository is required")
	errRelationshipRecordInvalid           = errors.New("friends relationship service: friend relationship record is invalid")
	errMalformedGeneratedID                = errors.New("friends relationship service: generated relationship id is malformed")
)
