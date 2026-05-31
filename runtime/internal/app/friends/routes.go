package friends

import (
	"github.com/iceiko/vibit/runtime/internal/app"
	friendsmodule "github.com/iceiko/vibit/runtime/internal/modules/friends"
)

const (
	CommandSendFriendRequest         = "SendFriendRequest"
	CommandAcceptFriendRequest       = "AcceptFriendRequest"
	CommandRejectFriendRequest       = "RejectFriendRequest"
	CommandRemoveFriend              = "RemoveFriend"
	CommandBlockPlayer               = "BlockPlayer"
	CommandUnblockPlayer             = "UnblockPlayer"
	QueryListFriendRelationships     = "ListFriendRelationships"
	QueryGetFriendRelationshipStatus = "GetFriendRelationshipStatus"
)

// Full route keys:
// - friends.SendFriendRequest
// - friends.AcceptFriendRequest
// - friends.RejectFriendRequest
// - friends.RemoveFriend
// - friends.BlockPlayer
// - friends.UnblockPlayer
// - friends.ListFriendRelationships
// - friends.GetFriendRelationshipStatus

func SendFriendRequestRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandSendFriendRequest}
}

func AcceptFriendRequestRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandAcceptFriendRequest}
}

func RejectFriendRequestRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandRejectFriendRequest}
}

func RemoveFriendRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandRemoveFriend}
}

func BlockPlayerRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandBlockPlayer}
}

func UnblockPlayerRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: friendsmodule.ModuleName, Name: CommandUnblockPlayer}
}

func ListFriendRelationshipsRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: friendsmodule.ModuleName, Name: QueryListFriendRelationships}
}

func GetFriendRelationshipStatusRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: friendsmodule.ModuleName, Name: QueryGetFriendRelationshipStatus}
}
