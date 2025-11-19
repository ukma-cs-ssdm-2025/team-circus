import { API_ENDPOINTS } from "../constants";
import { apiClient } from "./apiClient";
import type {
	CreateDocumentPayload,
	DeleteResponse,
	DocumentItem,
} from "../types";

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
