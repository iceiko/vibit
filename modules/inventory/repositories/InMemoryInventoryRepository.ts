export interface InventoryItemRecord {
  player_id: string;
  item_id: string;
  quantity: number;
}

export interface InventoryRepository {
  getItemQuantity(playerId: string, itemId: string): number;
  setItemQuantity(playerId: string, itemId: string, quantity: number): void;
  listItems(playerId: string): InventoryItemRecord[];
}

export class InMemoryInventoryRepository implements InventoryRepository {
  private readonly quantitiesByPlayer = new Map<string, Map<string, number>>();

  getItemQuantity(playerId: string, itemId: string): number {
    return this.quantitiesByPlayer.get(playerId)?.get(itemId) ?? 0;
  }

  setItemQuantity(playerId: string, itemId: string, quantity: number): void {
    let playerQuantities = this.quantitiesByPlayer.get(playerId);
    if (!playerQuantities) {
      playerQuantities = new Map<string, number>();
      this.quantitiesByPlayer.set(playerId, playerQuantities);
    }
    playerQuantities.set(itemId, quantity);
  }

  listItems(playerId: string): InventoryItemRecord[] {
    const playerQuantities = this.quantitiesByPlayer.get(playerId);
    if (!playerQuantities) {
      return [];
    }

    const items: InventoryItemRecord[] = [];
    for (const [itemId, quantity] of playerQuantities.entries()) {
      items.push({
        player_id: playerId,
        item_id: itemId,
        quantity,
      });
    }
    return items.sort((left, right) => left.item_id.localeCompare(right.item_id));
  }
}
