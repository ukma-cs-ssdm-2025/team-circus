import { API_ENDPOINTS } from "../constants";
import type {
	CreateDocumentPayload,
	DeleteResponse,
	DocumentItem,
	ShareLinkResponse,
} from "../types";
import { getApiUrl } from "../config/env";
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

export const generateShareLink = async (
	documentUUID: string,
	expirationDays: number,
): Promise<ShareLinkResponse> => {
	const response = await apiClient.post<ShareLinkResponse, { expiration_days: number }>(
		API_ENDPOINTS.DOCUMENTS.SHARE(documentUUID),
		{ expiration_days: expirationDays },
	);

	return response.data;
};

type PublicDocumentParams = {
	doc: string;
	sig: string;
	exp: string;
};

type PublicDocumentError = Error & { status?: number };

export const fetchPublicDocument = async (
	params: PublicDocumentParams,
): Promise<DocumentItem> => {
	const query = new URLSearchParams({
		doc: params.doc,
		sig: params.sig,
		exp: params.exp,
	});

	const response = await fetch(
		getApiUrl(`${API_ENDPOINTS.DOCUMENTS.PUBLIC}?${query.toString()}`),
		{
			method: "GET",
			credentials: "include",
		},
	);

	let parsed: unknown = null;
	try {
		parsed = await response.json();
	} catch {
		parsed = null;
	}

	if (!response.ok) {
		const errorMessage =
			(parsed &&
				typeof parsed === "object" &&
				"error" in (parsed as Record<string, unknown>) &&
				typeof (parsed as Record<string, unknown>).error === "string" &&
				(parsed as Record<string, unknown>).error) ||
			"Failed to load document";

		const error = new Error(errorMessage) as PublicDocumentError;
		error.status = response.status;
		throw error;
	}

	if (parsed && typeof parsed === "object") {
		return parsed as DocumentItem;
	}

	throw new Error("Invalid response");
};
