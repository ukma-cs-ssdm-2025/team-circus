// Entity types for the application

export type GroupRole = 'author' | 'coauthor' | 'reviewer';

export interface GroupItem {
  uuid: string;
  name: string;
  created_at: string;
  author_uuid: string;
  role?: GroupRole;
}

export interface GroupMember {
  group_uuid: string;
  user_uuid: string;
  role: GroupRole;
  created_at: string;
  user_login: string;
  user_email: string;
}

export interface DocumentItem {
  uuid: string;
  name: string;
  content: string;
  group_uuid: string;
  created_at: string;
}

export interface GroupsResponse {
  groups: GroupItem[];
}

export interface UsersResponse {
  users: Array<{
    uuid: string;
    login: string;
    email: string;
    created_at: string;
  }>;
}

export interface DocumentsResponse {
  documents: DocumentItem[];
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

// Mutation payload types
export interface CreateGroupPayload {
  name: string;
}

export type UpdateGroupPayload = CreateGroupPayload;

export interface CreateDocumentPayload {
  group_uuid: string;
  name: string;
  content: string;
}
