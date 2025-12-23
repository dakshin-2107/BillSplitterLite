import React, { useContext } from 'react';
import './TakersPanel.css';
import { ActionType } from '../../../common/interfaces';
import type { IAction } from '../../../common/interfaces';
import { SocketContextProvider } from '../action-manager/SocketContext';

interface TakersPanelProps {
    participants: Record<string, string>;
    activeItem: { id: string, name: string } | null;
}

const TakersPanel: React.FC<TakersPanelProps> = ({ participants, activeItem }) => {

    const SocketContext = useContext(SocketContextProvider);
    const OnClickTaker = (id: string) => {
        if (activeItem) {

            const action: IAction = {
                actionId: 0,
                actionType: ActionType.ADD_TAKER_FOR_ITEM,
                splitId: activeItem.id,
                itemId: activeItem.id,
                itemName: activeItem.name,
                takerId: id,
                price: 0,
                total: 0
            }

            if (SocketContext) {
                SocketContext.publishAction(action);
            }
        }
    }

    return (
        <div className={`takers-panel ${!activeItem ? 'disabled' : ''}`}>
            <h3 className="takers-panel-title">{activeItem ? `Selecting taker for : ${activeItem.name}` : "Select an item to add takers"}</h3>
            <div className="takers-grid">
                {Object.entries(participants).map(([id, name]) => (
                    <button
                        key={id}
                        className="taker-chip"
                        onClick={() => { OnClickTaker(id) }}
                    >
                        {name}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default TakersPanel;
