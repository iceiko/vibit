package storage

import (
	"github.com/iceiko/vibit/runtime/internal/app"
	storagemodule "github.com/iceiko/vibit/runtime/internal/modules/storage"
)

const (
	QueryGetOwnStorageObject      = "GetOwnStorageObject"
	QueryListOwnStorageObjects    = "ListOwnStorageObjects"
	CommandPutOwnStorageObject    = "PutOwnStorageObject"
	CommandDeleteOwnStorageObject = "DeleteOwnStorageObject"
)

// Full route keys:
// - storage.GetOwnStorageObject
// - storage.ListOwnStorageObjects
// - storage.PutOwnStorageObject
// - storage.DeleteOwnStorageObject

func GetOwnStorageObjectRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: storagemodule.ModuleName, Name: QueryGetOwnStorageObject}
}

func ListOwnStorageObjectsRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindQuery, Module: storagemodule.ModuleName, Name: QueryListOwnStorageObjects}
}

func PutOwnStorageObjectRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: storagemodule.ModuleName, Name: CommandPutOwnStorageObject}
}

func DeleteOwnStorageObjectRoute() app.RouteKey {
	return app.RouteKey{Kind: app.MessageKindCommand, Module: storagemodule.ModuleName, Name: CommandDeleteOwnStorageObject}
}
