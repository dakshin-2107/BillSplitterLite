import { toast } from 'sonner';

const useNotify = () => ({
    success: (msg: string) => toast.success(msg),
    error: (msg: string) => toast.error(msg),
    info: (msg: string) => toast.info(msg),
    warning: (msg: string) => toast.warning(msg),
});

export default useNotify;
