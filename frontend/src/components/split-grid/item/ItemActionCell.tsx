import { useContext, useState } from 'react';
import { Trash2, Pencil } from 'lucide-react';
import { SocketContextProvider } from '../action-manager/SocketContext';
import type { Item } from '../../../common/interfaces';
import { ActionType } from '../../../common/interfaces';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import ItemPopup from './ItemPopup';
import './ItemActionCell.css';

interface ItemActionCellProps {
    item: Item;
    billId: number;
}

const ItemActionCell = ({ item, billId }: ItemActionCellProps) => {
    const socketContext = useContext(SocketContextProvider);
    const [editOpen, setEditOpen] = useState(false);

    const handleDelete = () => {
        if (socketContext) {
            socketContext.publishAction({
                actionType: ActionType.DELETE_ITEM,
                billId: billId,
                itemId: item.id,
            });
        }
    };

    return (
        <div className="item-action-cell">
            <button className="action-btn edit-btn" onClick={() => setEditOpen(true)} title="Edit Item">
                <Pencil size={14} />
            </button>

            {editOpen && (
                <ItemPopup
                    billId={billId}
                    item={item}
                    onClose={() => setEditOpen(false)}
                />
            )}

            <AlertDialog>
                <AlertDialogTrigger asChild>
                    <button className="action-btn delete-btn" title="Delete Item">
                        <Trash2 size={14} />
                    </button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete item?</AlertDialogTitle>
                        <AlertDialogDescription>
                            "{item.name}" will be removed from this bill. This cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction variant="destructive" onClick={handleDelete}>Delete</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
};

export default ItemActionCell;
