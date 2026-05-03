import { useNavigate } from 'react-router-dom';
import './ErrorPage.css';

interface ErrorPageProps {
    code?: number | string;
    title?: string;
    message?: string;
}

const ErrorPage = ({
    code = 404,
    title = 'Page not found',
    message = "The page you're looking for doesn't exist or something went wrong.",
}: ErrorPageProps) => {
    const navigate = useNavigate();

    return (
        <div className="error-page">
            <div className="error-page__card">
                <span className="error-page__code">{code}</span>
                <div className="error-page__divider" />
                <h1 className="error-page__title">{title}</h1>
                <p className="error-page__message">{message}</p>
                <button className="error-page__btn" onClick={() => navigate('/')}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                        <line x1="19" y1="12" x2="5" y2="12" />
                        <polyline points="12 19 5 12 12 5" />
                    </svg>
                    Back to Dashboard
                </button>
            </div>
        </div>
    );
};

export default ErrorPage;
