// Entity types for the application
import type { MEMBER_ROLES } from "../constants";

export type MemberRole = (typeof MEMBER_ROLES)[number];

export interface GroupItem {
	uuid: string;
	name: string;
	created_at: string;
	role?: MemberRole;
}

export interface DocumentItem {
	uuid: string;
	name: string;
	content: string;
	group_uuid: string;
	created_at: string;
}

export interface CreateDocumentPayload {
	group_uuid: string;
	name: string;
	content: string;
}

export interface ShareLinkResponse {
	document_uuid: string;
	url: string;
	expires_at: string;
}

export interface GroupsResponse {
	groups: GroupItem[];
}

export interface DocumentsResponse {
	documents: DocumentItem[];
}

export interface CreateGroupPayload {
	name: string;
}

export interface MemberItem {
	group_uuid: string;
	user_uuid: string;
	role: MemberRole;
	created_at: string;
}

export interface MembersResponse {
	members: MemberItem[];
}

export interface CreateMemberPayload {
	user_uuid: string;
	role: MemberRole;
}

export interface UpdateMemberPayload {
	role: MemberRole;
}

export interface DeleteResponse {
	message: string;
}

export interface UserItem {
	uuid: string;
	login: string;
	email: string;
	created_at: string;
}

export interface UsersResponse {
	users: UserItem[];
}

// Filter and search types
export interface GroupOption {
	value: string;
	label: string;
}

export interface DocumentFilters {
	selectedGroup: string;
	searchTerm: string;
}
