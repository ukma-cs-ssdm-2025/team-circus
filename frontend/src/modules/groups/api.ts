import { apiClient } from '../../services/apiClient';
import { API_ENDPOINTS } from '../../constants';
import type { GroupMember, GroupRole } from '../../types';

interface GroupMembersResponse {
  members: GroupMember[];
}

interface AddMemberPayload {
  user_uuid: string;
  role: GroupRole;
}

interface UpdateMemberPayload {
  role: GroupRole;
}

const buildMembersEndpoint = (groupUUID: string, userUUID?: string): string => {
  const base = `${API_ENDPOINTS.GROUPS.BASE}/${groupUUID}/members`;
  if (userUUID) {
    return `${base}/${userUUID}`;
  }
  return base;
};

export const fetchGroupMembers = async (groupUUID: string): Promise<GroupMember[]> => {
  const response = await apiClient.get<GroupMembersResponse>(buildMembersEndpoint(groupUUID));
  return response.data.members;
};

export const addGroupMember = async (
  groupUUID: string,
  userUUID: string,
  role: GroupRole,
): Promise<GroupMember> => {
  const payload: AddMemberPayload = {
    user_uuid: userUUID,
    role,
  };
  const response = await apiClient.post<GroupMember>(buildMembersEndpoint(groupUUID), payload);
  return response.data;
};

export const updateGroupMemberRole = async (
  groupUUID: string,
  userUUID: string,
  role: GroupRole,
): Promise<GroupMember> => {
  const payload: UpdateMemberPayload = {
    role,
  };
  const response = await apiClient.put<GroupMember>(buildMembersEndpoint(groupUUID, userUUID), payload);
  return response.data;
};

export const removeGroupMember = async (groupUUID: string, userUUID: string): Promise<void> => {
  await apiClient.delete<null>(buildMembersEndpoint(groupUUID, userUUID));
};
