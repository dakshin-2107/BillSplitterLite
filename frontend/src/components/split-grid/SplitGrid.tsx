import './SplitGrid.css';
import React from 'react';
import type { BillData } from '../../common/interfaces';
import TakerCell from './item/TakerCell';
import TakersPanel from './people-panel/TakersPanel';
import ItemPopup from './item/ItemPopup';
import { SocketContextComponent } from './action-manager/SocketContext';
import { StatusIndicator } from './connection-status/StatusIndicator';
import ItemActionCell from './item/ItemActionCell';
import type { Tally } from '../../common/interfaces';
import TallyPanel from './tally-panel/TallyPanel';
import BillDebugger from '../debug/BillDebugger';
import BillActions from './bill-actions/BillActions';

interface SplitGridProps {
    billData: BillData;
    setBillData: React.Dispatch<React.SetStateAction<BillData | null>>;
}

const SplitGrid: React.FC<SplitGridProps> = ({ billData, setBillData }) => {

    const [tally, setTally] = React.useState<Tally | null>(null);
    const [activeItem, setActiveItem] = React.useState<{ id: string, name: string } | null>(null);
    const [isAddItemPopupOpen, setIsAddItemPopupOpen] = React.useState(false);

    React.useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            const target = event.target as HTMLElement;

            // Check if clicking outside takers panel AND outside any "Add" buttons
            const clickedTakersPanel = target.closest('.takers-panel');
            const clickedAddButton = target.closest('.btn-add-taker-inline');

            if (!clickedTakersPanel && !clickedAddButton && activeItem) {
                setActiveItem(null);
            }
        };

        if (activeItem) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [activeItem]);

    const handleAddTakerClick = (itemId: string, itemName: string) => {
        setActiveItem({ id: itemId, name: itemName });
    };

    return (
        <SocketContextComponent setBillData={setBillData} setTallyData={setTally}>
            {billData && (
                <div className="split-grid-container">
                    <div className="split-grid-header">
                        <h1 className="split-grid-title">Split the bill</h1>
                        <StatusIndicator />
                    </div>

                    <TakersPanel
                        participants={billData.participants}
                        activeItem={activeItem}
                    />

                    <div className="main-layout">
                        <div className="table-section">
                            <div className="table-wrapper">
                                <table className="split-table">
                                    <thead>
                                        <tr>
                                            <th>Actions</th>
                                            <th>Number</th>
                                            <th>Item Name</th>
                                            <th>Price</th>
                                            <th>Takers</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {Object.values(billData.items).map((item, index) => (
                                            <tr key={item.id || index}>
                                                <td>
                                                    <ItemActionCell
                                                        item={item}
                                                    />
                                                </td>
                                                <td>{index + 1}</td>
                                                <td>{item.name}</td>
                                                <td>{item.price.toFixed(2)}</td>
                                                <td className="takers-cell">
                                                    <TakerCell
                                                        takers={item.takers}
                                                        itemId={item.id}
                                                        itemName={item.name}
                                                        onAddTakerClick={handleAddTakerClick}
                                                    />
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                                <div className="grid-controls">
                                    <div className="left-controls">
                                        <button
                                            className="control-btn add-btn"
                                            onClick={() => setIsAddItemPopupOpen(true)}
                                        >
                                            + Add New Item
                                        </button>
                                    </div>
                                    <BillActions />
                                </div>

                                {isAddItemPopupOpen && (
                                    <ItemPopup
                                        onClose={() => setIsAddItemPopupOpen(false)}
                                        newItemId={Object.values(billData.items).length + 1}
                                    />
                                )}

                                {/* Debugger to visualize state updates */}
                                {/*<BillDebugger billData={billData} />*/}
                                {/*<LastMessage />*/}
                            </div>
                        </div>

                        <div className="tally-section">
                            <TallyPanel
                                tally={tally}
                                participants={billData.participants}
                                billData={billData}
                            />
                        </div>
                    </div>
                </div>
            )}
        </SocketContextComponent>
    )
}

export default SplitGrid;
