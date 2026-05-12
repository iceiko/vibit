export interface InventoryCapacityPolicy {
  canGrant(currentQuantity: number, grantQuantity: number): boolean;
}

export function createInventoryCapacityPolicy(maxQuantityPerItem: number): InventoryCapacityPolicy {
  return {
    canGrant(currentQuantity: number, grantQuantity: number): boolean {
      return currentQuantity + grantQuantity <= maxQuantityPerItem;
    },
  };
}
