import './SplitGrid.css';
import React from 'react';
import type { SplitData } from '../../common/interfaces';
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
    splitData: SplitData;
    setSplitData: React.Dispatch<React.SetStateAction<SplitData | null>>;
}

const SplitGrid: React.FC<SplitGridProps> = ({ splitData, setSplitData }) => {

    const [tally, setTally] = React.useState<Tally | null>(null);
    const [activeItem, setActiveItem] = React.useState<{ id: number, name: string, billId: number } | null>(null);
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

    const handleAddTakerClick = (itemId: number, itemName: string, billId: number) => {
        setActiveItem({ id: itemId, name: itemName, billId });
    };

    // Aggregate all items from all bills
    const allItems = Object.entries(splitData.bills).flatMap(([billId, bill]) =>
        Object.values(bill.items).map(item => ({ ...item, billId: Number(billId) }))
    );

    // Get the first bill ID for adding new items (default)
    const billsKeys = Object.keys(splitData.bills);
    const firstBillId = billsKeys.length > 0 ? Number(billsKeys[0]) : 0;

    return (
        <SocketContextComponent setSplitData={setSplitData} setTallyData={setTally}>
            {splitData && (
                <div className="split-grid-container">
                    <div className="split-grid-header">
                        <h1 className="split-grid-title">Split the bill</h1>
                        <StatusIndicator />
                    </div>

                    <TakersPanel
                        participants={splitData.participants}
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
                                        {Object.entries(splitData.bills).map(([billIdStr, bill]) => {
                                            const billId = Number(billIdStr);
                                            return (
                                                <React.Fragment key={billId}>
                                                    <tr className="bill-header-row">
                                                        <td colSpan={5}>
                                                            <div className="bill-header-content">
                                                                <span className="bill-location-tag">{bill.location}</span>
                                                                <span className="bill-date-tag">{bill.date}</span>
                                                            </div>
                                                        </td>
                                                    </tr>
                                                    {Object.values(bill.items).map((item, itemIndex) => (
                                                        <tr key={`${billId}-${item.id || itemIndex}`}>
                                                            <td>
                                                                <ItemActionCell
                                                                    item={item}
                                                                    billId={billId}
                                                                />
                                                            </td>
                                                            <td>{itemIndex + 1}</td>
                                                            <td>{item.name}</td>
                                                            <td>{item.price.toFixed(2)}</td>
                                                            <td className="takers-cell">
                                                                <TakerCell
                                                                    takers={item.takers}
                                                                    itemId={item.id}
                                                                    itemName={item.name}
                                                                    billId={billId}
                                                                    onAddTakerClick={handleAddTakerClick}
                                                                />
                                                            </td>
                                                        </tr>
                                                    ))}
                                                </React.Fragment>
                                            )
                                        })}
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
                                        newItemId={allItems.length + 1}
                                        billId={firstBillId}
                                    />
                                )}
                            </div>
                        </div>

                        <div className="tally-section">
                            <TallyPanel
                                tally={tally}
                                participants={splitData.participants}
                                splitData={splitData}
                            />
                        </div>
                    </div>
                </div>
            )}
        </SocketContextComponent>
    )
}


export default SplitGrid;
