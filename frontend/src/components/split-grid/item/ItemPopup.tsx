import { useState, useContext } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SocketContextProvider } from '../action-manager/SocketContext';
import type { Item } from '../../../common/interfaces';
import { ActionType } from '../../../common/interfaces';
import useNotify from '../../../hooks/useNotify';
import './ItemPopup.css';

interface ItemPopupProps {
    onClose: () => void;
    billId: number;
    item?: Item;
}

const ItemPopup = ({ onClose, billId, item }: ItemPopupProps) => {
    const notify = useNotify();
    const [name, setName] = useState(item?.name ?? '');
    const [price, setPrice] = useState(item ? String(item.price) : '');
    const socketContext = useContext(SocketContextProvider);

    const isEdit = item !== undefined;

    const handleConfirm = () => {
        if (!name || !price) {
            notify.error('Please enter both name and price');
            return;
        }

        const numericPrice = parseFloat(price);
        if (isNaN(numericPrice)) {
            notify.error('Please enter a valid price');
            return;
        }

        if (socketContext) {
            socketContext.publishAction(
                isEdit
                    ? {
                          actionType: ActionType.EDIT_ITEM,
                          billId: billId,
                          itemId: item.id,
                          itemName: name.trim(),
                          price: numericPrice,
                      }
                    : {
                          actionType: ActionType.ADD_NEW_ITEM,
                          billId: billId,
                          itemId: -1,
                          itemName: name,
                          price: numericPrice,
                      }
            );
        }

        onClose();
    };

    return (
        <div className="popup-overlay">
            <div className="popup-content">
                <h2>{isEdit ? 'Edit Item' : 'Add New Item'}</h2>
                <div className="input-group">
                    <label>Item Name</label>
                    <Input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="e.g. Pizza"
                        autoFocus
                        onKeyDown={(e) => e.key === 'Enter' && handleConfirm()}
                    />
                </div>
                <div className="input-group">
                    <label>Item Price</label>
                    <Input
                        type="number"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                        placeholder="0.00"
                        step="0.01"
                        onKeyDown={(e) => e.key === 'Enter' && handleConfirm()}
                    />
                </div>
                <div className="popup-actions">
                    <Button variant="outline" onClick={onClose}>Cancel</Button>
                    <Button onClick={handleConfirm}>{isEdit ? 'Save' : 'Confirm'}</Button>
                </div>
            </div>
        </div>
    );
};

export default ItemPopup;
