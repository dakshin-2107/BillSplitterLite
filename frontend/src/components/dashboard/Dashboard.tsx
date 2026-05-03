import type { SplitData } from '@utils/interfaces';
import './Dashboard.css';

interface DashboardProps {
    onStartNewBill: () => void;
    activeSession: SplitData | null;
    onResumeSession: (data: SplitData) => void;
    isLoadingSession: boolean;
}

const Dashboard = ({ onStartNewBill, activeSession, onResumeSession, isLoadingSession }: DashboardProps) => {
    const getSessionTitle = (session: SplitData): string => {
        const bills = Object.values(session.bills);
        if (bills.length > 0 && bills[0].location) {
            return bills[0].location;
        }
        return 'Bill Session';
    };

    const getParticipantNames = (session: SplitData): string => {
        const names = Object.values(session.participants);
        if (names.length === 0) return 'No participants';
        if (names.length <= 3) return names.join(', ');
        return `${names.slice(0, 3).join(', ')} +${names.length - 3}`;
    };

    const getBillCount = (session: SplitData): number => {
        return Object.keys(session.bills).length;
    };

    return (
        <div className="dashboard-root">
            {/* Hero Banner */}
            <section className="ds-banner">
                <div className="ds-banner-blob-tr" />
                <div className="ds-banner-blob-bl" />
                <button className="ds-banner-btn" onClick={onStartNewBill}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                        <circle cx="12" cy="12" r="10" />
                        <line x1="12" y1="8" x2="12" y2="16" />
                        <line x1="8" y1="12" x2="16" y2="12" />
                    </svg>
                    Start New Bill
                </button>
            </section>

            {/* Main grid */}
            <div className="ds-grid">
                {/* Active Sessions */}
                <div>
                    <h2 className="ds-section-title">Active Sessions</h2>
                    <p className="ds-section-subtitle">Sessions that can be resumed.</p>

                    {isLoadingSession ? (
                        <div className="ds-skeleton" />
                    ) : activeSession ? (
                        <div className="ds-card">
                            <div className="ds-card-strip" />
                            <div className="ds-card-icon">🧾</div>
                            <h3 className="ds-card-title">{getSessionTitle(activeSession)}</h3>
                            <p className="ds-card-subtitle">
                                {getParticipantNames(activeSession)} · {getBillCount(activeSession)} bill{getBillCount(activeSession) !== 1 ? 's' : ''}
                            </p>
                            <div className="ds-card-footer">
                                <span className="ds-status-chip">{getBillCount(activeSession)} Active Bill{getBillCount(activeSession) !== 1 ? 's' : ''}</span>
                                <button className="ds-resume-btn" onClick={() => onResumeSession(activeSession)}>
                                    Resume
                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                        <line x1="5" y1="12" x2="19" y2="12" />
                                        <polyline points="12 5 19 12 12 19" />
                                    </svg>
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div className="ds-empty-state">
                            <span className="ds-empty-state-icon">📋</span>
                            <span className="ds-empty-state-text">No active sessions. Start a new bill above.</span>
                        </div>
                    )}
                </div>

                {/* Recent Activity */}
                <div className="ds-activity-panel">
                    <h2 className="ds-activity-title">Recent Activity</h2>
                    <div className="ds-activity-empty">
                        <span className="ds-activity-empty-icon">🕐</span>
                        <span className="ds-activity-empty-text">No older transactions detected</span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
