// Entity types for the application

export interface GroupItem {
	uuid: string;
	name: string;
	created_at: string;
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
	role: "author" | "editor" | "viewer";
	created_at: string;
}

export interface MembersResponse {
	members: MemberItem[];
}

export interface CreateMemberPayload {
	user_uuid: string;
	role: "author" | "editor" | "viewer";
}

export interface UpdateMemberPayload {
	role: "author" | "editor" | "viewer";
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
