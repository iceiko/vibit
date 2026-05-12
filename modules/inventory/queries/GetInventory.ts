import {
  GetInventoryContract,
  type GetInventoryInput,
  type GetInventoryOutput,
  type GetInventoryErrorCode,
} from "../generated/contracts/GetInventory.generated.ts";
import type { InventoryPermissionPolicy } from "../policies/inventoryPermissionPolicy.ts";
import type { InventoryRepository } from "../repositories/InMemoryInventoryRepository.ts";

export interface GetInventoryContext {
  permissions: readonly string[];
}

export type GetInventoryResult =
  | {
      ok: true;
      value: GetInventoryOutput;
      events: [];
    }
  | {
      ok: false;
      error: GetInventoryErrorCode;
      events: [];
    };

export interface GetInventoryDependencies {
  repository: InventoryRepository;
  permissionPolicy: InventoryPermissionPolicy;
}

export function createGetInventoryHandler(dependencies: GetInventoryDependencies) {
  return function getInventory(input: GetInventoryInput, context: GetInventoryContext): GetInventoryResult {
    if (!dependencies.permissionPolicy.canReadInventory(context.permissions)) {
      return failure("INVENTORY_PERMISSION_DENIED");
    }

    return {
      ok: true,
      value: {
        player_id: input.player_id,
        items: dependencies.repository
          .listItems(input.player_id)
          .map((item) => ({
            item_id: item.item_id,
            quantity: item.quantity,
          })),
      },
      events: noEvents(),
    };
  };
}

function noEvents(): [] {
  if (GetInventoryContract.emits.length !== 0) {
    throw new Error("GetInventory contract must not emit events.");
  }
  return [];
}

function failure(error: GetInventoryErrorCode): GetInventoryResult {
  return {
    ok: false,
    error,
    events: [],
  };
}
