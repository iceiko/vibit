import test from "node:test";
import assert from "node:assert/strict";

import { createGetInventoryHandler } from "../queries/GetInventory.ts";
import { createInventoryPermissionPolicy } from "../policies/inventoryPermissionPolicy.ts";
import { InMemoryInventoryRepository } from "../repositories/InMemoryInventoryRepository.ts";

function createHarness() {
  const repository = new InMemoryInventoryRepository();
  const handler = createGetInventoryHandler({
    repository,
    permissionPolicy: createInventoryPermissionPolicy(),
  });

  return { repository, handler };
}

const readContext = {
  permissions: ["inventory_read"],
};

const validInput = {
  player_id: "player_1",
  requested_by: "test_actor",
};

test("GetInventory returns sorted inventory state without events", () => {
  const { repository, handler } = createHarness();
  repository.setItemQuantity("player_1", "z_potion", 1);
  repository.setItemQuantity("player_1", "a_elixir", 2);

  const result = handler(validInput, readContext);

  assert.equal(result.ok, true);
  assert.ok(result.ok);
  assert.deepEqual(result.value, {
    player_id: "player_1",
    items: [
      {
        item_id: "a_elixir",
        quantity: 2,
      },
      {
        item_id: "z_potion",
        quantity: 1,
      },
    ],
  });
  assert.deepEqual(result.events, []);
});

test("GetInventory returns an empty item list for an empty inventory", () => {
  const { handler } = createHarness();

  const result = handler(validInput, readContext);

  assert.equal(result.ok, true);
  assert.ok(result.ok);
  assert.deepEqual(result.value, {
    player_id: "player_1",
    items: [],
  });
});

test("GetInventory does not mutate inventory state", () => {
  const { repository, handler } = createHarness();
  repository.setItemQuantity("player_1", "potion", 3);

  const before = repository.listItems("player_1");
  const result = handler(validInput, readContext);
  const after = repository.listItems("player_1");

  assert.equal(result.ok, true);
  assert.deepEqual(after, before);
});

test("GetInventory keeps player and item identities distinct", () => {
  const { repository, handler } = createHarness();
  repository.setItemQuantity("player_1", "potion:rare", 2);
  repository.setItemQuantity("player_1:potion", "rare", 5);

  const result = handler(validInput, readContext);

  assert.equal(result.ok, true);
  assert.ok(result.ok);
  assert.deepEqual(result.value.items, [
    {
      item_id: "potion:rare",
      quantity: 2,
    },
  ]);
});

test("GetInventory rejects missing read permission", () => {
  const { repository, handler } = createHarness();
  repository.setItemQuantity("player_1", "potion", 3);

  const result = handler(validInput, { permissions: [] });

  assert.equal(result.ok, false);
  assert.equal(result.error, "INVENTORY_PERMISSION_DENIED");
  assert.deepEqual(result.events, []);
  assert.equal(repository.getItemQuantity("player_1", "potion"), 3);
});
