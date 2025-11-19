import { API_ENDPOINTS } from "../constants";
import { apiClient } from "./apiClient";
import type {
	CreateGroupPayload,
	DeleteResponse,
	GroupItem,
	GroupsResponse,
} from "../types";

export const createGroup = async (
	payload: CreateGroupPayload,
): Promise<GroupItem> => {
	const response = await apiClient.post<GroupItem, CreateGroupPayload>(
		API_ENDPOINTS.GROUPS.BASE,
		payload,
	);

	return response.data;
};

export const deleteGroup = async (groupUUID: string): Promise<void> => {
	await apiClient.delete<DeleteResponse>(
		API_ENDPOINTS.GROUPS.DETAIL(groupUUID),
	);
};

export const fetchGroups = async (): Promise<GroupsResponse> => {
	const response = await apiClient.get<GroupsResponse>(
		API_ENDPOINTS.GROUPS.BASE,
	);

	return response.data;
};
