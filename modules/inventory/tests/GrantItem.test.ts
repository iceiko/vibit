import test from "node:test";
import assert from "node:assert/strict";

import { createGrantItemHandler } from "../commands/GrantItem.ts";
import { createInventoryCapacityPolicy } from "../policies/inventoryCapacityPolicy.ts";
import { createInventoryPermissionPolicy } from "../policies/inventoryPermissionPolicy.ts";
import { InMemoryInventoryRepository } from "../repositories/InMemoryInventoryRepository.ts";

function createHarness(maxQuantityPerItem = 99) {
  const repository = new InMemoryInventoryRepository();
  const handler = createGrantItemHandler({
    repository,
    capacityPolicy: createInventoryCapacityPolicy(maxQuantityPerItem),
    permissionPolicy: createInventoryPermissionPolicy(),
    now: () => new Date("2026-05-12T00:00:00.000Z"),
    createEventId: () => "evt_test_grant_item",
  });

  return { repository, handler };
}

const validInput = {
  player_id: "player_1",
  item_id: "potion",
  quantity: 3,
  reason: "test_grant",
  requested_by: "test_actor",
};

const grantContext = {
  permissions: ["inventory_grant_item"],
};

test("GrantItem records a valid grant and emits ItemGranted once", () => {
  const { repository, handler } = createHarness();

  const result = handler(validInput, grantContext);

  assert.equal(result.ok, true);
  assert.ok(result.ok);
  assert.equal(repository.getItemQuantity("player_1", "potion"), 3);
  assert.deepEqual(result.value, {
    player_id: "player_1",
    item_id: "potion",
    quantity: 3,
    new_quantity: 3,
    event: "ItemGranted",
  });
  assert.equal(result.events.length, 1);
  assert.deepEqual(result.events[0], {
    event_id: "evt_test_grant_item",
    occurred_at: "2026-05-12T00:00:00.000Z",
    player_id: "player_1",
    item_id: "potion",
    quantity: 3,
    new_quantity: 3,
    reason: "test_grant",
  });
});

test("GrantItem rejects invalid quantity", () => {
  const { repository, handler } = createHarness();

  const result = handler({ ...validInput, quantity: 0 }, grantContext);

  assert.equal(result.ok, false);
  assert.equal(result.error, "INVALID_ITEM_QUANTITY");
  assert.equal(result.events.length, 0);
  assert.equal(repository.getItemQuantity("player_1", "potion"), 0);
});

test("GrantItem rejects grants that exceed capacity", () => {
  const { repository, handler } = createHarness(5);
  repository.setItemQuantity("player_1", "potion", 4);

  const result = handler({ ...validInput, quantity: 2 }, grantContext);

  assert.equal(result.ok, false);
  assert.equal(result.error, "INVENTORY_CAPACITY_EXCEEDED");
  assert.equal(result.events.length, 0);
  assert.equal(repository.getItemQuantity("player_1", "potion"), 4);
});

test("GrantItem rejects missing grant permission", () => {
  const { repository, handler } = createHarness();

  const result = handler(validInput, { permissions: [] });

  assert.equal(result.ok, false);
  assert.equal(result.error, "INVENTORY_PERMISSION_DENIED");
  assert.equal(result.events.length, 0);
  assert.equal(repository.getItemQuantity("player_1", "potion"), 0);
});

test("InMemoryInventoryRepository keeps player and item identities distinct", () => {
  const repository = new InMemoryInventoryRepository();

  repository.setItemQuantity("player_1", "potion:rare", 2);
  repository.setItemQuantity("player_1:potion", "rare", 5);

  assert.equal(repository.getItemQuantity("player_1", "potion:rare"), 2);
  assert.equal(repository.getItemQuantity("player_1:potion", "rare"), 5);
  assert.deepEqual(repository.listItems("player_1"), [
    {
      player_id: "player_1",
      item_id: "potion:rare",
      quantity: 2,
    },
  ]);
});
