import { ActionType } from "../../../common/interfaces";
import type { IAction, BillData } from "../../../common/interfaces";

export function processAction(action: IAction, billData: BillData): BillData {
    // Create a shallow copy of the top-level object
    const newBillData = { ...billData };

    // Ensure nested objects are also copied if they are going to be modified
    newBillData.items = { ...billData.items };
    newBillData.participants = { ...billData.participants };

    switch (action.actionType) {
        case ActionType.ADD_NEW_ITEM:
            newBillData.items[action.itemId] = {
                id: action.itemId,
                name: action.itemName ?? "Unknown Item",
                price: action.price ?? 0,
                takers: {}
            };
            break;
        case ActionType.DELETE_ITEM:
            delete newBillData.items[action.itemId];
            break;
        case ActionType.EDIT_ITEM:
            if (newBillData.items[action.itemId]) {
                newBillData.items[action.itemId] = {
                    ...newBillData.items[action.itemId],
                    name: action.itemName ?? newBillData.items[action.itemId].name,
                    price: action.price ?? newBillData.items[action.itemId].price
                };
            }
            break;
        case ActionType.ADD_TAKER_FOR_ITEM:
            if (newBillData.items[action.itemId] && action.takerId) {
                const item = newBillData.items[action.itemId];
                const newTakers = { ...item.takers };
                const current = newTakers[action.takerId] || 0;
                newTakers[action.takerId] = current + 1;

                newBillData.items[action.itemId] = {
                    ...item,
                    takers: newTakers
                };
            }
            break;
        case ActionType.DELETE_TAKER_FOR_ITEM:
            if (newBillData.items[action.itemId] && action.takerId && newBillData.items[action.itemId].takers[action.takerId]) {
                const item = newBillData.items[action.itemId];
                const newTakers = { ...item.takers };
                const current = newTakers[action.takerId];

                if (current > 1) {
                    newTakers[action.takerId] = current - 1;
                } else {
                    delete newTakers[action.takerId];
                }

                newBillData.items[action.itemId] = {
                    ...item,
                    takers: newTakers
                };
            }
            break;
        case ActionType.ADD_NEW_TAKER:
            if (action.takerId) {
                newBillData.participants[action.takerId] = action.itemName || action.takerId;
            }
            break;
        case ActionType.DELETE_TAKER:
            if (action.takerId) {
                delete newBillData.participants[action.takerId];
            }
            break;
        case ActionType.HELLO_THERE:
            break;
    }

    return newBillData;
}