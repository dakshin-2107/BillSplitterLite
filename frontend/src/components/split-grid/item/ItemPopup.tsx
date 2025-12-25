import React, { useState, useContext } from 'react';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ActionType } from '../../../common/interfaces';
import './ItemPopup.css';

interface ItemPopupProps {
    onClose: () => void;
    newItemId: number;
}

const ItemPopup: React.FC<ItemPopupProps> = ({ onClose, newItemId }) => {
    const [name, setName] = useState('');
    const [price, setPrice] = useState('');
    const socketContext = useContext(SocketContextProvider);

    const handleConfirm = () => {
        if (!name || !price) {
            alert('Please enter both name and price');
            return;
        }

        const numericPrice = parseFloat(price);
        if (isNaN(numericPrice)) {
            alert('Please enter a valid price');
            return;
        }

        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.ADD_NEW_ITEM,
                itemId: `${newItemId}`, // backend will generate it 
                itemName: name,
                price: numericPrice,
                takerId: '' // Optional/not needed for new item
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
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="e.g. Pizza"
                        autoFocus
                    />
                </div>
                <div className="input-group">
                    <label>Item Price</label>
                    <input
                        type="number"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                        placeholder="0.00"
                        step="0.01"
                    />
                </div>
                <div className="popup-actions">
                    <button className="popup-btn cancel" onClick={onClose}>Cancel</button>
                    <button className="popup-btn confirm" onClick={handleConfirm}>Confirm</button>
                </div>
            </div>
        </div>
    );
};

export default ItemPopup;
