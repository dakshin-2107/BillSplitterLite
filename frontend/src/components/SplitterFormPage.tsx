import { useNavigate } from 'react-router-dom';
import SplitForm from '@split-form/SplitForm';
import type { SplitData } from '@utils/interfaces';
import useNotify from '../hooks/useNotify';
import './components.css';

const SplitterFormPage = () => {
    const navigate = useNavigate();
    const notify = useNotify();

    const handleSuccess = (data: SplitData) => {
        notify.success('Session started successfully');
        navigate(`/splitter/${data.splitId}`, { state: { splitData: data } });
    };

    return (
        <div className="home-container">
            <SplitForm
                onSubmitSuccess={handleSuccess}
                onSubmitError={(msg) => notify.error(msg)}
            />
        </div>
    );
};

export default SplitterFormPage;
