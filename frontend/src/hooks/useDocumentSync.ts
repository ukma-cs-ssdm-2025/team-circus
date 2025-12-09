import { useCallback, useMemo } from "react";
import { API_ENDPOINTS } from "../constants";
import type { DocumentItem } from "../types";
import { useApi, useMutation } from "./useApi";

type SavePayload = {
	name: string;
	content: string;
};

type DocumentSyncReturn = {
	document: DocumentItem | null;
	loading: boolean;
	error: Error | null;
	saveDocument: (payload: SavePayload) => Promise<DocumentItem>;
	refetch: () => Promise<void>;
};

export function useDocumentSync(documentId: string): DocumentSyncReturn {
	const endpoint = useMemo(() => {
		if (!documentId) {
			return "";
		}

		if (typeof API_ENDPOINTS.DOCUMENTS.DETAIL === "function") {
			return API_ENDPOINTS.DOCUMENTS.DETAIL(documentId);
		}

		return `${API_ENDPOINTS.DOCUMENTS.BASE}/${documentId}`;
	}, [documentId]);

	const {
		data,
		loading,
		error,
		refetch: refetchDocument,
		mutate: updateCache,
	} = useApi<DocumentItem>(endpoint, {
		immediate: Boolean(documentId),
	});

	const { mutate: saveRequest } = useMutation<DocumentItem, SavePayload>(
		endpoint,
		"PUT",
	);

	const saveDocument = useCallback(
		async (payload: SavePayload) => {
			if (!endpoint) {
				throw new Error("Document endpoint is not defined.");
			}

			const saved = await saveRequest(payload);
			updateCache(saved);
			return saved;
		},
		[endpoint, saveRequest, updateCache],
	);

	const normalizedError = useMemo(() => {
		if (!error) {
			return null;
		}

		return new Error(error.message);
	}, [error]);

	const refetch = useCallback(async () => {
		if (!documentId) {
			return;
		}

		await refetchDocument();
	}, [documentId, refetchDocument]);

	return {
		document: data,
		loading,
		error: normalizedError,
		saveDocument,
		refetch,
	};
}
