import './SplitGrid.css';
import React from 'react';
import type { BillData } from '../../common/interfaces';
import TakerCell from './taker-cell/TakerCell';
import TakersPanel from './people-panel/TakersPanel';
import { SocketContextComponent } from './action-manager/SocketContext';
import { StatusIndicator } from './connection-status/StatusIndicator';
import BillDebugger from '../debug/BillDebugger';


/*
    - Create a grid that is rendered based on the split JSON object 
    - All actions performed are routed through the action mananger
*/

const SplitGrid: React.FC<{ billData: BillData }> = ({ billData: initialBillData }) => {

    const [billData, setBillData] = React.useState<BillData>(initialBillData);
    const [activeItem, setActiveItem] = React.useState<{ id: string, name: string } | null>(null);

    const handleAddTakerClick = (itemId: string, itemName: string) => {
        setActiveItem({ id: itemId, name: itemName });
    };

    return (
        <SocketContextComponent setBillData={setBillData}>
            <div className="split-grid-container">
                <div className="split-grid-header">
                    <h1 className="split-grid-title">Split Grid</h1>
                    <StatusIndicator />
                </div>

                <TakersPanel
                    participants={billData.participants}
                    activeItem={activeItem}
                />

                <div className="table-wrapper">
                    <table className="split-table">
                        <thead>
                            <tr>
                                <th>Number</th>
                                <th>Item Name</th>
                                <th>Price</th>
                                <th>Takers</th>
                            </tr>
                        </thead>
                        <tbody>
                            {Object.values(billData.items).map((item, index) => (
                                <tr key={item.id || index}>
                                    <td>{index + 1}</td>
                                    <td>{item.name}</td>
                                    <td>{item.price}</td>
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
                        {/* <tfoot>
                            <tr>
                                <td colSpan={3}>Total</td>
                                <td>{billData.total}</td>
                            </tr>
                        </tfoot> */}
                    </table>
                </div>

                {/* Debugger to visualize state updates */}
                {/*<BillDebugger billData={billData} />*/}
            </div>
        </SocketContextComponent>
    )
}

export default SplitGrid;
