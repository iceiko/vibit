export interface InventoryPermissionPolicy {
  canGrantItem(permissions: readonly string[]): boolean;
}

export function createInventoryPermissionPolicy(): InventoryPermissionPolicy {
  return {
    canGrantItem(permissions: readonly string[]): boolean {
      return permissions.includes("inventory_grant_item");
    },
  };
}
