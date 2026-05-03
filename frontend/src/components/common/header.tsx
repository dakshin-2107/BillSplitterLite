import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './header.css';

const Header: React.FC = () => {
    const siteTitle = import.meta.env.VITE_SITE_TITLE || 'Splitzy';
    const navigate = useNavigate();
    const location = useLocation();

    const isDashboard = location.pathname === '/';
    const isSplitter = location.pathname.startsWith('/splitter');

    return (
        <header className="app-header">
            <div className="header-content">
                <button className="header-logo" onClick={() => navigate('/')} aria-label="Go to dashboard">
                    {siteTitle}
                </button>
                <nav className="header-nav">
                    <button
                        className={`nav-link${isDashboard ? ' nav-link--active' : ''}`}
                        onClick={() => navigate('/')}
                    >
                        Dashboard
                    </button>
                    <button
                        className={`nav-link${isSplitter ? ' nav-link--active' : ''}`}
                        onClick={() => navigate('/splitter')}
                    >
                        Splitter
                    </button>
                </nav>
            </div>
        </header>
    );
};

export default Header;
