import { get, set, del } from 'idb-keyval';
import type { SplitFormData } from './interfaces';

const STORAGE_KEY = 'split-form-draft';

export const storageUtils = {
    async saveDraft(data: SplitFormData): Promise<void> {
        await set(STORAGE_KEY, data);
    },

    async getDraft(): Promise<SplitFormData | undefined> {
        return await get<SplitFormData>(STORAGE_KEY);
    },

    async clearDraft(): Promise<void> {
        await del(STORAGE_KEY);
    }
};
