import {
  GrantItemContract,
  type GrantItemInput,
  type GrantItemOutput,
  type GrantItemErrorCode,
} from "../generated/contracts/GrantItem.generated.ts";
import type { InventoryCapacityPolicy } from "../policies/inventoryCapacityPolicy.ts";
import type { InventoryPermissionPolicy } from "../policies/inventoryPermissionPolicy.ts";
import type { InventoryRepository } from "../repositories/InMemoryInventoryRepository.ts";

export interface GrantItemContext {
  permissions: readonly string[];
}

export interface ItemGrantedEvent {
  event_id: string;
  occurred_at: string;
  player_id: string;
  item_id: string;
  quantity: number;
  new_quantity: number;
  reason: string;
}

export type GrantItemResult =
  | {
      ok: true;
      value: GrantItemOutput;
      events: [ItemGrantedEvent];
    }
  | {
      ok: false;
      error: GrantItemErrorCode;
      events: [];
    };

export interface GrantItemDependencies {
  repository: InventoryRepository;
  capacityPolicy: InventoryCapacityPolicy;
  permissionPolicy: InventoryPermissionPolicy;
  now?: () => Date;
  createEventId?: () => string;
}

export function createGrantItemHandler(dependencies: GrantItemDependencies) {
  return function grantItem(input: GrantItemInput, context: GrantItemContext): GrantItemResult {
    if (!dependencies.permissionPolicy.canGrantItem(context.permissions)) {
      return failure("INVENTORY_PERMISSION_DENIED");
    }

    if (!Number.isInteger(input.quantity) || input.quantity <= 0) {
      return failure("INVALID_ITEM_QUANTITY");
    }

    const currentQuantity = dependencies.repository.getItemQuantity(input.player_id, input.item_id);
    if (!dependencies.capacityPolicy.canGrant(currentQuantity, input.quantity)) {
      return failure("INVENTORY_CAPACITY_EXCEEDED");
    }

    const newQuantity = currentQuantity + input.quantity;
    dependencies.repository.setItemQuantity(input.player_id, input.item_id, newQuantity);

    const event: ItemGrantedEvent = {
      event_id: dependencies.createEventId ? dependencies.createEventId() : "evt_inventory_item_granted",
      occurred_at: (dependencies.now ? dependencies.now() : new Date()).toISOString(),
      player_id: input.player_id,
      item_id: input.item_id,
      quantity: input.quantity,
      new_quantity: newQuantity,
      reason: input.reason,
    };

    return {
      ok: true,
      value: {
        player_id: input.player_id,
        item_id: input.item_id,
        quantity: input.quantity,
        new_quantity: newQuantity,
        event: GrantItemContract.emits[0],
      },
      events: [event],
    };
  };
}

function failure(error: GrantItemErrorCode): GrantItemResult {
  return {
    ok: false,
    error,
    events: [],
  };
}
