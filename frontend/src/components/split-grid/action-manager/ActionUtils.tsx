import { ActionType } from "../../../common/interfaces";
import type { IAction, SplitData, BillData } from "../../../common/interfaces";
import { URLProvider } from "../../../common/urlProvider";

export function processAction(action: IAction, splitData: SplitData | null): SplitData | null {

    if (!splitData) {
        return null;
    }

    const newSplitData: SplitData = {
        ...splitData,
        bills: { ...splitData.bills },
        participants: { ...splitData.participants }
    };

    const getBill = (billId?: number) => {
        if (billId === undefined) return null;
        return newSplitData.bills[billId];
    };

    const updateBill = (billId: number, updatedBill: BillData) => {
        newSplitData.bills[billId] = updatedBill;
    };

    switch (action.actionType) {
        case ActionType.ADD_NEW_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill) {
                const newBillItems = {
                    ...bill.items,
                    [action.itemId]: {
                        id: action.itemId,
                        name: action.itemName ?? "Unknown Item",
                        price: action.price ?? 0,
                        takers: {}
                    }
                };
                updateBill(action.billId, { ...bill, items: newBillItems });
            }
            break;
        }

        case ActionType.DELETE_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill) {
                const newItems = { ...bill.items };
                delete newItems[action.itemId];
                updateBill(action.billId, { ...bill, items: newItems });
            }
            break;
        }

        case ActionType.EDIT_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill && bill.items[action.itemId]) {
                const newItems = {
                    ...bill.items,
                    [action.itemId]: {
                        ...bill.items[action.itemId],
                        name: action.itemName ?? bill.items[action.itemId].name,
                        price: action.price ?? bill.items[action.itemId].price
                    }
                };
                updateBill(action.billId, { ...bill, items: newItems });
            }
            break;
        }

        case ActionType.ADD_TAKER_FOR_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill && bill.items[action.itemId] && action.takerId) {
                const item = bill.items[action.itemId];
                const newTakers = { ...item.takers };
                const current = newTakers[action.takerId] || 0;
                newTakers[action.takerId] = current + 1;

                const newItems = {
                    ...bill.items,
                    [action.itemId]: { ...item, takers: newTakers }
                };
                updateBill(action.billId, { ...bill, items: newItems });
            }
            break;
        }

        case ActionType.DELETE_TAKER_FOR_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill && bill.items[action.itemId] && action.takerId && bill.items[action.itemId].takers[action.takerId]) {
                const item = bill.items[action.itemId];
                const newTakers = { ...item.takers };
                const current = newTakers[action.takerId];

                if (current > 1) {
                    newTakers[action.takerId] = current - 1;
                } else {
                    delete newTakers[action.takerId];
                }

                const newItems = {
                    ...bill.items,
                    [action.itemId]: { ...item, takers: newTakers }
                };
                updateBill(action.billId, { ...bill, items: newItems });
            }
            break;
        }

        case ActionType.ADD_ALL_TAKERS_FOR_ITEM: {
            const bill = getBill(action.billId);
            if (action.billId !== undefined && bill && bill.items[action.itemId]) {
                const item = bill.items[action.itemId];
                const newTakers = { ...item.takers };

                // Add all participants with at least 1 share
                Object.keys(newSplitData.participants).forEach(takerId => {
                    if (!newTakers[takerId]) {
                        newTakers[takerId] = 1;
                    }
                });

                const newItems = {
                    ...bill.items,
                    [action.itemId]: { ...item, takers: newTakers }
                };
                updateBill(action.billId, { ...bill, items: newItems });
            }
            break;
        }

        case ActionType.ADD_NEW_TAKER:
            if (action.takerId) {
                newSplitData.participants[action.takerId] = action.itemName || action.takerId;
            }
            break;

        case ActionType.DELETE_TAKER:
            if (action.takerId) {
                delete newSplitData.participants[action.takerId];
            }
            break;

        case ActionType.HELLO_THERE:
            break;

        case ActionType.BYE_BYE:
            fetch(URLProvider.getCloseUrl(), {
                method: 'GET',
                credentials: 'include'
            }).finally(() => {
                window.location.href = '/';
            });
            break;
    }

    return newSplitData;
}
