import React, { useState, useEffect } from 'react';
import './components.css';
import SplitGrid from './split-grid/SplitGrid';
import { type SplitData, type ApiResponse } from '@utils/interfaces';
import { URLProvider } from '@utils/urlProvider';

const Example: React.FC = () => {
    const [splitData, setSplitData] = useState<SplitData | null>(null);
    const [notFound, setNotFound] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const load = async () => {
            try {
                const res = await fetch(URLProvider.getExampleUrl(), {
                    method: 'GET',
                    credentials: 'include',
                });

                if (res.status === 404) {
                    setNotFound(true);
                    return;
                }

                const result: ApiResponse = await res.json();
                if (result.success && result.split) {
                    setSplitData(result.split);
                } else {
                    setNotFound(true);
                }
            } catch {
                setNotFound(true);
            } finally {
                setIsLoading(false);
            }
        };

        load();
    }, []);

    if (isLoading) {
        return (
            <div className="home-container">
                <div className="loading-overlay">
                    <div className="loading-spinner">Loading...</div>
                </div>
            </div>
        );
    }

    if (notFound) {
        return (
            <div className="home-container">
                <p>Route not found.</p>
            </div>
        );
    }

    if (splitData) {
        return (
            <div className="home-container">
                <SplitGrid splitData={splitData} setSplitData={setSplitData} />
            </div>
        );
    }

    return null;
};

export default Example;
