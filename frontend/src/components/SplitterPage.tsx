import { useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import SplitGrid from './split-grid/SplitGrid';
import type { SplitData } from '@utils/interfaces';
import './components.css';

const SplitterPage = () => {
    const location = useLocation();
    const initialData = (location.state as { splitData?: SplitData } | null)?.splitData ?? null;
    const [splitData, setSplitData] = useState<SplitData | null>(initialData);

    if (!splitData) return <Navigate to="/" replace />;

    return (
        <div className="home-container">
            <SplitGrid splitData={splitData} setSplitData={setSplitData} />
        </div>
    );
};

export default SplitterPage;
