import React, { useContext } from 'react';
import './TakersPanel.css';
import { ActionType } from '../../../common/interfaces';
import type { IAction } from '../../../common/interfaces';
import { SocketContextProvider } from '../action-manager/SocketContext';

interface TakersPanelProps {
    participants: Record<string, string>;
    activeItem: { id: number, name: string, billId: number } | null;
}

const TakersPanel: React.FC<TakersPanelProps> = ({ participants, activeItem }) => {

    const SocketContext = useContext(SocketContextProvider);
    const OnClickTaker = (id: string) => {
        if (activeItem) {

            const action: IAction = {
                actionId: 0,
                actionType: ActionType.ADD_TAKER_FOR_ITEM,
                billId: activeItem.billId,
                itemId: activeItem.id,
                takerId: id
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
                {/*Add an all button*/}
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
