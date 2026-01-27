import React from 'react';
import './header.css';

const Header: React.FC = () => {
    const siteTitle = import.meta.env.VITE_SITE_TITLE || 'Bill Splitter';

    return (
        <header className="app-header">
            <div className="header-content">
                <h1 className="header-title">{siteTitle}</h1>
                <button className="settings-button" aria-label="Settings">
                    Settings
                </button>
            </div>
        </header>
    );
};

export default Header;
