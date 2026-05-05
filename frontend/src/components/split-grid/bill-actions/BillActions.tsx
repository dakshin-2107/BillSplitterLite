import { useContext, useState } from 'react';
import { Share2, Loader2, Check, AlertCircle, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
    AlertDialog,
    AlertDialogTrigger,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogAction,
    AlertDialogCancel,
} from '@/components/ui/alert-dialog';
import { SocketContextProvider } from '../action-manager/SocketContext';
import { ActionType } from '../../../common/interfaces';
import { URLProvider } from '../../../common/urlProvider';
import useNotify from '../../../hooks/useNotify';
import './BillActions.css';

interface ShareBillApiResponse {
    success: boolean;
    message: string;
    splitId: string;
}

const STATUS_RESET_DELAY_MS = 3000;

const SHARE_ICON: Record<string, React.ReactNode> = {
    idle:    <Share2 size={16} />,
    copying: <Loader2 size={16} className="spin" />,
    copied:  <Check size={16} />,
    error:   <AlertCircle size={16} />,
};

const SHARE_TITLE: Record<string, string> = {
    idle:    'Share Bill',
    copying: 'Generating link...',
    copied:  'Link copied!',
    error:   'Retry share',
};

const BillActions = () => {
    const notify = useNotify();
    const socketContext = useContext(SocketContextProvider);
    const [shareStatus, setShareStatus] = useState<'idle' | 'copying' | 'copied' | 'error'>('idle');

    const handleShareBill = async () => {
        if (shareStatus === 'copying' || shareStatus === 'copied') return;

        setShareStatus('copying');
        try {
            const generateUrl = URLProvider.getGenerateUrl();
            const response = await fetch(generateUrl, {
                method: 'GET',
                credentials: 'include',
            });

            if (!response.ok) throw new Error('Failed to generate session ID');

            const result: ShareBillApiResponse = await response.json();
            if (result.success && result.splitId) {
                const joinLink = `${URLProvider.getSiteOrigin()}/join/${result.splitId}`;
                await navigator.clipboard.writeText(joinLink);
                setShareStatus('copied');
                setTimeout(() => setShareStatus('idle'), STATUS_RESET_DELAY_MS);
            } else {
                throw new Error(result.message || 'Failed to generate session ID');
            }
        } catch {
            setShareStatus('error');
            notify.error('Failed to copy link');
            setTimeout(() => setShareStatus('idle'), STATUS_RESET_DELAY_MS);
        }
    };

    const handleRefresh = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.SYNC_BILL_STATE,
                itemId: 0,
            });
        }
    };

    const handleConfirmClose = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.BYE_BYE,
                itemId: 0,
            });
        }
    };

    return (
        <div className="right-controls">

            <Button
                className="control-btn refresh-btn"
                onClick={handleRefresh}
                variant="outline"
                title="Sync with server"
            >
                <RefreshCw size={16} />
            </Button>

            <Button
                className={`control-btn share-btn ${shareStatus}`}
                onClick={handleShareBill}
                disabled={shareStatus === 'copying'}
                variant="outline"
                title={SHARE_TITLE[shareStatus]}
            >
                {SHARE_ICON[shareStatus]}
            </Button>

            <AlertDialog>
                <AlertDialogTrigger asChild>
                    <Button className="control-btn close-btn" variant="outline">
                        Close Bill
                    </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Close Bill</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to close this bill session? All data will be deleted.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleConfirmClose}>Confirm</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
};

export default BillActions;
