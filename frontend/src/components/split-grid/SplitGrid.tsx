import './SplitGrid.css';
import React from 'react';
import type { SplitData, Tally } from '../../common/interfaces';
import TakersPanel from './people-panel/TakersPanel';
import { SocketContextComponent } from './action-manager/SocketContext';
import { StatusIndicator } from './connection-status/StatusIndicator';
import { BillGrid } from './bill/BillGrid';
import { BillPillCarousel } from './bill/BillPillCarousel';
import TallyPanel from './tally-panel/TallyPanel';
import ItemPopup from './item/ItemPopup';
import BillActions from './bill-actions/BillActions';

interface SplitGridProps {
    splitData: SplitData;
    setSplitData: React.Dispatch<React.SetStateAction<SplitData | null>>;
}

const SplitGrid: React.FC<SplitGridProps> = ({ splitData, setSplitData }) => {

    const [tally, setTally] = React.useState<Tally | null>(null);
    const [activeItem, setActiveItem] = React.useState<{ id: number, name: string, billId: number } | null>(null);
    const [isAddItemPopupOpen, setIsAddItemPopupOpen] = React.useState(false);

    // Get the first bill ID for initial state
    const billsKeys = Object.keys(splitData.bills);
    const initialBillId = billsKeys.length > 0 ? Number(billsKeys[0]) : 0;
    const [activeBillId, setActiveBillId] = React.useState<number>(initialBillId);

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

    return (
        <SocketContextComponent setSplitData={setSplitData} setTallyData={setTally}>
            {splitData && (
                <div className="split-grid-container">
                    <div className="split-grid-header">
                        <h1 className="split-grid-title">Split the bill</h1>
                        <StatusIndicator />
                    </div>

                    <div className="main-layout">
                        <div className="left-content">

                            <BillPillCarousel
                                bills={splitData.bills}
                                activeBillId={activeBillId}
                                onBillSelect={setActiveBillId}
                            />

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

                            <TakersPanel
                                participants={splitData.participants}
                                activeItem={activeItem}
                            />

                            <BillGrid
                                splitData={splitData}
                                onAddTakerClick={handleAddTakerClick}
                                activeBillId={activeBillId}
                            />
                        </div>

                        <div className="right-sidebar">
                            <TallyPanel
                                tally={tally}
                                participants={splitData.participants}
                                splitData={splitData}
                            />
                        </div>
                    </div>
                    {isAddItemPopupOpen && (
                        <ItemPopup
                            onClose={() => setIsAddItemPopupOpen(false)}
                            billId={activeBillId}
                        />
                    )}
                </div>
            )}
        </SocketContextComponent>
    )
}


export default SplitGrid;
