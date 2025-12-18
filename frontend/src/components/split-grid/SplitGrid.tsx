import React, { useState, useEffect, useMemo, createContext } from 'react';
import './SplitGrid.css';
import { ActionManager } from './ActionManager';


/*
    - Create a grid that is rendered based on the split JSON object 
    - All actions performed are routed through the action mananger

*/
const SplitGrid: React.FC = () => {

    const actionManager = useMemo(() => new ActionManager(), []);
    const ActionManagerContext = createContext<ActionManager | null>(null);

    return (
        <ActionManagerContext.Provider value={actionManager}>
            <div>
                <h1>Split Grid</h1>
            </div>
        </ActionManagerContext.Provider>
    )
}

export default SplitGrid;
