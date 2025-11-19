import { API_ENDPOINTS } from "../constants";
import { apiClient } from "./apiClient";
import type {
	CreateMemberPayload,
	DeleteResponse,
	MemberItem,
	MembersResponse,
	UpdateMemberPayload,
} from "../types";

export const getGroupMembers = async (
	groupUUID: string,
): Promise<MembersResponse> => {
	const response = await apiClient.get<MembersResponse>(
		API_ENDPOINTS.GROUPS.MEMBERS(groupUUID),
	);

	return response.data;
};

export const addGroupMember = async (
	groupUUID: string,
	payload: CreateMemberPayload,
): Promise<MemberItem> => {
	const response = await apiClient.post<MemberItem, CreateMemberPayload>(
		API_ENDPOINTS.GROUPS.MEMBERS(groupUUID),
		payload,
	);

	return response.data;
};

export const updateGroupMemberRole = async (
	groupUUID: string,
	userUUID: string,
	payload: UpdateMemberPayload,
): Promise<MemberItem> => {
	const response = await apiClient.put<MemberItem, UpdateMemberPayload>(
		API_ENDPOINTS.GROUPS.MEMBER(groupUUID, userUUID),
		payload,
	);

	return response.data;
};

export const removeGroupMember = async (
	groupUUID: string,
	userUUID: string,
): Promise<void> => {
	await apiClient.delete<DeleteResponse>(
		API_ENDPOINTS.GROUPS.MEMBER(groupUUID, userUUID),
	);
};
