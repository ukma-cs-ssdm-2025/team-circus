import { API_ENDPOINTS } from "../constants";
import type {
	CreateDocumentPayload,
	DeleteResponse,
	DocumentItem,
} from "../types";
import { apiClient } from "./apiClient";

export const createDocument = async (
	payload: CreateDocumentPayload,
): Promise<DocumentItem> => {
	const response = await apiClient.post<DocumentItem, CreateDocumentPayload>(
		API_ENDPOINTS.DOCUMENTS.BASE,
		payload,
	);

	return response.data;
};

export const deleteDocument = async (documentUUID: string): Promise<void> => {
	await apiClient.delete<DeleteResponse>(
		API_ENDPOINTS.DOCUMENTS.DETAIL(documentUUID),
	);
};
