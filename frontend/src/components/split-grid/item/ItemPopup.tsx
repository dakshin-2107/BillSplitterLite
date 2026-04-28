import { useState, useContext } from 'react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ActionType } from '../../../common/interfaces';
import './ItemPopup.css';

interface ItemPopupProps {
    onClose: () => void;
    billId: number;
}

const ItemPopup = ({ onClose, billId }: ItemPopupProps) => {
    const [name, setName] = useState('');
    const [price, setPrice] = useState('');
    const socketContext = useContext(SocketContextProvider);

    const handleConfirm = () => {
        if (!name || !price) {
            toast.error('Please enter both name and price');
            return;
        }

        const numericPrice = parseFloat(price);
        if (isNaN(numericPrice)) {
            toast.error('Please enter a valid price');
            return;
        }

        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.ADD_NEW_ITEM,
                billId: billId,
                itemId: -1,
                itemName: name,
                price: numericPrice,
            });
        }

        onClose();
    };

    return (
        <div className="popup-overlay">
            <div className="popup-content">
                <h2>Add New Item</h2>
                <div className="input-group">
                    <label>Item Name</label>
                    <Input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="e.g. Pizza"
                        autoFocus
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
                    />
                </div>
                <div className="popup-actions">
                    <Button variant="outline" onClick={onClose}>Cancel</Button>
                    <Button onClick={handleConfirm}>Confirm</Button>
                </div>
            </div>
        </div>
    );
};

export default ItemPopup;
