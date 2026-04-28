import { useContext } from 'react';
import { SocketContextProvider } from '../action-manager/SocketContext';
import type { Item } from '../../../common/interfaces';
import { ActionType } from '../../../common/interfaces';
import './ItemActionCell.css';

interface ItemActionCellProps {
    item: Item;
    billId: number;
}

const ItemActionCell = ({ item, billId }: ItemActionCellProps) => {
    const socketContext = useContext(SocketContextProvider);

    const handleDelete = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.DELETE_ITEM,
                billId: billId,
                itemId: item.id,
            });
        }
    };

    const handleEdit = () => {
        // Empty for now
    };

    return (
        <div className="item-action-cell">
            <button className="action-btn edit-btn" onClick={handleEdit} title="Edit Item">
                Edit
            </button>
            <button className="action-btn delete-btn" onClick={handleDelete} title="Delete Item">
                Delete
            </button>
        </div>
    );
};

export default ItemActionCell;
