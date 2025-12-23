
import React, { useContext } from 'react';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ActionType } from '../../../common/interfaces';
import './TakerCell.css';

interface TakerCellProps {
    takers: { [key: string]: number };
    itemId: string;
    itemName: string;
    onAddTakerClick: (itemId: string, itemName: string) => void;
}

const TakerCell: React.FC<TakerCellProps> = ({ takers, itemId, itemName, onAddTakerClick }) => {
    const socketContext = useContext(SocketContextProvider);

    // to add a taker's share to the item
    const handleAddTakerShare = (takerId: string) => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.ADD_TAKER_FOR_ITEM,
                itemId: itemId,
                takerId: takerId
            });
        }
    }

    // to remove a taker's share from the item
    const handleRemoveTakerShare = (takerId: string) => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.DELETE_TAKER_FOR_ITEM,
                itemId: itemId,
                takerId: takerId
            });
        }
    }

    // to add a new taker to the item
    const handleAddTaker = () => {
        onAddTakerClick(itemId, itemName);
    }

    return (
        <div className="taker-cell">
            <div className="takers-list">
                {Object.entries(takers).filter(([takerId, count]) => count >= 1).map(([takerId, count]) => (
                    <div key={takerId} className="taker-pill">
                        <button className="pill-btn minus" onClick={() => handleRemoveTakerShare(takerId)}>−</button>
                        <div className="pill-content">
                            <span className="taker-id">{takerId}</span>
                            {count > 1 && <span className="taker-count">x{count}</span>}
                        </div>
                        <button className="pill-btn plus" onClick={() => handleAddTakerShare(takerId)}>+</button>
                    </div>
                ))}
                <button className="btn-add-taker-inline" onClick={() => handleAddTaker()}>
                    + Add
                </button>
            </div>
        </div>
    );
};

export default TakerCell;
